run: build
	@./bin/heyecho

build: templ
	@go build -o bin/heyecho main.go

templ:
	@templ generate

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