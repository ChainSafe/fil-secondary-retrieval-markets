PROJECTNAME=$(shell basename "$(PWD)")
GOLANGCI := $(GOPATH)/bin/golangci-lint

.PHONY: help lint test build client provider install
all: help
help: Makefile
	@echo
	@echo " Choose a make command to run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

$(GOLANGCI):
	if [ ! -f ./bin/golangci-lint ]; then \
		wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest; \
	fi;

lint: $(GOLANGCI)
	./bin/golangci-lint run ./... --timeout 5m0s

test:
	go test ./...

license:
	./scripts/add_license.sh

build: client provider

client:
	go build -o ./build/retrieval-client ./cmd/retrieval-client

provider:
	go build -o ./build/retrieval-provider ./cmd/retrieval-provider

install:
	cd cmd/retrieval-client && go install