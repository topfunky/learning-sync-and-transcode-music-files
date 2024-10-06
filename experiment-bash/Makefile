# Variables
SHELL := /bin/bash
TEST_SCRIPT := ./run_tests.sh
LINT_DIR := .

# Default target
all: test lint

# Run tests
test:
	@echo "Running tests..."
	@$(TEST_SCRIPT)

# Lint code
lint:
	@echo "Linting code..."
	@find $(LINT_DIR) -type f -name '*-func.sh' -exec shellcheck -x {} \;

# Clean target (optional)
clean:
	@echo "Cleaning up..."


# Phony targets
.PHONY: all test lint clean