# Makefile for Proto Event Deserializer

.PHONY: deserialize build clean help

# Default target
help:
	@echo "Proto Event Deserializer Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make deserialize JSON=<json_file_path>  - Deserialize events from JSON file"
	@echo "  make build                              - Build the deserialize_proto binary"
	@echo "  make clean                              - Remove built binary"
	@echo "  make help                               - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make run JSON=events.json"
	@echo "  make run JSON=cloudwatch_logs.json"

# Run the deserializer with the specified JSON file
run:
	@if [ -z "$(JSON)" ]; then \
		echo "Error: JSON file path is required"; \
		echo "Usage: make run JSON=<json_file_path>"; \
		exit 1; \
	fi
	@echo "Deserializing events from: $(JSON)"
	@go run deserialize_proto.go $(JSON)

# Build the binary
build:
	@echo "Building deserialize_proto binary..."
	@go build -o deserialize_proto deserialize_proto.go
	@echo "Binary built successfully: deserialize_proto"

# Clean up built binary
clean:
	@echo "Cleaning up..."
	@rm -f deserialize_proto
	@echo "Clean complete"
