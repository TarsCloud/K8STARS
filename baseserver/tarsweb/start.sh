#!/bin/bash

source ./default.env

sed -i "s/{{DB_HOST}}/$DB_HOST/g" /web/config/webConf.js
sed -i "s/{{DB_PORT}}/$DB_PORT/g" /web/config/webConf.js
sed -i "s/{{DB_USER}}/$DB_USER/g" /web/config/webConf.js
sed -i "s/{{DB_PASS}}/$DB_PASS/g" /web/config/webConf.js

sed -i "s/{{TARS_LOCATOR}}/$TARS_LOCATOR/g" /web/config/tars.conf


# start server
cd /web && npm run start
