FROM ccr.ccs.tencentyun.com/tarsbase/tarscli:latest

# please set SERVER from docker build --build-arg SERVER=xxx
ARG SERVER=please_build_by_make_img

ENV TARS_SERVER=${SERVER}
ENV TARS_PATH=/tars
WORKDIR ${TARS_PATH}

RUN mkdir -p $TARS_PATH/bin/../conf/../data/../log

ENV TARS_MERGE_CONF=${TARS_PATH}/bin/${SERVER}.conf
COPY ${SERVER} ${SERVER}.conf default.env start.sh  _server_meta.yaml ${TARS_PATH}/bin/

CMD tarscli supervisor

