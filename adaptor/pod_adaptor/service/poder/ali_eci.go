package poder

import (
	"context"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	alieci "github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
	"github.com/bitly/go-simplejson"
	"github.com/golang/glog"
	"strconv"
	"sync"

	"github.com/JCCE-nudt/PCM/common/tenanter"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbpod"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbtenant"
	"github.com/pkg/errors"
)

var aliClientMutex sync.Mutex

type AliEci struct {
	cli      *alieci.Client
	region   tenanter.Region
	tenanter tenanter.Tenanter
}

func (eci *AliEci) GetPodRegion(ctx context.Context, req *pbpod.GetPodRegionReq) (*pbpod.GetPodRegionResp, error) {

	regions := make([]*pbtenant.Region, 0)

	requestRegion := requests.NewCommonRequest()
	requestRegion.Method = "POST"
	requestRegion.Scheme = "https" // https | http
	requestRegion.Domain = "eci.aliyuncs.com"
	requestRegion.Version = "2018-08-08"
	requestRegion.ApiName = "DescribeRegions"
	//这里需要给一个空字符串的RegionId
	requestRegion.QueryParams["RegionId"] = ""

	resp, err := eci.cli.ProcessCommonRequest(requestRegion)

	var respRegion *simplejson.Json
	respRegion, err = simplejson.NewJson([]byte(resp.GetHttpContentString()))
	if err != nil {
		panic("解析失败")
	}

	i := 0
	for i < len(respRegion.Get("Regions").MustArray()) {
		regionsJson := respRegion.Get("Regions").GetIndex(i)
		regionName, _ := regionsJson.Get("RegionId").String()
		regionId, _ := tenanter.GetAliRegionId(regionName)
		regionPod := &pbtenant.Region{
			Id:   regionId,
			Name: regionName,
		}
		regions = append(regions, regionPod)
		i++
	}

	return &pbpod.GetPodRegionResp{Regions: regions}, nil

}

func newAliEciClient(region tenanter.Region, tenant tenanter.Tenanter) (Poder, error) {
	var (
		client *alieci.Client
		err    error
	)

	switch t := tenant.(type) {
	case *tenanter.AccessKeyTenant:
		aliClientMutex.Lock()
		client, err = alieci.NewClientWithAccessKey(region.GetName(), t.GetId(), t.GetSecret())
		aliClientMutex.Unlock()
	default:
	}

	if err != nil {
		return nil, errors.Wrap(err, "init ali ecs client error")
	}

	return &AliEci{
		cli:      client,
		region:   region,
		tenanter: tenant,
	}, nil
}

func (eci *AliEci) CreatePod(ctx context.Context, req *pbpod.CreatePodReq) (*pbpod.CreatePodResp, error) {
	request := alieci.CreateCreateContainerGroupRequest()
	request.RegionId = eci.region.GetName()
	request.ContainerGroupName = req.PodName
	requestContainer := make([]alieci.CreateContainerGroupContainer, 1)
	requestContainer[0].Image = req.ContainerImage
	requestContainer[0].Name = req.ContainerName
	requestContainer[0].Cpu = requests.Float(req.CpuPod)
	requestContainer[0].Memory = requests.Float(req.MemoryPod)
	request.Container = &requestContainer

	resp, err := eci.cli.CreateContainerGroup(request)
	if err != nil {
		return nil, errors.Wrap(err, "Aliyun CreatePod error")
	}

	isFinished := false
	if len(resp.ContainerGroupId) > 0 {
		isFinished = true
	}
	glog.Infof("--------------------Aliyun ECI Instance created--------------------")

	return &pbpod.CreatePodResp{
		Finished:  isFinished,
		RequestId: "Create Ali pod request ID:" + resp.RequestId,
		PodId:     resp.ContainerGroupId,
		PodName:   req.PodName,
	}, nil
}

func (eci *AliEci) DeletePod(ctx context.Context, req *pbpod.DeletePodReq) (*pbpod.DeletePodResp, error) {

	deleteContainerGroupRequest := alieci.CreateDeleteContainerGroupRequest()
	deleteContainerGroupRequest.RegionId = eci.region.GetName()
	deleteContainerGroupRequest.ContainerGroupId = req.PodId

	resp, err := eci.cli.DeleteContainerGroup(deleteContainerGroupRequest)
	isFinished := true
	if err != nil {
		isFinished = false
		return nil, errors.Wrap(err, "Aliyun DeletePod error")
	}

	glog.Infof("--------------------Aliyun ECI Instance deleted--------------------")

	return &pbpod.DeletePodResp{
		Finished:  isFinished,
		RequestId: "Delete Ali pod request ID:" + resp.RequestId,
		PodId:     req.PodId,
		PodName:   req.PodName,
	}, nil
}

func (eci *AliEci) UpdatePod(ctx context.Context, req *pbpod.UpdatePodReq) (*pbpod.UpdatePodResp, error) {

	updateContainerGroupRequest := alieci.CreateUpdateContainerGroupRequest()
	updateContainerGroupRequest.RegionId = eci.region.GetName()
	updateContainerGroupRequest.ContainerGroupId = req.PodId

	updateContainerRequestContainer := make([]alieci.UpdateContainerGroupContainer, 1)
	updateContainerRequestContainer[0].Image = req.ContainerImage
	updateContainerRequestContainer[0].Name = req.ContainerName
	updateContainerRequestContainer[0].Cpu = requests.Float(req.CpuPod)
	updateContainerRequestContainer[0].Memory = requests.Float(req.MemoryPod)
	updateContainerGroupRequest.Container = &updateContainerRequestContainer
	updateContainerGroupRequest.RestartPolicy = req.RestartPolicy

	resp, err := eci.cli.UpdateContainerGroup(updateContainerGroupRequest)
	isFinished := true
	if err != nil {
		isFinished = false
		return nil, errors.Wrap(err, "Aliyun UpdatePod error")
	}

	glog.Infof("--------------------Aliyun ECI Instance updated--------------------")

	return &pbpod.UpdatePodResp{
		Finished:  isFinished,
		RequestId: "Update Ali pod request ID:" + resp.RequestId,
		PodId:     req.PodId,
		PodName:   req.PodName,
	}, nil
}

func (eci *AliEci) ListPodDetail(ctx context.Context, req *pbpod.ListPodDetailReq) (*pbpod.ListPodDetailResp, error) {
	request := alieci.CreateDescribeContainerGroupsRequest()
	request.NextToken = req.NextToken
	resp, err := eci.cli.DescribeContainerGroups(request)

	if err != nil {
		return nil, errors.Wrap(err, "Aliyun ListDetail error")
	}

	var ecies = make([]*pbpod.PodInstance, len(resp.ContainerGroups))
	for k, v := range resp.ContainerGroups {
		ecies[k] = &pbpod.PodInstance{
			Provider:        pbtenant.CloudProvider_ali,
			AccountName:     eci.tenanter.AccountName(),
			PodId:           v.ContainerGroupId,
			PodName:         v.ContainerGroupName,
			RegionId:        eci.region.GetId(),
			RegionName:      v.RegionId,
			ContainerImage:  v.Containers[0].Image,
			ContainerName:   v.Containers[0].Name,
			CpuPod:          strconv.FormatFloat(float64(v.Cpu), 'f', 6, 64),
			MemoryPod:       strconv.FormatFloat(float64(v.Memory), 'f', 6, 64),
			SecurityGroupId: v.SecurityGroupId,
			SubnetId:        v.InternetIp,
			VpcId:           v.VpcId,
			Namespace:       "",
			Status:          v.Status,
		}

	}

	isFinished := false
	if len(ecies) < int(req.PageSize) {
		isFinished = true
	}

	return &pbpod.ListPodDetailResp{
		Pods:       ecies,
		Finished:   isFinished,
		PageNumber: req.PageNumber + 1,
		PageSize:   req.PageSize,
		NextToken:  resp.NextToken,
		RequestId:  resp.RequestId,
	}, nil
}
