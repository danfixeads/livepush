package senders

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gopkg.in/guregu/null.v3"

	"github.com/danfixeads/livepush/models"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/vjeantet/jodaTime"
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
		go iosworker(db, client, pushes)
	}

	for _, token := range i.Push.Tokens {

		iospush := IOSPush{}

		iospush.notification = &apns2.Notification{}
		iospush.notification.DeviceToken = token.String
		iospush.notification.Topic = i.client.BundleIdentifier.String

		pLoad, _ := json.Marshal(i.Push.Payload)
		iospush.notification.Payload = []byte(pLoad)

		iospush.push = models.Push{}
		iospush.push.ClientID = i.Push.ClientID
		iospush.push.Token = token
		iospush.push.Platform = null.String{NullString: sql.NullString{
			String: "ios",
			Valid:  true,
		}}
		iospush.push.Payload = null.String{NullString: sql.NullString{
			String: string(pLoad),
			Valid:  true,
		}}

		// add to the worker array
		pushes <- &iospush
	}

	close(pushes)

	return nil
}

func iosworker(db *sql.DB, client *apns2.Client, pushes <-chan *IOSPush) {

	for p := range pushes {
		res, err := client.Push(p.notification)
		if err != nil {
			log.Fatal("Push Error: ", err)
		}
		p.push.Response = null.String{NullString: sql.NullString{
			String: fmt.Sprintf("%d %s", res.StatusCode, res.Reason),
			Valid:  true,
		}}

		// add the sent datetime if statuscode is 200
		if res.StatusCode == 200 {
			p.push.Sent = null.String{NullString: sql.NullString{
				String: jodaTime.Format("YYYY-MM-dd HH:mm:ss", time.Now()),
				Valid:  true,
			}}
		}

		p.push.Create(db)

		fmt.Printf("DeviceToken: %v StatusCode: %v ApnsID: %v Reason: %v\n", p.notification.DeviceToken, res.StatusCode, res.ApnsID, res.Reason)
	}

}
