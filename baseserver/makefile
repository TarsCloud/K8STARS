SERVER ?= tarsstat
IMG_REPO ?= ccr.ccs.tencentyun.com/tarsbase
PWD = $(shell pwd)

deploy:
	mkdir -p deploy
	docker run --entrypoint '' --rm -v ${PWD}/deploy:/deploy tarscloud/framework:v2.4.3 cp -rf /usr/local/tars/cpp/deploy/ /deploy
	rm -rf ./deploy/deploy/web//node_modules/.bin/


img: deploy
	TARS_SERVER=${SERVER} sh replace_obj.sh
	mkdir -p build
	cp Dockerfile conf/${SERVER}.conf deploy/deploy/framework/servers/${SERVER}/bin/${SERVER}  default.env start.sh  _server_meta.yaml build
	cd build && docker build --build-arg SERVER=${SERVER} -t $(IMG_REPO)/${SERVER}:latest .

registry: deploy
	TARS_SERVER=tarsregistry sh replace_obj.sh
	cp deploy/deploy/framework/servers/tarsregistry/bin/tarsregistry build/tarscppregistry
	GOOS=linux go build -o build/tarsregistry github.com/tarscloud/k8stars/tarsregistry
	docker build -f tarsregistry/Dockerfile -t $(IMG_REPO)/tarsregistry:latest .
	docker build -f cppregistry/Dockerfile  -t $(IMG_REPO)/cppregistry:latest  .

web:
	cp default.env tarsweb/
	cd tarsweb && make img

clean:
	rm -rf deploy
