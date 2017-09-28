package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"

	null "gopkg.in/guregu/null.v3"
)

// Push struct
type Push struct {
	ID       int          `json:"id"`
	ClientID null.Int     `json:"clientid"`
	Token    null.String  `json:"token"`
	Platform null.String  `json:"platform"`
	Title    null.String  `json:"title"`
	Subtitle null.String  `json:"subtitle"`
	Body     null.String  `json:"body"`
	Badge    null.Int     `json:"badge"`
	Image    null.String  `json:"image"`
	Sound    null.String  `json:"sound"`
	Type     null.Int     `json:"type"`
	TheID    null.String  `json:"the_id"`
	UserID   null.Int     `json:"user_id"`
	Track    null.String  `json:"track"`
	Main     null.String  `json:"main"`
	Options  []PushOption `json:"options"`
	Inserted null.String  `json:"inserted"`
	Sent     null.String  `json:"sent"`
	Response null.String  `json:"response"`
	Attempts null.Int     `json:"attempts"`
}

// Get method
func (p *Push) Get(db *sql.DB, id int) error {

	//jsonOptions, _ := json.Marshal(&p.Options)
	//, options
	//, bytes.NewBuffer(jsonOptions)

	// TODO: Need to resolve the bloody options!

	return db.QueryRow("SELECT id, clientid, token, platform, title, subtitle, body, badge, image, sound, type, the_id, user_id, track, main, inserted, sent, response, attempts FROM push WHERE id = ?", id).Scan(&p.ID, &p.ClientID, &p.Token, &p.Platform, &p.Title, &p.Subtitle, &p.Body, &p.Badge, &p.Image, &p.Sound, &p.Type, &p.TheID, &p.UserID, &p.Track, &p.Main, &p.Inserted, &p.Sent, &p.Response, &p.Attempts)
}

// Create function
func (p *Push) Create(db *sql.DB) error {
	var err error

	// check the required fields
	err = p.validateFields()
	if err != nil {
		// println("Required field err:", err.Error())
		return err
	}

	jsonOptions, _ := json.Marshal(&p.Options)

	res, err := db.Exec("INSERT INTO push (clientid, token, platform, title, subtitle, body, badge, image, sound, type, the_id, user_id, track, main, options, inserted, sent, response) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,NOW(),?,?)", &p.ClientID, &p.Token, &p.Platform, &p.Title, &p.Subtitle, &p.Body, &p.Badge, &p.Image, &p.Sound, &p.Type, &p.TheID, &p.UserID, &p.Track, &p.Main, fmt.Sprint(bytes.NewBuffer(jsonOptions)), &p.Sent, &p.Response)
	if err != nil {
		println("Exec err:", err.Error())
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		println("Error:", err.Error())
		return err
	}

	//println("LastInsertId:", id)
	err = p.Get(db, int(id))
	if err != nil {
		println("Error:", err.Error())
		return err
	}

	return err
}

// -----------------------
// HELPERS -------------
// -----------------------

func (p *Push) validateFields() error {

	var err error
	if !p.ClientID.Valid {
		err = ErrMissingClientID
	}

	// fmt.Printf("Notification object: %v err: %v", n, err)

	return err
}
