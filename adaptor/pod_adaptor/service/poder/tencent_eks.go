package poder

import (
	"context"
	"strconv"
	"sync"

	"github.com/golang/glog"

	"github.com/zhouqunjie-cs/PCM/lan_trans/idl/pbtenant"

	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencenteks "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/zhouqunjie-cs/PCM/common/tenanter"
	"github.com/zhouqunjie-cs/PCM/lan_trans/idl/pbpod"
)

var tencentClientMutex sync.Mutex

type TencentEks struct {
	cli      *tencenteks.Client
	region   tenanter.Region
	tenanter tenanter.Tenanter
}

func (eks TencentEks) GetPodRegion(ctx context.Context, req *pbpod.GetPodRegionReq) (*pbpod.GetPodRegionResp, error) {

	regions := make([]*pbtenant.Region, 0)

	request := tencenteks.NewDescribeEKSContainerInstanceRegionsRequest()
	resp, err := eks.cli.DescribeEKSContainerInstanceRegions(request)
	if err != nil {
		return nil, errors.Wrap(err, "tencent eks describe region error")
	}
	for _, eksRegion := range resp.Response.Regions {

		regionPod := &pbtenant.Region{
			Id:   int32(*eksRegion.RegionId),
			Name: *eksRegion.RegionName,
		}
		regions = append(regions, regionPod)
	}
	return &pbpod.GetPodRegionResp{Regions: regions}, nil

}

func newTencentEksClient(region tenanter.Region, tenant tenanter.Tenanter) (Poder, error) {
	var (
		client *tencenteks.Client
		err    error
	)

	switch t := tenant.(type) {
	case *tenanter.AccessKeyTenant:
		tencentClientMutex.Lock()

		credential := common.NewCredential(
			t.GetId(),
			t.GetSecret(),
		)
		cpf := profile.NewClientProfile()
		client, err = tencenteks.NewClient(credential, region.GetName(), cpf)
		tencentClientMutex.Unlock()
	default:
	}

	if err != nil {
		return nil, errors.Wrap(err, "init tencent eks client error")
	}

	return &TencentEks{
		cli:      client,
		region:   region,
		tenanter: tenant,
	}, nil
}

func (eks TencentEks) CreatePod(ctx context.Context, req *pbpod.CreatePodReq) (*pbpod.CreatePodResp, error) {

	request := tencenteks.NewCreateEKSContainerInstancesRequest()

	eksCiName := req.PodName
	containerName := req.ContainerName
	containerImage := req.ContainerImage
	eksCpu := req.CpuPod
	eksMemory := req.MemoryPod
	securityGroupId := req.SecurityGroupId
	securityGroupIds := make([]*string, 1)
	securityGroupIds[0] = &securityGroupId
	subNetId := req.SubnetId
	vpcId := req.VpcId

	request.EksCiName = &eksCiName
	container := make([]*tencenteks.Container, 1)
	container[0] = new(tencenteks.Container)
	container[0].Name = &containerName
	container[0].Image = &containerImage

	request.Containers = container
	eksCpu64, err := strconv.ParseFloat(eksCpu, 64)
	eksMemory64, err := strconv.ParseFloat(eksMemory, 64)
	request.Cpu = &eksCpu64
	request.Memory = &eksMemory64
	request.SecurityGroupIds = securityGroupIds
	request.SubnetId = &subNetId
	request.VpcId = &vpcId

	resp, err := eks.cli.CreateEKSContainerInstances(request)
	if err != nil {
		return nil, errors.Wrap(err, "Tencent CreatePod error")
	}

	isFinished := false
	if resp.Response.RequestId != nil {
		isFinished = true
	}

	glog.Infof("--------------------K8S Pod Instance created--------------------")

	return &pbpod.CreatePodResp{
		Finished:  isFinished,
		RequestId: "tencent pod create request id:" + *resp.Response.RequestId,
		PodId:     *resp.Response.EksCiIds[0],
		PodName:   req.PodName,
	}, nil

}

func (eks *TencentEks) DeletePod(ctx context.Context, req *pbpod.DeletePodReq) (*pbpod.DeletePodResp, error) {
	request := tencenteks.NewDeleteEKSContainerInstancesRequest()
	request.EksCiIds = make([]*string, 1)
	request.EksCiIds[0] = &req.PodId
	resp, err := eks.cli.DeleteEKSContainerInstances(request)

	isFinished := true
	if err != nil {
		isFinished = false
		return nil, errors.Wrap(err, "Tencent DeletePod error")
	}

	glog.Infof("--------------------K8S Pod Instance deleted--------------------")

	return &pbpod.DeletePodResp{
		Finished:  isFinished,
		RequestId: "tencent pod delete request id:" + *resp.Response.RequestId,
		PodId:     req.PodId,
		PodName:   req.PodName,
	}, nil
}

func (eks *TencentEks) UpdatePod(ctx context.Context, req *pbpod.UpdatePodReq) (*pbpod.UpdatePodResp, error) {
	//创建更新pod请求
	request := tencenteks.NewUpdateEKSContainerInstanceRequest()
	request.EksCiId = &req.PodId
	request.RestartPolicy = &req.RestartPolicy
	request.Name = &req.PodName
	cpu, err := strconv.ParseFloat(req.CpuPod, 64)
	memory, err := strconv.ParseFloat(req.MemoryPod, 64)
	request.Containers = make([]*tencenteks.Container, 1)
	request.Containers[0] = new(tencenteks.Container)
	request.Containers[0].Name = &req.ContainerName
	request.Containers[0].Image = &req.ContainerImage
	request.Containers[0].Cpu = &cpu
	request.Containers[0].Memory = &memory
	resp, err := eks.cli.UpdateEKSContainerInstance(request)
	isFinished := true
	if err != nil {
		isFinished = false
		return nil, errors.Wrap(err, "Tencent UpdatePod error")
	}

	glog.Infof("--------------------K8S Pod Instance deleted--------------------")

	return &pbpod.UpdatePodResp{
		Finished:  isFinished,
		RequestId: "tencent pod update request id:" + *resp.Response.RequestId,
		PodId:     req.PodId,
		PodName:   req.PodName,
	}, nil
}

func (eks TencentEks) ListPodDetail(ctx context.Context, req *pbpod.ListPodDetailReq) (*pbpod.ListPodDetailResp, error) {
	request := tencenteks.NewDescribeEKSContainerInstancesRequest()
	resp, err := eks.cli.DescribeEKSContainerInstances(request)
	if err != nil {
		return nil, errors.Wrap(err, "Tencent ListDetail pod error")
	}
	var ekspods = make([]*pbpod.PodInstance, len(resp.Response.EksCis))
	for k, v := range resp.Response.EksCis {
		ekspods[k] = &pbpod.PodInstance{
			Provider:        pbtenant.CloudProvider_tencent,
			AccountName:     eks.tenanter.AccountName(),
			PodId:           *v.EksCiId,
			PodName:         *v.EksCiName,
			RegionId:        eks.region.GetId(),
			RegionName:      eks.region.GetName(),
			ContainerImage:  *v.Containers[0].Image,
			ContainerName:   *v.Containers[0].Name,
			CpuPod:          strconv.FormatFloat(*v.Cpu, 'f', 6, 64),
			MemoryPod:       strconv.FormatFloat(*v.Memory, 'f', 6, 64),
			SecurityGroupId: *v.SecurityGroupIds[0],
			SubnetId:        *v.SubnetId,
			VpcId:           *v.VpcId,
			Namespace:       "",
			Status:          *v.Status,
		}
	}
	isFinished := false
	if len(ekspods) < int(req.PageSize) {
		isFinished = true
	}

	glog.Infof("--------------------K8S Pod Instance deleted--------------------")

	return &pbpod.ListPodDetailResp{
		Pods:       ekspods,
		Finished:   isFinished,
		PageNumber: req.PageNumber + 1,
		PageSize:   req.PageSize,
		RequestId:  *resp.Response.RequestId,
	}, nil
}
