# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: run-deploy-example clean lint vet help
.DEFAULT_GOAL := help

run-deploy-example: ## Run the example program
	go run deploy-example.go

clean: ## Clean up files that aren't checked into version control
	go clean
	$(RM) $(BINARIES)

lint: ## Run golint linter
	golint -set_exit_status $(shell go list ./...)

vet: ## Run go vet linter
	go vet $(shell go list ./...)

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
