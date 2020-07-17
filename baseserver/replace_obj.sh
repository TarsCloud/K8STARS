#!/bin/bash

cp _server_meta_tpl.yaml _server_meta.yaml

if [[ "$OSTYPE" == "darwin"* ]]; then
    alias sed="sed -i ''"
else
    alias sed="sed -i"
fi

sed "s/tars_server/${TARS_SERVER}/g" _server_meta.yaml
case ${TARS_SERVER} in
    tarslog)
        sed "s/MainObj/LogObj/g" _server_meta.yaml
        ;;
    tarsconfig)
        sed "s/MainObj/ConfigObj/g" _server_meta.yaml
        ;;
    tarsnotify)
        sed "s/MainObj/NotifyObj/g" _server_meta.yaml
        ;;
    tarsproperty)
        sed "s/MainObj/PropertyObj/g" _server_meta.yaml
        ;;
    tarsstat)
        sed "s/MainObj/StatObj/g" _server_meta.yaml
        ;;
    tarsquerystat)
        sed "s/MainObj/QueryObj/g" _server_meta.yaml
        ;;
    tarsqueryproperty)
        sed "s/MainObj/QueryObj/g" _server_meta.yaml
        ;;
    tarsregistry)
        sed "s/MainObj/Registry/g" _server_meta.yaml
        ;;
    tarsjaeger)
        sed "s/MainObj/JaegerObj/g" _server_meta.yaml
        ;;
    *)
        echo "Server ${TARS_SERVER} not found"
        exit 1
        ;;
esac