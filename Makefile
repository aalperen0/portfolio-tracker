include .env

# =====================================================================================

# HELPER

# =====================================================================================

## help: to get information print this help message with 'make help'
.PHONY: help
help:
	@echo 'Usage: '
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ] 


# =====================================================================================

# BUILD

# =====================================================================================

## build/app: build the cmd/app application
.PHONY: build/app
build/api:
	@echo 'Building cmd/app...'
	@go build -o=./bin/app ./cmd/app



# =====================================================================================

# DEVELOPMENT

# =====================================================================================

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@echo 'Running application...'
	@go run ./cmd/app/main.go



## docker/up: docker compose up
.PHONY: docker/up
docker/up:
	@echo 'Starting Application...'
	@docker compose up --build 


## docker/down: docker compose down
.PHONY: docker/down
docker/down:
	@echo 'Stopping Application...'
	@docker compose down 

## docker/db: open db in terminal
.PHONY: docker/db
docker/db:
	@echo "Opening db from terminal..."
	docker compose exec postgres psql -U ptracker -d ptracker

## db/migrations/new: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo "Creating migration files for ${name}..."
	migrate create -seq -ext sql -dir ./migrations ${name}


## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo "Running up migrations..."
	docker run -v "${MIGRATION_PATH}":/migrations --network host --env-file .env migrate/migrate -path=/migrations/ -database "${DB_DSN_ZD}" up


## db/migrations/down: apply down database migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo "Down version of migrations..."
	docker run -v "${MIGRATION_PATH}":/migrations --network host --env-file .env migrate/migrate -path=/migrations/ -database "${DB_DSN_ZD}" down 1
 

# ============================================================================== #

# QUALITY 

# ============================================================================== #

## audit: tidy and vendor dependecies and format, vet and test all code
.PHONY: audit 
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running test...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies'
	go mod vendor
