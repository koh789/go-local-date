.PHONY: help
help: ## help
	@echo "------- Task list ------"
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36mmake %-20s\033[0m %s\n", $$1, $$2}'

.PHONY: install-all
install-all:
	make install-golangci-lint

.PHONY: install-golangci-lint
install-golangci-lint: ## install golangci-lint
	which golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.52.1

.PHONY: test-all
test-all: ## 静的解析＆テスト
	make lint
	make test

.PHONY: lint
lint: ## 静的解析 golangci-lint
	make install-golangci-lint
	golangci-lint cache clean
	golangci-lint run

.PHONY: test
test: ## test. 一応カバレジ結果c.outとして出力します
	@echo "+ go test..."
	go clean ./... && go test -v ./...
	@echo "+ go test clear."
