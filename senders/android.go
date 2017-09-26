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
	Message message `json:"message"`
}

type message struct {
	Data  data2  `json:"data"`
	Alert string `json:"alert"`
}

type data2 struct {
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

// SetUpClient function
func (a *Android) SetUpClient() error {

	return nil
}

// SendMessage function
func (a *Android) SendMessage(db *sql.DB) error {

	c := fcm.NewFcmClient(a.client.FCMAuthKey.String)
	//c.SetDryRun(true)

	/*
		data := map[string]string{
			"msg":        "Hello World1",
			"type":       "0",
			"main":       "\\",
			"alert":      "hello",
			"resulttype": "test",
			"title":      "Hello!",
		}*/

	data := data{
		Message: message{
			Data: data2{
				Type: 0,
				Actions: actions{
					Main: "/",
				},
			},
			Alert: "Hello World!",
		},
	}
	jsonByte, _ := json.Marshal(data)
	fmt.Print(jsonByte)
	fmt.Print(bytes.NewBuffer(jsonByte))

	ids := []string{
		"fRE69G6iGx0:APA91bGJZBlY-2Ljor-WeDEWZghcA0yY5SC5pJeNtQp_OHnlktCy_2uQTacceaRUp5ieIiW6CLk6DXndBJeAReHLVvV1DgA4cpOyUaBU0Wb6CNJ86vOo9RnG0U9h9PFuAdi4nSNbc1qH",
	}

	c.NewFcmRegIdsMsg(ids, data)

	status, err := c.Send()

	if err == nil {
		status.PrintResults()
	} else {
		fmt.Println(err)
	}

	return nil
}
