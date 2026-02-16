# HiveClaw Makefile
VERSION := 0.1.0
BINARY := hiveclaw
BUILD_DIR := build

.PHONY: all build clean test frontend run install

all: frontend build

build:
	@echo "🔨 Building HiveClaw..."
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BINARY) ./cmd/hiveclaw
	@echo "✅ Built: ./$(BINARY)"

build-all: frontend
	@echo "🔨 Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY)-linux-amd64 ./cmd/hiveclaw
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY)-linux-arm64 ./cmd/hiveclaw
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY)-darwin-amd64 ./cmd/hiveclaw
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY)-darwin-arm64 ./cmd/hiveclaw
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe ./cmd/hiveclaw
	@echo "✅ Binaries in $(BUILD_DIR)/"

frontend:
	@echo "🎨 Building frontend..."
	cd web/frontend && npm install && npm run build
	@echo "✅ Frontend built"

clean:
	rm -f $(BINARY)
	rm -rf $(BUILD_DIR)
	rm -rf web/frontend/dist
	rm -rf web/frontend/node_modules

test:
	go test -v ./...

run: build
	./$(BINARY) start

install: build
	@echo "📦 Installing to /usr/local/bin..."
	sudo cp $(BINARY) /usr/local/bin/
	@echo "✅ Installed: hiveclaw"

dev:
	@echo "🚀 Starting dev mode..."
	cd web/frontend && npm run dev &
	go run ./cmd/hiveclaw start --port 8080
