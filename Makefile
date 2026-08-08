run:
	docker-compose up --build --abort-on-container-exit
cover:
	go test -count=1 ./... -coverprofile cover.out && go tool cover -func cover.out
clean:
	rm cover.out && rm agent && rm server
build:
	go build -o agent cmd/agent/main.go && go build -o server cmd/server/main.go
gen:
	go generate ./...

all: gen build


