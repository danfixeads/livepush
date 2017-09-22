package models

import (
	"os"
	"strings"
)

type dbconfig struct {
	DBHost string
	DBPort string
	DBName string
	DBUser string
	DBPass string
}

func ReturnConfig() dbconfig {

	var dbName = "livepush"
	if strings.Contains(os.Args[0], "/_test/") {
		dbName = "livepush_test"
	}

	return dbconfig{
		DBHost: "192.168.31.254",
		DBPort: "3306",
		DBName: dbName,
		DBUser: "tester",
		DBPass: "test",
	}
}
