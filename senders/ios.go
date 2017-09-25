package senders

import (
	"crypto/tls"
	"database/sql"
	"fmt"

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

	for _, token := range i.Push.Tokens {

		notification := &apns2.Notification{}
		notification.DeviceToken = token.String
		notification.Topic = i.client.BundleIdentifier.String
		//notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`) // See Payload section below

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

		client := apns2.NewClient(i.cert).Development()
		res, err := client.Push(notification)

		if err != nil {
			//log.Fatal("Error:", err)
			return err
		}

		fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)

	}

	return nil
}

/*
func main() {

	cert, err := certificate.FromP12File("../cert.p12", "")
	if err != nil {
		log.Fatal("Cert Error:", err)
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7"
	notification.Topic = "com.sideshow.Apns2"
	notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`) // See Payload section below

	client := apns2.NewClient(cert).Production()
	res, err := client.Push(notification)

	if err != nil {
		log.Fatal("Error:", err)
	}

	fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
}

*/
