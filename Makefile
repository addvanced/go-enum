.DEFAULT_GOAL := generate

APP_NAME=enum-gen

.PHONY: build
build:
	@echo "Building..."
	@go build -o $(APP_NAME) cmd/$(APP_NAME)/main.go && chmod +x $(APP_NAME)
	@echo "Build complete"


.PHONY: generate
generate: build
	@echo "Generating..."
	@go generate ./...
	@echo "Generate complete"