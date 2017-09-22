package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// -----------------------
// CREATE ----------------
// -----------------------

func TestCreateClient(t *testing.T) {

	payload := []byte(`{"clientid":1,"fcmtoken":"abc"}`)

	req, _ := http.NewRequest("POST", "/client", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] == 0 || m["id"] == nil {
		t.Errorf("Expected the ID to be greater than '0'. Got '%v' instead!", m["id"])
	}

	if m["clientid"] == 0 || m["clientid"] == nil {
		t.Errorf("Expected the ClientID to be greater than '0'. Got '%v' instead!", m["clientid"])
	}

}

func TestCreateClientWithMissingClientID(t *testing.T) {

	payload := []byte(`{"missing":1}`)

	req, _ := http.NewRequest("POST", "/client", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

func TestCreateClientWithInvalidPayload(t *testing.T) {

	payload := []byte(`{"clientid",1}`)

	req, _ := http.NewRequest("POST", "/client", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

// -----------------------
// UPDATE ----------------
// -----------------------

func TestUpdateClient(t *testing.T) {

	clearTestClients()
	addTestClients(12)

	req, _ := http.NewRequest("GET", "/client/2", nil)
	response := executeRequest(req)
	var originalClient map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalClient)

	payload := []byte(`{"clientid":22,"passphrase":"abc","pemfile":"a_file.pem"}`)

	req, _ = http.NewRequest("PUT", "/client/2", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalClient["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalClient["id"], m["id"])
	}

	if m["clientid"] == originalClient["clientid"] {
		t.Errorf("Expected the clientid to change from '%v' to '%v'. Got '%v'", originalClient["clientid"], m["clientidname"], m["clientid"])
	}

}

func TestUpdateClientWithMissingClientID(t *testing.T) {

	clearTestClients()
	addTestClients(12)

	req, _ := http.NewRequest("GET", "/client/2", nil)
	response := executeRequest(req)
	var originalClient map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalClient)

	payload := []byte(`{"passphrase":"abc","pemfile":"a_file.pem"}`)

	req, _ = http.NewRequest("PUT", "/client/2", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusInternalServerError, response.Code)
}

func TestUpdateClientWithInvalidPayload(t *testing.T) {

	clearTestClients()
	addTestClients(12)

	req, _ := http.NewRequest("GET", "/client/2", nil)
	response := executeRequest(req)
	var originalClient map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalClient)

	payload := []byte(`{"passphrase","abc","pemfile","a_file.pem"}`)

	req, _ = http.NewRequest("PUT", "/client/2", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

// -----------------------
// GET -------------------
// -----------------------

func TestGetClient(t *testing.T) {

	clearTestClients()
	addTestClients(12)

	req, _ := http.NewRequest("GET", "/client/9", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] == 0 || m["id"] == nil {
		t.Errorf("Expected the ID to be greater than '0'. Got '%v' instead!", m["id"])
	}
}

func TestGetClientWithInvalidID(t *testing.T) {

	clearTestClients()
	addTestClients(1)

	req, _ := http.NewRequest("GET", "/client/123", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// -----------------------
// DELETE ----------------
// -----------------------

func TestDeleteClient(t *testing.T) {

	clearTestClients()
	addTestClients(12)

	req, _ := http.NewRequest("DELETE", "/client/8", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteClientWithInvalidID(t *testing.T) {

	clearTestClients()
	addTestClients(12)

	req, _ := http.NewRequest("DELETE", "/client/501", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusInternalServerError, response.Code)
}

// -----------------------
// LIST ------------------
// -----------------------

func TestListClients(t *testing.T) {

	clearTestClients()
	addTestClients(122)

}

// -----------------------
// HELPERS ---------------
// -----------------------

func clearTestClients() {
	_, err := a.Database.Exec("TRUNCATE clients")
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
		values[i] = fmt.Sprintf("('%v',NOW())", i+1)
	}

	var query = fmt.Sprintf("INSERT INTO clients (clientid,inserted) VALUES %v", strings.Join(values, ", "))
	//log.Print(query)
	a.Database.Exec(query)
}
