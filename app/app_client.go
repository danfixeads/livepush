package app

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danfixeads/livepush/models"
	"github.com/gorilla/mux"
)

func (a *App) clientCreate(w http.ResponseWriter, r *http.Request) {

	var c models.Client
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		a.respondWithError(w, r, http.StatusBadRequest, models.ErrInvalidPayload.Error())
		return
	}
	defer r.Body.Close()

	if err := c.Create(a.Database); err != nil {
		switch err {
		default:
			a.respondWithError(w, r, http.StatusBadRequest, err.Error())
		}
		return
	}

	a.respondWithJSON(w, r, http.StatusCreated, c)

}

func (a *App) clientUpdate(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var c models.Client
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		a.respondWithError(w, r, http.StatusBadRequest, models.ErrInvalidPayload.Error())
		return
	}
	defer r.Body.Close()
	c.ID = id

	if err := c.Update(a.Database); err != nil {
		a.respondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondWithJSON(w, r, http.StatusOK, c)

}

func (a *App) clientDelete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	c := models.Client{ID: id}
	if err := c.Delete(a.Database); err != nil {
		a.respondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondWithJSON(w, r, http.StatusOK, map[string]string{"result": "success"})

}

func (a *App) clientList(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	limit, _ := strconv.Atoi(vars["limit"])
	start, _ := strconv.Atoi(vars["start"])

	if limit > 10 || limit < 1 {
		limit = 10
	}

	clients, _ := models.ListClients(a.Database, start, limit)

	a.respondWithJSON(w, r, http.StatusOK, clients)

}

func (a *App) clientGet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	c := models.Client{}
	if err := c.Get(a.Database, id); err != nil {
		switch err {
		default:
			a.respondWithError(w, r, http.StatusNotFound, "Client not found")
		}
		return
	}

	a.respondWithJSON(w, r, http.StatusOK, c)
}
