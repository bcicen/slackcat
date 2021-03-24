NAME=slackcat
VERSION=$(shell cat VERSION)
BUILD=$(shell git rev-parse --short HEAD)
LDFLAGS="-s -X main.version=$(VERSION) -X main.build=$(BUILD)"

clean:
	rm -rf _build/ _release/ _aur/

deps:
	go mod download

build: deps
	go build -tags osusergo,netgo -ldflags $(LDFLAGS) -o slackcat

build-all: deps
	mkdir -p _build
	GOOS=darwin  GOARCH=amd64 go build -tags osusergo,netgo -ldflags $(LDFLAGS) -o _build/slackcat-$(VERSION)-darwin-amd64
	GOOS=linux   GOARCH=amd64 go build -tags osusergo,netgo -ldflags $(LDFLAGS) -o _build/slackcat-$(VERSION)-linux-amd64
	GOOS=linux   GOARCH=arm   go build -tags osusergo,netgo -ldflags $(LDFLAGS) -o _build/slackcat-$(VERSION)-linux-arm
	GOOS=freebsd GOARCH=amd64 go build -tags osusergo,netgo -ldflags $(LDFLAGS) -o _build/slackcat-$(VERSION)-freebsd-amd64
	cd _build; sha256sum * > sha256sums.txt

release:
	mkdir _release
	cp _build/* _release/
	cd _release; sha256sum --quiet --check sha256sums.txt && \
	gh release create $(VERSION) -d -t v$(VERSION) *

aur:
	git clone ssh://aur@aur.archlinux.org/slackcat.git _aur
	cd _aur && \
		sed -i "/^pkgver=/c\pkgver=$(VERSION)" PKGBUILD && \
		makepkg --printsrcinfo > .SRCINFO
