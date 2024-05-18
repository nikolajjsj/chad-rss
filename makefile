.PHONY: all info backend-build frontend-build install-template-dependencies clean

# Default target
all: frontend-build backend-build

# Print information about available commands
info:
	$(info ------------------------------------------)
	$(info -           ChadRSS                      -)
	$(info ------------------------------------------)
	$(info This Makefile helps you manage your projects.)
	$(info )
	$(info Available commands:)
	$(info - backend-build:  Build the Golang project.)
	$(info - frontend-build:  Build the SvelteKit project.)
	$(info - all:  Run all commands (SvelteBuild, GoBuild).)
	$(info )
	$(info Usage: make <command>)

# Build the Golang project
backend-build:
	@echo "=== Building Golang Project ==="
	@go build -o app -v

# Build the SvelteKit project
frontend-build: install-template-dependencies
	@echo "=== Building Reacy Project ==="
	@if command -v pnpm >/dev/null; then \
		pnpm run -C ./frontend build; \
	else \
		npm run --prefix ./frontend build; \
	fi

# Install template dependencies
install-template-dependencies:
	@if command -v pnpm >/dev/null; then \
		pnpm install -C ./frontend; \
	else \
		npm install --prefix ./frontend; \
	fi

# Clean build artifacts
clean:
	@echo "=== Cleaning build artifacts ==="
	@rm -f app
	@if [ -d "./template/build" ]; then \
		rm -rf ./frontend/build; \
	fi
