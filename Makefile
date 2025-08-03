# Go项目Makefile

# 变量定义
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
BINARY_NAME = agent
BINARY_PATH = ./bin/$(BINARY_NAME)
COVERAGE_FILE = coverage.out
COVERAGE_HTML = coverage.html

# 默认目标
.PHONY: all
all: test build

# 构建项目
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) -v .

# 运行测试
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

# 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

# 显示覆盖率统计
.PHONY: coverage-stats
coverage-stats: test-coverage
	@echo "Coverage statistics:"
	@$(GOCMD) tool cover -func=$(COVERAGE_FILE) | grep total

# 运行基准测试
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# 代码格式化
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GOCMD) fmt ./...

# 代码整理
.PHONY: tidy
tidy:
	@echo "Tidying modules..."
	@$(GOMOD) tidy

# 代码检查
.PHONY: vet
vet:
	@echo "Running go vet..."
	@$(GOCMD) vet ./...

# 运行golangci-lint
.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	@$(shell go env GOPATH)/bin/golangci-lint run

# 安装golangci-lint
.PHONY: install-lint
install-lint:
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin

# 运行安全检查
.PHONY: security
security:
	@echo "Running security checks..."
	@gosec ./...

# 安装gosec
.PHONY: install-security
install-security:
	@echo "Installing gosec..."
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# 完整的CI检查
.PHONY: ci
ci: fmt vet lint test-coverage
	@echo "All CI checks passed!"

# 清理构建文件
.PHONY: clean
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -f $(BINARY_PATH)
	@rm -f $(COVERAGE_FILE)
	@rm -f $(COVERAGE_HTML)
	@rm -rf bin/

# 下载依赖
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@$(GOMOD) download
	@$(GOMOD) verify

# 更新依赖
.PHONY: update-deps
update-deps:
	@echo "Updating dependencies..."
	@$(GOGET) -u ./...
	@$(GOMOD) tidy

# 运行应用
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BINARY_PATH)

# 开发模式运行
.PHONY: dev
dev:
	@echo "Running in development mode..."
	@$(GOCMD) run main.go

# 查看覆盖率
.PHONY: coverage-view
coverage-view: test-coverage
	@echo "Opening coverage report in browser..."
	@open $(COVERAGE_HTML) || xdg-open $(COVERAGE_HTML)

# 生成覆盖率徽章
.PHONY: coverage-badge
coverage-badge: test-coverage
	@echo "Generating coverage badge..."
	@COVERAGE=$$($(GOCMD) tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$COVERAGE >= 80" | bc -l) -eq 1 ]; then \
		COLOR="brightgreen"; \
	elif [ $$(echo "$$COVERAGE >= 60" | bc -l) -eq 1 ]; then \
		COLOR="yellow"; \
	else \
		COLOR="red"; \
	fi; \
	curl -s "https://img.shields.io/badge/coverage-$${COVERAGE}%25-$$COLOR" > coverage-badge.svg; \
	echo "Coverage badge saved as coverage-badge.svg ($$COVERAGE%)"

# 显示帮助信息
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  coverage-stats - Show coverage statistics"
	@echo "  coverage-view  - Open coverage report in browser"
	@echo "  coverage-badge - Generate coverage badge"
	@echo "  bench          - Run benchmarks"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Tidy modules"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run golangci-lint"
	@echo "  install-lint   - Install golangci-lint"
	@echo "  security       - Run security checks"
	@echo "  install-security - Install gosec"
	@echo "  ci             - Run all CI checks"
	@echo "  clean          - Clean build files"
	@echo "  deps           - Download dependencies"
	@echo "  update-deps    - Update dependencies"
	@echo "  run            - Build and run the application"
	@echo "  dev            - Run in development mode"
	@echo "  help           - Show this help message"
