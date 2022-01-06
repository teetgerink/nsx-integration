MAKEFLAGS := --no-print-directory --silent

default: help

help:
	@echo "Please use 'make <target>' where <target> is one of"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z\._-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

r: run
run: fmt ## Run the program, alias: r
	@echo "Initializing swagger documentation"
	swag init --output swagger --parseDependency --parseDepth=4
	HTTP_SERVER_ADDRESS=0.0.0.0:8080 \
	DATABASE_CONNECTION_URL="host=localhost port=5432 user=postgres dbname=nsx_api password=changeme" \
		go run main.go

dr: docker.run
docker.run: ## Run the containers in docker
	cd docker-compose && docker-compose -p nsxapi up -d
	@echo "---------------------------------------------"
	@echo "NSX API      http://localhost:8088"
	@echo "---------------------------------------------"


stop: docker.stop
docker.stop: ## Stop the docker containers	
	cd docker-compose && docker-compose -p nsxapi down

b: docker.build
docker.build: ## Build the docker container
	# Read the access token from the user
	@docker build . 

t: test
test: fmt ## Run unit tests, alias: t
	@echo "Running tests"
	go test ./... --cover -timeout=70s -parallel=8

fmt: ## Format go code
	@go mod tidy
	@go fmt ./...

tools: ## Fetch the required tools
	GO111MODULE=off go get -u github.com/swaggo/swag/cmd/swag
