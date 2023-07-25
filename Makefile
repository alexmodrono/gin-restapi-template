build:
	go build -o bin/api cmd/gin-restapi-template/main.go

run:
	go run cmd/gin-restapi-template/main.go

test:
	go test ./...
