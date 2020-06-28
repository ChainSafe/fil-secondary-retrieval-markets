PROJECTNAME=$(shell basename "$(PWD)")
GOLANGCI := $(GOPATH)/bin/golangci-lint

.PHONY: help lint test
all: help
help: Makefile
	@echo
	@echo " Choose a make command to run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

$(GOLANGCI):
	wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest

lint: $(GOLANGCI)
	./bin/golangci-lint run ./... --timeout 5m0s

test:
	go test ./...

license:
	./scripts/add_license.sh

