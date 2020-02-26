.SILENT :

export GO111MODULE=on

# App name
APPNAME=httpsig

# Go configuration
GOOS?=linux
GOARCH?=amd64

# Add exe extension if windows target
is_windows:=$(filter windows,$(GOOS))
EXT:=$(if $(is_windows),".exe","")

# Archive name
ARCHIVE=$(APPNAME)-$(GOOS)-$(GOARCH).tgz

# Executable name
EXECUTABLE=$(APPNAME)$(EXT)

all: build

clean:
	-rm -rf release
.PHONY: clean

## Build executable
build:
	-mkdir -p release
	echo "Building: $(EXECUTABLE) $(VERSION) for $(GOOS)-$(GOARCH) ..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o release/$(EXECUTABLE)
.PHONY: build

release/$(EXECUTABLE): build

key:
	openssl genrsa -out key.pem 2048
	openssl rsa -in key.pem -outform PEM -pubout -out key-pub.pem
.PHONY: key
