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
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	hwecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	hwregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/region"
	hwiam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	iammodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	iamregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
	"github.com/pkg/errors"
)

type HuaweiEcs struct {
	cli      *hwecs.EcsClient
	region   tenanter.Region
	tenanter tenanter.Tenanter
}

func newHuaweiEcsClient(region tenanter.Region, tenant tenanter.Tenanter) (Ecser, error) {
	var (
		client *hwecs.EcsClient
		err    error
	)

	switch t := tenant.(type) {
	case *tenanter.AccessKeyTenant:
		auth := basic.NewCredentialsBuilder().WithAk(t.GetId()).WithSk(t.GetSecret()).Build()
		rName := region.GetName()
		cli := hwiam.IamClientBuilder().WithRegion(iamregion.ValueOf(rName)).WithCredential(auth).Build()
		c := hwiam.NewIamClient(cli)
		request := new(iammodel.KeystoneListProjectsRequest)
		request.Name = &rName
		r, err := c.KeystoneListProjects(request)
		if err != nil || len(*r.Projects) == 0 {
			return nil, errors.Wrapf(err, "Huawei KeystoneListProjects regionName %s", rName)
		}
		projectId := (*r.Projects)[0].Id

		auth = basic.NewCredentialsBuilder().WithAk(t.GetId()).WithSk(t.GetSecret()).WithProjectId(projectId).Build()
		hcClient := hwecs.EcsClientBuilder().WithRegion(hwregion.ValueOf(rName)).WithCredential(auth).Build()
		client = hwecs.NewEcsClient(hcClient)
	default:
	}

	if err != nil {
		return nil, errors.Wrap(err, "init huawei ecs client error")
	}
	return &HuaweiEcs{
		cli:      client,
		region:   region,
		tenanter: tenant,
	}, nil
}

func (ecs *HuaweiEcs) CreateEcs(ctx context.Context, req *pbecs.CreateEcsReq) (*pbecs.CreateEcsResp, error) {
	subnetIds := string_.Split(&req.SubnetId, tea.String(","), tea.Int(-1))
	Nics := make([]model.PrePaidServerNic, 0)
	for _, nic := range subnetIds {
		Nics = append(Nics, model.PrePaidServerNic{
			SubnetId: *nic,
		})
	}
	Volumetype := &model.PrePaidServerRootVolume{}
	switch req.SystemDisk.Category {
	case "SATA":
		Volumetype.Volumetype = model.GetPrePaidServerRootVolumeVolumetypeEnum().SATA
	case "SAS":
		Volumetype.Volumetype = model.GetPrePaidServerRootVolumeVolumetypeEnum().SAS
	case "SSD":
		Volumetype.Volumetype = model.GetPrePaidServerRootVolumeVolumetypeEnum().SSD
	case "GPSSD":
		Volumetype.Volumetype = model.GetPrePaidServerRootVolumeVolumetypeEnum().GPSSD
	case "co-p1":
		Volumetype.Volumetype = model.GetPrePaidServerRootVolumeVolumetypeEnum().CO_P1
	case "uh-l1":
		Volumetype.Volumetype = model.GetPrePaidServerRootVolumeVolumetypeEnum().UH_L1
	case "ESSD":
		Volumetype.Volumetype = model.GetPrePaidServerRootVolumeVolumetypeEnum().ESSD
	}

	PrePaidServerExtendParam := &model.PrePaidServerExtendParam{}
	prePaid := model.GetPrePaidServerExtendParamChargingModeEnum().PRE_PAID
	if req.InstanceChargeType == "PrePaid" {
		PrePaidServerExtendParam.ChargingMode = &prePaid
	}
	request := &model.CreateServersRequest{
		Body: &model.CreateServersRequestBody{
			DryRun: util.EqualString(&req.DryRun, tea.String("true")),
			Server: &model.PrePaidServer{
				ImageRef:    req.GetImageId(),
				FlavorRef:   req.InstanceType,
				Name:        req.InstanceName,
				Vpcid:       req.VpcId,
				Nics:        Nics,
				RootVolume:  Volumetype,
				Count:       &req.Amount,
				Extendparam: PrePaidServerExtendParam,
			},
		}}
	resp, err := ecs.cli.CreateServers(request)
	if err != nil {
		return nil, errors.Wrap(err, "Huawei create ecs error")
	}
	glog.Infof("--------------------华为ECS实例创建成功--------------------")
	glog.Infof(resp.String())
	isFinished := false
	if len(*resp.ServerIds) > 0 {
		isFinished = true
	}
	//订单ID。该参数只有创建包年包月ECS实例（请求参数InstanceChargeType=PrePaid）时有返回值。
	OrderId := ""
	if req.InstanceChargeType == "PrePaid" {
		OrderId = *resp.OrderId
	}
	InstanceIds := make([]string, 0)
	for _, v := range *resp.ServerIds {
		InstanceIds = append(InstanceIds, v)
	}
	return &pbecs.CreateEcsResp{
		OrderId:        OrderId,
		RequestId:      "Huawei ECS RequestId: " + *resp.JobId,
		InstanceIdSets: InstanceIds,
		Finished:       isFinished,
	}, nil
}

func (ecs *HuaweiEcs) DeleteEcs(ctx context.Context, req *pbecs.DeleteEcsReq) (*pbecs.DeleteEcsResp, error) {
	if req.GetInstanceIds() == "" {
		return nil, errors.New("InstanceId is empty")
	}
	deleteReq := &model.DeleteServersRequest{}
	InstanceIds := string_.Split(&req.InstanceIds, tea.String(","), tea.Int(-1))
	Servers := make([]model.ServerId, 0)
	for _, v := range InstanceIds {
		Servers = append(Servers, model.ServerId{
			Id: *v,
		})
	}
	deleteReq.Body = &model.DeleteServersRequestBody{
		DeletePublicip: util.EqualString(&req.DeletePublicip, tea.String("true")),
		DeleteVolume:   util.EqualString(&req.DeleteVolume, tea.String("true")),
		Servers:        Servers,
	}
	resp, err := ecs.cli.DeleteServers(deleteReq)
	if err != nil {
		return nil, errors.Wrap(err, "Huawei Delete ECS error")
	}
	glog.Infof("--------------------华为ECS实例删除成功--------------------")
	glog.Infof(resp.String())
	return &pbecs.DeleteEcsResp{
		RequestId: *resp.JobId,
	}, nil
}

func (ecs *HuaweiEcs) UpdateEcs(ctx context.Context, req *pbecs.UpdateEcsReq) (*pbecs.UpdateEcsResp, error) {
	HostName := req.GetHostName()
	InstanceName := req.GetInstanceName()
	Description := req.GetDescription()
	UpdateReq := &model.UpdateServerRequest{}
	if req.GetInstanceIds() == "" {
		return nil, errors.New("InstanceId is empty")
	}
	Server := &model.UpdateServerOption{}
	UpdateReq.ServerId = req.GetInstanceIds()
	if HostName != "" {
		Server.Hostname = &HostName
	}
	if InstanceName != "" {
		Server.Name = &InstanceName
	}
	if Description != "" {
		Server.Description = &Description
	}
	UpdateReq.Body = &model.UpdateServerRequestBody{
		Server: Server,
	}
	resp, err := ecs.cli.UpdateServer(UpdateReq)
	if err != nil {
		return nil, errors.Wrap(err, "Huawei Update ECS error")
	}
	glog.Infof("--------------------华为ECS实例修改成功--------------------")
	glog.Infof(resp.String())
	return &pbecs.UpdateEcsResp{
		RequestId: resp.Server.Id,
	}, nil
}

func (ecs *HuaweiEcs) ListDetail(ctx context.Context, req *pbecs.ListDetailReq) (*pbecs.ListDetailResp, error) {
	request := new(model.ListServersDetailsRequest)
	offset := (req.PageNumber - 1) * req.PageSize
	request.Offset = &offset
	limit := req.PageSize
	request.Limit = &limit

	resp, err := ecs.cli.ListServersDetails(request)
	if err != nil {
		return nil, errors.Wrap(err, "Huawei ListDetail error")
	}

	servers := *resp.Servers
	var ecses = make([]*pbecs.EcsInstance, len(servers))
	for k, v := range servers {
		vCpu, err := strconv.ParseInt(v.Flavor.Vcpus, 10, 32)
		vMemory, err := strconv.ParseInt(v.Flavor.Vcpus, 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "Huawei ListDetail error")
		}
		PublicIps := make([]string, 0)
		InnerIps := make([]string, 0)
		for s := range v.Addresses {
			for _, a := range v.Addresses[s] {
				// 判断是内网ip还是公网ip
				if *a.OSEXTIPStype == model.GetServerAddressOSEXTIPStypeEnum().FIXED {
					InnerIps = append(InnerIps, a.Addr)
				} else {
					PublicIps = append(PublicIps, a.Addr)
				}
			}
		}
		ecses[k] = &pbecs.EcsInstance{
			Provider:           pbtenant.CloudProvider_huawei,
			AccountName:        ecs.tenanter.AccountName(),
			InstanceId:         v.Id,
			InstanceName:       v.Name,
			RegionName:         ecs.region.GetName(),
			InstanceType:       v.Flavor.Name,
			PublicIps:          PublicIps,
			InnerIps:           InnerIps,
			Cpu:                strconv.FormatInt(vCpu, 10),
			Memory:             strconv.FormatInt(vMemory, 10),
			Description:        *v.Description,
			Status:             v.Status,
			CreationTime:       v.Created,
			ExpireTime:         v.OSSRVUSGterminatedAt,
			InstanceChargeType: v.Metadata["charging_mode"],
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
		RequestId:  "",
	}, nil
}

func (ecs *HuaweiEcs) ActionEcs(ctx context.Context, req *pbecs.ActionReq) (resp *pbecs.ActionResp, err error) {
	return nil, nil
}
