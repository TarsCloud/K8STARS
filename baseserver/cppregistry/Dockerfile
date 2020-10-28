FROM ccr.ccs.tencentyun.com/tarsbase/tarscli:latest

ARG SERVER_VERSION=unstable
ENV SERVER_VERSION=${SERVER_VERSION}

ENV TARS_SERVER=tarsregistry
ENV TARS_PATH=/tars
WORKDIR ${TARS_PATH}

RUN mkdir -p $TARS_PATH/bin/../conf/../data/../log

ENV TARS_MERGE_CONF=${TARS_PATH}/bin/tarsregistry.conf
COPY cppregistry/_server_meta.yaml cppregistry/tarsregistry.conf default.env cppregistry/start.sh ${TARS_PATH}/bin/
COPY build/tarscppregistry ${TARS_PATH}/bin/tarsregistry

CMD  ${TARS_PATH}/bin/start.sh