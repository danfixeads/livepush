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
	Client   models.Client
	cert     tls.Certificate
}

// IOSPush struct
type IOSPush struct {
	Push         models.Push `json:"push"`
	notification *apns2.Notification
}

// setup the waitgroup variable
var wgIOS sync.WaitGroup

// GetClient function
func (i *IOS) GetClient(db *sql.DB) error {

	if err := i.Client.GetByClientID(db, i.ClientID); err != nil {
		return err
	}

	if !i.Client.BundleIdentifier.Valid || !i.Client.PassPhrase.Valid || (!i.Client.PemFile.Valid && !i.Client.P12File.Valid) {
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

	// check and load the PEM file
	if i.Client.PemFile.Valid {
		i.cert, err = certificate.FromPemFile(fmt.Sprintf("%s/files/%s", basepath, i.Client.PemFile.String), i.Client.PassPhrase.String)
	}
	// if no PEM file, then try to load the P12 file
	if i.cert.Certificate == nil {
		if i.Client.P12File.Valid {
			i.cert, err = certificate.FromP12File(fmt.Sprintf("%s/files/%s", basepath, i.Client.P12File.String), i.Client.PassPhrase.String)
		}
	}

	// no files?  Then return an error
	if i.cert.Certificate == nil && err != nil {
		err = models.ErrFailedToLoadPEMFile
	}

	return err
}

// SendMessage function
func (i *IOS) SendMessage(db *sql.DB) ([]models.Push, error) {

	totalPushes := len(i.Push.Tokens)
	// make the worker arrays
	pushes := make(chan *IOSPush, totalPushes)
	responses := make(chan *IOSPush, totalPushes)

	// update the waitgroup with the amount of goroutines
	wgIOS.Add(totalPushes)

	client := apns2.NewClient(i.cert)

	// determine which APNS service to use
	if i.Client.UseSandboxIOS.Bool {
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
		iospush.notification.Topic = i.Client.BundleIdentifier.String

		pLoad, _ := json.Marshal(i.Push.Payload)
		iospush.notification.Payload = []byte(pLoad)

		iospush.Push = models.Push{}
		iospush.Push.ClientID = i.Push.ClientID
		iospush.Push.Token = token
		iospush.Push.Platform = null.String{NullString: sql.NullString{
			String: "ios",
			Valid:  true,
		}}
		iospush.Push.Payload = null.String{NullString: sql.NullString{
			String: string(pLoad),
			Valid:  true,
		}}

		// add to the worker array
		pushes <- &iospush
	}

	// check if there are "errors" and do callback (if applicable)
	failed, err := ioscheckAndCallback(totalPushes, responses, i)

	close(pushes)
	close(responses)
	wgIOS.Wait()

	return failed, err
}

func ioscheckAndCallback(total int, responses <-chan *IOSPush, i *IOS) ([]models.Push, error) {

	failed := make([]models.Push, 0)
	for i := 0; i < total; i++ {
		res := <-responses
		//fmt.Printf("response: %v\n", res)
		if !res.Push.Sent.Valid {
			failed = append(failed, res.Push)
		}
	}
	// fmt.Print(failed)

	// if all of the pushes failed (even if it's just one push sent)
	// then return an error message
	// NOTE: If there were various pushes sent successfully, and just a few that failed
	// then it's not an "error"
	var err error
	if len(failed) == total {
		err = models.ErrFailedToSendPush
	}

	return failed, err
}

func iosworker(db *sql.DB, client *apns2.Client, pushes <-chan *IOSPush, responses chan<- *IOSPush) {

	for p := range pushes {
		// push the message
		res, err := client.Push(p.notification)
		if err != nil {
			log.Fatal("IOS Push Error: ", err)
		}
		// check and save the response
		p.Push.Response = null.String{NullString: sql.NullString{
			String: fmt.Sprintf("%d %s", res.StatusCode, res.Reason),
			Valid:  true,
		}}

		// add the sent datetime if statuscode is 200
		if res.StatusCode == 200 {
			p.Push.Sent = null.String{NullString: sql.NullString{
				String: jodaTime.Format("YYYY-MM-dd HH:mm:ss", time.Now()),
				Valid:  true,
			}}
		}

		// save the push record
		p.Push.Create(db)

		// now update the responses
		responses <- p

		// this waitgroup routine is "done"
		wgIOS.Done()

		//fmt.Printf("DeviceToken: %v StatusCode: %v ApnsID: %v Reason: %v\n", p.notification.DeviceToken, res.StatusCode, res.ApnsID, res.Reason)
	}

}
