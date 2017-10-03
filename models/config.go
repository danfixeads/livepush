package models

import (
	"os"
	"strings"
)

// Config struct
type Config struct {
	DBHost       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPass       string
	NewRelicName string
	NewRelicKey  string
	MQHost       string
	MQUser       string
	MQPass       string
	MQPort       string
}

// ReturnConfig function
func ReturnConfig() Config {

	// determine the DB required
	var dbName = "livepush"
	if strings.Contains(os.Args[0], "/_test/") {
		dbName = "livepush_test"
	}

	return Config{
		// DB
		DBHost: "192.168.31.254",
		DBPort: "3306",
		DBName: dbName,
		DBUser: "tester",
		DBPass: "test",
		// New Relic
		NewRelicName: "LivePush",
		NewRelicKey:  "504d6bc51121bc1e3b88cdc654b1411456979237",
		// Rabbit MQ
		MQHost: "rabbit.storiaro.fixeads.com",
		MQUser: "rabbit",
		MQPass: "administrator",
		MQPort: "5672",
	}
}
