package api

import (
	"encoding/json"
	"log"
	"net/http"

	"backend-go/internal/dto"
  "backend-go/internal/api/handlers"
)

func (s *Server) createMailingList(w http.ResponseWriter, r *http.Request) {
	var newMailingList dto.MailingList

	if err := json.NewDecoder(r.Body).Decode(&newMailingList); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

  w.Header().Set("Content-Type", "application/json")

  err, mailingList := handlers.HandleCreate(newMailingList)
  if err != nil {
		w.WriteHeader(http.StatusBadRequest)
    errorResponse := map[string]map[string]string{
      "error": {
        "message": err.Error(),
      },
    }
		json.NewEncoder(w).Encode(errorResponse)
		return
  }

	w.WriteHeader(http.StatusCreated)
  err = json.NewEncoder(w).Encode(mailingList)
	if err != nil {
		log.Default().Print(err)
	}
}
