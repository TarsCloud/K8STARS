#!/bin/bash

cd $(dirname $0) || exit 1

source ./default.env

tarscli genconf

sed -i "s/{{DB_HOST}}/$DB_HOST/g" ${TARS_PATH}/conf/${TARS_SERVER}.conf
sed -i "s/{{DB_PORT}}/$DB_PORT/g" ${TARS_PATH}/conf/${TARS_SERVER}.conf
sed -i "s/{{DB_USER}}/$DB_USER/g" ${TARS_PATH}/conf/${TARS_SERVER}.conf
sed -i "s/{{DB_PASS}}/$DB_PASS/g" ${TARS_PATH}/conf/${TARS_SERVER}.conf


# start server
cd ${TARS_PATH}
${TARS_PATH}/bin/${TARS_SERVER} --config=${TARS_PATH}/conf/${TARS_SERVER}.conf
