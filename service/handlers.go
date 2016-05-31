package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func updateRealityHandler(formatter *render.Render, repo realityRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		gameID := vars["gameId"]
		payload, _ := ioutil.ReadAll(req.Body)
		var newReality reality
		err := json.Unmarshal(payload, &newReality)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse game state: %s\n", err.Error()))
			return
		}

		err = repo.updateReality(gameID, newReality)
		if err != nil {
			formatter.Text(w, http.StatusInternalServerError, fmt.Sprintf("Error saving updated state: %s\n", err.Error()))
			return
		}

		formatter.JSON(w, http.StatusOK, nil)
		return
	}
}

func getRealityHandler(formatter *render.Render, repo realityRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		gameID := vars["gameId"]

		gameReality, err := repo.getReality(gameID)
		if err != nil {
			fmt.Println(err.Error())
			formatter.Text(w, http.StatusNotFound, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, &gameReality)
		return
	}
}
