
.PHONY: wire
wire: ## generated file based in ./di/inject_*.go
	wire gen ./di

.PHONY: docs
docs: ## generate docs
	@echo "Generating docs..."
	@swag init -g cmd/server/main.go -o ./docs