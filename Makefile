.PHONY: up reset stop start server test

up:
	docker compose up -d
	sleep 3

reset:
	docker compose down
	make up

stop:
	docker compose stop

start:
	docker compose start

server:
	go run cmd/fwt/main.go

test:
	go test ./...
