#!/bin/bash

kubectl delete deployments tars-registry
kubectl delete deployments tars-web
kubectl delete deployments tarsconfig
kubectl delete deployments tarslog
kubectl delete deployments tarslog
kubectl delete deployments tarsproperty
kubectl delete deployments tarsqueryproperty
kubectl delete deployments tarsquerystat
kubectl delete deployments tarsstat
kubectl delete svc tars-registry
kubectl delete svc tars-web
kubectl delete configmaps tars-db-config
