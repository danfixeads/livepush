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

	// Need the MySql driver
	_ "github.com/go-sql-driver/mysql"
)

// App struct
type App struct {
	Router   *mux.Router
	Database *sql.DB
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

	tableNotificationsQuery := `CREATE TABLE IF NOT EXISTS notifications (
					id int(10) unsigned NOT NULL AUTO_INCREMENT,
					token varchar(255) DEFAULT NULL,
					platform enum('android','ios') DEFAULT NULL,
					title varchar(100) NOT NULL,
					subtitle varchar(100) DEFAULT NULL,
					message varchar(255) DEFAULT NULL,
					badge int(11) DEFAULT NULL,
					inserted datetime DEFAULT NULL,
					sent datetime DEFAULT NULL,
					response varchar(255) DEFAULT NULL,
					attempts int(10) DEFAULT NULL,
					PRIMARY KEY (id), 
					KEY platform (platform)
			   ) ENGINE=InnoDB DEFAULT CHARSET=latin1`

	_, err = a.Database.Query(tableNotificationsQuery)

	tableAppLogQuery := `CREATE TABLE IF NOT EXISTS applog (
					id int(10) unsigned NOT NULL AUTO_INCREMENT,
					inserted datetime DEFAULT NULL,
					platform enum('android','ios') DEFAULT NULL,
					description varchar(255) DEFAULT NULL,
					PRIMARY KEY (id), 
					KEY platform (platform)
			   ) ENGINE=InnoDB DEFAULT CHARSET=latin1`

	_, err = a.Database.Query(tableAppLogQuery)

	tableClientsQuery := `CREATE TABLE IF NOT EXISTS clients (
					id int(10) unsigned NOT NULL AUTO_INCREMENT,
					clientid int(10) DEFAULT NULL,
					pemfile varchar(30) DEFAULT NULL,
					p12file varchar(30) DEFAULT NULL,
					passphrase varchar(30) DEFAULT NULL,
					fcmtoken varchar(255) DEFAULT NULL,
					active int(1) DEFAULT 0,
					inserted datetime DEFAULT NULL,
					updated datetime DEFAULT NULL,
					PRIMARY KEY (id), 
					KEY platform (clientid)
			   ) ENGINE=InnoDB DEFAULT CHARSET=latin1`

	_, err = a.Database.Query(tableClientsQuery)

	return err
}

// SetUpRouter func
func (a *App) SetUpRouter() error {
	a.Router = mux.NewRouter()

	//a.Router.HandleFunc("/push/ios", a.createPushIOS).Methods("POST")
	//a.Router.HandleFunc("/push/android", a.createPushAndroid).Methods("POST")

	// client handling
	a.Router.HandleFunc("/clients", a.clientList).Methods("GET")
	a.Router.HandleFunc("/client", a.clientCreate).Methods("POST")
	a.Router.HandleFunc("/client/{id:[0-9]+}", a.clientUpdate).Methods("PUT")
	a.Router.HandleFunc("/client/{id:[0-9]+}", a.clientDelete).Methods("DELETE")
	a.Router.HandleFunc("/client/{id:[0-9]+}", a.clientGet).Methods("GET")

	if strings.Contains(os.Args[0], "/_test/") {
		return nil
	}
	return http.ListenAndServe(":8080", a.Router)
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
