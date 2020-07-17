#!/bin/bash

export pod=tars-db-all-in-one-d747fbb89-qpxvw


kubectl exec ${pod} -- mysql -uroot -ppass -e "create database db_tars" || exit 1

kubectl exec ${pod} -- mysql -uroot -ppass -e "create database db_user_system"
kubectl exec ${pod} -- mysql -uroot -ppass -e "create database db_tars_web"
kubectl exec ${pod} -- mysql -uroot -ppass -e "create database db_cache_web"
kubectl exec ${pod} -- mysql -uroot -ppass -e "create database tars_stat"
kubectl exec ${pod} -- mysql -uroot -ppass -e "create database tars_property"

kubectl exec -i ${pod} -- mysql -uroot -ppass db_tars < ./deploy/deploy/framework/sql/db_tars.sql
kubectl exec -i ${pod} -- mysql -uroot -ppass db_user_system < ./deploy/deploy/web/demo/sql/db_user_system.sql
kubectl exec -i ${pod} -- mysql -uroot -ppass db_tars_web < ./deploy/deploy/web/sql/db_tars_web.sql
kubectl exec -i ${pod} -- mysql -uroot -ppass db_cache_web < ./deploy/deploy/web/sql/db_cache_web.sql
kubectl exec -i ${pod} -- mysql -uroot -ppass tars_stat < ./deploy/deploy/web/sql/db_cache_web.sql

kubectl exec ${pod} -- mysql -uroot -ppass db_tars -e "show tables"
