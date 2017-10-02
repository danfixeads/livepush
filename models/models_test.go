package models_test

import (
	"fmt"
	"os"
	"strings"
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

// -----------------------
// HELPERS ---------------
// -----------------------

func clearTestClients() {
	_, err := a.Database.Exec("TRUNCATE client")
	if err != nil {
		panic(err)
	}
}

func addTestClients(count int) {

	if count < 1 {
		count = 1
	}

	var values = make([]string, count)

	for i := 0; i < count; i++ {
		values[i] = fmt.Sprintf("('%v',1,NOW())", i+1)
	}

	var query = fmt.Sprintf("INSERT INTO client (clientid,active,inserted) VALUES %v", strings.Join(values, ", "))
	//log.Print(query)
	a.Database.Exec(query)
}

func clearTestPushes() {
	_, err := a.Database.Exec("TRUNCATE push")
	if err != nil {
		panic(err)
	}
}

func addTestPushes(count int) {

	if count < 1 {
		count = 1
	}

	var values = make([]string, count)

	for i := 0; i < count; i++ {
		values[i] = fmt.Sprintf("('1','token_%v','ios',NOW())", i+1)
	}

	var query = fmt.Sprintf("INSERT INTO push (clientid,token,platform,inserted) VALUES %v", strings.Join(values, ", "))
	//log.Print(query)
	a.Database.Exec(query)
}
