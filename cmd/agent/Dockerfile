FROM golang:1.21

WORKDIR /src
COPY cmd/agent/main.go cmd/agent/
COPY internal/agent internal/agent/
COPY internal/shared internal/shared/
COPY go.mod .
COPY go.sum .

RUN go mod download
CMD go run cmd/agent/main.go

