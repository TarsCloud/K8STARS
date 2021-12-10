GOBUILD := GOOS=linux go build
VERSION ?= $(shell git describe --tags)

ALL:
	$(GOBUILD) -ldflags "-X main.VERSION=$(VERSION)" -o tarscli

img: ALL
	docker build -t ccr.ccs.tencentyun.com/tarsbase/tarscli:${VERSION} .

clean:
	rm -rf tarscli