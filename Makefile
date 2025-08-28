PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .

.PHONY: generate-proto
generate-proto:
	@mkdir -p $(GO_OUT)
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)


# Run trip-service directly
run-trip-service:
	go run ./services/trip-service/cmd/main.go

run-test: 
	go test -v ./...

# build commend
build-driver-service:
	go build -o build/driver-service ./services/driver-service/cmd/main.go
build-trip-service:
	go build -o build/trip-service ./services/trip-service/cmd/main.go
build-getway:
	go build -o build/api-gateway ./services/api-gateway