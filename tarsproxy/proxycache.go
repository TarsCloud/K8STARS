package tarsproxy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/endpointf"
	"github.com/tarscloud/k8stars/algorithm/retry"
	"github.com/tarscloud/k8stars/consts"
	"github.com/tarscloud/k8stars/logger"
)

var (
	comm = tars.NewCommunicator()
	log  = logger.GetLogger()

	rpcTimeout = 1000

	cacheTime = time.Second * 10
	tryOpt    = retry.New(retry.MaxTimeoutOpt(time.Second*3, time.Second*1))
)

// StringToProxy sets the servant of ProxyPrx p with file-cached endpoints
func StringToProxy(locator, obj string, p tars.ProxyPrx) error {
	obj = getEndpoints(locator, obj)
	if obj == "" || !strings.Contains(obj, "@") {
		return fmt.Errorf("get object failed %s", obj)
	}
	comm.StringToProxy(obj, p)
	return nil
}

func getEndpoints(locator, obj string) string {
	if strings.Contains(obj, "@") {
		return obj
	}
	retObj := ""
	os.MkdirAll(consts.TarsPath+"/data", 0755)
	cacheFile := filepath.Join(consts.TarsPath, "data", "cache_"+obj)
	if bs, err := ioutil.ReadFile(cacheFile); err == nil {
		retObj = string(bs)
		if st, _ := os.Stat(cacheFile); st != nil && st.ModTime().Add(cacheTime).After(time.Now()) {
			return retObj // use cache
		}
	}
	log.Debugf("Start get endpoints for %s", obj)
	client := GetQueryClient(locator)

	activeEp := make([]endpointf.EndpointF, 0)
	inactiveEp := make([]endpointf.EndpointF, 0)
	err := tryOpt(func() error {
		ret, err := client.FindObjectByIdInSameGroup(obj, &activeEp, &inactiveEp)
		if retObj != "" { // use cache
			if err != nil {
				log.Errorf("FindObjectByIdInSameGroup error(use cache) %v", err)
			}
			return nil
		}
		if err != nil || ret != 0 {
			return fmt.Errorf("findObjectByIdInSameGroup error %d %v", ret, err)
		}
		if len(activeEp) == 0 {
			return fmt.Errorf("empty active endpoints")
		}
		return nil
	})
	if err != nil || len(activeEp) == 0 {
		log.Errorf("Get enpoints failed %v", err)
		return obj
	}
	retObj = obj + "@"
	for i, ep := range activeEp {
		protocol := "udp"
		if ep.Istcp == 1 {
			protocol = "tcp"
		}
		if i > 0 {
			retObj += ":"
		}
		retObj += fmt.Sprintf("%s -h %s -p %d -t %d",
			protocol, ep.Host, ep.Port, ep.Timeout)
	}
	_ = ioutil.WriteFile(cacheFile, []byte(retObj), 0644)
	return retObj
}
