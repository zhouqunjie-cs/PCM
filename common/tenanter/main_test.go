package tenanter

import (
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbtenant"
	"os"
	"testing"
)

var (
	aliTenant []Tenanter
	tcTenant  []Tenanter
	hwTenant  []Tenanter
)

func TestMain(m *testing.M) {
	err := LoadCloudConfigs("../../config.yaml")
	if err != nil {
		panic(err)
	}

	if aliTenant, err = GetTenanters(pbtenant.CloudProvider_ali); err != nil {
		panic("get aliTenant failed")
	}
	if tcTenant, err = GetTenanters(pbtenant.CloudProvider_tencent); err != nil {
		panic("get tcTenantr failed")
	}
	if hwTenant, err = GetTenanters(pbtenant.CloudProvider_huawei); err != nil {
		panic("get hwTenant failed")
	}

	os.Exit(m.Run())
}
