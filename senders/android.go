package senders

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"

	fcm "github.com/NaySoftware/go-fcm"
	"github.com/danfixeads/livepush/models"
)

// Android Struct
type Android struct {
	ClientID int
	Push     models.MultiplePush
	client   models.Client
}

// AndroidPush struct
type AndroidPush struct {
	push         models.Push
	notification fcm.NotificationPayload
}

type data struct {
	Message  message `json:"message"`
	Title    string  `json:"title"`
	Subtitle string  `json:"subtitle"`
	Body     string  `json:"body"`
	Badge    int     `json:"badge"`
	Image    string  `json:"image"`
}

type message struct {
	Data     dataType `json:"data"`
	Alert    string   `json:"alert"`
	Sound    int      `json:"sound"`
	ImageURL string   `json:"image_url"`
}

type dataType struct {
	Type    int     `json:"type"`
	Actions actions `json:"actions"`
}

type actions struct {
	Main string `json:"main"`
}

// GetClient function
func (a *Android) GetClient(db *sql.DB) error {

	if err := a.client.GetByClientID(db, a.ClientID); err != nil {
		return err
	}

	if !a.client.FCMAuthKey.Valid {
		return models.ErrMissingVitalFields
	}

	return nil
}

// SendMessage function
func (a *Android) SendMessage(db *sql.DB) error {

	c := fcm.NewFcmClient(a.client.FCMAuthKey.String)
	//c.SetDryRun(true)

	data := data{
		Message: message{
			Data: dataType{
				Type: 0,
				Actions: actions{
					Main: "/",
				},
			},
			Alert:    "Hello Green world!",
			Sound:    1,
			ImageURL: "https://imovirtualpt-images.akamaized.net///images_imovirtualpt///5343829_1_655x491_moradia-na-praia-barbecue-jardim-e-bricolage-vila-nova-de-gaia.jpg",
		},
		Title:    "",
		Subtitle: "",
		Body:     "",
		Badge:    0,
		Image:    "",
	}
	jsonByte, _ := json.Marshal(data)
	// fmt.Print(jsonByte)
	fmt.Print(bytes.NewBuffer(jsonByte))

	for _, token := range a.Push.Tokens {

		ids := []string{
			token.String,
		}

		c.NewFcmRegIdsMsg(ids, data)

		status, err := c.Send()

		if err == nil {
			status.PrintResults()
		} else {
			fmt.Println(err)
		}

	}

	return nil
}
