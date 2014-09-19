.SILENT :
.PHONY : build clean

TAG:=`git describe --abbrev=0 --tags`
LDFLAGS:=-X main.buildVersion $(TAG)
APPNAME:=webhookd

all: build

build:
	echo "Building $(APPNAME)..."
	go build -ldflags "$(LDFLAGS)" -o bin/$(APPNAME) ./src

clean: clean-dist
	rm -rf bin

clean-dist:
	rm -rf dist

dist: clean-dist
	mkdir -p dist/linux/amd64 && GOOS=linux GOARCH=amd64 go build -o dist/linux/amd64/$(APPNAME) ./src
#	mkdir -p dist/linux/i386  && GOOS=linux GOARCH=386 go build -o dist/linux/i386/$(APPNAME) ./src

release: dist
#	godep restore
	tar -cvzf $(APPNAME)-linux-amd64-$(TAG).tar.gz -C dist/linux/amd64 $(APPNAME)
#	tar -cvzf $(APPNAME)-linux-i386-i386$(TAG).tar.gz -C dist/linux/i386 $(APPNAME)

