package poder

import (
	"context"

	"github.com/JCCE-nudt/PCM/common/tenanter"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbpod"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbtenant"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

var (
	ErrPodListNotSupported = errors.New("cloud not supported pod list")
	ErrPoderPanic          = errors.New("pod init panic")
)

type Poder interface {
	ListPodDetail(ctx context.Context, req *pbpod.ListPodDetailReq) (resp *pbpod.ListPodDetailResp, err error)
	CreatePod(ctx context.Context, req *pbpod.CreatePodReq) (resp *pbpod.CreatePodResp, err error)
	DeletePod(ctx context.Context, req *pbpod.DeletePodReq) (*pbpod.DeletePodResp, error)
	UpdatePod(ctx context.Context, req *pbpod.UpdatePodReq) (*pbpod.UpdatePodResp, error)
	GetPodRegion(ctx context.Context, req *pbpod.GetPodRegionReq) (*pbpod.GetPodRegionResp, error)
}

func NewPodClient(provider pbtenant.CloudProvider, region tenanter.Region, tenant tenanter.Tenanter) (poder Poder, err error) {
	// 部分sdk会在内部panic
	defer func() {
		if err1 := recover(); err1 != nil {
			glog.Errorf("NewPodClient panic %v", err1)
			err = errors.WithMessagef(ErrPoderPanic, "%v", err1)
		}
	}()

	switch provider {
	case pbtenant.CloudProvider_ali:
		return newAliEciClient(region, tenant)
	case pbtenant.CloudProvider_tencent:
		return newTencentEksClient(region, tenant)
	case pbtenant.CloudProvider_huawei:
		return newHuaweiCciClient(region, tenant)
	case pbtenant.CloudProvider_k8s:
		return newK8SClient(tenant)
		//TODO aws
		//case pbtenant.CloudProvider_aws:
		//	return newAwsPodClient(region, tenant)
	}

	err = errors.WithMessagef(ErrPodListNotSupported, "cloud provider %v region %v", provider, region)
	return
}
