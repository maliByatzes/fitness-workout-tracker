.PHONY: up reset stop start server test pgcli

DB_URL="postgresql://admin:secret@localhost:5432/fwtdb?sslmode=disable"

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

create-migration:
	migrate create -ext sql -dir postgres/migrations -seq $(name)

run-migration:
	migrate -database $(DB_URL) -path postgres/migrations up

down-migration:
	migrate -database $(DB_URL) -path postgres/migrations down $(v)

force-migration:
	migrate -database $(DB_URL) -path postgres/migrations force $(version)

pgcli:
	pgcli -h localhost -p 5432 -U admin -W -d fwtdb
