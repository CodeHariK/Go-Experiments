info:
	@echo openapi
	@echo proto
	@echo makehat
	@echo server
	@echo client

openapi:
	swag init --parseDependency --parseInternal --ot json --dir cmd/server

proto:
	protoc --go_out=. --twirp_out=. rpc/haberdasher/service.proto

makehat:
	curl --request "POST" \
		--header "Content-Type: application/json" \
		--data '{"inches": 12}' \
		http://localhost:8080/dash/twirphat.haberdasher.Haberdasher/MakeHat

server:
	go run cmd/server/main.go

client:
	go run cmd/client/main.go
