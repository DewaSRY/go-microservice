package http

import (
	"encoding/json"
	"log"
	"net/http"
	grpcClient "ride-sharing/services/api-gateway/grpc_client"
	"ride-sharing/shared/contracts"

	"ride-sharing/services/api-gateway/pkg/types"
	"ride-sharing/services/api-gateway/pkg/util"
)

func HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reqBody types.PreviewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed_to_parse_json_data", http.StatusBadRequest)
		return
	}

	tripService, err := grpcClient.NewTripServiceClient()

	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), reqBody.ToTripProtoTripRequest())

	if err != nil {
		log.Printf("failed_to_preview_a_trip %v", err)
		return
	}

	response := contracts.APIResponse{
		Data: tripPreview,
	}

	util.WriteJSON(w, http.StatusCreated, response)
}

func HandleTripStart(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reqBody types.StartTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed_to_parse_json_data", http.StatusBadRequest)
		return
	}

	tripService, err := grpcClient.NewTripServiceClient()

	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	trip, err := tripService.Client.CreateTrip(r.Context(), reqBody.ToTripProtoCreateTrip())

	if err != nil {
		log.Printf("failed_to_preview_a_trip %v", err)
		return
	}

	response := contracts.APIResponse{
		Data: trip,
	}

	util.WriteJSON(w, http.StatusCreated, response)
}
