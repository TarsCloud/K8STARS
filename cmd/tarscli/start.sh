#!/bin/bash

if [ -z "$TARS_PATH" ]
then
    export TARS_PATH=/tars
fi

if [ -d ${TARS_PATH}/vol-conf ]; then
    # sync conf to dir
    TARS_SYNC_DIRECTORY=${TARS_PATH}/vol-conf tarscli syncdir
fi

# generate sever config
tarscli genconf

# start server
cd ${TARS_PATH}
${TARS_PATH}/bin/${TARS_BUILD_SERVER} --config=${TARS_PATH}/conf/${TARS_BUILD_SERVER}.conf
