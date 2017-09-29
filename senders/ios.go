package senders

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sync"
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

// setup the waitgroup variable
var wg sync.WaitGroup

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

	// determine the basepath value
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(b)
	)

	//fmt.Println(basepath)

	if i.client.PemFile.Valid {
		i.cert, err = certificate.FromPemFile(fmt.Sprintf("%s/files/%s", basepath, i.client.PemFile.String), i.client.PassPhrase.String)
	}
	if i.cert.Certificate == nil {
		if i.client.P12File.Valid {
			i.cert, err = certificate.FromP12File(fmt.Sprintf("%s/files/%s", basepath, i.client.P12File.String), i.client.PassPhrase.String)
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
	// make the worker arrays
	pushes := make(chan *IOSPush, totalPushes)
	responses := make(chan *IOSPush, totalPushes)

	// update the waitgroup with the amount of goroutines
	wg.Add(totalPushes)

	client := apns2.NewClient(i.cert)

	// determine which APNS service to use
	if i.client.UseSandboxIOS.Bool {
		client.Development() // the Sandbox
	} else {
		client.Production() // or the Production
	}

	for i := 0; i < totalPushes; i++ {
		go iosworker(db, client, pushes, responses)
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

	/*
		for i := 0; i < totalPushes; i++ {
			res := <-responses
			fmt.Printf("res: %v\n", res)
		}*/

	err := checkAndCallback(totalPushes, responses)

	close(pushes)
	close(responses)
	wg.Wait()

	return err
}

func checkAndCallback(total int, responses <-chan *IOSPush) error {

	failed := make([]IOSPush, 0)
	for i := 0; i < total; i++ {
		res := <-responses
		//fmt.Printf("response: %v\n", res)
		if !res.push.Sent.Valid {
			failed = append(failed, *res)
			//fmt.Printf("failed: %v\n", failed)
		}
	}

	// if all of the pushes failed (even if it's just one push sent)
	// then return an error message
	// NOTE: If there were various pushes sent successfully, and just a few that failed
	// then it's not an "error"
	var err error
	if len(failed) == total {
		err = models.ErrFailedToSendPush
	}

	return err
}

func iosworker(db *sql.DB, client *apns2.Client, pushes <-chan *IOSPush, responses chan<- *IOSPush) {

	for p := range pushes {
		// push the message
		res, err := client.Push(p.notification)
		if err != nil {
			log.Fatal("Push Error: ", err)
		}
		// check and save the response
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

		// save the push record
		p.push.Create(db)

		responses <- p

		// this waitgroup routine is "done"
		wg.Done()

		//fmt.Printf("DeviceToken: %v StatusCode: %v ApnsID: %v Reason: %v\n", p.notification.DeviceToken, res.StatusCode, res.ApnsID, res.Reason)
	}

}
