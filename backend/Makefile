.PHONY: help run build test clean dev install

help: ## 显示帮助信息
	@echo "可用命令:"
	@echo "  make run      - 运行应用"
	@echo "  make build    - 编译应用"
	@echo "  make dev      - 开发模式（热重载）"
	@echo "  make test     - 运行测试"
	@echo "  make clean    - 清理编译文件"
	@echo "  make install  - 安装依赖"

run: ## 运行应用
	go run main.go

build: ## 编译应用
	go build -o bin/server main.go

dev: ## 开发模式（需要安装 air）
	air

test: ## 运行测试
	go test -v ./...

clean: ## 清理编译文件
	rm -rf bin/
	go clean

install: ## 安装依赖
	go mod download
	go mod tidy
