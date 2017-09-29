package senders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	fcm "github.com/NaySoftware/go-fcm"
	"github.com/danfixeads/livepush/models"
	"github.com/vjeantet/jodaTime"
	null "gopkg.in/guregu/null.v3"
)

// Android Struct
type Android struct {
	ClientID int
	Push     models.MultiplePush
	client   models.Client
}

// AndroidPush struct
type AndroidPush struct {
	push    models.Push
	token   string
	payload map[string]interface{}
}

// setup the waitgroup variable
var wgAndroid sync.WaitGroup

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

	totalPushes := len(a.Push.Tokens)
	// make the worker array
	pushes := make(chan *AndroidPush, totalPushes)
	responses := make(chan *AndroidPush, totalPushes)

	// update the waitgroup with the amount of goroutines
	wgAndroid.Add(totalPushes)

	for i := 0; i < totalPushes; i++ {
		go androidworker(db, a.client.FCMAuthKey.String, pushes, responses)
	}

	for _, token := range a.Push.Tokens {

		androidpush := AndroidPush{}

		androidpush.token = token.String
		androidpush.payload = a.Push.Payload

		androidpush.push = models.Push{}
		androidpush.push.ClientID = a.Push.ClientID
		androidpush.push.Token = token
		androidpush.push.Platform = null.String{NullString: sql.NullString{
			String: "android",
			Valid:  true,
		}}

		pLoad, _ := json.Marshal(a.Push.Payload)

		androidpush.push.Payload = null.String{NullString: sql.NullString{
			String: string(pLoad),
			Valid:  true,
		}}

		pushes <- &androidpush

	}

	// check if there are "errors" and do callback (if applicable)
	err := androidcheckAndCallback(totalPushes, responses)

	close(pushes)
	close(responses)
	wgAndroid.Wait()

	return err
}

func androidcheckAndCallback(total int, responses <-chan *AndroidPush) error {

	failed := make([]AndroidPush, 0)
	for i := 0; i < total; i++ {
		res := <-responses
		//fmt.Printf("response: %v\n", res)
		if !res.push.Sent.Valid {
			failed = append(failed, *res)
			fmt.Printf("failed: %v\n", failed)
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

func androidworker(db *sql.DB, authKey string, pushes <-chan *AndroidPush, responses chan<- *AndroidPush) {

	for p := range pushes {

		ids := []string{
			p.token,
		}

		client := fcm.NewFcmClient(authKey)

		client.NewFcmRegIdsMsg(ids, p.payload) // p.payload

		res, err := client.Send()
		if err != nil {
			log.Fatal("FCM Push Error: ", err)
		}

		results, _ := json.Marshal(res.Results)
		p.push.Response = null.String{NullString: sql.NullString{
			String: fmt.Sprintf("%d %v", res.StatusCode, string(results)),
			Valid:  true,
		}}

		// add the sent datetime if statuscode is 200
		if res.StatusCode == 200 && res.Success == 1 {
			p.push.Sent = null.String{NullString: sql.NullString{
				String: jodaTime.Format("YYYY-MM-dd HH:mm:ss", time.Now()),
				Valid:  true,
			}}
		}

		// add to the database
		p.push.Create(db)

		// now update the responses
		responses <- p

		/*
			if err == nil {
				res.PrintResults()
			} else {
				fmt.Println(err)
			}*/

		// this waitgroup routine is "done"
		wgAndroid.Done()

	}

}
