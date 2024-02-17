# Compiler options
GO := go
GOCMD := $(GO) build

# Directories
BUILD_DIR := build
TEMPLATES_DIR := templates
STATIC_DIR := static

# Source files
SOURCES := $(wildcard *.go)

# Targets
TARGET := $(BUILD_DIR)/lethimcook

# Phony targets
.PHONY: all clean

# Default target
all: $(TARGET) copy_templates copy_static

# Build target
$(TARGET): $(SOURCES)
	$(GOCMD) -o $(TARGET)

# Copy templates
copy_templates:
	mkdir -p $(BUILD_DIR)/$(TEMPLATES_DIR)
	cp -r $(TEMPLATES_DIR)/* $(BUILD_DIR)/$(TEMPLATES_DIR)

# Copy static files
copy_static:
	mkdir -p $(BUILD_DIR)/$(STATIC_DIR)
	cp -r $(STATIC_DIR)/* $(BUILD_DIR)/$(STATIC_DIR)

# Clean target
clean:
	rm -rf $(BUILD_DIR)

