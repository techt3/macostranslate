# Makefile for macostranslate App

.PHONY: build run clean install deps test

# Build the application
build:
	go build -o macostranslate

# Build with optimizations for distribution
build-release:
	go build -ldflags "-s -w" -o macostranslate

# Run the application
run: build
	./macostranslate

# Clean build artifacts
clean:
	rm -f macostranslate

# Install dependencies
deps:
	go mod tidy

# Install the app to /usr/local/bin (optional)
install: build-release
	sudo cp macostranslate /usr/local/bin/

# Uninstall from /usr/local/bin
uninstall:
	sudo rm -f /usr/local/bin/macostranslate

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-release - Build with optimizations"
	@echo "  run           - Build and run the application"
	@echo "  clean         - Remove build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  install       - Install to /usr/local/bin"
	@echo "  uninstall     - Remove from /usr/local/bin"
	@echo "  help          - Show this help message"
