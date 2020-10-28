FROM centos:7

ENV TARS_PATH=/tars
WORKDIR ${TARS_PATH}
RUN mkdir -p $TARS_PATH/bin/../conf/../data/../log
RUN ln -nsf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY start.sh ${TARS_PATH}/bin
COPY tarscli /usr/bin/

CMD ${TARS_PATH}/bin/start.sh
