package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	defer r.Body.Close()

	var reqBody PreviewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed_to_parse_josn_data", http.StatusBadRequest)
		return
	}

	log.Println("SUCCESS")
	writeJSON(w, http.StatusCreated, reqBody)
}
