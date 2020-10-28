FROM ccr.ccs.tencentyun.com/tarsbase/tarscli:latest

# please set SERVER from docker build --build-arg SERVER=xxx
ARG SERVER_VERSION=unstable

ENV SERVER_VERSION=${SERVER_VERSION}
ENV TARS_SERVER=tarsregistry

ENV TARS_PATH=/tars
WORKDIR ${TARS_PATH}

RUN mkdir -p $TARS_PATH/bin/../conf/../data/../log
COPY _server_meta.yaml default.env tarsregistry/start.sh ${TARS_PATH}/bin/
COPY build/tarsregistry ${TARS_PATH}/bin/tarsregistry


CMD ${TARS_PATH}/bin/start.sh
