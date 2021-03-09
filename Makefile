NAME=slackcat
VERSION=$(shell cat VERSION)
BUILD=$(shell git rev-parse --short HEAD)

clean:
	rm -rf build/ release/ arch-release/

deps:
	go mod download

build: deps
	go build -tags osusergo,netgo -ldflags "-s -X main.version=$(VERSION) -X main.build=$(BUILD)" -o slackcat

build-all: deps
	mkdir -p build
	GOOS=darwin GOARCH=amd64 go build -tags osusergo,netgo -ldflags "-s -X main.version=$(VERSION) -X main.build=$(BUILD)" -o build/slackcat-$(VERSION)-darwin-amd64
	GOOS=linux GOARCH=amd64 go build -tags osusergo,netgo -ldflags "-s -X main.version=$(VERSION) -X main.build=$(BUILD)" -o build/slackcat-$(VERSION)-linux-amd64
	GOOS=linux GOARCH=arm go build -tags osusergo,netgo -ldflags "-s -X main.version=$(VERSION) -X main.build=$(BUILD)" -o build/slackcat-$(VERSION)-linux-arm
	GOOS=freebsd GOARCH=amd64 go build -tags osusergo,netgo -ldflags "-s -X main.version=$(VERSION) -X main.build=$(BUILD)" -o build/slackcat-$(VERSION)-freebsd-amd64

release:
	mkdir release
	cp build/* release
	cd release; sha256sum --quiet --check sha256sums.txt && \
	gh release create $(VERSION) -d -t v$(VERSION) *

arch-release:
	mkdir -p arch-release
	go get github.com/seletskiy/go-makepkg/...
	cd arch-release && \
		go-makepkg -p version "Commandline utility for posting snippets to Slack" git://github.com/bcicen/slackcat.git; \
		git clone ssh://aur@aur.archlinux.org/slackcat.git; \
		cp build/* slackcat/
	cd arch-release/slackcat/ && \
		mksrcinfo
