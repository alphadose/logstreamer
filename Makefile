PROJECTNAME := $(shell basename "$(PWD)")
PACKAGES := $(shell go list ./... | grep -v vendor)

# Go related variables.
PROJECTROOT := $(shell pwd)
GOBIN := $(PROJECTROOT)/bin
GOFILES := $(PROJECTROOT)/*.go

# Shell script related variables.
BUILDIR := $(PROJECTROOT)/scripts

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: default
default: help

## install: Install missing dependencies
install:
	@printf "🔨 Installing project dependencies to vendor\n" 
	@GOBIN=$(GOBIN) go get ./... && go mod vendor
	@printf "👍 Done\n"

## build: Build the project binary
build:
	@printf "🔨 Building binary $(GOBIN)/$(PROJECTNAME)\n" 
	@go build -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)
	@printf "👍 Done\n"

## tools: Install development tools
tools:
	@$(BUILDIR)/install_air.sh
	@$(BUILDIR)/install_golint.sh

## start: Start in development mode with hot-reload enabled
start: tools
	@$(PROJECTROOT)/bin/air -c .air.toml

## clean: Clean build files
clean:
	@printf "🔨 Cleaning build cache\n" 
	@go clean $(PACKAGES)
	@printf "👍 Done\n"
	@-rm $(GOBIN)/$(PROJECTNAME) 2>/dev/null

## fmt: Format entire codebase
fmt:
	@printf "🔨 Formatting\n" 
	@gofmt -l $(shell find . -type f -name '*.go'| grep -v "/vendor/")
	@printf "👍 Done\n"

## vet: Vet entire codebase
vet:
	@printf "🔨 Vetting\n" 
	@go vet $(PACKAGES)
	@printf "👍 Done\n"

## lint: Check codebase for style mistakes
lint:
	@printf "🔨 Linting\n"
	@golint -set_exit_status $(PACKAGES)
	@printf "👍 Done\n"

## deps: Sync all dependencies
deps:
	@printf "🔨 Syncing Dependencies\n"
	@go get ./...; go mod vendor
	@printf "👍 Done\n"

## update: Updates all dependencies
update:
	@printf "🔨 Updating Dependencies\n"
	@go get -d
	@printf "👍 Done\n"

## precommit: Formats, vets and lints the codebase before commit
precommit: fmt vet lint

## test: Run tests
test:
	@printf "🔨 Testing\n"
	@go test -race -coverprofile=coverage.txt -covermode=atomic
	@printf "👍 Done\n"

## help: Display this help
help: Makefile
	@printf "\n Rimworld: The dark side of EzFlow\n\n"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@printf ""
