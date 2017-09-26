package senders

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"

	"github.com/danfixeads/livepush/models"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

// IOS Struct
type IOS struct {
	ClientID int
	Push     models.MultiplePush
	client   models.Client
	cert     tls.Certificate
}

// IOSPush struct
type IOSPush struct {
	push         models.Push
	notification apns2.Notification
	response     apns2.Response
}

// GetClient function
func (i *IOS) GetClient(db *sql.DB) error {

	if err := i.client.Get(db, i.ClientID); err != nil {
		return err
	}

	if !i.client.BundleIdentifier.Valid || !i.client.PassPhrase.Valid || (!i.client.PemFile.Valid && !i.client.P12File.Valid) {
		return models.ErrMissingVitalFields
	}

	return nil
}

// GetCertificate function
func (i *IOS) GetCertificate() error {

	var err error

	if i.client.PemFile.Valid {
		i.cert, err = certificate.FromPemFile(fmt.Sprintf("./senders/files/%s", i.client.PemFile.String), i.client.PassPhrase.String)
	}
	if i.cert.Certificate == nil {
		if i.client.P12File.Valid {
			i.cert, err = certificate.FromP12File(fmt.Sprintf("./senders/files/%s", i.client.P12File.String), i.client.PassPhrase.String)
		}
	}

	if i.cert.Certificate == nil && err != nil {
		err = models.ErrFailedToLoadPEMFile
	}

	return err
}

// SendMessage function
func (i *IOS) SendMessage() error {

	total := len(i.Push.Tokens)

	notifications := make(chan *apns2.Notification, total)
	responses := make(chan *apns2.Response, total)

	client := apns2.NewClient(i.cert).Development()

	for i := 0; i < total; i++ {
		go worker(client, notifications, responses)
	}

	for _, token := range i.Push.Tokens {

		notification := &apns2.Notification{}
		notification.DeviceToken = token.String
		notification.Topic = i.client.BundleIdentifier.String

		var p = payload.NewPayload()
		p.AlertTitle(i.Push.Title.String)
		p.AlertSubtitle(i.Push.Subtitle.String)
		p.AlertBody(i.Push.Body.String)
		p.Badge(int(i.Push.Badge.Int64))
		p.AlertLaunchImage(i.Push.Image.String)
		p.ContentAvailable()
		p.MutableContent()

		image := map[string]string{
			"attachment-url": i.Push.Image.String,
		}
		p.Custom("data", image)

		notification.Payload = p

		notifications <- notification

		/*
			res, err := client.Push(notification)

			if err != nil {
				//log.Fatal("Error:", err)
				return err
			}

			fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
		*/
	}

	/*
		for i := 0; i < total; i++ {
			res := <-responses
			fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
		}*/

	close(notifications)
	//close(responses)

	return nil
}

func worker(client *apns2.Client, notifications <-chan *apns2.Notification, responses chan<- *apns2.Response) {
	for n := range notifications {
		res, err := client.Push(n)
		if err != nil {
			log.Fatal("Push Error: ", err)
		}
		fmt.Printf("DeviceToken: %v StatusCode: %v ApnsID: %v Reason: %v\n", n.DeviceToken, res.StatusCode, res.ApnsID, res.Reason)
		responses <- res
	}
}
