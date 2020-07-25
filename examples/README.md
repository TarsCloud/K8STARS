# 部署服务
这里举了两个例子：
- 部署示例simpleserver
- 基于Go语言的 HelloWorld 程序 TestApp.HelloGo
## 部署示例simpleserver
1. 执行：
   ```cd examples/simple && kubectl apply -f simpleserver.yaml```
   
   示例说明：
   - 镜像由examples/simple/Dockerfile文件制作，基础镜像由cmd/tarscli/Dockerfile制作
   - start.sh中的`tarscli genconf`用于生成tars服务启动配置
   - _server_meta.yaml文件用于配置服务的元数据，字段信息可以参考app/genconf/config.go中的`ServerConf`结构体。endpoint默认为`tcp -h ${local_ip} -p ${random_port}`，支持自动填入IP和随机端口。
   
2. 验证部署
   - 登录db_tars，然后执行`select * from t_server_conf\G`可以看到simpleserver的节点信息已自动注册。
   - 访问TarsWeb查看服务已部署成功。

## 基于Go语言的 HelloWorld 程序 TestApp.HelloGo
参照 [TarsGo开发环境部署](https://tarscloud.github.io/TarsDocs/env/tarsgo.html) 部署好开发环境

参照 [TarsGo快速入门](https://tarscloud.github.io/TarsDocs/hello-world/tarsgo.html) 的说明创建服务

### 创建服务

1. 运行create_tars_server.sh脚本，自动创建服务必须的文件

   ```shell
   sh $GOPATH/src/github.com/TarsCloud/TarsGo/tars/tools/create_tars_server.sh [App] [Server] [Servant]
     例如： 
   sh $GOPATH/src/github.com/TarsCloud/TarsGo/tars/tools/create_tars_server.sh TestApp HelloGo SayHello
   ```

2. 命令执行后将生成代码至GOPATH中，并以`APP/Server`命名目录，生成代码中也有提示具体路径。

   ```
   [root@1-1-1-1 src]# sh $GOPATH/src/github.com/TarsCloud/TarsGo/tars/tools/create_tars_server.sh TestApp HelloGo SayHello
   [create server: TestApp.HelloGo ...]
   [mkdir: /data/gopath/src/TestApp/HelloGo/]
   >>>Now doing:./config.conf >>>>
   >>>Now doing:./main.go >>>>
   >>>Now doing:./Servant_imp.go >>>>
   >>>Now doing:./makefile >>>>
   >>>Now doing:./Servant.tars >>>>
   >>>Now doing:./start.sh >>>>
   >>>Now doing:client/client.go >>>>
   >>>Now doing:vendor/vendor.json >>>>
   >>>Now doing:debugtool/dumpstack.go >>>>
   >>> Great！Done! You can jump in /data/gopath/src/TestApp/HelloGo
   >>> Tips: After editing the Tars file, execute the following cmd to automatically generate golang files.
   >>>       /data/gopath/bin/tars2go *.tars
   ```

### 定义接口文件

我们在 SayHello 接口下定义两个函数，一个`Add`和一个`Sub`，客户端请求参数是两个整型a和b，服务端计算a和b的运算结果赋值给c，并返回c的值。 

```go
module TestApp
{
	interface SayHello
	{
	    int Add(int a,int b,out int c); // Some example function
	    int Sub(int a,int b,out int c); // Some example function
	};
};
```

**注意** ： 参数中 **out** 修饰关键字标识输出参数。

### 服务端开发

1. 首先把tars协议文件转化为Golang语言形式

   ```shell
   cd $GOPATH/src/TestApp/HelloGo/
   $GOPATH/bin/tars2go SayHello.tars
   ```

2. 现在开始实现服务端的逻辑：

   ```vim $GOPATH/src/TestApp/HelloGo/sayhello_imp.go```

  - 在 `Add` 函数里计算 `a + b`，并赋值给c
  - 在 `Sub` 函数里计算 `a - b`，并赋值给c

      ```go
    package main
    
    import (
      "context"
    )
    
    // SayHelloImp servant implementation
    type SayHelloImp struct {
    }
    
    // Init servant init
    func (imp *SayHelloImp) Init() (error) {
      //initialize servant here:
      //...
      return nil
    }
    
    // Destroy servant destory
    func (imp *SayHelloImp) Destroy() {
      //destroy servant here:
      //...
    }
    
    func (imp *SayHelloImp) Add(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
      //Doing something in your function
      *c = a + b
      return 0, nil
    }
    func (imp *SayHelloImp) Sub(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
      //Doing something in your function
      *c = a - b
      return 0, nil
    }
      ```

**注意** ： 这里函数名要大写，Go语言方法导出规定。

3. 编辑main函数，初始代码已经由TARS框架实现了。

   ```vim $GOPATH/src/TestApp/HelloGo/main.go```

   ```go
   package main
   
   import (
     "fmt"
     "os"
   
     "github.com/TarsCloud/TarsGo/tars"
   
     "TestApp/HelloGo/TestApp" //Edit like this
   )
   
   func main() {
     // Get server config
     cfg := tars.GetServerConfig()
   
     // New servant imp
     imp := new(SayHelloImp)
     err := imp.Init()
     if err != nil {
       fmt.Printf("SayHelloImp init fail, err:(%s)\n", err)
       os.Exit(-1)
     }
     // New servant
     app := new(TestApp.SayHello)
     // Register Servant
     app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".SayHelloObj")
   
     // Run application
     tars.Run()
   }
   ```

**注意** ：将 `import` 里的 `TestApp` 改为 `TestApp/HelloGo/TestApp`

4. 将 [HelloGo的Demo](TestApp/HelloGo) 下的

    `Dockerfile`、`makefile`、`_server_meta.yaml`、`simpleserver.yaml`、`start.sh` 

   五个文件拷贝到`$GOPATH/src/TestApp/HelloGo/`下

5. 修改`_server_meta.yaml`文件，根据实际情况填写 `application`、`server`、`object` 字段的值。

   `vim $GOPATH/src/TestApp/HelloGo/_server_meta.yaml`

   ```yaml
   version: v1
   application: TestApp
   server: HelloGo
   adapters:
     - object: SayHelloObj
   ```

6. 修改 `start.sh`，根据实际情况修改服务名及配置名

   `vim $GOPATH/src/TestApp/HelloGo/start.sh`

   ```shell
   #!/bin/bash
   
   # start server
   ${TARS_PATH}/bin/HelloGo --config=${TARS_PATH}/conf/HelloGo.conf
   ```

7. 修改 `makefile`，根据实际情况填写 `APP`、`TARGET`

   `vim $GOPATH/src/TestApp/HelloGo/makefile`

   ```makefile
   APP    := TestApp
   TARGET := HelloGo
   
   GOBUILD      := go build
   DOCKER_BUILD := docker build
   
   REPO         ?= ccr.ccs.tencentyun.com/tarsbase
   VERSION  ?= $(shell date "+%Y%m%d%H%M%S") #Edit if you need
   
   LOWCASE_TARGET := $(shell echo $(TARGET) | tr '[:upper:]' '[:lower:]')
   IMG_REPO       := $(REPO)/$(LOWCASE_TARGET)
   
   build:
   	GOOS=linux $(GOBUILD) -o $(TARGET)
   
   img: build
   	$(DOCKER_BUILD) --build-arg SERVER=$(TARGET) -t $(IMG_REPO):$(VERSION) .
   
   tgz: build
   	tar czf $(TARGET).tgz $(TARGET) _server_meta.yaml
   
   patch: tgz
   	curl --data-binary @$(TARGET).tgz "${TARS_EP}/patch?server=$(TARGET)&version=$(VERSION)"
   	
   stdout:
   	@curl "${TARS_EP}/stdout?server=$(TARGET)"
   
   listlog:
   	@curl "${TARS_EP}/listlog?app=$(APP)&server=$(TARGET)"
   
   tailog:
   	@curl "${TARS_EP}/tailog?app=$(APP)&server=$(TARGET)&filename=$(LOG_NAME)"
   
   clean:
   	rm -rf $(TARGET)
   ```

8. 制作镜像

   ```shell
   cd $GOPATH/src/TestApp/HelloGo/ && make img
   ```

   ```
   [root@1-1-1-1 HelloGo]# make img
   GOOS=linux go build -o HelloGo
   docker build --build-arg SERVER=HelloGo -t ccr.ccs.tencentyun.com/tarsbase/hellogo:20200725145313 .
   Sending build context to Docker daemon  23.44MB
   Step 1/5 : FROM ccr.ccs.tencentyun.com/tarsbase/tarscli:latest
    ---> ef39bfa9b5e6
   Step 2/5 : ARG SERVER=please_build_by_make_img
    ---> Using cache
    ---> c3a89b1e1be5
   Step 3/5 : ENV TARS_BUILD_SERVER ${SERVER}
    ---> Using cache
    ---> 5f8fc936380e
   Step 4/5 : COPY $SERVER _server_meta.yaml start.sh $TARS_PATH/bin/
    ---> Using cache
    ---> 05f31e9e6e4a
   Step 5/5 : CMD tarscli supervisor
    ---> Using cache
    ---> be8b4a482834
   Successfully built be8b4a482834
   Successfully tagged ccr.ccs.tencentyun.com/tarsbase/hellogo:20200725145313
   ```

   记住制作好的镜像名，本例为： `ccr.ccs.tencentyun.com/tarsbase/hellogo:20200725145313`

9. 修改 `simpleserver.yaml` 文件，并将`TestApp.HelloGo`部署至k8s集群

   `vim $GOPATH/src/TestApp/HelloGo/simpleserver.yaml` 修改`image`的值为上一步制作好的镜像名

   ```yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: hellogo
   spec:
     selector:
       matchLabels:
         app: hellogo
     replicas: 1
     template:
       metadata:
         labels:
           app: hellogo
       spec:
         containers:
         - name: hellogo
           image: ccr.ccs.tencentyun.com/tarsbase/hellogo:20200725145313 #Edit here to the image you just build
           imagePullPolicy: IfNotPresent 
           lifecycle:
             preStop:
               exec:
                 command: ["tarscli", "prestop"]
           readinessProbe:
             exec:
               command: ["tarscli", "hzcheck"]
             initialDelaySeconds: 2
             timeoutSeconds: 8
             periodSeconds: 6
   ```

   执行 `kubectl apply -f simpleserver.yaml`

9. 验证部署
   - 登录db_tars，然后执行`select * from t_server_conf\G`可以看到HelloGo的节点信息已自动注册。
   - 访问TarsWeb查看服务已部署成功。

### 客户端开发

[client](TestApp/HelloGo/client) 目录下提供了一个客户端的demo

```go
package main

import (
	"fmt"

	"github.com/TarsCloud/TarsGo/tars"

	"TestApp/HelloGo/TestApp"
)

func main() {
	comm := tars.NewCommunicator()
  //If your server "HelloGo" has been registered to tarsregistry
  obj := fmt.Sprintf("TestApp.HelloGo.SayHelloObj")
  //If your server "HelloGo" hasn't been registerd to tarsregistry
  //obj := fmt.Sprintf("TestApp.HelloNew.SayHelloObj@tcp -h 127.0.0.1 -p 10015 -t 60000")
  //"127.0.0.1" and "10015" should be change to your HelloGo's IP and Port.
	app := new(TestApp.SayHello)

	//If your server has been registered to tarsregistry
  comm.SetProperty("locator", "tars.tarsregistry.QueryObj@tcp -h 192.168.0.55 -p 17890")
  //"192.168.0.55" should be change to your tarsregistry's IP

	comm.StringToProxy(obj, app)
	var out, i int32
	i = 123
	ret, err := app.Add(i, i*2, &out)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret, out)
}
```

编译测试

```shell
[root@1-1-1-1 client]# go build client.go
[root@1-1-1-1 client]# ./client 
0 369
```




