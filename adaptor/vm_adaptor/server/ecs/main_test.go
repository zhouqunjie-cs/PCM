package ecs

import (
	"github.com/JCCE-nudt/PCM/common/tenanter"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbtenant"
	"os"
	"testing"
)

var (
	aliTenant, tcTenant, hwTenant, k8sTenant []tenanter.Tenanter
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
	if k8sTenant, err = tenanter.GetTenanters(pbtenant.CloudProvider_k8s); err != nil {
		panic("get awsTenant failed")
	}
	os.Exit(m.Run())
}
