package senders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

	for i := 0; i < totalPushes; i++ {
		go androidworker(db, a.client.FCMAuthKey.String, pushes)
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

	return nil
}

func androidworker(db *sql.DB, authKey string, pushes <-chan *AndroidPush) {

	for p := range pushes {

		ids := []string{
			p.token,
		}

		client := fcm.NewFcmClient(authKey)

		//payload := []byte(`{"message":{"data":{"type":0,"actions":{"main":"/"}},"alert":"Hello World!"}}`)

		//shit := bytes.NewBuffer(payload)

		client.NewFcmRegIdsMsg(ids, p.payload) // p.payload

		res, err := client.Send()
		if err != nil {
			log.Fatal("Push Error: ", err)
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

		if err == nil {
			res.PrintResults()
		} else {
			fmt.Println(err)
		}

	}

}
