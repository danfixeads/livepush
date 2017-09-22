package main

import (
	"log"

	"github.com/danfixeads/livepush/app"
)

// App variable
var App = app.App{}

func main() {
	setUpApp()
}

func setUpApp() {
	checkSetUpError(App.SetUpDatabase())
	checkSetUpError(App.SetUpDatabaseTables())
	checkSetUpError(App.SetUpRouter())
}

func checkSetUpError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
