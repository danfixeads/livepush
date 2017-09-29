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
	Payload  null.String `json:"payload"`
	Inserted null.String `json:"inserted"`
	Sent     null.String `json:"sent"`
	Response null.String `json:"response"`
	Attempts null.Int    `json:"attempts"`
}

// Get method
func (p *Push) Get(db *sql.DB, id int) error {
	return db.QueryRow("SELECT id, clientid, token, platform, payload, inserted, sent, response, attempts FROM push WHERE id = ?", id).Scan(&p.ID, &p.ClientID, &p.Token, &p.Platform, &p.Payload, &p.Inserted, &p.Sent, &p.Response, &p.Attempts)
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

	res, err := db.Exec("INSERT INTO push (clientid, token, platform, payload, inserted, sent, response) VALUES (?,?,?,?,NOW(),?,?)", &p.ClientID, &p.Token, &p.Platform, &p.Payload, &p.Sent, &p.Response)
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

// Delete function
func (p *Push) Delete(db *sql.DB) error {
	rows, err := db.Exec("DELETE FROM push WHERE id = ?", p.ID)

	if affected, _ := rows.RowsAffected(); affected == 0 {
		err = errors.New("No records were deleted")
	}

	return err
}

// ListPushes function
func ListPushes(db *sql.DB, start, limit int) ([]Push, error) {

	rows, err := db.Query("SELECT id, clientid, token, platform, payload, inserted, sent, response, attempts FROM push ORDER BY id DESC LIMIT ? OFFSET ?", limit, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	pushes := []Push{}

	for rows.Next() {
		var p Push
		if err := rows.Scan(&p.ID, &p.ClientID, &p.Token, &p.Platform, &p.Payload, &p.Inserted, &p.Sent, &p.Response, &p.Attempts); err != nil {
			return nil, err
		}
		pushes = append(pushes, p)
	}

	return pushes, nil
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
