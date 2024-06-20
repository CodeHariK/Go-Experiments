proto:
	protoc --go_out=. --twirp_out=. rpc/haberdasher/service.proto

makehat:
	curl --request "POST" \
		--header "Content-Type: application/json" \
		--data '{"inches": 12}' \
		http://localhost:8080/twirp/twirphat.haberdasher.Haberdasher/MakeHat
