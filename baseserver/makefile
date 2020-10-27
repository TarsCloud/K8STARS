SERVER ?= tarsstat
IMG_REPO ?= ccr.ccs.tencentyun.com/tarsbase
PWD = $(shell pwd)
VERSION ?= $(shell date "+%Y%m%d%H%M%S")

deploy:
	mkdir -p deploy
	docker run --entrypoint '' --rm -v ${PWD}/deploy:/deploy tarscloud/framework:v2.4.3 cp -rf /usr/local/tars/cpp/deploy/ /deploy
	rm -rf ./deploy/deploy/web//node_modules/.bin/


img: deploy
	TARS_SERVER=${SERVER} sh replace_obj.sh
	mkdir -p build
	cp Dockerfile conf/${SERVER}.conf deploy/deploy/framework/servers/${SERVER}/bin/${SERVER}  default.env start.sh  _server_meta.yaml build
	cd build && docker build --build-arg SERVER_VERSION=${VERSION} --build-arg SERVER=${SERVER} -t $(IMG_REPO)/${SERVER}:${VERSION} .

registry: deploy
	TARS_SERVER=tarsregistry sh replace_obj.sh
	GOOS=linux go build -o build/tarsregistry ../tarsregistry
	docker build -f tarsregistry/Dockerfile --build-arg SERVER_VERSION=${VERSION} -t $(IMG_REPO)/tarsregistry:${VERSION} .
	
	cp deploy/deploy/framework/servers/tarsregistry/bin/tarsregistry build/tarscppregistry
	docker build -f cppregistry/Dockerfile  --build-arg SERVER_VERSION=${VERSION} -t $(IMG_REPO)/cppregistry:${VERSION}  .

web:
	cp default.env tarsweb/
	cd tarsweb && make img

clean:
	rm -rf deploy
