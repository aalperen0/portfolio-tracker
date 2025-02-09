include .env


# =====================================================================================

# HELPER

# =====================================================================================

## help: print this help message
.PHONY: help
help:
	@echo 'Usage: '
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y] 


# =====================================================================================

# BUILD

# =====================================================================================

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	@go build -o=./bin/api ./cmd/api



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




# ============================================================================== #

# QUALITY CONTROL

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
	@echo 'Tidying and verifying module dependecies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies'
	go mod vendor
