run: build
	@./bin/heyecho

build:
	@go build -o bin/heyecho cmd/app/main.go

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