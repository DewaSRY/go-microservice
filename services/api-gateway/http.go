package main

import (
	"encoding/json"
	"log"
	"net/http"
	grpcclient "ride-sharing/services/api-gateway/grpc_client"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reqBody PreviewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed_to_parse_josn_data", http.StatusBadRequest)
		return
	}

	tripService, err := grpcclient.NewTripServiceClient()

	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), reqBody.mappedToProto())

	if err != nil {
		log.Printf("failed_to_preview_a_trip %v", err)
		return
	}

	response := contracts.APIResponse{
		Data: tripPreview,
	}

	writeJSON(w, http.StatusCreated, response)
}

func handleTripSatart(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reqBody startTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed_to_parse_josn_data", http.StatusBadRequest)
		return
	}

	tripService, err := grpcclient.NewTripServiceClient()

	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	trip, err := tripService.Client.CreateTrip(r.Context(), reqBody.toProto())

	if err != nil {
		log.Printf("failed_to_preview_a_trip %v", err)
		return
	}

	response := contracts.APIResponse{
		Data: trip,
	}

	writeJSON(w, http.StatusCreated, response)
}
