package ecser

import (
	"context"
	"strconv"

	"github.com/JCCE-nudt/PCM/common/tenanter"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbecs"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbtenant"
	string_ "github.com/alibabacloud-go/darabonba-string/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type TencentCvm struct {
	cli      *cvm.Client
	region   tenanter.Region
	tenanter tenanter.Tenanter
}

func newTencentCvmClient(region tenanter.Region, tenant tenanter.Tenanter) (Ecser, error) {
	var (
		client *cvm.Client
		err    error
	)

	switch t := tenant.(type) {
	case *tenanter.AccessKeyTenant:
		client, err = cvm.NewClient(common.NewCredential(t.GetId(), t.GetSecret()), region.GetName(), profile.NewClientProfile())
	default:
	}

	if err != nil {
		return nil, errors.Wrap(err, "init tencent cvm client error")
	}
	return &TencentCvm{
		cli:      client,
		region:   region,
		tenanter: tenant,
	}, nil
}

func (ecs *TencentCvm) CreateEcs(ctx context.Context, req *pbecs.CreateEcsReq) (*pbecs.CreateEcsResp, error) {
	ImageId := req.GetImageId()
	InstanceType := req.GetInstanceType()
	InstanceName := req.GetInstanceName()
	Amount := int64(req.GetAmount())
	DryRun := req.GetDryRun()
	InstanceChargeType := req.GetInstanceChargeType()
	ZoneId := req.GetZoneId()
	request := cvm.NewRunInstancesRequest()
	request.ImageId = &ImageId
	request.InstanceType = &InstanceType
	request.InstanceName = &InstanceName
	request.DryRun = util.EqualString(&DryRun, tea.String("true"))
	request.InstanceChargeType = &InstanceChargeType
	request.InstanceCount = &Amount
	request.Placement = &cvm.Placement{
		Zone: &ZoneId,
	}
	resp, err := ecs.cli.RunInstances(request)
	if err != nil {
		return nil, errors.Wrap(err, "Tencent Create ECS error")
	}
	InstanceIds := make([]string, 0)
	for _, v := range resp.Response.InstanceIdSet {
		InstanceIds = append(InstanceIds, *v)
	}
	isFinished := false
	if len(resp.Response.InstanceIdSet) > 0 {
		isFinished = true
	}
	glog.Infof("--------------------腾讯ECS实例创建成功--------------------")
	glog.Infof(*util.ToJSONString(util.ToMap(resp)))
	return &pbecs.CreateEcsResp{
		RequestId:      "Tencent ECS RequestId: " + *resp.Response.RequestId,
		InstanceIdSets: InstanceIds,
		Finished:       isFinished,
	}, nil
}

func (ecs *TencentCvm) DeleteEcs(ctx context.Context, req *pbecs.DeleteEcsReq) (*pbecs.DeleteEcsResp, error) {
	idStr := req.GetInstanceIds()
	InstanceIds := string_.Split(&idStr, tea.String(","), tea.Int(-1))
	//腾讯云支持批量操作，每次请求批量实例的上限为100
	if len(InstanceIds) > 100 {
		return nil, errors.New("Tencent Delete ECS error InstanceIds > 100")
	}
	deleteReq := cvm.NewTerminateInstancesRequest()
	deleteReq.InstanceIds = InstanceIds
	resp, err := ecs.cli.TerminateInstances(deleteReq)
	if err != nil {
		return nil, errors.Wrap(err, "Tencent Delete ECS error")
	}
	glog.Infof("--------------------腾讯ECS实例释放成功--------------------")
	glog.Infof(*util.ToJSONString(util.ToMap(resp)))
	return &pbecs.DeleteEcsResp{
		Provider:    req.Provider,
		RequestId:   *resp.Response.RequestId,
		AccountName: req.AccountName,
		RegionId:    req.RegionId,
	}, nil
}

func (ecs *TencentCvm) UpdateEcs(ctx context.Context, req *pbecs.UpdateEcsReq) (*pbecs.UpdateEcsResp, error) {
	InstanceName := req.GetInstanceName()
	InstanceId := req.GetInstanceIds()
	UpdateReq := cvm.NewModifyInstancesAttributeRequest()
	if req.GetInstanceIds() == "" {
		return nil, errors.New("InstanceId is empty")
	}
	UpdateReq.InstanceIds = string_.Split(&InstanceId, tea.String(","), tea.Int(-1))
	if InstanceName != "" {
		UpdateReq.InstanceName = &InstanceName
	}
	resp, err := ecs.cli.ModifyInstancesAttribute(UpdateReq)
	if err != nil {
		return nil, errors.Wrap(err, "Tencent Update ECS error")
	}
	glog.Infof("--------------------腾讯ECS实例修改成功--------------------")
	glog.Infof(*util.ToJSONString(util.ToMap(resp)))
	return &pbecs.UpdateEcsResp{
		RequestId: *resp.Response.RequestId,
	}, nil
}

func (ecs *TencentCvm) ListDetail(ctx context.Context, req *pbecs.ListDetailReq) (*pbecs.ListDetailResp, error) {
	request := cvm.NewDescribeInstancesRequest()
	request.Offset = common.Int64Ptr(int64((req.PageNumber - 1) * req.PageSize))
	request.Limit = common.Int64Ptr(int64(req.PageSize))
	resp, err := ecs.cli.DescribeInstances(request)
	if err != nil {
		return nil, errors.Wrap(err, "Tencent ListDetail error")
	}

	var ecses = make([]*pbecs.EcsInstance, len(resp.Response.InstanceSet))
	for k, v := range resp.Response.InstanceSet {
		ExpiredTime := ""
		if v.ExpiredTime != nil {
			ExpiredTime = *v.ExpiredTime
		}
		ecses[k] = &pbecs.EcsInstance{
			Provider:           pbtenant.CloudProvider_tencent,
			AccountName:        ecs.tenanter.AccountName(),
			InstanceId:         *v.InstanceId,
			InstanceName:       *v.InstanceName,
			RegionName:         ecs.region.GetName(),
			PublicIps:          make([]string, len(v.PublicIpAddresses)),
			InstanceType:       *v.InstanceType,
			Cpu:                strconv.FormatInt(*v.CPU, 10),
			Memory:             strconv.FormatInt(*v.Memory, 10),
			Description:        "",
			Status:             *v.InstanceState,
			CreationTime:       *v.CreatedTime,
			ExpireTime:         ExpiredTime,
			InnerIps:           make([]string, len(v.PrivateIpAddresses)),
			VpcId:              *v.VirtualPrivateCloud.VpcId,
			ResourceGroupId:    "",
			InstanceChargeType: *v.InstanceChargeType,
		}
		for k1, v1 := range v.PublicIpAddresses {
			ecses[k].PublicIps[k1] = *v1
		}
		for k1, v1 := range v.PrivateIpAddresses {
			ecses[k].InnerIps[k1] = *v1
		}
	}

	isFinished := false
	if len(ecses) < int(req.PageSize) {
		isFinished = true
	}

	return &pbecs.ListDetailResp{
		Ecses:      ecses,
		Finished:   isFinished,
		NextToken:  "",
		PageNumber: req.PageNumber + 1,
		PageSize:   req.PageSize,
		RequestId:  *resp.Response.RequestId,
	}, nil
}

func (ecs *TencentCvm) ActionEcs(ctx context.Context, req *pbecs.ActionReq) (resp *pbecs.ActionResp, err error) {
	return nil, nil
}
