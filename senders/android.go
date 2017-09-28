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
	Data  dataType `json:"data"`
	Alert string   `json:"alert"`
	Sound int      `json:"sound"`
}

type dataType struct {
	Type     int     `json:"type"`
	ID       string  `json:"id"`
	UserID   int     `json:"user_id"`
	Track    string  `json:"track"`
	Actions  actions `json:"actions"`
	ImageURL string  `json:"image_url"`
}

type actions struct {
	Main string              `json:"main"`
	Opt  []models.PushOption `json:"opt"`
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

	o := make([]models.PushOption, 0)
	o = append(o, models.PushOption{Label: "Ver anúncio", Path: "/ads/8803927"})
	o = append(o, models.PushOption{Label: "Ver resultados", Path: "/saved-searches/@15063799774512"})

	data := data{
		Message: message{
			Data: dataType{
				Type:   int(a.Push.Type.Int64),
				ID:     a.Push.TheID.String,
				UserID: int(a.Push.UserID.Int64),
				Track:  a.Push.Track.String,
				Actions: actions{
					Main: a.Push.Main.String,
					Opt:  o,
				},
				ImageURL: a.Push.Image.String,
			},
			Alert: a.Push.Body.String,
			Sound: 1,
		},
		Title:    a.Push.Title.String,
		Subtitle: a.Push.Subtitle.String,
		Body:     a.Push.Body.String,
		Badge:    int(a.Push.Badge.Int64),
		Image:    a.Push.Image.String,
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
