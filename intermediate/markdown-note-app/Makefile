# Define the Go binary name: it will be created in the current directory, next to the Makefile, and will be named server for now
BINARY_NAME=server

# Define the path to the main.go file
MAIN_FILE=cmd/api/main.go

.PHONY: dev build run clean

# Development run: runs the app directly without building
dev:
	go run $(MAIN_FILE)

# Production build
build:
	go build -o $(BINARY_NAME) $(MAIN_FILE)

# Run the built binary (production run)
run: build
	./$(BINARY_NAME)

# Clean up the binary
clean:
	rm -f $(BINARY_NAME)