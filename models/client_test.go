package models_test

import (
	"database/sql"
	"testing"

	"gopkg.in/guregu/null.v3/zero"

	"github.com/danfixeads/livepush/models"
	null "gopkg.in/guregu/null.v3"
)

// -----------------------
// GET -------------------
// -----------------------

func TestGetClient(t *testing.T) {

	clearTestClients()
	addTestClients(5)

	var client models.Client
	err := client.Get(a.Database, 1)
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}

}

func TestGetByClientID(t *testing.T) {

	clearTestClients()
	addTestClients(5)

	var client models.Client
	err := client.GetByClientID(a.Database, "xpto")
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}

}

func TestGetClientByInvalidClientID(t *testing.T) {

	clearTestClients()
	addTestClients(5)

	var client models.Client
	err := client.GetByClientID(a.Database, "fake")
	if err == nil {
		t.Error("Should not have returned any rows")
	}

}

// -----------------------
// CREATE ----------------
// -----------------------
func TestCreateClient(t *testing.T) {

	clearTestClients()

	var client models.Client
	client.ClientID = null.String{NullString: sql.NullString{
		String: "xpto",
		Valid:  true,
	}}
	client.P12File = null.String{NullString: sql.NullString{
		String: "p12.pem",
		Valid:  true,
	}}
	client.BundleIdentifier = null.String{NullString: sql.NullString{
		String: "com.fixeads.anApp",
		Valid:  true,
	}}
	client.PassPhrase = null.String{NullString: sql.NullString{
		String: "a_passphrase",
		Valid:  true,
	}}
	client.Active = zero.Bool{NullBool: sql.NullBool{
		Bool:  true,
		Valid: true,
	}}

	err := client.Create(a.Database)
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}

}

func TestCreateClientWithEmptyValues(t *testing.T) {

	var client models.Client
	err := client.Create(a.Database)
	if err == nil {
		t.Error("Should have returned validation errors")
	}
}

func TestCreateClientWithEmptyPassPhrase(t *testing.T) {

	var client models.Client
	client.ClientID = null.String{NullString: sql.NullString{
		String: "xpto",
		Valid:  true,
	}}
	client.P12File = null.String{NullString: sql.NullString{
		String: "p12.pem",
		Valid:  true,
	}}
	client.BundleIdentifier = null.String{NullString: sql.NullString{
		String: "com.fixeads.anApp",
		Valid:  true,
	}}
	err := client.Create(a.Database)
	if err == nil {
		t.Error("Should have returned validation errors")
	}
}

func TestCreateClientWithEmptyBundleIdentifier(t *testing.T) {

	var client models.Client
	client.ClientID = null.String{NullString: sql.NullString{
		String: "xpto",
		Valid:  true,
	}}
	client.P12File = null.String{NullString: sql.NullString{
		String: "p12.pem",
		Valid:  true,
	}}
	client.PassPhrase = null.String{NullString: sql.NullString{
		String: "a_phrase",
		Valid:  true,
	}}
	err := client.Create(a.Database)
	if err == nil {
		t.Error("Should have returned validation errors")
	}
}

// -----------------------
// UPDATE ----------------
// -----------------------

func TestUpdate(t *testing.T) {

	clearTestClients()
	addTestClients(5)

	var client models.Client
	client.Get(a.Database, 2)

	client.ClientID = null.String{NullString: sql.NullString{
		String: "abc",
		Valid:  true,
	}}
	client.P12File = null.String{NullString: sql.NullString{
		String: "p12.pem",
		Valid:  true,
	}}
	client.BundleIdentifier = null.String{NullString: sql.NullString{
		String: "com.fixeads.anApp",
		Valid:  true,
	}}
	client.PassPhrase = null.String{NullString: sql.NullString{
		String: "a_passphrase",
		Valid:  true,
	}}
	client.Active = zero.Bool{NullBool: sql.NullBool{
		Bool:  true,
		Valid: true,
	}}

	err := client.Update(a.Database)
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}

	var client2 models.Client
	client2.Get(a.Database, 2)

	if client.ClientID != client2.ClientID {
		t.Error("Did not update the ClientID")
	}
}

func TestUpdateClientWithInvalidValues(t *testing.T) {
	clearTestClients()
	addTestClients(5)

	var client models.Client
	client.Get(a.Database, 2)

	client.ClientID = null.String{NullString: sql.NullString{
		String: "abc",
		Valid:  true,
	}}
	client.FCMAuthKey = null.String{NullString: sql.NullString{
		String: "",
		Valid:  false,
	}}

	err := client.Update(a.Database)
	if err == nil {
		t.Error("Should have returned validation errors")
	}
}

func TestUpdateNonExistingClient(t *testing.T) {

	clearTestClients()
	addTestClients(5)

	var client models.Client

	client.ClientID = null.String{NullString: sql.NullString{
		String: "abc",
		Valid:  true,
	}}
	client.P12File = null.String{NullString: sql.NullString{
		String: "p12.pem",
		Valid:  true,
	}}
	client.BundleIdentifier = null.String{NullString: sql.NullString{
		String: "com.fixeads.anApp",
		Valid:  true,
	}}
	client.PassPhrase = null.String{NullString: sql.NullString{
		String: "a_passphrase",
		Valid:  true,
	}}
	client.Active = zero.Bool{NullBool: sql.NullBool{
		Bool:  true,
		Valid: true,
	}}

	err := client.Update(a.Database)
	if err == nil {
		t.Error("Should have returned validation errors")
	}
}

// -----------------------
// DELETE ----------------
// -----------------------

func TestDeleteClient(t *testing.T) {

	clearTestClients()
	addTestClients(5)

	var client models.Client
	client.ID = 2
	err := client.Delete(a.Database)
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}
}

func TestDeleteNonExistingClient(t *testing.T) {

	clearTestClients()
	addTestClients(5)

	var client models.Client
	client.ID = 200
	err := client.Delete(a.Database)
	if err == nil {
		t.Error("Should have returned validation errors")
	}
}

// -----------------------
// LIST ------------------
// -----------------------

func TestListClients(t *testing.T) {

	clearTestClients()
	addTestClients(25)

	clients, err := models.ListClients(a.Database, 0, 50)
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}

	if len(clients) != 25 {
		t.Error("Should have returned 25 mock clients")
	}

}

func TestListClientsPagination(t *testing.T) {

	clearTestClients()
	addTestClients(40)

	clients, err := models.ListClients(a.Database, 0, 20)
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}

	if len(clients) != 20 {
		t.Error("Should have returned 20 mock clients and not 40!")
	}

}
