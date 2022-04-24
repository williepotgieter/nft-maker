dev:
	go run app/*.go

swagger:
	swag init -g app/server.go