# Grpc Gateway

go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

buf generate

    ```
        curl -X POST -k http://localhost:8090/v1/sayhello -d '{"name": "hello"}'

        curl -X GET -k http://localhost:8090/v1/users

        curl -X GET -k http://localhost:8090/v1/users/4bf5ef86-af0d-11ee-84f2-6e66b6bac5af

        curl -X 'POST' \
            'http://0.0.0.0:8090/v1/users' \
            -H 'accept: application/json' \
            -H 'Content-Type: application/json' \
            -d '{
            "id": "string"
            }'

            9a0a7aa6-af09-11ee-837b-6e66b6bac5af
    ```
