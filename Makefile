# Variables
GO=go
BINARY_NAME=MagiChessServer

# Default target
.PHONY: all
all: build  # Default action when "make" is called without arguments

# Build the binary
.PHONY: build
build:
	$(GO) build -o $(BINARY_NAME)

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

# Run the binary
.PHONY: run
run: build
	./$(BINARY_NAME)

# Run tests
.PHONY: test
test:
	$(GO) test ./...

# Fetch Go dependencies
.PHONY: deps
deps:
	$(GO) mod tidy
