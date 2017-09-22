package models

import (
	"os"
	"strings"
)

// DBconfig struct
type DBconfig struct {
	DBHost string
	DBPort string
	DBName string
	DBUser string
	DBPass string
}

// ReturnConfig function
func ReturnConfig() DBconfig {

	var dbName = "livepush"
	if strings.Contains(os.Args[0], "/_test/") {
		dbName = "livepush_test"
	}

	return DBconfig{
		DBHost: "192.168.31.254",
		DBPort: "3306",
		DBName: dbName,
		DBUser: "tester",
		DBPass: "test",
	}
}
