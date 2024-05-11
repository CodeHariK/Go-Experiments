# gRPC Hello World

1. Get the code:
    go mod tidy

2. Run

go run greeter_server/main.go
go run greeter_client/main.go

3.Generate proto files

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative helloproto/helloworld.proto
