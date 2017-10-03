package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/danfixeads/livepush/models"
	"github.com/danfixeads/livepush/senders"
	"github.com/gorilla/mux"
)

func (a *App) createPushIOS(w http.ResponseWriter, r *http.Request) {

	var mp models.MultiplePush
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&mp); err != nil {
		a.respondWithError(w, r, http.StatusBadRequest, models.ErrInvalidPayload.Error())
		return
	}
	defer r.Body.Close()

	var ios senders.IOS
	ios.ClientID = mp.ClientID.String
	ios.Push = mp
	if err := ios.GetClient(a.Database); err != nil {
		a.respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := ios.GetCertificate(); err != nil {
		a.respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	failed, err := ios.SendMessage(a.Database)
	// send webhook (if applicable)
	go sendWebHook(&ios.Client, failed, "ios")

	if err != nil {
		a.respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	a.respondWithJSON(w, r, http.StatusCreated, mp)
}

func (a *App) createPushAndroid(w http.ResponseWriter, r *http.Request) {

	var mp models.MultiplePush
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&mp); err != nil {
		a.respondWithError(w, r, http.StatusBadRequest, models.ErrInvalidPayload.Error())
		return
	}
	defer r.Body.Close()

	var android senders.Android
	android.ClientID = mp.ClientID.String
	android.Push = mp
	if err := android.GetClient(a.Database); err != nil {
		fmt.Print(err.Error())
		a.respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	failed, err := android.SendMessage(a.Database)
	// send webhook (if applicable)
	go sendWebHook(&android.Client, failed, "android")

	if err != nil {
		a.respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	a.respondWithJSON(w, r, http.StatusCreated, mp)
}

func (a *App) pushDelete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p := models.Push{ID: id}
	if err := p.Delete(a.Database); err != nil {
		a.respondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondWithJSON(w, r, http.StatusOK, map[string]string{"result": "success"})

}

func (a *App) pushList(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	limit, _ := strconv.Atoi(vars["limit"])
	start, _ := strconv.Atoi(vars["start"])

	if limit > 10 || limit < 1 {
		limit = 10
	}

	clients, _ := models.ListPushes(a.Database, start, limit)

	a.respondWithJSON(w, r, http.StatusOK, clients)

}

func (a *App) pushGet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p := models.Push{}
	if err := p.Get(a.Database, id); err != nil {
		switch err {
		default:
			a.respondWithError(w, r, http.StatusNotFound, "Push not found")
		}
		return
	}

	a.respondWithJSON(w, r, http.StatusOK, p)

}

func sendWebHook(c *models.Client, p []models.Push, s string) {

	// fmt.Print("-------- HUH --------")

	// if a webhook address exists, then use it to inform the client that something went wrong
	if len(p) > 0 && c.WebHook.Valid {

		// payload struct to send to the webhook receiver
		type payload struct {
			Service string
			Failed  []models.Push
		}

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(payload{Service: s, Failed: p})
		req, _ := http.NewRequest("POST", c.WebHook.String, b)

		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/x-www-form-urlencoded")

		res, _ := http.DefaultClient.Do(req)

		body, _ := ioutil.ReadAll(res.Body)
		//fmt.Println(res)
		defer res.Body.Close()
		fmt.Println(string(body))

	}
}
