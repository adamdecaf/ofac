VERSION := $(shell grep -Eo '(v[0-9]+[\.][0-9]+[\.][0-9]+(-[a-zA-Z0-9]*)?)' version.go)

.PHONY: build build-server build-examples docker release check

build: check build-server build-examples

build-server:
	CGO_ENABLED=1 go build -o ./bin/server github.com/moov-io/ofac/cmd/server

build-examples: build-search-example build-webhook-example

build-search-example:
	CGO_ENABLED=0 go build -o ./bin/search-example github.com/moov-io/ofac/examples/search

build-webhook-example:
	CGO_ENABLED=0 go build -o ./bin/webhook-example github.com/moov-io/ofac/examples/webhook

check:
	go fmt ./...
	@mkdir -p ./bin/

.PHONY: client
client:
# Versions from https://github.com/OpenAPITools/openapi-generator/releases
	@chmod +x ./openapi-generator
	@rm -rf ./client
	OPENAPI_GENERATOR_VERSION=4.0.0-beta2 ./openapi-generator generate -i openapi.yaml -g go -o ./client
	go fmt ./client
	go build github.com/moov-io/ofac/client
	go test ./client

.PHONY: clean
clean:
	@rm -rf client/
	@rm -f openapi-generator-cli-*.jar

docker:
# Main OFAC server Docker image
	docker build --pull -t moov/ofac:$(VERSION) -f Dockerfile .
	docker tag moov/ofac:$(VERSION) moov/ofac:latest
# search example
	docker build --pull -t moov/ofac-search-example:$(VERSION) -f ./examples/search/Dockerfile .
	docker tag moov/ofac-search-example:$(VERSION) moov/ofac-search-example:latest
# webhook example
	docker build --pull -t moov/ofac-webhook-example:$(VERSION) -f ./examples/webhook/Dockerfile .
	docker tag moov/ofac-webhook-example:$(VERSION) moov/ofac-webhook-example:latest

release: docker AUTHORS
	go vet ./...
	go test -coverprofile=cover-$(VERSION).out ./...
	git tag -f $(VERSION)

release-push:
	docker push moov/ofac:$(VERSION)
	docker push moov/ofac-webhook-example:$(VERSION)

.PHONY: cover-test cover-web
cover-test:
	go test -coverprofile=cover.out ./...
cover-web:
	go tool cover -html=cover.out

# From https://github.com/genuinetools/img
.PHONY: AUTHORS
AUTHORS:
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@echo "$(shell git log --format='\n%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf)" >> $@
