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

arch-release:
	rm -rf arch-release && mkdir -p arch-release
	go get github.com/seletskiy/go-makepkg/...
	cd arch-release && \
		go-makepkg -p version "Commandline utility for posting snippets to Slack" git://github.com/vektorlab/slackcat.git; \
		git clone ssh://aur@aur.archlinux.org/slackcat.git; \
		cp build/* slackcat/
	cd arch-release/slackcat/ && \
		mksrcinfo; \
		git commit -a -m "Update to $(VERSION)"

.PHONY: release arch-release
