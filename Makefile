swagger:
	swag init -g cmd/app/main.go -o docs

run: swagger
	go run cmd/app/main.go