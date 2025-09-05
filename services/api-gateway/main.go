package main

import (
	"log"
	"net/http"

	internalHttp "ride-sharing/services/api-gateway/internal/infrastructure/http"
	"ride-sharing/services/api-gateway/internal/infrastructure/ws"
	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	log.Println("Starting API Gateway")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from API Gateway"))
	})

	mux.HandleFunc("POST /trip/preview", enableCORS(internalHttp.HandleTripPreview))
	mux.HandleFunc("POST /trip/start", enableCORS(internalHttp.HandleTripStart))

	mux.HandleFunc("/ws/riders", ws.HandlerRidersWebSocket)
	mux.HandleFunc("/ws/drivers", ws.HandleDriverWebSocket)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("Http server error : %v  ", err)
	} else {
		log.Println("Server run on port : 8081")
	}
}
