.SILENT:
.DEFAULT_GOAL := help

GO ?= go
GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin
GOLINT ?= $(GOBIN)/golint
GOSEC ?= $(GOBIN)/gosec

VILICUS_API_BIN ?= vilicus-api-bin
VILICUS_CLIENT_BIN ?= vilicus-client-bin
VILICUS_MIGRATION_BIN ?= vilicus-migration-bin
CMD_API ?= ./cmd/api/main.go
CMD_CLIENT ?= ./cmd/client/main.go
CMD_MIGRATION ?= ./cmd/migration/main.go

COLOR_RESET=\033[0;39;49m
COLOR_BOLD=\033[1m
COLOR_ULINE=\033[4m
COLOR_BOLD_OFF=\033[0;21m
COLOR_ULINE_OFF=\033[0;24m
COLOR_NORM=\033[0;39m
COLOR_GREN=\033[38;5;118m
COLOR_BLUE=\033[38;5;81m
COLOR_RED=\033[38;5;161m
COLOR_PURP=\033[38;5;135m
COLOR_ORNG=\033[38;5;208m
COLOR_YELO=\033[38;5;227m
COLOR_GRAY=\033[38;5;245m
COLOR_WHIT=\033[38;5;15m

PROJECT := Vilicus

TAG := $(shell git rev-parse --abbrev-ref HEAD)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT := $(shell git rev-parse $(TAG))
LDFLAGS := '-X "main.version=$(TAG)" -X "main.commit=$(COMMIT)" -X "main.date=$(DATE)" -w -s'

## Builds all project binaries
build: build-api build-client build-migration

## Builds all project binaries using linux architecture
build-linux: build-api-linux build-client-linux build-migration-linux

## Builds API code into a binary
build-api:
	$(GO) build -ldflags $(LDFLAGS) -o "$(VILICUS_API_BIN)" $(CMD_API)

## Builds API code using linux architecture into a binary
build-api-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o "$(VILICUS_API_BIN)" $(CMD_API)

## Builds API code into a binary
build-client:
	$(GO) build -ldflags $(LDFLAGS) -o "$(VILICUS_CLIENT_BIN)" $(CMD_CLIENT)

## Builds API code using linux architecture into a binary
build-client-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o "$(VILICUS_CLIENT_BIN)" $(CMD_CLIENT)

## Builds all image locally with the latest tags
build-images:
	chmod +x scripts/build-.sh
	./scripts/build-images.sh

## Builds preset images locally with the latest tags
build-preset-images:
	chmod +x scripts/preset/build-images.sh
	./scripts/preset/build-images.sh no_updater

## Builds preset images(with database updater) locally with the latest tags
build-preset-images-updater:
	chmod +x scripts/preset/build-images.sh 
	./scripts/preset/build-images.sh updater

## Builds API migration code into a binary
build-migration:
	$(GO) build -ldflags $(LDFLAGS) -o "$(VILICUS_MIGRATION_BIN)" $(CMD_MIGRATION)

## Builds API migration code using linux architecture into a binary
build-migration-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(LDFLAGS) -o "$(VILICUS_MIGRATION_BIN)" $(CMD_MIGRATION)

## Checks dependencies
check-deps:
	$(GO) mod tidy && $(GO) mod verify

## Runs a security static analysis using Gosec
check-sec:
	GO111MODULE=off $(GO) get -u github.com/securego/gosec/cmd/gosec
	$(GOSEC) ./...

## Prints help message
help:

	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

	printf "\n${COLOR_YELO}${PROJECT}\n-------\n${COLOR_RESET}"
	printf "${COLOR_BOLD}  Usage:${COLOR_RESET}"
	printf "\n"
	printf "${COLOR_BLUE}    make <target>${COLOR_RESET}"
	printf "\n"
	printf "\n"
	printf "${COLOR_BOLD}  Targets:${COLOR_RESET}"
	printf "\n"
	awk '/^[a-zA-Z\-\_0-9\.%]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "${COLOR_BLUE}    make %-30s${COLOR_RESET} %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST) | sort
	printf "\n"

## Runs all Vilicus lint
lint:
	GO111MODULE=off $(GO) get -u golang.org/x/lint/golint
	$(GOLINT) ./...

## Push images to hub.docker
push-images:
	chmod +x scripts/push-images.sh
	./scripts/push-images.sh

## Push preset images to hub.docker
push-preset-images:
	chmod +x scripts/preset/push-images.sh
	./scripts/preset/push-images.sh

## Builds and push images with the latest tags
update-images: build-images push-images

## Builds and push preset images with the latest tags
update-preset-images: build-preset-images push-preset-images