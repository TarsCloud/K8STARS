#!/bin/bash
set -e

# 创建tars-system命名空间
kubectl create namespace tars-system

# 设置默认的命名空间
kubectl config set-context --current --namespace=tars-system

git clone https://github.com/TarsCloud/K8STARS

# 进入baseserver目录，获取tars的deploy相关文件
cd K8STARS/baseserver
make deploy

# 创建一个MySQL数据库用于体验(在生产环境中，建议使用云db实例)
kubectl apply -f yaml/db_all_in_one.yaml

# 等待状态正常
kubectl wait --timeout=30s --for=condition=available deployment/tars-db-all-in-one

# 获取pod名
export db_pod=$(kubectl get pod -l  app=tars-db-all-in-one -o jsonpath='{.items[0].metadata.name}')

# 基于上面的pod
sh db/install_db_k8s.sh

kubectl apply -f yaml/registry.yaml
kubectl apply -f yaml/tarsweb.yaml
for server in "tarsnotify" "tarsconfig" "tarslog" "tarsstat" "tarsproperty" "tarsquerystat" "tarsqueryproperty"; do
    kubectl apply -f yaml/$server.yaml 
done 

# 安装完成后恢复到默认的namespace：
kubectl config set-context --current --namespace=default

echo "tarsweb url:"
echo "http://node_ip:30000"