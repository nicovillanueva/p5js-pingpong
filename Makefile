all:
	@echo "available: start, swagger"

start: swagger
	go run cli/pingpong-server/main.go

swagger:
	swag init -g server/server.go -o server/docs

