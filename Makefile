# Simple Makefile to set version, build and create cross-platform binaries.
NAME := kcskit
PACKAGE := github.com/arturscheiner/kcskit/cmd
DIST := dist

# default values (overridable when invoking make)
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X '$(PACKAGE).Version=$(VERSION)' -X '$(PACKAGE).Commit=$(COMMIT)' -X '$(PACKAGE).Date=$(DATE)'

PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	windows/amd64 \
	windows/arm64 \
	darwin/amd64 \
	darwin/arm64

.PHONY: all release build clean single

all: build

# build for the current host (build the main package in repo root)
build:
	@echo "building $(NAME) (version=$(VERSION)) for host"
	@go build -ldflags "$(LDFLAGS)" -o $(NAME) .

# produce cross-platform binaries into $(DIST)
release: clean
	@mkdir -p $(DIST)
	@echo "creating release artifacts version=$(VERSION)"
	@for plat in $(PLATFORMS); do \
	    OS=$${plat%/*}; ARCH=$${plat#*/}; \
	    EXT=""; BINNAME="$(NAME)-$(VERSION)-$${OS}-$${ARCH}"; \
	    if [ "$${OS}" = "windows" ]; then EXT=".exe"; fi; \
	    echo "  -> $${BINNAME}$${EXT}"; \
	    GOOS=$${OS} GOARCH=$${ARCH} CGO_ENABLED=0 \
	      go build -trimpath -ldflags "$(LDFLAGS)" -o "$(DIST)/$${BINNAME}$${EXT}" . || exit 1; \
	done
	@echo "artifacts placed in $(DIST)/"

clean:
	@rm -rf $(DIST) $(NAME)

# convenience: build a single platform, e.g.:
# make PLATFORM=linux/amd64 single
single:
ifndef PLATFORM
	$(error PLATFORM is required, e.g. PLATFORM=linux/amd64)
endif
	@mkdir -p $(DIST)
	@OS=$${PLATFORM%/*}; ARCH=$${PLATFORM#*/}; \
	EXT=""; BINNAME="$(NAME)-$(VERSION)-$${OS}-$${ARCH}"; \
	if [ "$${OS}" = "windows" ]; then EXT=".exe"; fi; \
	echo "building $${BINNAME}$${EXT}"; \
	GOOS=$${OS} GOARCH=$${ARCH} CGO_ENABLED=0 \
	  go build -trimpath -ldflags "$(LDFLAGS)" -o "$(DIST)/$${BINNAME}$${EXT}" .
