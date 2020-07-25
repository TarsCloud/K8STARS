APP    := TestApp
TARGET := HelloGo

GOBUILD      := go build
DOCKER_BUILD := docker build

REPO         ?= ccr.ccs.tencentyun.com/tarsbase
#VERSION  ?= $(shell date "+%Y%m%d%H%M%S")
VERSION ?= "latest"

LOWCASE_TARGET := $(shell echo $(TARGET) | tr '[:upper:]' '[:lower:]')
IMG_REPO       := $(REPO)/$(LOWCASE_TARGET)

build:
	GOOS=linux $(GOBUILD) -o $(TARGET)

img: build
	$(DOCKER_BUILD) --build-arg SERVER=$(TARGET) -t $(IMG_REPO):$(VERSION) .

tgz: build
	tar czf $(TARGET).tgz $(TARGET) _server_meta.yaml

patch: tgz
	curl --data-binary @$(TARGET).tgz "${TARS_EP}/patch?server=$(TARGET)&version=$(VERSION)"
	
stdout:
	@curl "${TARS_EP}/stdout?server=$(TARGET)"

listlog:
	@curl "${TARS_EP}/listlog?app=$(APP)&server=$(TARGET)"

tailog:
	@curl "${TARS_EP}/tailog?app=$(APP)&server=$(TARGET)&filename=$(LOG_NAME)"

clean:
	rm -rf $(TARGET)
