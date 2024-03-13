GO := go
GOCMD := $(GO) build

BUILD_DIR := build
STATIC_DIR := static

SOURCES := $(wildcard *.go)

TARGET := $(BUILD_DIR)/lethimcook

.PHONY: all clean

all: $(TARGET) copy_static copy_env

$(TARGET): $(SOURCES)
	$(GOCMD) -o $(TARGET)

copy_static:
	mkdir -p $(BUILD_DIR)/$(STATIC_DIR)
	cp -r $(STATIC_DIR)/* $(BUILD_DIR)/$(STATIC_DIR)

copy_env:
	cp .env.example $(BUILD_DIR)/.env

clean:
	rm -rf $(BUILD_DIR)

