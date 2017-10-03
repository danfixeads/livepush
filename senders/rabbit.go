package senders

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/danfixeads/livepush/models"
	"github.com/streadway/amqp"
	"github.com/vjeantet/jodaTime"
)

// Rabbit struct
type Rabbit struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

type payload struct {
	Date       string `json:"date"`
	FormatDate string `json:"format_date"`
	Status     int    `json:"status"`
	Service    string `json:"service"`
	Operation  string `json:"operation"`
	IPClient   string `json:"ip_client"`
	IPServer   string `json:"ip_server"`
	ClientID   string `json:"clientid"`
	Data       string `json:"data"`
}

// SetUp function
func (r *Rabbit) SetUp() error {

	config := models.ReturnConfig()
	if len(config.MQHost) == 0 || strings.Contains(os.Args[0], "/_test/") {
		return nil
	}

	var err error

	conn, errConnect := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.MQUser, config.MQPass, config.MQHost, config.MQPort))
	err = errConnect

	if err == nil {
		ch, errChannel := conn.Channel()
		err = errChannel
		r.channel = ch

		errExchange := ch.ExchangeDeclare(
			"livepush", // name
			"topic",    // type
			true,       // durable
			false,      // auto-deleted
			false,      // internal
			false,      // no-wait
			nil,        // arguments
		)
		err = errExchange

		if err == nil {
			q, errQueue := ch.QueueDeclare(
				"livepush", // name
				true,       // durable
				false,      // delete when unused
				false,      // exclusive
				false,      // no-wait
				nil,        // arguments
			)
			r.queue = q
			err = errQueue
		}

	}

	return err
}

// Send function
func (r *Rabbit) Send(clientID string, operation string, code int, request *http.Request, data string, attempt int) error {

	if r.channel != nil {

		body, err := json.Marshal(returnPayload(clientID, operation, code, request, data))
		if err != nil {
			return err
		}

		err = r.channel.Publish(
			"livepush",   // exchange
			r.queue.Name, // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(string(body)),
			})

		fmt.Printf("Sending msg: %v - Error: %v", string(body), err)

		return err

	} else if r.channel == nil && attempt < 5 {
		err := r.SetUp()
		if err == nil {
			r.Send(clientID, operation, code, request, data, attempt+1)
		}
	}
	return nil
}

func returnPayload(clientID, operation string, code int, r *http.Request, data string) payload {
	p := payload{}
	p.Date = jodaTime.Format("YYYY-MM-dd HH:mm:ss", time.Now())
	p.FormatDate = "UTC"
	p.Service = "livepush"
	p.Status = code
	p.Operation = operation
	p.IPClient = r.Referer()
	p.IPServer = getServerIP()
	p.ClientID = clientID
	p.Data = data
	return p
}

func getServerIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
