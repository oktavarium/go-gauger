FROM golang:1.21

WORKDIR /src
COPY cmd/server/main.go cmd/server/
COPY internal/server internal/server/
COPY internal/shared internal/shared/
COPY go.mod .
COPY go.sum .

RUN go mod download
CMD go run cmd/server/main.go

