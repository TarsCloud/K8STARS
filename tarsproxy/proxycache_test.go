package tarsproxy

import (
	"os"
	"testing"

	"github.com/tarscloud/k8stars/consts"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/endpointf"
	"github.com/stretchr/testify/assert"
)

func TestGetEndpoints(t *testing.T) {
	consts.TarsPath = "tmp_TestGetEndpoints"
	defer os.RemoveAll("tmp_TestGetEndpoints")
	mockQuery = &mockQueryClient{}
	obj := getEndpoints("12abc@aa", "obj1")
	assert.Equal(t, "obj1@udp -h localhost -p 9987 -t 50000", obj)
}

type mockQueryClient struct{}

func (m *mockQueryClient) FindObjectByIdInSameGroup(Id string, ActiveEp *[]endpointf.EndpointF,
	InactiveEp *[]endpointf.EndpointF, _opt ...map[string]string) (ret int32, err error) {
	if Id == "obj1" {
		eps := make([]endpointf.EndpointF, 1)
		eps[0].Host = "localhost"
		eps[0].Istcp = 0
		eps[0].Port = 9987
		eps[0].Timeout = 50000
		*ActiveEp = eps
		return
	}
	return 0, nil
}
