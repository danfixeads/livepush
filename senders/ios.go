package senders

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gopkg.in/guregu/null.v3"

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
	notification *apns2.Notification
	response     apns2.Response
}

// GetClient function
func (i *IOS) GetClient(db *sql.DB) error {

	if err := i.client.GetByClientID(db, i.ClientID); err != nil {
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
func (i *IOS) SendMessage(db *sql.DB) error {

	totalPushes := len(i.Push.Tokens)
	// make the worker array
	pushes := make(chan *IOSPush, totalPushes)

	client := apns2.NewClient(i.cert)

	// determine which APNS service to use
	if i.client.UseSandboxIOS.Bool {
		client.Development() // the Sandbox
	} else {
		client.Production() // or the Production
	}

	for i := 0; i < totalPushes; i++ {
		go worker(db, client, pushes)
	}

	for _, token := range i.Push.Tokens {

		iospush := IOSPush{}

		iospush.notification = &apns2.Notification{}
		iospush.notification.DeviceToken = token.String
		iospush.notification.Topic = i.client.BundleIdentifier.String

		pLoad := payload.NewPayload()
		pLoad.AlertTitle(i.Push.Title.String)
		pLoad.AlertSubtitle(i.Push.Subtitle.String)
		pLoad.AlertBody(i.Push.Body.String)
		pLoad.Badge(int(i.Push.Badge.Int64))
		pLoad.Sound(i.Push.Sound.String)
		//pLoad.AlertLaunchImage(i.Push.Image.String)
		pLoad.ContentAvailable()
		pLoad.MutableContent()

		// add image if applicable
		if i.Push.Image.Valid {
			image := map[string]string{
				"attachment-url": i.Push.Image.String,
			}
			pLoad.Custom("data", image)
		}

		iospush.notification.Payload = pLoad

		iospush.push = models.Push{}
		iospush.push.ClientID = i.Push.ClientID
		iospush.push.Token = token
		iospush.push.Platform = null.String{NullString: sql.NullString{
			String: "ios",
			Valid:  true,
		}}
		iospush.push.Title = i.Push.Title
		iospush.push.Subtitle = i.Push.Subtitle
		iospush.push.Body = i.Push.Body
		iospush.push.Badge = i.Push.Badge
		iospush.push.Image = i.Push.Image
		iospush.push.Sound = i.Push.Sound

		// add to the worker array
		pushes <- &iospush
	}

	close(pushes)

	return nil
}

func worker(db *sql.DB, client *apns2.Client, pushes <-chan *IOSPush) {

	for p := range pushes {
		res, err := client.Push(p.notification)
		if err != nil {
			log.Fatal("Push Error: ", err)
		}
		p.push.Response = null.String{NullString: sql.NullString{
			String: fmt.Sprintf("%d %s", res.StatusCode, res.Reason),
			Valid:  true,
		}}

		if res.StatusCode == 200 {

			t := time.Now()

			p.push.Sent = null.String{NullString: sql.NullString{
				String: fmt.Sprintf("%v", t),
				Valid:  true,
			}}
		}

		p.push.Create(db)

		fmt.Printf("DeviceToken: %v StatusCode: %v ApnsID: %v Reason: %v\n", p.notification.DeviceToken, res.StatusCode, res.ApnsID, res.Reason)
	}

}
