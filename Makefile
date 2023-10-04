.SILENT :

export GO111MODULE=on

# App name
APPNAME=webhookd

# Go configuration
GOOS?=$(shell go env GOHOSTOS)
GOARCH?=$(shell go env GOHOSTARCH)

# Add exe extension if windows target
is_windows:=$(filter windows,$(GOOS))
EXT:=$(if $(is_windows),".exe","")

# Archive name
ARCHIVE=$(APPNAME)-$(GOOS)-$(GOARCH).tgz

# Executable name
EXECUTABLE=$(APPNAME)$(EXT)

# Extract version infos
PKG_VERSION:=github.com/ncarlier/$(APPNAME)/pkg/version
VERSION:=`git describe --always --dirty`
GIT_COMMIT:=`git rev-list -1 HEAD --abbrev-commit`
BUILT:=`date`
define LDFLAGS
-X '$(PKG_VERSION).Version=$(VERSION)' \
-X '$(PKG_VERSION).GitCommit=$(GIT_COMMIT)' \
-X '$(PKG_VERSION).Built=$(BUILT)'
endef

all: build

# Include common Make tasks
root_dir:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
makefiles:=$(root_dir)/makefiles
include $(makefiles)/help.Makefile

## Clean built files
clean:
	-rm -rf release
.PHONY: clean

## Build executable
build:
	-mkdir -p release
	echo ">>> Building: $(EXECUTABLE) $(VERSION) for $(GOOS)-$(GOARCH) ..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o release/$(EXECUTABLE)
.PHONY: build

release/$(EXECUTABLE): build

# Check code style
check-style:
	echo ">>> Checking code style..."
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...
.PHONY: check-style

# Check code criticity
check-criticity:
	echo ">>> Checking code criticity..."
	go run github.com/go-critic/go-critic/cmd/gocritic@latest check -enableAll ./...
.PHONY: check-criticity

# Check code security
check-security:
	echo ">>> Checking code security..."
	go run github.com/securego/gosec/v2/cmd/gosec@latest -quiet ./...
.PHONY: check-security

## Code quality checks
checks: check-style check-criticity
.PHONY: checks

## Run tests
test: 
	go test ./...
.PHONY: test

## Install executable
install: release/$(EXECUTABLE)
	echo "Installing $(EXECUTABLE) to ${HOME}/.local/bin/$(EXECUTABLE) ..."
	cp release/$(EXECUTABLE) ${HOME}/.local/bin/$(EXECUTABLE)
.PHONY: install

## Create Docker image
image:
	echo "Building Docker image ..."
	docker build --rm --target slim -t ncarlier/$(APPNAME) .
.PHONY: image

# Generate changelog
CHANGELOG.md:
	standard-changelog --first-release

## Create archive
archive: release/$(EXECUTABLE) CHANGELOG.md
	echo "Creating release/$(ARCHIVE) archive..."
	tar czf release/$(ARCHIVE) README.md LICENSE CHANGELOG.md -C release/ $(EXECUTABLE)
	rm release/$(EXECUTABLE)
.PHONY: archive

## Create distribution binaries
distribution:
	GOARCH=amd64 make build archive
	GOARCH=arm64 make build archive
	GOARCH=arm make build archive
	GOOS=darwin make build archive
.PHONY: distribution
