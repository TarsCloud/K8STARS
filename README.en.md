# k8stars
K8stars is a solution to deploy TARS in a Kubernetes environment and it has the following features: 
* Maintain the same functionality of TARS Framework
* Support TARS naming service and configuration management 
* Support a smooth transfer of TARS services into container orchestrations such as K8S. 
* Non-invasive; there isn’t a coupling between the deployment environment and the services
## The solution
1. Add 3 interfaces in tarsregistry for naming service, heartbeat report, and disabling nodes. 
2. Build a `tarscli` (command line tool) for assigning ports, generating configurations, reporting heartbeats, and retiring nodes. 
## Deployment example
1. Please refer to [baseserver](https://github.com/TarsCloud/K8STARS/blob/master/baseserver/README.md) for deploying tarsregistry 
2. Run:
  `cd examples/simple && kubectl apply -f simpleserver.yaml`
* example description：
  * `examples/simple/Dockerfile` will create the container image, and `cmd/tarscli/Dockerfile` will build the fundamental image
  * `tarscli genconf` in start.sh will be used to generate tars services’ starting configuration. 
  * The file _server_meta.yaml is used to configure services’ metadata. You can refer to the structure of ServerConf in `app/genconf/config.go` for the field information. The endpoint is `tcp -h ${local_ip} -p ${random_port}` by default and supports automatically filling in IP and random ports. 
  * You can also find a demo HelloWorld program TestApp.Hello in Go language in the repository
    * Please refer to [examples/README.md](https://github.com/TarsCloud/K8STARS/tree/master/examples)
3. To confirm if deployment is successfully, log into db_tars and execute `select * from t_server_conf\G`. You should see that the simpleserver’s node information is registered. 
## TARS Deployment Directories
`tarscli` manages services based on the environment variables in `TARS_PATH`, the function of each directory is as follows: 
* `${TARS_PATH}/bin`: to start scripts and binary files
* `${TARS_PATH}/conf`: to configure files
* `${TARS_PATH}/log`: to store log files
* `${TARS_PATH}/data`: to provide execution status and cache files
 
## tarscli description
`tarscli` provides a set of command tools to help with tars services’ containerized deployment. If you want to designate the parameters, please refer to `tarscli help`. Below are the sub commands tarscli supports: 
* `genconf` is used to create tars services’ starting configuration files. The supported environment variables include: 
  * `TARS_APPLICATION` : the application name, which comes from _server_meta.yaml by default. 
  * `TARS_SERVER`: the server name, which is read from _server_meta.yaml. 
  * `TARS_BUILD_SERVER`: the server name when compiling. It is used when the compiled server name is different from the name of the server that is running. 
  * `TARS_LOCATOR` can designate the address of registry, and the default path is tars-registry.default.svc.cluster.local -p 17890 (service’s address) 
  * `TARS_SET_ID`: can designate the service set. 
  * `TARS_MERGE_CONF` can configure the template file, and merge it into the service’s starting configuration file. 
* `Supervisor`, by default, executes genconf commands first, and then starts/monitors services. The supported environment variables include: 
  * `TARS_START_PATH`: Services’ starting script and it is TARS_PATH/bin/start.sh by default
  * `TARS_STOP_PATH`: services’ stopping script, and by default kills all the processes under $TARS_PATH 
  * `TARS_REPORT_INTERVAL` is the time it takes to report heartbeat to registry 
  * `TARS_DISABLE_FLOW`is disabled by default
  * `TARS_CHECK_INTERVAL`: the time it takes to check the server status. If the status changes, it will be synced to registry.
  * `TARS_BEFORE_CHECK_SCRIPT`: it is the shell script that runs before each examination. 
  * `TARS_CHECK_SCRIPT_TIMEOUT`: the timeout before each shell script execution. 
  * `TARS_PRESTOP_WAITTIME`: the wait time period before shutting-stopping services and it makes sure nothing is lost afterward. The default is set to 80 seconds
* `hzcheck` is used to synced service status and the status of k8s’ pod and it needs to set pod’s readiness probe to tarscli hzcheck.
* `prestop` is used to delete the corresponding configurations before services end
  * `TARS_PRESTOP_WAITTIME`: the wait time period before shutting-stopping services and it makes sure nothing is lost afterward. The default is set to 80 seconds
* `notify` is used to send management commands. An example of commonly used commands is tars.setloglevel/tars.pprof.
 
 
