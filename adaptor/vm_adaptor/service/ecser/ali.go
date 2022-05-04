package ecser

import (
	"context"
	"sync"

	"github.com/JCCE-nudt/PCM/common/tenanter"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbecs"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbtenant"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	string_ "github.com/alibabacloud-go/darabonba-string/client"
	aliecs "github.com/alibabacloud-go/ecs-20140526/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

var aliClientMutex sync.Mutex

type AliEcs struct {
	cli      *aliecs.Client
	region   tenanter.Region
	tenanter tenanter.Tenanter
}

func newAliEcsClient(region tenanter.Region, tenant tenanter.Tenanter) (Ecser, error) {
	var (
		client *aliecs.Client
		err    error
	)
	switch t := tenant.(type) {
	case *tenanter.AccessKeyTenant:
		// 阿里云的sdk有一个 map 的并发问题，go test 加上-race 能检测出来，所以这里加一个锁
		aliClientMutex.Lock()
		config := &openapi.Config{}
		AccessKeyId := t.GetId()
		AccessKeySecret := t.GetSecret()
		RegionId := region.GetName()
		config.AccessKeyId = &AccessKeyId
		config.AccessKeySecret = &AccessKeySecret
		config.RegionId = &RegionId
		client, err = aliecs.NewClient(config)
		aliClientMutex.Unlock()
	default:
		return nil, errors.New("unsupported tenant type")
	}

	if err != nil {
		return nil, errors.Wrap(err, "init ali ecs client error")
	}

	return &AliEcs{
		cli:      client,
		region:   region,
		tenanter: tenant,
	}, nil
}
func (ecs *AliEcs) CreateEcs(ctx context.Context, req *pbecs.CreateEcsReq) (*pbecs.CreateEcsResp, error) {
	RegionId := ecs.region.GetName()
	ImageId := req.GetImageId()
	InstanceType := req.GetInstanceType()
	SecurityGroupId := req.GetSecurityGroupId()
	InstanceName := req.GetInstanceName()
	Description := req.GetDescription()
	ZoneId := req.GetZoneId()
	VSwitchId := req.GetVSwitchId()
	Amount := req.GetAmount()
	DryRun := req.GetDryRun()
	Category := req.GetCategory()
	InstanceChargeType := req.GetInstanceChargeType()
	request := &aliecs.RunInstancesRequest{
		RegionId:        &RegionId,
		InstanceType:    &InstanceType,
		ImageId:         &ImageId,
		SecurityGroupId: &SecurityGroupId,
		InstanceName:    &InstanceName,
		Description:     &Description,
		ZoneId:          &ZoneId,
		VSwitchId:       &VSwitchId,
		Amount:          &Amount,
		DryRun:          util.EqualString(&DryRun, tea.String("true")),
		SystemDisk: &aliecs.RunInstancesRequestSystemDisk{
			Category: &Category,
		},
		InstanceChargeType: &InstanceChargeType,
	}
	// 创建并运行实例
	resp, err := ecs.cli.RunInstances(request)
	if err != nil {
		return nil, errors.Wrap(err, "Ali Create ECS error")
	}
	isFinished := false

	if len(resp.Body.InstanceIdSets.InstanceIdSet) > 0 {
		isFinished = true
	}
	glog.Infof("--------------------阿里ECS实例创建成功--------------------")
	glog.Infof(*util.ToJSONString(util.ToMap(resp)))
	requestId := *resp.Body.RequestId
	//订单ID。该参数只有创建包年包月ECS实例（请求参数InstanceChargeType=PrePaid）时有返回值。
	OrderId := ""
	if req.InstanceChargeType == "PrePaid" {
		OrderId = *resp.Body.OrderId
	}
	TradePrice := float32(0)
	if resp.Body.TradePrice != nil {
		TradePrice = *resp.Body.TradePrice
	}
	InstanceIds := make([]string, 0)
	for _, v := range resp.Body.InstanceIdSets.InstanceIdSet {
		InstanceIds = append(InstanceIds, *v)
	}
	return &pbecs.CreateEcsResp{
		OrderId:        OrderId,
		TradePrice:     TradePrice,
		RequestId:      "Ali ECS RequestId: " + requestId,
		InstanceIdSets: InstanceIds,
		Finished:       isFinished,
	}, nil
}

func (ecs *AliEcs) DeleteEcs(ctx context.Context, req *pbecs.DeleteEcsReq) (*pbecs.DeleteEcsResp, error) {
	RegionId := ecs.region.GetName()
	InstanceIds := req.GetInstanceIds()
	DryRun := req.GetDryRun()
	Force := req.GetForce()
	TerminateSubscription := req.GetTerminateSubscription()
	deleteReq := &aliecs.DeleteInstancesRequest{
		RegionId:              &RegionId,
		InstanceId:            string_.Split(&InstanceIds, tea.String(","), tea.Int(-1)),
		Force:                 util.EqualString(&Force, tea.String("true")),
		DryRun:                util.EqualString(&DryRun, tea.String("true")),
		TerminateSubscription: util.EqualString(&TerminateSubscription, tea.String("true")),
	}
	resp, err := ecs.cli.DeleteInstances(deleteReq)
	if err != nil {
		return nil, errors.Wrap(err, "Ali Delete ECS error")
	}
	glog.Infof("--------------------阿里ECS实例释放成功--------------------")
	glog.Infof(*util.ToJSONString(util.ToMap(resp)))
	return &pbecs.DeleteEcsResp{
		RequestId:   *resp.Body.RequestId,
		AccountName: req.AccountName,
		RegionId:    req.RegionId,
	}, nil
}

func (ecs *AliEcs) UpdateEcs(ctx context.Context, req *pbecs.UpdateEcsReq) (*pbecs.UpdateEcsResp, error) {
	Password := req.GetPassword()
	HostName := req.GetHostName()
	InstanceName := req.GetInstanceName()
	Description := req.GetDescription()
	InstanceId := req.GetInstanceIds()
	UpdateReq := &aliecs.ModifyInstanceAttributeRequest{}
	if req.GetInstanceIds() == "" {
		return nil, errors.New("InstanceId is empty")
	}
	UpdateReq.InstanceId = &InstanceId
	if Password != "" {
		UpdateReq.Password = &Password
	}
	if HostName != "" {
		UpdateReq.HostName = &HostName
	}
	if InstanceName != "" {
		UpdateReq.InstanceName = &InstanceName
	}
	if Description != "" {
		UpdateReq.Description = &Description
	}
	resp, err := ecs.cli.ModifyInstanceAttribute(UpdateReq)
	if err != nil {
		return nil, errors.Wrap(err, "Ali Update ECS error")
	}
	glog.Infof("--------------------阿里ECS实例修改成功--------------------")
	glog.Infof(*util.ToJSONString(util.ToMap(resp)))
	return &pbecs.UpdateEcsResp{
		RequestId: *resp.Body.RequestId,
	}, nil
}

func (ecs *AliEcs) ListDetail(ctx context.Context, req *pbecs.ListDetailReq) (*pbecs.ListDetailResp, error) {
	request := &aliecs.DescribeInstancesRequest{}
	request.PageNumber = &req.PageNumber
	request.PageSize = &req.PageSize
	request.NextToken = &req.NextToken
	request.RegionId = ecs.cli.RegionId
	resp, err := ecs.cli.DescribeInstances(request)
	if err != nil {
		return nil, errors.Wrap(err, "Ali ListDetail error")
	}
	var ecses = make([]*pbecs.EcsInstance, len(resp.Body.Instances.Instance))
	for k, v := range resp.Body.Instances.Instance {
		publicIps := make([]string, 0)
		for _, vv := range v.PublicIpAddress.IpAddress {
			publicIps = append(publicIps, *vv)
		}
		InnerIps := make([]string, 0)
		for _, vv := range v.VpcAttributes.PrivateIpAddress.IpAddress {
			InnerIps = append(InnerIps, *vv)
		}
		ecses[k] = &pbecs.EcsInstance{
			Provider:           pbtenant.CloudProvider_ali,
			AccountName:        ecs.tenanter.AccountName(),
			InstanceId:         *v.InstanceId,
			InstanceName:       *v.InstanceName,
			RegionName:         ecs.region.GetName(),
			PublicIps:          publicIps,
			InstanceType:       *v.InstanceType,
			Cpu:                string(*v.Cpu),
			Memory:             string(*v.Memory),
			Description:        *v.Description,
			Status:             *v.Status,
			CreationTime:       *v.CreationTime,
			ExpireTime:         *v.ExpiredTime,
			InnerIps:           InnerIps,
			VpcId:              *v.VpcAttributes.VpcId,
			ResourceGroupId:    *v.ResourceGroupId,
			InstanceChargeType: *v.InstanceChargeType,
		}
	}
	isFinished := false
	if len(ecses) < int(req.PageSize) {
		isFinished = true
	}

	return &pbecs.ListDetailResp{
		Ecses:      ecses,
		Finished:   isFinished,
		PageNumber: req.PageNumber + 1,
		PageSize:   req.PageSize,
		NextToken:  *resp.Body.NextToken,
		RequestId:  *resp.Body.RequestId,
	}, nil
}

func (ecs *AliEcs) ActionEcs(ctx context.Context, req *pbecs.ActionReq) (resp *pbecs.ActionResp, err error) {
	return nil, nil
}
