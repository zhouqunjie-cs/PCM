package ecser

import (
	"context"

	"github.com/JCCE-nudt/PCM/common/tenanter"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbecs"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbtenant"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

var (
	ErrEcsListNotSupported = errors.New("cloud not supported ecs list")
	ErrEcserPanic          = errors.New("ecs init panic")
)

type Ecser interface {
	CreateEcs(ctx context.Context, req *pbecs.CreateEcsReq) (resp *pbecs.CreateEcsResp, err error)    //创建ecs
	DeleteEcs(ctx context.Context, req *pbecs.DeleteEcsReq) (resp *pbecs.DeleteEcsResp, err error)    //批量删除ecs
	UpdateEcs(ctx context.Context, req *pbecs.UpdateEcsReq) (resp *pbecs.UpdateEcsResp, err error)    //修改ecs
	ListDetail(ctx context.Context, req *pbecs.ListDetailReq) (resp *pbecs.ListDetailResp, err error) //查询ecs详情
	ActionEcs(ctx context.Context, req *pbecs.ActionReq) (resp *pbecs.ActionResp, err error)          //操作ecs
}

func NewEcsClient(provider pbtenant.CloudProvider, region tenanter.Region, tenant tenanter.Tenanter) (ecser Ecser, err error) {
	// 部分sdk会在内部panic
	defer func() {
		if err1 := recover(); err1 != nil {
			glog.Errorf("NewEcsClient panic %v", err1)
			err = errors.WithMessagef(ErrEcserPanic, "%v", err1)
		}
	}()

	switch provider {
	case pbtenant.CloudProvider_ali:
		return newAliEcsClient(region, tenant)
	case pbtenant.CloudProvider_tencent:
		return newTencentCvmClient(region, tenant)
	case pbtenant.CloudProvider_huawei:
		return newHuaweiEcsClient(region, tenant)
		//TODO aws
	case pbtenant.CloudProvider_harvester:
		return newHarvesterClient(tenant)
	}

	err = errors.WithMessagef(ErrEcsListNotSupported, "cloud provider %v region %v", provider, region)
	return
}
