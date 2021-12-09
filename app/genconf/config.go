package genconf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tarscloud/k8stars/consts"
	"github.com/tarscloud/k8stars/mergeconf"
	"gopkg.in/yaml.v2"
)

// Adapter is the config of tars adapter
type Adapter struct {
	Object       string `yaml:"object"`
	Endpoint     string `yaml:"endpoint"`
	MaxConns     int    `yaml:"maxconns"`
	Protocol     string `yaml:"protocol"`
	QueueCap     int    `yaml:"queuecap"`
	QueueTimeout int    `yaml:"queuetimeout"`
	Threads      int    `yaml:"threads"`
}

// ServerConf is the config of tars server
type ServerConf struct {
	// global config
	Application string `yaml:"application"`
	Server      string `yaml:"server"`
	SetID       string `yaml:"set_id"`

	// client config
	Locator string `yaml:"locator"`

	SampleRate              int    `yaml:"sample_rate"`
	MaxSampleCount          int    `yaml:"max_sample_count"`
	StatObj                 string `yaml:"statObj"`
	PropertyObj             string `yaml:"propertyObj"`
	AsyncThreadNum          int    `yaml:"asyncThreadNum"`
	SyncInvokeTimeout       int    `yaml:"sync_invoke_timeout"`
	AsyncInvokeTimeout      int    `yaml:"async_invoke_timeout"`
	ReportInterval          int    `yaml:"report_interval"`
	RefreshEndpointInterval int    `yaml:"refresh_endpoint_interval"`

	// server config
	Adapters            []Adapter `yaml:"adapters"`
	LocalEndpoint       string    `yaml:"local"`
	LogSize             string    `yaml:"logSize"`
	LogNum              int       `yaml:"logNum"`
	ConfigObj           string    `yaml:"configObj"`
	NotifyObj           string    `yaml:"notifyObj"`
	LogObj              string    `yaml:"logObj"`
	LogLevel            string    `yaml:"logLevel"`
	DeactivatingTimeout int       `yaml:"deactivating_timeout"`
}

// GlobalConf is all of config for server
type GlobalConf struct {
	Conf      ServerConf `yaml:"conf"`
	EnableSet string     `yaml:"-"`
	Servant   string     `yaml:"-"`
	LocalIP   string     `yaml:"local_ip"`
	TarsPath  string     `yaml:"tars_path"`
}

// GetGlobalConf returns global config
func GetGlobalConf() (*GlobalConf, error) {
	// get config from file
	sConfFile := filepath.Join(consts.TarsPath, "data", consts.ServerInfoFile)
	sConfBytes, err := ioutil.ReadFile(sConfFile)
	if err != nil {
		return nil, err
	}
	gConf := &GlobalConf{}
	if err := yaml.Unmarshal(sConfBytes, gConf); err != nil {
		return nil, err
	}
	return gConf, nil
}

func parseServerConf(f string, conf *ServerConf) error {
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bs, conf)
	return err
}
func defaultServerConf() ServerConf {
	return ServerConf{
		// client config
		Locator:                 "tars.tarsregistry.QueryObj@tcp -h tars-registry.${namespace}.svc.cluster.local -p 17890",
		SampleRate:              100000,
		MaxSampleCount:          50,
		StatObj:                 "tars.tarsstat.StatObj",
		RefreshEndpointInterval: 60000,
		PropertyObj:             "tars.tarsproperty.PropertyObj",
		AsyncThreadNum:          3,
		SyncInvokeTimeout:       3000,
		AsyncInvokeTimeout:      5000,
		ReportInterval:          6000,
		// server config
		LogSize:             "100M",
		LogNum:              5,
		ConfigObj:           "tars.tarsconfig.ConfigObj",
		NotifyObj:           "tars.tarsnotify.NotifyObj",
		LogObj:              "tars.tarslog.LogObj",
		LogLevel:            "DEBUG",
		DeactivatingTimeout: 3000,
		LocalEndpoint:       "tcp -h 127.0.0.1 -p ${random_port} -t 5000",
	}
}

// UnmarshalYAML defines a custom unmarsher for setting default values
func (s *Adapter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawObj Adapter
	raw := rawObj{
		MaxConns:     200000,
		Protocol:     "tars",
		QueueCap:     10000,
		QueueTimeout: 60000,
		Threads:      5,
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	*s = Adapter(raw)
	return nil
}

func generateConf(sConf *ServerConf, buildServer, mergeConf string) error {
	// fill endpoint
	for i := range sConf.Adapters {
		sv := &sConf.Adapters[i]
		if sv.Endpoint == "" {
			sv.Endpoint = "tcp -h ${local_ip} -p ${random_port} -t 30000"
		}
		if strings.Contains(sv.Endpoint, "${local_ip}") {
			sv.Endpoint = strings.ReplaceAll(sv.Endpoint, "${local_ip}", consts.LocalIP)
		}
		if strings.Contains(sv.Endpoint, "${random_port}") {
			port, err := getRandomPort(sv.Endpoint, consts.RandPortMin, consts.RandPortMax)
			if err != nil {
				return err
			}
			sv.Endpoint = strings.ReplaceAll(sv.Endpoint, "${random_port}", port)
		}
	}
	// fill locator
	if strings.Contains(sConf.Locator, "${namespace}") {
		sConf.Locator = strings.ReplaceAll(sConf.Locator, "${namespace}", consts.NameSpace)
	}

	// fill local endpoint
	if strings.Contains(sConf.LocalEndpoint, "${random_port}") {
		port, err := getRandomPort(sConf.LocalEndpoint, consts.RandPortMin, consts.RandPortMax)
		if err != nil {
			return err
		}
		sConf.LocalEndpoint = strings.ReplaceAll(sConf.LocalEndpoint, "${random_port}", port)
	}

	setName := sConf.SetID
	enableSet := "y"
	if setName == "" {
		setName = "NULL"
		enableSet = "n"
	}
	servant := ""
	for _, sv := range sConf.Adapters {
		obj := fmt.Sprintf("%s.%s.%s", sConf.Application, sConf.Server, sv.Object)
		servant += fmt.Sprintf(`        <%sAdapter>
		    allow
			handlegroup=%sAdapter
			servant=%s
			endpoint=%s
			maxconns=%d
			protocol=%s
			queuecap=%d
			queuetimeout=%d
			threads=%d
        </%sAdapter>
`, obj, obj, obj, sv.Endpoint, sv.MaxConns, sv.Protocol, sv.QueueCap, sv.QueueTimeout, sv.Threads, obj)
	}

	tpl, err := template.New("tpl").Parse(templateStr)
	if err != nil {
		return err
	}
	gConf := GlobalConf{
		EnableSet: enableSet,
		Servant:   servant,
		Conf:      *sConf,
		TarsPath:  consts.TarsPath,
		LocalIP:   consts.LocalIP,
	}
	outBuf := bytes.NewBuffer(nil)
	if err := tpl.Execute(outBuf, gConf); err != nil {
		return err
	}
	if mergeConf != "" {
		bs, err := ioutil.ReadFile(mergeConf)
		if err != nil {
			return err
		}
		out, err := mergeconf.MergeConf(outBuf.Bytes(), bs)
		if err != nil {
			return err
		}
		outBuf.Reset()
		outBuf.Write([]byte(out))
	}
	outConf := filepath.Join(consts.TarsPath, "conf", buildServer+".conf")
	err = ioutil.WriteFile(outConf, outBuf.Bytes(), 0644)

	gConfFile := filepath.Join(consts.TarsPath, "data", consts.ServerInfoFile)
	outBytes, err := yaml.Marshal(&gConf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(gConfFile, outBytes, 0644)
	return err
}

var templateStr = `<tars>
<application>
    #是否启用SET分组
	enableset={{.EnableSet}}
	
    #SET分组的全名.(mtt.s.1)
    setdivision={{.Conf.SetID}}
    <client>
        # registry地址
		locator={{.Conf.Locator}}
		
        #同步调用超时时间,缺省3s(毫秒)
		sync-invoke-timeout={{.Conf.SyncInvokeTimeout}}
		
        #异步超时时间,缺省5s(毫秒)
		async-invoke-timeout={{.Conf.AsyncInvokeTimeout}}
		
        #重新获取服务列表时间间隔(毫秒)
		refresh-endpoint-interval={{.Conf.RefreshEndpointInterval}}
		
        #模块间调用服务[可选]
        stat ={{.Conf.StatObj}}

        #属性上报服务[可选]
        property={{.Conf.PropertyObj}}

        #上报间隔时间,默认60s(毫秒)
        report-interval={{.Conf.ReportInterval}}

        #stat采样比1:n 例如sample-rate为1000时 采样比为千分之一
        sample-rate={{.Conf.SampleRate}}

        #1分钟内stat最大采样条数
        max-sample-count={{.Conf.MaxSampleCount}}

        #网络异步回调线程个数
        asyncthread={{.Conf.AsyncThreadNum}}

        #模块名称
        modulename={{.Conf.Application}}.{{.Conf.Server}}
    </client>
        
    <server>
        #应用名称
        app={{.Conf.Application}}

        #服务名称
        server={{.Conf.Server}}

        #本地ip
		localip={{.LocalIP}}
		
        #本地ip
        local={{.Conf.LocalEndpoint}}

        # servant列表信息
{{.Servant}}

        # 服务程序目录
        basepath={{.TarsPath}}/bin/
        
        # 服务的数据目录,可执行文件,配置文件等
        datapath={{.TarsPath}}/data/

        # 日志路径
        logpath ={{.TarsPath}}/log/

        #滚动日志等级默认值
		logLevel={{.Conf.LogLevel}}
		
        # 日志大小
        logsize={{.Conf.LogSize}}

        # 日志数量
        lognum={{.Conf.LogNum}}

        # 配置中心的地址[可选]
        config={{.Conf.ConfigObj}}

        # 信息中心的地址[可选]
        notify={{.Conf.NotifyObj}}

        # 远程LogServer[可选]
        log={{.Conf.LogObj}}

        #关闭服务时等待时间
        deactivating-timeout={{.Conf.DeactivatingTimeout}}
    </server>          
</application>
</tars>
`
