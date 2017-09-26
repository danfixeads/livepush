package models

import (
	"database/sql"

	null "gopkg.in/guregu/null.v3"
)

// Push struct
type Push struct {
	ID       int         `json:"id"`
	ClientID null.Int    `json:"clientid"`
	Token    null.String `json:"token"`
	Platform null.String `json:"platform"`
	Title    null.String `json:"title"`
	Subtitle null.String `json:"subtitle"`
	Body     null.String `json:"body"`
	Badge    null.Int    `json:"badge"`
	Image    null.String `json:"image"`
	Sound    null.String `json:"sound"`
	Inserted null.String `json:"inserted"`
	Response null.String `json:"response"`
	Attempts null.Int    `json:"attempts"`
}

// Get method
func (p *Push) Get(db *sql.DB, id int) error {
	return db.QueryRow("SELECT id, clientid, token, platform, title, subtitle, body, badge, image, sound, inserted, response, attempts FROM push WHERE id = ?", id).Scan(&p.ID, &p.ClientID, &p.Token, &p.Platform, &p.Title, &p.Subtitle, &p.Body, &p.Badge, &p.Image, &p.Sound, &p.Inserted, &p.Response, &p.Attempts)
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

	res, err := db.Exec("INSERT INTO push (clientid, token, platform, title, subtitle, body, badge, image, sound, inserted, response) VALUES (?,?,?,?,?,?,?,?,?,NOW(),?)", &p.ClientID, &p.Token, &p.Platform, &p.Title, &p.Subtitle, &p.Body, &p.Badge, &p.Image, &p.Sound, &p.Response)
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
