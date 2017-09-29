package app

import (
	"encoding/json"
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
		respondWithError(w, http.StatusBadRequest, models.ErrInvalidPayload.Error())
		return
	}
	defer r.Body.Close()

	var ios senders.IOS
	ios.ClientID = int(mp.ClientID.Int64)
	ios.Push = mp
	if err := ios.GetClient(a.Database); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ios.GetCertificate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ios.SendMessage(a.Database); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, mp)
}

func (a *App) createPushAndroid(w http.ResponseWriter, r *http.Request) {

	var mp models.MultiplePush
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&mp); err != nil {
		respondWithError(w, http.StatusBadRequest, models.ErrInvalidPayload.Error())
		return
	}
	defer r.Body.Close()

	var android senders.Android
	android.ClientID = int(mp.ClientID.Int64)
	android.Push = mp
	if err := android.GetClient(a.Database); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := android.SendMessage(a.Database); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, mp)
}

func (a *App) pushDelete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p := models.Push{ID: id}
	if err := p.Delete(a.Database); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

func (a *App) pushList(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	limit, _ := strconv.Atoi(vars["limit"])
	start, _ := strconv.Atoi(vars["start"])

	if limit > 10 || limit < 1 {
		limit = 10
	}

	clients, _ := models.ListPushes(a.Database, start, limit)

	respondWithJSON(w, http.StatusOK, clients)

}

func (a *App) pushGet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p := models.Push{}
	if err := p.Get(a.Database, id); err != nil {
		switch err {
		default:
			respondWithError(w, http.StatusNotFound, "Push not found")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)

}
