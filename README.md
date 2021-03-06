[中文版](./README_cn.md) 

# K8STARS
K8STARS is a convenient solution to run TARS services in kubernetes.

## Characteristics
- Maintain the native development capability of TARS
- Automatic registration and configuration deletion of name service for TARS
- Support smooth migration of original TARS services to K8S and other container platforms
- Non intrusive design, no coupling relationship with operating environment

## How it works
1. Three interfaces are added in the tarsregistry, which are used for automatic registration, heartbeat reporting and node offline. For details, please refer to [interface definition](./tarsregistry/protocol/tarsregistry.tars)。

2. A 'tarscli' command-line tool is provided to allocate ports, generate configuration, report heartbeat and node offline.

## Deployment examples
1. Deployment tars basic service
```
curl https://raw.githubusercontent.com/TarsCloud/K8STARS/master/baseserver/install_all.sh | sh
```

2. Deployment service example
    -Deploy sample simpleserver

     ```cd examples/simple && kubectl apply -f  simpleserver.yaml```

     Example description:
     - The image is created by the `examples/simple/dockerfile` file, and the basic image is created by `cmd/tarscli/dockerfile`
     - start.sh: `tarscli genconf` in is used to generate the tars service startup configuration
     - server_ meta.yaml The file is used to configure the metadata of the service. For field information, please refer to `app/genconf/config.go` structure  `ServerConf` . Endpoint defaults to `tcp -h ${local_ip} -p ${random_port}` , supports automatic filling of IP and random ports.
     -ased on Golang HelloWorld program TestApp.HelloGo
     See [examples/README.md](examples)
     
3. Verify the deployment
Login `db_tars` , then execute `select * from t_server_conf\G` The node information of simpleserver has been registered automatically.

## Tars deployment directory structure
`tarscli` based on environment variable `TARS_PATH`(default/tars) to manage services. The directory functions are as follows:
   - `${TARS_PATH}/bin`：Startup scripts and binaries
   - `${TARS_PATH}/conf`：Configuration file
   - `${TARS_PATH}/log`： Log file
   - `${TARS_PATH}/data`：Runtime, Cache file

## About tarscli
`tarscli` provides a set of command tools to facilitate container deployment of TARS services. Parameters can be specified through environment variables. For details, see `tarscli help`.

Here are the sub commands supported by tarscli
- `genconf` is used to generate the startup configuration file of the TARS service. The supported environment variables are:
  - `TARS_APPLICATION` the application name specified. By default, the`_ server_ meta.yaml `Read from
  - `TARS_SERVER` is the service name specified by the`_ server_ meta.yaml `Read from
  - `TARS_BUILD_SERVER` the service name at compile time. It will be used when the compiled service name is different from the running service name
  - `TARS_LOCATOR` can specify the address of registry. The default is `tars.tarsregistry.QueryObj@tcp -h tars-registry.tars-system.svc.cluster.local -p 17890` (address of service)
  - `TARS_SET_ID` can specify service set
  - `TARS_MERGE_Conf` can specify the configuration template file and merge the configuration into the service startup configuration file

- `supervisor` executes the `genconf` command by default, and then starts and monitors the service. The supported environment variables are:
  - `TARS_START_PATH` The startup script of the service `$TARS_PATH/bin/start.sh`
  - `TARS_STOP_PATH` The stop script, by default, kill all service processes under path `$TARS_PATH`
  - `TARS_REPORT_INTERVAL` reports the interval heartbeat to registry
  - `TARS_DISABLE_FLOW` whether to enable traffic when registering with registry. If it is not empty, it means it is off. It is enabled by default
  - `TARS_CHECK_INTERVAL` check the service status interval. If the status changes, it will be synchronized to the registry in real time
  - `TARS_BEFORE_CHECK_SCRIPT` the shell command that runs before each check
  - `TARS_CHECK_SCRIPT_TIMEOUT` the timeout to run the shell command before each check
  - `TARS_PRESTOP_WAITTIME` turn off traffic - the waiting time before stopping the service. It is used for lossless changes. The default value is 80 seconds

- `hzcheck` is used to synchronize the service status and the pod status of k8s. You need to set the `readiness Probe` of pod to `tarscli hzcheck`  command
- `prestop` is used to delete the configuration corresponding to the registry before the service exits
    - `TARS_PRESTOP_WAITTIME`  turn off traffic - the waiting time before stopping the service. It is used for lossless changes. The default value is 80 seconds
- `notify` is used to send management commands. The common commands are: tars.setloglevel/tars.pprof, etc

## Basic services
TARS related basic services provide rich service governance functions. Please refer to [baseserver](./baseserver) for deployment.
