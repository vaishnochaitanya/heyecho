run: build
	@./bin/heyecho

build: templ
	@go build -o bin/heyecho cmd/app/main.go

templ:
	@templ generate

up:
	@go run cmd/migrate/main.go up

down:
	@go run cmd/migrate/main.go down

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

drop:
	@go run cmd/drop/main.go

seed:
	@go run cmd/seed/main.go