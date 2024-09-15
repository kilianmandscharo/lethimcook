BUILD_DIR := build
STATIC_DIR := static

.PHONY: all clean

all: app copy_static copy_env

app: $(wildcard *.go)
	go build -o $(BUILD_DIR)/lethimcook

copy_static:
	mkdir -p $(BUILD_DIR)/$(STATIC_DIR)
	cp -r $(STATIC_DIR)/* $(BUILD_DIR)/$(STATIC_DIR)

copy_env:
	cp .env.example $(BUILD_DIR)/.env

clean:
	rm -rf $(BUILD_DIR)

