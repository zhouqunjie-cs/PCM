package poder

import (
	"os"
	"testing"

	"github.com/zhouqunjie-cs/PCM/common/tenanter"
	"github.com/zhouqunjie-cs/PCM/lan_trans/idl/pbtenant"
)

var (
	aliTenant, tcTenant, hwTenant, awsTenant []tenanter.Tenanter
)

func TestMain(m *testing.M) {
	err := tenanter.LoadCloudConfigs("../../../config.yaml")
	if err != nil {
		panic(err)
	}
	if aliTenant, err = tenanter.GetTenanters(pbtenant.CloudProvider_ali); err != nil {
		panic("get aliTenant failed")
	}
	if tcTenant, err = tenanter.GetTenanters(pbtenant.CloudProvider_tencent); err != nil {
		panic("get tcTenant failed")
	}
	if hwTenant, err = tenanter.GetTenanters(pbtenant.CloudProvider_huawei); err != nil {
		panic("get hwTenant failed")
	}
	if awsTenant, err = tenanter.GetTenanters(pbtenant.CloudProvider_k8s); err != nil {
		panic("get awsTenant failed")
	}
	os.Exit(m.Run())
}
