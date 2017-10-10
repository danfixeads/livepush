package models_test

import (
	"database/sql"
	"testing"

	"github.com/danfixeads/livepush/models"
	null "gopkg.in/guregu/null.v3"
)

var testAuthorization = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzdHMiLCJleHAiOjc1MDcxMzE5OTYsImlhdCI6MTUwNzEzMTk5NiwiaXNzIjoiYWRtaW5Ac2VydmljZXMub2x4LmNvbSIsIm5iZiI6MTUwNzEzMTk5Nn0.Zm8pSlZahvWPhsl9eoVhAKgEtezjn-ht5H91XmE8W3pKenTkYT4bC-bn97dO-6NACQ43RTtNIlU327cBlpAuM22kUWVR2qd5X6Bo3-PeRb39s_uJNTDcngBxzwLaAItePVUa5fEcSz4_PTn8jVW4-m3K9kq0_Ql3rElvjZVT_6c"
var testClientID = "admin@services.olx.com"

// -----------------------
// GET -------------------
// -----------------------

func TestGetPush(t *testing.T) {

	clearTestPushes()
	addTestPushes(5)

	var push models.Push
	push.ClientID = null.String{NullString: sql.NullString{
		String: testClientID,
		Valid:  true,
	}}
	err := push.Get(a.Database, 1)
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}

}

func TestGetPushByInvalidID(t *testing.T) {

	clearTestPushes()
	addTestPushes(5)

	var push models.Push
	push.ClientID = null.String{NullString: sql.NullString{
		String: testClientID,
		Valid:  true,
	}}
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
	push.ClientID = null.String{NullString: sql.NullString{
		String: testClientID,
		Valid:  true,
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
	push.ClientID = null.String{NullString: sql.NullString{
		String: testClientID,
		Valid:  true,
	}}
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
	push.ClientID = null.String{NullString: sql.NullString{
		String: testClientID,
		Valid:  true,
	}}
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

	pushes, err := models.ListPushes(a.Database, 0, 50, testClientID)
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

	pushes, err := models.ListPushes(a.Database, 0, 20, testClientID)
	if err != nil {
		t.Errorf("Following error occured: %v", err)
	}

	if len(pushes) != 20 {
		t.Error("Should have returned 20 mock pushes and not 40!")
	}

}
