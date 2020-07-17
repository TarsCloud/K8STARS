FROM ccr.ccs.tencentyun.com/tarsbase/tarscli:latest

ARG SERVER=tarsregistry

ENV TARS_SERVER=${SERVER}
ENV TARS_PATH=/tars
WORKDIR ${TARS_PATH}

RUN mkdir -p $TARS_PATH/bin/../conf/../data/../log

ENV TARS_MERGE_CONF=${TARS_PATH}/bin/${SERVER}.conf
COPY cppregistry/_server_meta.yaml cppregistry/tarsregistry.conf default.env cppregistry/start.sh ${TARS_PATH}/bin/
COPY build/tarscppregistry ${TARS_PATH}/bin/tarsregistry

CMD  ${TARS_PATH}/bin/start.sh