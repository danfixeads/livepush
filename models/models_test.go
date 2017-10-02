package models_test

import (
	"os"
	"testing"

	"github.com/danfixeads/livepush/app"
)

var a = app.App{}

func TestMain(m *testing.M) {

	a.SetUpDatabase()
	a.SetUpDatabaseTables()
	a.SetUpRouter()

	code := m.Run()

	os.Exit(code)
}
