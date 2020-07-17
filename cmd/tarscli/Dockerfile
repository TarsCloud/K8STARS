FROM centos:7

ENV TARS_PATH=/tars
WORKDIR ${TARS_PATH}
RUN mkdir -p $TARS_PATH/bin/../conf/../data/../log

COPY start.sh ${TARS_PATH}/bin
COPY tarscli /usr/bin/

CMD ${TARS_PATH}/bin/start.sh
