GOBUILD := GOOS=linux go build
VERSION := $(shell date "+%Y%m%d%H%M%S")

ALL:
	$(GOBUILD) -ldflags "-X main.VERSION=$(VERSION)" -o tarscli

img: ALL
	docker build -t ccr.ccs.tencentyun.com/tarsbase/tarscli .

clean:
	rm -rf tarscli