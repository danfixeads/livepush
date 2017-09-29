package app_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// -----------------------
// GET -------------------
// -----------------------

func TestGetPush(t *testing.T) {

	clearTestPushes()
	addTestPushes(12)

	req, _ := http.NewRequest("GET", "/push/9", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] == 0 || m["id"] == nil {
		t.Errorf("Expected the ID to be greater than '0'. Got '%v' instead!", m["id"])
	}
}

func TestGetPushWithInvalidID(t *testing.T) {

	clearTestPushes()
	addTestPushes(1)

	req, _ := http.NewRequest("GET", "/push/123", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// -----------------------
// DELETE ----------------
// -----------------------

func TestDeletePush(t *testing.T) {

	clearTestPushes()
	addTestPushes(12)

	req, _ := http.NewRequest("DELETE", "/push/8", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeletePushWithInvalidID(t *testing.T) {

	clearTestPushes()
	addTestPushes(12)

	req, _ := http.NewRequest("DELETE", "/push/501", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusInternalServerError, response.Code)
}

// -----------------------
// LIST ------------------
// -----------------------

func TestListPushes(t *testing.T) {

	clearTestPushes()
	addTestPushes(122)

	req, _ := http.NewRequest("GET", "/pushes", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected an array. Got %s", body)
	}

}

func TestListPushesEmptyResponse(t *testing.T) {

	clearTestPushes()

	req, _ := http.NewRequest("GET", "/pushes", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}

}

func TestListPushesPagination(t *testing.T) {

	clearTestPushes()
	addTestPushes(122)

	req, _ := http.NewRequest("GET", "/pushes/2/7", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected an array. Got %s", body)
	}
}

func TestListPushesPaginationWithInvalidStartValue(t *testing.T) {

	clearTestPushes()
	addTestPushes(122)

	req, _ := http.NewRequest("GET", "/pushes/0/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected an array. Got %s", body)
	}
}

// -----------------------
// HELPERS ---------------
// -----------------------

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
