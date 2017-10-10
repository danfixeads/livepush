package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

var testAuthorization = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzdHMiLCJleHAiOjc1MDcxMzE5OTYsImlhdCI6MTUwNzEzMTk5NiwiaXNzIjoiYWRtaW5Ac2VydmljZXMub2x4LmNvbSIsIm5iZiI6MTUwNzEzMTk5Nn0.Zm8pSlZahvWPhsl9eoVhAKgEtezjn-ht5H91XmE8W3pKenTkYT4bC-bn97dO-6NACQ43RTtNIlU327cBlpAuM22kUWVR2qd5X6Bo3-PeRb39s_uJNTDcngBxzwLaAItePVUa5fEcSz4_PTn8jVW4-m3K9kq0_Ql3rElvjZVT_6c"

// -----------------------
// CREATE (IOS)-----------
// -----------------------

func TestCreatePushIOS(t *testing.T) {

	clearTestPushes()
	clearTestClients()
	addTestClientValues()

	payload := []byte(`{
		"tokens": ["f7534f19f6103e1a7ee26de615f2b8c8d3eeb63dc3da9922388ebfbf2b4d7717",
			"944f5c35533d770566901cf533aceef1111d8bd86d8e081e4f3603ddbb928875",
			"e9d03c7b63f950944eb5e34e4b875d7ad4918bc9ca71926afe11b1e30ec235c3",
			"87242750b176aeef6cd0c342b485560c46c092786ae17136404541855a8a5c59",
			"rubbish",
			"06bafb503c7168c57d5db46d238d76479f1a3466a808dff34c9e6e1c2834a6ff",
			"e9a736109990723b7bf9a33c7fb566601c375a8e64dacd07442a399361246f82"
		],
	
		"payload": {
			"aps": {
				"alert": {
					"title": "A test message",
					"subtitle": "with a subtitle",
					"body": "and a very nice body"
				},
				"badge": 1,
				"mutable-content": 1,
				"content-available": 1
			},
			"data": {
				"attachment-url": "https://media.mnn.com/assets/images/2015/08/sunset-sunrise-tips-water-reflections0.jpg.990x0_q80_crop-smart.jpg"
			}
		}
	}`)

	req, _ := http.NewRequest("POST", "/push/ios", bytes.NewBuffer(payload))
	req.Header.Add("authorization", testAuthorization)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

}

func TestCreatePushIOSWithInvalidPayload(t *testing.T) {

	payload := []byte(`{"clientid",1}`)

	req, _ := http.NewRequest("POST", "/push/ios", bytes.NewBuffer(payload))
	req.Header.Add("authorization", testAuthorization)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

func TestCreatePushIOSWithMissingClientID(t *testing.T) {

	payload := []byte(`{"missing":1}`)

	req, _ := http.NewRequest("POST", "/push/ios", bytes.NewBuffer(payload))
	req.Header.Add("authorization", testAuthorization)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

func TestCreatePushIOSWithIncorrectCertificates(t *testing.T) {

	clearTestClients()
	addTestClientValuesIncorrectCertificates()

	payload := []byte(`{
		"tokens": [
			"rubbish"
		],
	
		"payload": {
			"aps": {
				"alert": {
					"title": "A test message",
					"subtitle": "with a subtitle",
					"body": "and a very nice body"
				},
				"badge": 1,
				"mutable-content": 1,
				"content-available": 1
			},
			"data": {
				"attachment-url": "https://media.mnn.com/assets/images/2015/08/sunset-sunrise-tips-water-reflections0.jpg.990x0_q80_crop-smart.jpg"
			}
		}
	}`)

	req, _ := http.NewRequest("POST", "/push/ios", bytes.NewBuffer(payload))
	req.Header.Add("authorization", testAuthorization)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

func TestCreatePushIOSWithInvalidDevices(t *testing.T) {

	clearTestClients()
	addTestClientValues()

	payload := []byte(`{
		"tokens": [
			"expired_token...",
			"rubbish"
		],
	
		"payload": {
			"aps": {
				"alert": {
					"title": "A test message",
					"subtitle": "with a subtitle",
					"body": "and a very nice body"
				},
				"badge": 1,
				"mutable-content": 1,
				"content-available": 1
			},
			"data": {
				"attachment-url": "https://media.mnn.com/assets/images/2015/08/sunset-sunrise-tips-water-reflections0.jpg.990x0_q80_crop-smart.jpg"
			}
		}
	}`)

	req, _ := http.NewRequest("POST", "/push/ios", bytes.NewBuffer(payload))
	req.Header.Add("authorization", testAuthorization)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

// -----------------------
// CREATE (Android)-------
// -----------------------

func TestCreatePushAndroid(t *testing.T) {

	clearTestPushes()
	clearTestClients()
	addTestClientValues()

	payload := []byte(`{
			"tokens": ["fRE69G6iGx0:APA91bGJZBlY-2Ljor-WeDEWZghcA0yY5SC5pJeNtQp_OHnlktCy_2uQTacceaRUp5ieIiW6CLk6DXndBJeAReHLVvV1DgA4cpOyUaBU0Wb6CNJ86vOo9RnG0U9h9PFuAdi4nSNbc1qH",
				"rubbish",
				"dCB_XXqF-NU:APA91bGbqW5v_qd9gaAaVvhITgsohGhUHIp-pHxxFMAzSmvRKIqJPnjMZcqMVAZX4O8PSW9iYZcd-JRHSNKMf0Mb9JWxYY1llOtxN0dx1_fhxSjKPo0-SuObdfqPw3ZpNo7_AndKSq7P"
			],
			"payload": {
				"message": {
					"alert": "message",
					"sound": 1,
					"data": {
						"type": 0,
						"actions": {
							"main": "/"
						}
					}
				}
			}
		}`)

	req, _ := http.NewRequest("POST", "/push/android", bytes.NewBuffer(payload))
	req.Header.Add("authorization", testAuthorization)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

}

func TestCreatePushAndroidWithInvalidPayload(t *testing.T) {

	payload := []byte(`{"clientid","` + testClientID + `"}`)

	req, _ := http.NewRequest("POST", "/push/android", bytes.NewBuffer(payload))
	req.Header.Add("authorization", testAuthorization)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

func TestCreatePushAndroidWithMissingClientID(t *testing.T) {

	payload := []byte(`{"missing":1}`)

	req, _ := http.NewRequest("POST", "/push/android", bytes.NewBuffer(payload))
	req.Header.Add("authorization", testAuthorization)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

func TestCreatePushAndroidWithInvalidDevices(t *testing.T) {

	clearTestClients()
	addTestClientValues()

	payload := []byte(`{
			"tokens": [
				"expired_token...",
				"rubbish"
			],
		
			"payload": {
				"message": {
					"alert": "message",
					"sound": 1,
					"data": {
						"type": 0,
						"actions": {
							"main": "/"
						}
					}
				}
			}
		}`)

	req, _ := http.NewRequest("POST", "/push/android", bytes.NewBuffer(payload))
	req.Header.Add("authorization", testAuthorization)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

// -----------------------
// GET -------------------
// -----------------------

func TestGetPush(t *testing.T) {

	clearTestPushes()
	addTestPushes(12)

	req, _ := http.NewRequest("GET", "/push/9", nil)
	req.Header.Add("authorization", testAuthorization)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] == 0 || m["id"] == nil {
		t.Errorf("Expected the ClientID to contain a value. Got '%v' instead!", m["id"])
	}
}

func TestGetPushWithInvalidID(t *testing.T) {

	clearTestPushes()
	addTestPushes(1)

	req, _ := http.NewRequest("GET", "/push/123", nil)
	req.Header.Add("authorization", testAuthorization)
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
	req.Header.Add("authorization", testAuthorization)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeletePushWithInvalidID(t *testing.T) {

	clearTestPushes()
	addTestPushes(12)

	req, _ := http.NewRequest("DELETE", "/push/501", nil)
	req.Header.Add("authorization", testAuthorization)
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
	req.Header.Add("authorization", testAuthorization)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected an array. Got %s", body)
	}

}

func TestListPushesEmptyResponse(t *testing.T) {

	clearTestPushes()

	req, _ := http.NewRequest("GET", "/pushes", nil)
	req.Header.Add("authorization", testAuthorization)
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
	req.Header.Add("authorization", testAuthorization)
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
	req.Header.Add("authorization", testAuthorization)
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
		values[i] = fmt.Sprintf("('%s','token_%v','ios',NOW())", testClientID, i+1)
	}

	var query = fmt.Sprintf("INSERT INTO push (clientid,token,platform,inserted) VALUES %v", strings.Join(values, ", "))
	//log.Print(query)
	a.Database.Exec(query)
}
