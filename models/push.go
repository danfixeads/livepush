package models

import (
	"database/sql"
	"errors"

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
	Sent     null.String `json:"sent"`
	Response null.String `json:"response"`
	Attempts null.Int    `json:"attempts"`
}

func (p *Push) Create(db *sql.DB) error {
	return errors.New("Not implemented")
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
