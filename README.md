# k8stars
K8stars is a convenient solution to run the tar service in kubernetes.
##Characteristics
-Maintain the native development framework capability of tars
-Automatic registration and configuration deletion of name service for tars
-Support smooth migration of original tar services to k8s and other container platforms
-Non intrusive design, no coupling relationship with operating environment
##Implementation scheme
1. Three interfaces are added in the tarsregistry, which are used for automatic registration / heartbeat reporting and node offline. For details, please refer to [interface definition] (. / tar registry / protocol)/ tarsregistry.tars )。
2. A 'tarscli' command-line tool is provided to allocate ports / generate configuration / report heartbeat and node offline.
##Deployment examples
1. Tarsregistry (name service)
Refer to the instructions in [baseserver] (. / baseserver) to deploy tar registry
2. Deployment service example
-Deploy sample simpleserver
```cd examples/simple &amp;&amp; kubectl apply -f  simpleserver.yaml ` ` ` `
Example description:
-The image is created by the 'examples / simple / dockerfile' file, and the basic image is created by 'CMD / tarscli / dockerfile'
- start.sh The 'tarscli genconf' in is used to generate the tars service startup configuration
- server_ meta.yaml The file is used to configure the metadata of the service. For field information, please refer to 'app / genconf'/ config.go `The 'serverconf' structure in. Endpoint defaults to 'TCP - h ${local_ ip} -p ${random_ Port} 'supports automatic filling of IP and random ports.
-HelloWorld program based on go language TestApp.HelloGo
See [examples] for details/ README.md ](examples)
3. Verify the deployment
Login ` DB_ Then execute 'select * from t'_ server_ The node information of simpleserver has been registered automatically.
##Tar deployment directory structure
`Tarscli ` based on environment variable 'tars_ Path '(default / tars) to manage services. The directory functions are as follows:
- `${TARS_ Path} / bin ': startup scripts and binaries
- `${TARS_ Path} / conf ': configuration file
- `${TARS_ Path} / log ': log file
- `${TARS_ Path} / data ': Runtime / cache file
##About tarscli
`Tarscli 'provides a set of command tools to facilitate container deployment of tar services. Parameters can be specified through environment variables. For details, see' tarscli help '.
Here are the sub commands supported by tarscli
-'genconf' is used to generate the startup configuration file of the TAS service. The supported environment variables are:
- `TARS_ Application 'the application name specified. By default, the`_ server_ meta.yaml `Read from
- `TARS_ Server 'is the service name specified by the`_ server_ meta.yaml `Read from
- `TARS_ BUILD_ Server 'the service name at compile time. It will be used when the compiled service name is different from the running service name
- `TARS_ Location 'can specify the address of registry. The default is' tars'- registry.default.svc . cluster.local  -P 17890 '(address of service)
- `TARS_ SET_ ID 'can specify service set
- `TARS_ MERGE_ Conf ` can specify the configuration template file and merge the configuration into the service startup configuration file
-The 'supervisor' executes the 'genconf' command by default, and then starts and monitors the service. The supported environment variables are:
- `TARS_ START_ The startup script of the path 'service is' $tars by default_ PATH/bin/ start.sh `
- `TARS_ STOP_ The stop script of the path 'service. By default, kill is in' $tars_ All service processes under path '
- `TARS_ REPORT_ Interval 'reports the interval between heartbeat and registry
- `TARS_ DISABLE_ Flow 'whether to enable traffic when registering with registry. If it is not empty, it means it is off. It is enabled by default
- `TARS_ CHECK_ Interval 'check the service status interval. If the status changes, it will be synchronized to the registry in real time
- `TARS_ BEFORE_ CHECK_ Script ` the shell command that runs before each check
- `TARS_ CHECK_ SCRIPT_ Timeout ` the timeout to run the shell command before each check
- `TARS_ PRESTOP_ Waittime 'turn off traffic - the waiting time before stopping the service. It is used for lossless changes. The default value is 80 seconds
-Hzcheck is used to synchronize the service status and the pod status of k8s. You need to set the 'readiness probe' of pod to tarscli 'hzcheck' command
-Prestop is used to delete the configuration corresponding to the registry before the service exits
- `TARS_ PRESTOP_ Waittime 'turn off traffic - the waiting time before stopping the service. It is used for lossless changes. The default value is 80 seconds
-Notify is used to send management commands. The common commands are: tars.setloglevel/tars . pprof, etc
##Basic services
Tars related basic services provide rich service governance functions. Please refer to [baseserver] (. / baseserver) for deployment.
