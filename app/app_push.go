package app

import (
	"encoding/json"
	"net/http"

	"github.com/danfixeads/livepush/models"
	"github.com/danfixeads/livepush/senders"
)

func (a *App) createPushIOS(w http.ResponseWriter, r *http.Request) {

	var mp models.MultiplePush
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&mp); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
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

	if err := ios.SendMessage(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, mp)
}

func (a *App) createPushAndroid(w http.ResponseWriter, r *http.Request) {}
