build:
	go build -o bin/api cmd/gin-restapi-template/main.go

run:
	go run cmd/gin-restapi-template/main.go

.PHONY: test
test:
	go test $(if $(VERBOSE),-v,) ./pkg/...
