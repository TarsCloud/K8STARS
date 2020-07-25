FROM ccr.ccs.tencentyun.com/tarsbase/tarscli:latest

# please set SERVER from docker build --build-arg SERVER=xxx
ARG SERVER=please_build_by_make_img

ENV TARS_BUILD_SERVER ${SERVER}

# copy server and meta files
COPY $SERVER _server_meta.yaml start.sh $TARS_PATH/bin/

CMD tarscli supervisor