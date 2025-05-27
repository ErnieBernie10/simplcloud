SRC=./src/main.go
BUILD_DIR=build
BIN=$(BUILD_DIR)/main
TEMPLATES_SRC=src/internal/web/templates
TEMPLATES_DST=$(BUILD_DIR)/templates
STATIC_SRC=src/internal/web/static
STATIC_DST=$(BUILD_DIR)/static

TAILWIND_CMD=npx @tailwindcss/cli -i ./src/internal/web/input.css -o ./src/internal/web/static/output.css --watch

.PHONY: build
build: clean $(BIN) copy-templates

$(BIN): $(SRC)
	@echo "👉 Building Go binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BIN) $(SRC)

.PHONY: copy-templates
copy-templates:
	@echo "👉 Copying templates and static files..."
	@mkdir -p $(TEMPLATES_DST)
	@mkdir -p $(STATIC_DST)
	@cp -r $(TEMPLATES_SRC)/* $(TEMPLATES_DST)/
	@cp -r $(STATIC_SRC)/* $(STATIC_DST)/
	@cp -r .env $(BUILD_DIR)/.env

.PHONY: clean
clean:
	@echo "🧹 Cleaning build directory..."
	@rm -rf $(BUILD_DIR)

.PHONY: run
run:
	@echo "🚀 Running app..."
	@$(BIN) serve

.PHONY: tailwind-watch
tailwind-watch:
	@echo "🎨 Starting Tailwind CSS watcher..."
	@$(TAILWIND_CMD)

.PHONY: dev
dev: build
	@$(MAKE) -j 2 run tailwind-watch
