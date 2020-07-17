
BUILD_IMG ?= ccr.ccs.tencentyun.com/tarsbase/tarsweb:latest

web-files:
	mkdir -p web
	cp -rf ../deploy/deploy/web/* web/

img: web-files
	docker build -t ${BUILD_IMG} .

clean:
	rm -rf web