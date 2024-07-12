build:
	@echo "Building binary..."
	@go build -o main ./cmd/argus

test:
	@echo "Testing..."
	@go test -v -cover -mod=mod ./... -coverprofile=coverage.out

k6:
	@echo "Load testing..."
	@k6 run test/k6/script.js
