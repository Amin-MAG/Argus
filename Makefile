build:
	@echo "Building binary..."
	@go build -o main ./cmd/argus

test:
	@echo "Testing..."
	@go test -v -cover -mod=mod ./...