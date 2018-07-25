# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

HYPERIONCLI_BIN     = bin/hyperion-cli
HYPERIONCLI_PKG     = ./cmd/hyperion-cli
BINARIES            = $(HYPERIONCLI_BIN)
BINARY_PKGS         = $(HYPERIONCLI_PKG)
TEST_APP_NAME       = hyperion-cli-test-$(shell date +'%Y%m%d%H%M%S')
TEST_APP_IMAGE      = k8s.gcr.io/echoserver:1.4
TEST_APP_COUNT      = 1

.PHONY: clean build cli-smoketest check docker-image lint test test-cover test-cover-html test-race top-cyclo vet html help
.DEFAULT_GOAL := help

clean: ## Clean up files that aren't checked into version control
	go clean -cache -testcache
	$(RM) $(BINARIES)

build: $(HYPERIONCLI_BIN) ## Build all the things

$(HYPERIONCLI_BIN): ## Build hyperion-cli
	go build -o $(HYPERIONCLI_BIN) $(HYPERIONCLI_PKG)
	CLICOLOR=1 ls -l $(HYPERIONCLI_BIN)

cli-smoketest: ## Quickly exercise hyperion-cli for Marathon and Kubernetes
	$(MAKE) cli-smoketest-marathon
	@echo; echo
	$(MAKE) cli-smoketest-kubernetes

cli-smoketest-marathon: ## Quickly exercise hyperion-cli for Marathon
	HYPERIONCLI_ENV=local_marathon $(MAKE) _cli-smoketest

cli-smoketest-kubernetes: ## Quickly exercise hyperion-cli for Marathon
	HYPERIONCLI_ENV=kubeconfig $(MAKE) _cli-smoketest

_cli-smoketest: $(HYPERIONCLI_BIN)
	@echo
	@echo "--------------------------------------------------------------------------------"
	@echo "Deploying service in $(HYPERIONCLI_ENV) ..."
	@echo "--------------------------------------------------------------------------------"
	@echo
	$(HYPERIONCLI_BIN) svc deploy --svc-id=$(TEST_APP_NAME) --image=$(TEST_APP_IMAGE) --count=$(TEST_APP_COUNT)
	@echo
	@sleep 5
	@echo "--------------------------------------------------------------------------------"
	@echo "Destroying service in $(HYPERIONCLI_ENV) ..."
	@echo "--------------------------------------------------------------------------------"
	@echo
	$(HYPERIONCLI_BIN) svc destroy --svc-id=$(TEST_APP_NAME)

check: test-race vet lint ## Run tests and linters

test: ## Run tests
	HYPERIONCLI_ENV=minikube go test ./...

test-cover: ## Generate test coverage report
	HYPERIONCLI_ENV=minikube scripts/coverage

test-cover-html: ## Generate HTML test coverage report
	go test -coverprofile=coverage.out -covermode=count
	go tool cover -func=coverage.out

test-race: ## Run tests with race detector
	go test -race ./...

lint: ## Run golint linter
	golint -set_exit_status $(shell go list ./...)

vet: ## Run go vet linter
	go vet $(shell go list ./...)

top-cyclo: ## Display function with most cyclomatic complexity
	gocyclo -top 10 $(shell find . \( -name vendor -o -name .git -o -name bin \) -prune -o -type d -print)

metalinter: ## Run gometalinter, which does a bunch of checks
	@echo "Running: gometalinter --config=gometalinter.json ./..."
	@gometalinter --config=gometalinter.json ./... && echo "All gometalinter checks passed!"

docker-image: ## Build a docker image with hyperion-cli
	docker build -t hyperion-cli .

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
