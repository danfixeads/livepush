package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/danfixeads/livepush/models"
	"github.com/gorilla/mux"
	newrelic "github.com/newrelic/go-agent"

	// Need the MySql driver
	_ "github.com/go-sql-driver/mysql"
)

// App struct
type App struct {
	Router      *mux.Router
	Database    *sql.DB
	newrelicapp newrelic.Application
}

// SetUpDatabase function
func (a *App) SetUpDatabase() error {

	dbConfig := models.ReturnConfig()

	// try and connect
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig.DBUser, dbConfig.DBPass, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
	//fmt.Printf("Database connection string: %s", dsn)
	a.Database, _ = sql.Open("mysql", dsn)
	return a.Database.Ping()
}

// CloseDatabaseConnnection will close the database connection
func (a *App) CloseDatabaseConnnection() {
	a.Database.Close()
}

// SetUpDatabaseTables function
func (a *App) SetUpDatabaseTables() error {

	var err error

	tablePushQuery := `CREATE TABLE IF NOT EXISTS push (
					id int(10) unsigned NOT NULL AUTO_INCREMENT,
					clientid int(10) DEFAULT NULL,
					token varchar(255) DEFAULT NULL,
					platform enum('android','ios') DEFAULT NULL,
					payload mediumtext DEFAULT NULL,
					inserted datetime DEFAULT NULL,
					sent datetime DEFAULT NULL,
					response varchar(255) DEFAULT NULL,
					attempts int(10) DEFAULT NULL,
					PRIMARY KEY (id), 
					KEY platform (platform)
			   ) ENGINE=InnoDB DEFAULT CHARSET=latin1`

	_, err = a.Database.Query(tablePushQuery)

	tableAppLogQuery := `CREATE TABLE IF NOT EXISTS applog (
					id int(10) unsigned NOT NULL AUTO_INCREMENT,
					inserted datetime DEFAULT NULL,
					platform enum('android','ios') DEFAULT NULL,
					description varchar(255) DEFAULT NULL,
					PRIMARY KEY (id), 
					KEY platform (platform)
			   ) ENGINE=InnoDB DEFAULT CHARSET=latin1`

	_, err = a.Database.Query(tableAppLogQuery)

	tableClientQuery := `CREATE TABLE IF NOT EXISTS client (
					id int(10) unsigned NOT NULL AUTO_INCREMENT,
					clientid int(10) DEFAULT NULL,
					pemfile varchar(30) DEFAULT NULL,
					p12file varchar(30) DEFAULT NULL,
					passphrase varchar(30) DEFAULT NULL,
					bundleidentifier varchar(40) DEFAULT NULL,
					usesandboxios int(1) DEFAULT NULL,
					fcmauthkey varchar(255) DEFAULT NULL,
					webhook varchar(255) DEFAULT NULL,
					active int(1) DEFAULT NULL,
					inserted datetime DEFAULT NULL,
					updated datetime DEFAULT NULL,
					PRIMARY KEY (id), 
					KEY platform (clientid)
			   ) ENGINE=InnoDB DEFAULT CHARSET=latin1`

	_, err = a.Database.Query(tableClientQuery)

	return err
}

// SetUpRouter func
func (a *App) SetUpRouter() error {
	a.Router = mux.NewRouter()

	// client handling
	a.routerFunc("/clients", a.clientList).Methods("GET")
	a.routerFunc("/clients/{start:[0-9]+}", a.clientList).Methods("GET")
	a.routerFunc("/clients/{limit:[0-9]+}", a.clientList).Methods("GET")
	a.routerFunc("/clients/{start:[0-9]+}/{limit:[0-9]+}", a.clientList).Methods("GET")
	a.routerFunc("/client", a.clientCreate).Methods("POST")
	a.routerFunc("/client/{id:[0-9]+}", a.clientUpdate).Methods("PUT")
	a.routerFunc("/client/{id:[0-9]+}", a.clientDelete).Methods("DELETE")
	a.routerFunc("/client/{id:[0-9]+}", a.clientGet).Methods("GET")

	// push handling
	a.routerFunc("/push/ios", a.createPushIOS).Methods("POST")
	a.routerFunc("/push/android", a.createPushAndroid).Methods("POST")

	a.routerFunc("/pushes", a.pushList).Methods("GET")
	a.routerFunc("/pushes/{start:[0-9]+}", a.pushList).Methods("GET")
	a.routerFunc("/pushes/{limit:[0-9]+}", a.pushList).Methods("GET")
	a.routerFunc("/pushes/{start:[0-9]+}/{limit:[0-9]+}", a.pushList).Methods("GET")
	a.routerFunc("/push/{id:[0-9]+}", a.pushDelete).Methods("DELETE")
	a.routerFunc("/push/{id:[0-9]+}", a.pushGet).Methods("GET")

	// return and start the server (if not test)
	if strings.Contains(os.Args[0], "/_test/") {
		return nil
	}
	return http.ListenAndServe(":8080", a.Router)
}

func (a *App) routerFunc(path string, f func(http.ResponseWriter,
	*http.Request)) *mux.Route {

	// if New Relic is active, then wrap the function to enable the transactions for this route request
	if a.newrelicapp != nil {
		a.Router.HandleFunc(newrelic.WrapHandleFunc(a.newrelicapp, path, f))
	}

	return a.Router.HandleFunc(path, f)
}

// SetUpNewRelic function
func (a *App) SetUpNewRelic() error {

	if strings.Contains(os.Args[0], "/_test/") {
		return nil
	}

	var err error

	cfg := models.ReturnConfig()

	if len(cfg.NewRelicKey) > 0 {
		config := newrelic.NewConfig(cfg.NewRelicName, cfg.NewRelicKey)
		app, errNewRelic := newrelic.NewApplication(config)
		err = errNewRelic
		a.newrelicapp = app
	}

	return err
}

// -----------------------
// HELPERS -------------
// -----------------------

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
