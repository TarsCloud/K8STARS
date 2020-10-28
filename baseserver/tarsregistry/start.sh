#!/bin/bash

cd $(dirname $0) || exit 1

source ./default.env
tarscli genconf

export DB_DSN="$DB_USER:$DB_PASS@tcp($DB_HOST:$DB_PORT)/db_tars"

# start server
cd ${TARS_PATH}
${TARS_PATH}/bin/${TARS_SERVER} --config=${TARS_PATH}/conf/${TARS_SERVER}.conf
