dev:
	@air

build:
	@go build -o bin/pro-backend cmd/api/main.go 

start: build
	@./bin/pro-backend -config config/local.yaml

seed: 
	@go build -o bin/pro-backend cmd/seed/main.go
	@./bin/pro-backend -config config/local.yaml

test:
	go test ./... -v