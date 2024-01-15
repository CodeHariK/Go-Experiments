FROM golang:alpine3.18

ENV PROJECT_DIR=/app GO111MODULE=on CGO_ENABLED=0

WORKDIR /app
RUN mkdir "/build"
COPY . .

EXPOSE 3000

RUN go build -o /build/app

CMD ["/build/app"]