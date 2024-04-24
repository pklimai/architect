#======================================#
# VARIABLES
#======================================#

## GENERAL
GOLANGCI_TAG := 1.57.1

## BINS
LOCAL_BIN := $(CURDIR)/bin
GOIMPORTS_BIN := $(LOCAL_BIN)/goimports
GOLANGCI_BIN := $(LOCAL_BIN)/golangci-lint

#======================================#
# INSTALLATION
#======================================#

.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps: ## install custom necessary bins
	$(info Installing custom bins for project...)    
	tmp=$$(mktemp -d) && cd $$tmp && go mod init temp && \
		go install github.com/pressly/goose/v3/cmd/goose@latest
		go install github.com/gojuno/minimock/v3/cmd/minimock@latest
		go install golang.org/x/tools/cmd/goimports@latest
	rm -rf $$tmp

bin-deps: .bin-deps ## install necessary bin dependencies

.install-lint: export GOBIN := $(LOCAL_BIN)
.install-lint:
ifeq (,$(wildcard $(GOLANGCI_BIN)))
	$(info Installing golangci-lint v$(GOLANGCI_TAG))
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_TAG)
else
	$(info Golangci-lint is already installed to $(GOLANGCI_BIN))
endif

#======================================#
# TEST
#======================================#

.test:
	go test -race -count 100 ./...

test: .test # run all test in project

#======================================#
# CHECK
#======================================#

.lint: .install-lint
	$(info Running lint against changed files...)
	$(GOLANGCI_BIN) run \
		--new-from-rev=origin/main \
		--config=.golangci.yaml \
		./...

lint: .lint # run golangci-lint against changed files from main

.lint-full: .install-lint
	$(info Running lint all project files...)
	$(GOLANGCI_BIN) run \
		--config=.golangci.yaml \
		./...

lint-full: .lint-full # run golangci-lint against all project files

.pre-push:
	@$(GOIMPORTS_BIN) -w .
	go mod tidy
	@make .lint
	@make test

pre-push: .pre-push # execute checks before push

# Declare that the current commands are not files and
# instruct Makefile not to look for changes to the filesystem.
.PHONY: \
    .bin-deps \
	bin-deps \
	.install-lint \
	.test \
	test \
	.lint \
	lint \
	.lint-full \
	lint-full \
	.pre-push \
	pre-push