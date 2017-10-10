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

	payload := []byte(`{"clientid":"` + testClientID + `","fcmauthkey":"abc"}`)

	req, _ := http.NewRequest("POST", "/client", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] == 0 || m["id"] == nil {
		t.Errorf("Expected the ID to be greater than '0'. Got '%v' instead!", m["id"])
	}

	if m["clientid"] == "" || m["clientid"] == nil {
		t.Errorf("Expected the ClientID to contain a value. Got '%v' instead!", m["clientid"])
	}

}

func TestCreateClientWithMissingClientID(t *testing.T) {

	payload := []byte(`{"missing":1}`)

	req, _ := http.NewRequest("POST", "/client", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

func TestCreateClientWithInvalidPayload(t *testing.T) {

	payload := []byte(`{"clientid","` + testClientID + `"}`)

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

	payload := []byte(`{"clientid":"` + testClientID + `","passphrase":"abc","pemfile":"a_file.pem","bundleidentifier":"com.fixeads.imo"}`)

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

	req, _ := http.NewRequest("GET", "/clients", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected an array. Got %s", body)
	}

}

func TestListClientsEmptyResponse(t *testing.T) {

	clearTestClients()

	req, _ := http.NewRequest("GET", "/clients", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}

}

func TestListClientsPagination(t *testing.T) {

	clearTestClients()
	addTestClients(122)

	req, _ := http.NewRequest("GET", "/clients/2/7", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected an array. Got %s", body)
	}
}

func TestListClientsPaginationWithInvalidStartValue(t *testing.T) {

	clearTestClients()
	addTestClients(122)

	req, _ := http.NewRequest("GET", "/clients/0/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected an array. Got %s", body)
	}
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
		values[i] = fmt.Sprintf("('%v',NOW())", i+1)
	}

	var query = fmt.Sprintf("INSERT INTO client (clientid,inserted) VALUES %v", strings.Join(values, ", "))
	//log.Print(query)
	a.Database.Exec(query)
}

func addTestClientValues() {
	_, err := a.Database.Exec("INSERT INTO `client` (`id`, `clientid`, `pemfile`, `p12file`, `passphrase`, `bundleidentifier`, `usesandboxios`, `fcmauthkey`, `webhook`, `active`, `inserted`, `updated`) VALUES (2, '" + testClientID + "', 'dev_imo.pem', NULL, 'bragaBoss8', 'com.fixeads.imo.Imovirtual', 1, 'AAAAZnheqwk:APA91bHMkh--KR3BYy-l6WX5cRzjelGGskcJy0p-LFnWdP0AsAc7HGvvmE7Aih6MwVd1ObMNkfpbu4vMYoABi5Y25cP2-c09wHOhkQWh-03XreyCXt-AKYCqKo6hY9Ru34iumLP5DQY5', 'http://mockbin.org/bin/42fc077c-2536-4ce3-b1ea-f3da3aa90971', 1, '2017-09-25 11:20:57', '2017-09-25 14:44:05')")
	if err != nil {
		panic(err)
	}
}

func addTestClientValuesIncorrectCertificates() {
	_, err := a.Database.Exec("INSERT INTO `client` (`id`, `clientid`, `pemfile`, `passphrase`, `bundleidentifier`, `usesandboxios`, `fcmauthkey`, `active`, `inserted`, `updated`) VALUES (2, '" + testClientID + "', 'incorrect.pem', 'rubbish', 'com.fixeads.imo.Imovirtual', 1, 'AAAAZnheqwk:APA91bHMkh--KR3BYy-l6WX5cRzjelGGskcJy0p-LFnWdP0AsAc7HGvvmE7Aih6MwVd1ObMNkfpbu4vMYoABi5Y25cP2-c09wHOhkQWh-03XreyCXt-AKYCqKo6hY9Ru34iumLP5DQY5', 1, '2017-09-25 11:20:57', '2017-09-25 14:44:05')")
	if err != nil {
		panic(err)
	}
}
