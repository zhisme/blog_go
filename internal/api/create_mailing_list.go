package api

import (
	"encoding/json"
	"log"
	"net/http"
  "errors"
  "io"

	"backend-go/internal/dto"
  "backend-go/internal/api/handlers"
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
    json.NewEncoder(w).Encode(errorResponse)
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
  if err := json.NewEncoder(w).Encode(mailingList); err != nil {
    log.Default().Print(err)
  }
}
