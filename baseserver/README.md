# TARS基础服务部署

## 部署步骤

1. 前置工作
   - 安装kubernetes，可以使用kubectl或其他方式来管理集群
   - 能执行docker命令行的终端

2. 部署tars db
   ```
   // 进入baseserver目录，获取tars的deploy相关文件
   cd baseserver
   make deploy
   
   // 创建一个MySQL数据库用于体验(在生产环境中，建议使用云db实例)
   kubectl apply -f yaml/db_all_in_one.yaml

   // 获取新仓库的db名
   kubectl get pods | grep tars-db

   // 修改 db/install_db_k8s.sh中的Pod名，然后导入数据
   sh db/install_db_k8s.sh
   ```
   对于已有的tars db，再请执行sql文件会清空原有的数据，只需要导入缺少的db即可。

3. tars registry
   使用`kubectl apply -f yaml/registry.yaml`部署tars registry。
   如果没用k8s创建的db，请修改`registry.yaml`中的数据地址。

4. tarsweb
   使用`kubectl apply -f yaml/tarsweb.yaml`部署。
   tarsweb默认使用3000端口，可以配置对应的访问方式，通过浏览器来访问。
   说明：当前tarsweb版本未兼容k8s中的场景，页面中有重启/停止等入口，但是操作会失败。

5. 部署其他服务
   以`tarsnotify`为例，使用`kubectl apply -f yaml/tarsnotify.yaml`来部署。其中数据库相关配置可以按需要替换。
   其他服务可以使用同一方式来部署，将tarsnotify替换成其他服务即可。其他服务有：
   1. tarslog
   2. tarsconfig
   3. tarsproperty
   4. tarsstat
   5. tarsquerystat
   6. tarsqueryproperty

## 镜像生成说明

`make registry` 生成registry的镜像
`make web` 生成tarsweb的镜像
`make img SERVER=XXX` 生成基础服务XXX的镜像

说明：cppregistry是原主控，后续可以合并到registry中。
