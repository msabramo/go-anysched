# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

BINARIES        = hyperion-cli
BINARY_PKGS     = ./cmd/hyperion-cli
HYPERION_CLI    = ./hyperion-cli
TEST_APP_NAME   = hyperion-cli-test-$(shell date +'%Y%m%d%H%M%S')
TEST_APP_IMAGE  = k8s.gcr.io/echoserver:1.4
TEST_APP_COUNT  = 1

.PHONY: clean build cli-smoketest check lint test test-race vet test-cover-html help
.DEFAULT_GOAL := help

clean: ## Clean up files that aren't checked into version control
	go clean
	$(RM) $(BINARIES)

build: ## Build all the things
	go build $(BINARY_PKGS)
	CLICOLOR=1 ls -l $(BINARIES)

cli-smoketest: ## Quickly exercise hyperion-cli for Marathon and Kubernetes
	HYPERIONCLI_ENV=local_marathon $(MAKE) _cli-smoketest
	@echo; echo
	HYPERIONCLI_ENV=kubeconfig $(MAKE) _cli-smoketest

_cli-smoketest:
	@echo
	@echo "--------------------------------------------------------------------------------"
	@echo "Deploying service in $(HYPERIONCLI_ENV) ..."
	@echo "--------------------------------------------------------------------------------"
	@echo
	$(HYPERION_CLI) svc deploy --svc-id=$(TEST_APP_NAME) --image=$(TEST_APP_IMAGE) --count=$(TEST_APP_COUNT)
	@echo
	@sleep 5
	@echo "--------------------------------------------------------------------------------"
	@echo "Destroying service in $(HYPERIONCLI_ENV) ..."
	@echo "--------------------------------------------------------------------------------"
	@echo
	$(HYPERION_CLI) svc destroy --svc-id=$(TEST_APP_NAME)

check: test-race vet lint ## Run tests and linters

test: ## Run tests
	go test ./...

test-race: ## Run tests with race detector
	go test -race ./...

lint: ## Run golint linter
	golint -set_exit_status $(shell go list ./...)

vet: ## Run go vet linter
	go vet $(shell go list ./...)

test-cover-html: ## Generate test coverage report
	go test -coverprofile=coverage.out -covermode=count
	go tool cover -func=coverage.out

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
