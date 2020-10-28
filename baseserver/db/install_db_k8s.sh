#!/bin/bash

export pod=$(db_pod)
export user=root
export pass=pass
export host=localhost

kubectl exec ${pod} -- mysql -h ${host} -u${user} -p${pass} -e "create database db_tars" || exit 1

kubectl exec ${pod} -- mysql -h ${host} -u${user} -p${pass} -e "create database db_user_system"
kubectl exec ${pod} -- mysql -h ${host} -u${user} -p${pass} -e "create database db_tars_web"
kubectl exec ${pod} -- mysql -h ${host} -u${user} -p${pass} -e "create database db_cache_web"
kubectl exec ${pod} -- mysql -h ${host} -u${user} -p${pass} -e "create database tars_stat"
kubectl exec ${pod} -- mysql -h ${host} -u${user} -p${pass} -e "create database tars_property"

kubectl exec -i ${pod} -- mysql -h ${host} -u${user} -p${pass} db_tars < ./deploy/deploy/framework/sql/db_tars.sql
kubectl exec -i ${pod} -- mysql -h ${host} -u${user} -p${pass} db_user_system < ./deploy/deploy/web/demo/sql/db_user_system.sql
kubectl exec -i ${pod} -- mysql -h ${host} -u${user} -p${pass} db_tars_web < ./deploy/deploy/web/sql/db_tars_web.sql
kubectl exec -i ${pod} -- mysql -h ${host} -u${user} -p${pass} db_cache_web < ./deploy/deploy/web/sql/db_cache_web.sql
kubectl exec -i ${pod} -- mysql -h ${host} -u${user} -p${pass} tars_stat < ./deploy/deploy/web/sql/db_cache_web.sql

kubectl exec ${pod} -- mysql -h ${host} -u${user} -p${pass} db_tars -e "show tables"
