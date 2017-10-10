package app_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/danfixeads/livepush/app"
)

var a = app.App{}
var testClientID = "admin@services.olx.com"

func TestMain(m *testing.M) {

	a.SetUpDatabase()
	a.SetUpDatabaseTables()
	a.SetUpRouter()

	code := m.Run()

	os.Exit(code)
}

func TestCloseDatabaseConnection(t *testing.T) {
	a.CloseDatabaseConnnection()

	err := a.Database.Ping()
	if err == nil {
		t.Error("Database connection should be closed")
	} else {
		// reopen the connection
		a.SetUpDatabase()
	}
}

func TestSetUpDatabaseTables(t *testing.T) {

	var err error

	err = a.SetUpDatabaseTables()
	if err != nil {
		t.Errorf("Error setting up the tables: %v", err)
	}

	// check if the table exists
	_, err = a.Database.Query("SELECT 1 FROM push LIMIT 1")
	if err != nil {
		t.Errorf("Database query error: %v", err)
	}

	_, err = a.Database.Query("SELECT 1 FROM client LIMIT 1")
	if err != nil {
		t.Errorf("Database query error: %v", err)
	}

}

func TestSetUpRouter(t *testing.T) {
	err := a.SetUpRouter()
	if err != nil {
		t.Errorf("Setup router error: %v", err)
	}
}

// -----------------------
// HELPERS -------------
// -----------------------

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}
