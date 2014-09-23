.SILENT :
.PHONY : build clean

TAG:=`git describe --abbrev=0 --tags`
LDFLAGS:=-X main.buildVersion $(TAG)
APPNAME:=webhookd
ROOTPKG:=github.com/ncarlier
PKGDIR:=$(GOPATH)/src/$(ROOTPKG)


all: build

prepare:
	rm -rf $(PKGDIR)
	mkdir -p $(PKGDIR)
	ln -s $(PWD)/src $(PKGDIR)/$(APPNAME)

build: prepare
	echo "Building $(APPNAME)..."
	go build -ldflags "$(LDFLAGS)" -o bin/$(APPNAME) ./src

clean: clean-dist
	rm -f bin/$(APPNAME)

clean-dist:
	rm -rf dist

dist: clean-dist
#	godep restore
	mkdir -p dist/linux/amd64 && GOOS=linux GOARCH=amd64 go build -o dist/linux/amd64/$(APPNAME) ./src
	tar -cvzf dist/$(APPNAME)-linux-amd64-$(TAG).tar.gz -C dist/linux/amd64 $(APPNAME)
#	mkdir -p dist/linux/i386  && GOOS=linux GOARCH=386 go build -o dist/linux/i386/$(APPNAME) ./src
#	tar -cvzf dist/$(APPNAME)-linux-i386-i386$(TAG).tar.gz -C dist/linux/i386 $(APPNAME)

