APP_NAME ?= ctx
OUT_DIR  ?= bin
GOOS     ?= $(shell go env GOOS)
GOARCH   ?= $(shell go env GOARCH)
LDFLAGS   = -s -w

.PHONY: build clean

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(OUT_DIR)/$(APP_NAME) .

clean:
	rm -rf $(OUT_DIR)
