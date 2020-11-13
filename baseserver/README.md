# TARS基础服务部署

## 部署步骤

### 前置工作
   - 安装kubernetes，可以使用kubectl或其他方式来管理集群
   - 能执行docker命令行的终端

### 安装基础服务步骤

运行一键安装脚本：
```
curl https://raw.githubusercontent.com/TarsCloud/K8STARS/master/baseserver/install_all.sh | sh
```

或按以下步骤执行安装:

1. 部署tars db
   ```
   git clone https://github.com/TarsCloud/K8STARS

   // 创建tars-system命名空间
   kubectl create namespace tars-system

   // 设置默认的命名空间
   kubectl config set-context --current --namespace=tars-system
   
   // 进入baseserver目录，获取tars的deploy相关文件
   cd K8STARS/baseserver
   make deploy

   // 创建一个MySQL数据库用于体验(在生产环境中，建议使用云db实例)
   kubectl apply -f yaml/db_all_in_one.yaml

   // 确认pod的状态是正常的
   kubectl get pods

   // 获取pod名
   export db_pod=$(kubectl get pod -l  app=tars-db-all-in-one -o jsonpath='{.items[0].metadata.name}')

   // 基于上面的pod
   sh db/install_db_k8s.sh
   ```
   如果使用了外部数据库，可以通过db/all_db.sql将表结构导入。
   对于已有的tars db，再请执行sql文件会清空原有的数据，只需要导入缺少的db即可。
   

2. 安装tars registry
   
   使用`kubectl apply -f yaml/registry.yaml`部署tars registry。
   如果没用k8s创建的db，请修改`registry.yaml`中的数据库地址。

3. 安装tarsweb
   
   使用`kubectl apply -f yaml/tarsweb.yaml`部署。
   
   tarsweb默认使用3000端口，可以根据实际使用情况，配置对应的访问方式，通过浏览器来访问。（ 当前使用NodePort方式映射至外网30000端口，通过 http://tarsweb所在Node的外网IP:30000 访问tarsweb）
   
   说明：当前tarsweb版本未兼容k8s中的场景，页面中有重启/停止等入口，但是操作会失败。

4. 部署其他服务
   
   以`tarsnotify`为例，使用`kubectl apply -f yaml/tarsnotify.yaml`来部署。其中数据库相关配置可以按需要替换。
   其他服务可以使用同一方式来部署，将tarsnotify替换成其他服务即可。其他服务有：
   1. tarslog
   2. tarsconfig
   3. tarsproperty
   4. tarsstat
   5. tarsquerystat
   6. tarsqueryproperty


   安装完成后恢复到默认的namespace：
   kubectl config set-context --current --namespace=default

## 镜像生成说明

`make registry` 生成registry的镜像

`make web` 生成tarsweb的镜像

`make img SERVER=XXX` 生成基础服务XXX的镜像

说明：cppregistry是原主控，后续可以合并到registry中。

## 清理所有tars相关基础服务
```
kubectl delete all --all -n tars-system
kubectl delete namespace tars-system
```