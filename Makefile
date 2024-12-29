run: build
	@CGO_ENABLED=1 ./bin/heyecho

build: templ sqlc
	@go build -o bin/heyecho main.go

templ:
	@templ generate

sqlc:
	@sqlc generate

up:
	@go run cmd/migrate/up.go up

down:
	@go run cmd/migrate/down.go down

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

drop:
	@go run cmd/drop/drop.go

seed:
	@go run cmd/seed/seed.go