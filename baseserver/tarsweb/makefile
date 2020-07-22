
BUILD_IMG ?= ccr.ccs.tencentyun.com/tarsbase/tarsweb:latest

TarsWeb:
	rm -rf TarsWeb && mkdir -p TarsWeb
	git clone -b k8stars https://github.com/defool/TarsWeb

img: TarsWeb
	docker build -t ${BUILD_IMG} .

clean:
	rm -rf TarsWeb