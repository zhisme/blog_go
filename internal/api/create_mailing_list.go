package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"backend-go/internal/api/handlers"
	"backend-go/internal/dto"
)

func (s *Server) createMailingList(w http.ResponseWriter, r *http.Request) {
	var newMailingList dto.MailingList

	err := json.NewDecoder(r.Body).Decode(&newMailingList)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		var msg string
		if errors.Is(err, io.EOF) {
			msg = "request body is empty"
		} else {
			msg = "invalid JSON: " + err.Error()
		}

		errorResponse := map[string]map[string]string{
			"error": {
				"message": msg,
			},
		}
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			log.Default().Print(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	mailingList, err := handlers.HandleCreate(newMailingList)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := map[string]map[string]string{
			"error": {
				"message": err.Error(),
			},
		}
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			log.Default().Print(err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(mailingList); err != nil {
		log.Default().Print(err)
	}
}
