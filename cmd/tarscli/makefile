GOBUILD := GOOS=linux go build

ALL:
	$(GOBUILD) -o tarscli

img: ALL
	docker build -t ccr.ccs.tencentyun.com/tarsbase/tarscli .

clean:
	rm -rf tarscli