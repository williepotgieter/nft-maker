dev:
	go run app/*.go

build:
	go build -o bin/nftmaker app/*.go

swagger:
	swag init -g app/server.go