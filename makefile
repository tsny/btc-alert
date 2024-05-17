dev:
	go run cmd/alert/*.go

build:
	mkdir -p dist
	go build -o ./dist ./cmd/alert