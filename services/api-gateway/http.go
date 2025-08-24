package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	defer r.Body.Close()

	var reqBody PreviewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed_to_parse_josn_data", http.StatusBadRequest)
		return
	}

	jsonBody, _ := json.Marshal(reqBody)
	reader := bytes.NewReader(jsonBody)

	// tripService, err := grpcclient.NewTripServiceClient()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer tripService.Close()

	res, err := http.Post("http://trip-service:8083/trip/preview", "application/json", reader)
	if err != nil {
		log.Println(err)
		return
	}

	defer res.Body.Close()
	var tripRes any
	if err := json.NewDecoder(res.Body).Decode(&tripRes); err != nil {
		http.Error(w, "failed_to_parse_trip_service_response", http.StatusBadGateway)
		return
	}

	response := contracts.APIResponse{
		Data: tripRes,
	}

	writeJSON(w, http.StatusCreated, response)
}
