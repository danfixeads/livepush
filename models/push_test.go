package models_test

import (
	"database/sql"
	"testing"

	"github.com/danfixeads/livepush/models"
	null "gopkg.in/guregu/null.v3"
)

// -----------------------
// GET -------------------
// -----------------------

func TestGetPush(t *testing.T) {

	clearTestPushes()
	addTestPushes(5)

	var push models.Push
	err := push.Get(a.Database, 1)
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}

}

func TestGetPushByInvalidID(t *testing.T) {

	clearTestPushes()
	addTestPushes(5)

	var push models.Push
	err := push.Get(a.Database, 100)
	if err == nil {
		t.Error("An error should have occured")
	}

}

// -----------------------
// CREATE ----------------
// -----------------------

func TestCreatePush(t *testing.T) {

	clearTestPushes()

	var push models.Push
	push.ClientID = null.Int{NullInt64: sql.NullInt64{
		Int64: 2,
		Valid: true,
	}}
	push.Platform = null.String{NullString: sql.NullString{
		String: "android",
		Valid:  true,
	}}
	push.Payload = null.String{NullString: sql.NullString{
		String: "{\"message\":{\"alert\":\"message\",\"data\":{\"actions\":{\"main\":\"/\"},\"type\":0},\"sound\":1}}",
		Valid:  true,
	}}

	err := push.Create(a.Database)
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}

}

func TestCreatePushWithEmptyValues(t *testing.T) {

	clearTestPushes()

	var push models.Push

	err := push.Create(a.Database)
	if err == nil {
		t.Error("Should have returned validation errors")
	}

}

// -----------------------
// DELETE ----------------
// -----------------------

func TestDeletePush(t *testing.T) {

	clearTestPushes()
	addTestPushes(5)

	var push models.Push
	push.ID = 2
	err := push.Delete(a.Database)
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}
}

func TestDeleteNonExistingPush(t *testing.T) {

	clearTestPushes()
	addTestPushes(5)

	var push models.Push
	push.ID = 200
	err := push.Delete(a.Database)
	if err == nil {
		t.Error("Should have returned validation errors")
	}
}

// -----------------------
// LIST ------------------
// -----------------------

func TestListPushes(t *testing.T) {

	clearTestPushes()
	addTestPushes(25)

	pushes, err := models.ListPushes(a.Database, 0, 50)
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}

	if len(pushes) != 25 {
		t.Error("Should have returned 25 mock pushes")
	}

}

func TestListPushesPagination(t *testing.T) {

	clearTestPushes()
	addTestPushes(40)

	pushes, err := models.ListPushes(a.Database, 0, 20)
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}

	if len(pushes) != 20 {
		t.Error("Should have returned 20 mock pushes and not 40!")
	}

}
