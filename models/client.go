package models

import (
	"database/sql"
	"errors"

	null "gopkg.in/guregu/null.v3"
	"gopkg.in/guregu/null.v3/zero"
)

// Client struct
type Client struct {
	ID               int         `json:"id"`
	ClientID         null.Int    `json:"clientid"`
	PemFile          null.String `json:"pemfile"`
	P12File          null.String `json:"p12file"`
	PassPhrase       null.String `json:"passphrase"`
	BundleIdentifier null.String `json:"bundleidentifier"`
	UseSandboxIOS    zero.Bool   `json:"usesandboxios"`
	FCMAuthKey       null.String `json:"fcmauthkey"`
	WebHook          null.String `json:"webhook"`
	Active           zero.Bool   `json:"active"`
	Inserted         null.String `json:"inserted"`
	Updated          null.String `json:"updated"`
}

// Get method
func (c *Client) Get(db *sql.DB, id int) error {
	return db.QueryRow("SELECT id, clientid, pemfile, p12file, passphrase, bundleidentifier, usesandboxios, fcmauthkey, webhook, active, inserted, updated FROM client WHERE id = ?", id).Scan(&c.ID, &c.ClientID, &c.PemFile, &c.P12File, &c.PassPhrase, &c.BundleIdentifier, &c.UseSandboxIOS, &c.FCMAuthKey, &c.WebHook, &c.Active, &c.Inserted, &c.Updated)
}

// GetByClientID function
func (c *Client) GetByClientID(db *sql.DB, clientid int) error {
	if err := db.QueryRow("SELECT id FROM client WHERE id = ? AND active = 1", clientid).Scan(&c.ID); err != nil {
		return err
	}
	return c.Get(db, c.ID)
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

	var isActive = false
	if c.Active.Valid {
		isActive = c.Active.Bool
	}

	res, errExec := db.Exec("INSERT INTO client (clientid, pemfile, p12file, passphrase, bundleidentifier, usesandboxios, fcmauthkey, webhook, active, inserted) VALUES(?,?,?,?,?,?,?,?,?,NOW())", &c.ClientID, &c.PemFile, &c.P12File, &c.PassPhrase, &c.BundleIdentifier, &c.UseSandboxIOS, &c.FCMAuthKey, &c.WebHook, isActive)
	err = errExec

	if err == nil {
		id, errLastInsertID := res.LastInsertId()
		err = errLastInsertID

		//println("LastInsertId:", id)

		if err == nil {
			err = c.Get(db, int(id))
		}

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

	// check to see if the record exists or not
	var check Client
	if err := check.Get(db, int(c.ID)); err != nil {
		return ErrRecordNotFound
	}

	_, err :=
		db.Exec("UPDATE client SET clientid=?, pemfile=?, p12file=?, passphrase=?, bundleidentifier=?, usesandboxios=?, fcmauthkey=?, webhook=?, active=?, updated=NOW() WHERE id = ?",
			c.ClientID, c.PemFile, c.P12File, c.PassPhrase, c.BundleIdentifier, c.UseSandboxIOS, c.FCMAuthKey, c.WebHook, c.Active, c.ID)

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

	var err error
	rows, errQuery := db.Query("SELECT id, clientid, pemfile, p12file, bundleidentifier, usesandboxios, fcmauthkey, webhook, active, inserted, updated FROM client ORDER BY id DESC LIMIT ? OFFSET ?", limit, start)
	err = errQuery

	defer rows.Close()

	clients := []Client{}

	for rows.Next() {
		var c Client
		errScan := rows.Scan(&c.ID, &c.ClientID, &c.PemFile, &c.P12File, &c.BundleIdentifier, &c.UseSandboxIOS, &c.FCMAuthKey, &c.WebHook, &c.Active, &c.Inserted, &c.Updated)
		err = errScan
		clients = append(clients, c)
	}

	return clients, err
}

// -----------------------
// HELPERS -------------
// -----------------------

func (c *Client) validateFields() error {

	var err error
	if !c.ClientID.Valid {
		err = ErrMissingClientID
	}
	if !c.PemFile.Valid && !c.P12File.Valid && !c.FCMAuthKey.Valid {
		err = ErrMissingVitalFields
	}
	if (c.PemFile.Valid || c.P12File.Valid) && !c.PassPhrase.Valid {
		err = ErrMissingPassPhrase
	}
	if (c.PemFile.Valid || c.P12File.Valid) && !c.BundleIdentifier.Valid {
		err = ErrMissingBundleIdentifier
	}

	// fmt.Printf("Client object: %v err: %v", c, err)

	return err
}
