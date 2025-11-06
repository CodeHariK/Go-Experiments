FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /booktest

FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM alpine:latest

WORKDIR /

COPY --from=build-stage /booktest /booktest

EXPOSE 8080

CMD ["./booktest"]