package models

import (
	"database/sql"
	"errors"

	null "gopkg.in/guregu/null.v3"
)

// Client struct
type Client struct {
	ID         int         `json:"id"`
	ClientID   null.Int    `json:"clientid"`
	PemFile    null.String `json:"pemfile"`
	P12File    null.String `json:"p12file"`
	PassPhrase null.String `json:"passphrase"`
	FCMToken   null.String `json:"fcmtoken"`
	Active     bool        `json:"ative"`
	Inserted   null.String `json:"inserted"`
	Updated    null.String `json:"updated"`
}

// Get method
func (c *Client) Get(db *sql.DB, id int) error {
	return db.QueryRow("SELECT id, clientid, pemfile, p12file, fcmtoken, active, inserted, updated FROM client WHERE id = ?", id).Scan(&c.ID, &c.ClientID, &c.PemFile, &c.P12File, &c.FCMToken, &c.Active, &c.Inserted, &c.Updated)
}

// Create function
func (c *Client) Create(db *sql.DB) error {

	var err error

	// check the required fields
	err = c.validateFields()
	if err != nil {
		// println("Required field err:", err.Error())
		return err
	}

	res, err := db.Exec("INSERT INTO client (clientid) VALUES(?)", &c.ClientID)
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
	err = c.Get(db, int(id))
	if err != nil {
		println("Error:", err.Error())
		return err
	}

	return err
}

// Update function
func (c *Client) Update(db *sql.DB) error {

	// check the required fields
	if err := c.validateFields(); err != nil {
		// println("Required field err:", err.Error())
		return err
	}

	_, err :=
		db.Exec("UPDATE client SET clientid=?, pemfile=?, p12file=?, passphrase=?, fcmtoken=?, active=?, updated=NOW() WHERE id = ?",
			c.ClientID, c.PemFile, c.P12File, c.PassPhrase, c.FCMToken, c.Active, c.ID)

	return err
}

// Delete function
func (c *Client) Delete(db *sql.DB) error {
	rows, err := db.Exec("DELETE FROM client WHERE id = ?", c.ID)

	if affected, _ := rows.RowsAffected(); affected == 0 {
		err = errors.New("No records were deleted")
	}

	return err
}

// ListClients function
func ListClients(db *sql.DB, start, limit int) ([]Client, error) {

	rows, err := db.Query("SELECT id, clientid, pemfile, p12file, fcmtoken, active, inserted, updated FROM client LIMIT ? OFFSET ?", limit, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	clients := []Client{}

	for rows.Next() {
		var c Client
		if err := rows.Scan(&c.ID, &c.ClientID, &c.PemFile, &c.P12File, &c.FCMToken, &c.Active, &c.Inserted, &c.Updated); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}

	return clients, nil
}

// -----------------------
// HELPERS -------------
// -----------------------

func (c *Client) validateFields() error {

	var err error
	if !c.ClientID.Valid {
		err = ErrMissingClientID
	}
	if !c.PemFile.Valid && !c.P12File.Valid && !c.FCMToken.Valid {
		err = ErrMissingVitalFields
	}
	if (c.PemFile.Valid || c.P12File.Valid) && !c.PassPhrase.Valid {
		err = ErrMissingPassPhrase
	}

	// fmt.Printf("Notification object: %v err: %v", n, err)

	return err
}
