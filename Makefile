NAME=slackcat
VERSION=$(shell cat VERSION)

build:
	mkdir -p build
	go get -v -d
	GOOS=darwin GOARCH=amd64 go build -o build/slackcat-$(VERSION)-darwin-amd64
	GOOS=linux GOARCH=amd64 go build -o build/slackcat-$(VERSION)-linux-amd64

release:
	rm -rf release && mkdir release
	go get github.com/progrium/gh-release/...
	cp build/* release
	gh-release create vektorlab/$(NAME) $(VERSION) \
		$(shell git rev-parse --abbrev-ref HEAD) $(VERSION)

.PHONY: release
