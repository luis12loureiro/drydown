.DEFAULT_GOAL := build

fmt:
	@echo "Formatting Go code..."
	go fmt ./...
.PHONY:fmt

lint: fmt
	@echo "Linting Go code..."
	golint ./...
.PHONY:lint

vet: lint
	@echo "Vetting Go code..."
	go vet ./...
.PHONY:vet

build: vet
	@echo "Building application..."
	go build -o bin/drydown ./main.go
.PHONY:build