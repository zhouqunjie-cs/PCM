package ali

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/zhouqunjie-cs/PCM/adaptor/pod_adaptor/server/pod"
	"github.com/zhouqunjie-cs/PCM/common/tenanter"
	"github.com/zhouqunjie-cs/PCM/lan_trans/idl/pbpod"
	"github.com/zhouqunjie-cs/PCM/lan_trans/idl/pbtenant"
)

// CreateContainerGroup invokes the eci.CreateContainerGroup API synchronously
// api document: https://help.aliyun.com/api/eci/createcontainergroup.html
func CreateContainerGroup(request *CreateContainerGroupRequest) (response *CreateContainerGroupResponse, err error) {

	provider := pbtenant.CloudProvider(request.ProviderId)
	tenanters, err := tenanter.GetTenanters(provider)
	regionId, err := tenanter.GetAliRegionId(request.RegionId)
	container := *request.Container
	containerImage := container[0].Image
	containerName := container[0].Name
	containerPod := container[0].Cpu
	memoryPod := container[0].Memory

	requestPCM := &pbpod.CreatePodReq{
		Provider:        provider,
		AccountName:     tenanters[0].AccountName(),
		PodName:         request.ContainerGroupName,
		RegionId:        regionId,
		ContainerImage:  containerImage,
		ContainerName:   containerName,
		CpuPod:          string(containerPod),
		MemoryPod:       string(memoryPod),
		SecurityGroupId: "sg-6qlun7hd",
		SubnetId:        "subnet-mnwfg2fk",
		VpcId:           "vpc-rkwt40g5",
		Namespace:       "pcm",
	}

	resp, err := pod.CreatePod(nil, requestPCM)

	response = &CreateContainerGroupResponse{
		BaseResponse:     nil,
		RequestId:        resp.RequestId,
		ContainerGroupId: resp.PodId,
	}

	return response, nil
}

// CreateContainerGroupRequest is the request struct for api CreateContainerGroup
type CreateContainerGroupRequest struct {
	*requests.RpcRequest
	/*********PCM param************/
	ProviderId  int32  `position:"Query" name:"ProviderId"`
	AccountName string `position:"Query" name:"AccountName"`
	Namespace   string `position:"Query" name:"Namespace"`
	/*********PCM param************/
	OwnerId                       requests.Integer                               `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount          string                                         `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId               requests.Integer                               `position:"Query" name:"ResourceOwnerId"`
	OwnerAccount                  string                                         `position:"Query" name:"OwnerAccount"`
	RegionId                      string                                         `position:"Query" name:"RegionId"`
	ZoneId                        string                                         `position:"Query" name:"ZoneId"`
	SecurityGroupId               string                                         `position:"Query" name:"SecurityGroupId"`
	VSwitchId                     string                                         `position:"Query" name:"VSwitchId"`
	ContainerGroupName            string                                         `position:"Query" name:"ContainerGroupName"`
	RestartPolicy                 string                                         `position:"Query" name:"RestartPolicy"`
	Tag                           *[]CreateContainerGroupTag                     `position:"Query" name:"Tag" type:"Repeated"`
	ImageRegistryCredential       *[]CreateContainerGroupImageRegistryCredential `position:"Query" name:"ImageRegistryCredential" type:"Repeated"`
	Container                     *[]CreateContainerGroupContainer               `position:"Query" name:"Container" type:"Repeated"`
	Volume                        *[]CreateContainerGroupVolume                  `position:"Query" name:"Volume" type:"Repeated"`
	EipInstanceId                 string                                         `position:"Query" name:"EipInstanceId"`
	InitContainer                 *[]CreateContainerGroupInitContainer           `position:"Query" name:"InitContainer" type:"Repeated"`
	Cpu                           requests.Float                                 `position:"Query" name:"Cpu"`
	Memory                        requests.Float                                 `position:"Query" name:"Memory"`
	ResourceGroupId               string                                         `position:"Query" name:"ResourceGroupId"`
	DnsPolicy                     string                                         `position:"Query" name:"DnsPolicy"`
	ClientToken                   string                                         `position:"Query" name:"ClientToken"`
	HostAliase                    *[]CreateContainerGroupHostAliase              `position:"Query" name:"HostAliase" type:"Repeated"`
	Arn                           *[]CreateContainerGroupArn                     `position:"Query" name:"Arn" type:"Repeated"`
	InstanceType                  string                                         `position:"Query" name:"InstanceType"`
	SlsEnable                     requests.Boolean                               `position:"Query" name:"SlsEnable"`
	ImageSnapshotId               string                                         `position:"Query" name:"ImageSnapshotId"`
	RamRoleName                   string                                         `position:"Query" name:"RamRoleName"`
	NtpServer                     []string                                       `position:"Query" name:"NtpServer" type:"Repeated"`
	TerminationGracePeriodSeconds requests.Integer                               `position:"Query" name:"TerminationGracePeriodSeconds"`
	AutoMatchImageCache           requests.Boolean                               `position:"Query" name:"AutoMatchImageCache"`
	VkClientVersion               string                                         `position:"Query" name:"VkClientVersion"`
	Ipv6AddressCount              requests.Integer                               `position:"Query" name:"Ipv6AddressCount"`
	ActiveDeadlineSeconds         requests.Integer                               `position:"Query" name:"ActiveDeadlineSeconds"`
	SpotStrategy                  string                                         `position:"Query" name:"SpotStrategy"`
	SpotPriceLimit                requests.Float                                 `position:"Query" name:"SpotPriceLimit"`
	VSwitchStrategy               string                                         `position:"Query" name:"VSwitchStrategy"`
	DnsConfig                     CreateContainerGroupDnsConfig                  `position:"Query" name:"DnsConfig" type:"Struct"`
	SecurityContext               CreateContainerGroupSecurityContext            `position:"Query" name:"SecurityContext" type:"Struct"`
}

type CreateContainerGroupTag struct {
	Key   string `name:"Key"`
	Value string `name:"Value"`
}

type CreateContainerGroupImageRegistryCredential struct {
	Server   string `name:"Server"`
	UserName string `name:"UserName"`
	Password string `name:"Password"`
}

type CreateContainerGroupContainer struct {
	Image                                      string                                                            `name:"Image"`
	Name                                       string                                                            `name:"Name"`
	Cpu                                        requests.Float                                                    `name:"Cpu"`
	Memory                                     requests.Float                                                    `name:"Memory"`
	WorkingDir                                 string                                                            `name:"WorkingDir"`
	ImagePullPolicy                            string                                                            `name:"ImagePullPolicy"`
	Command                                    []string                                                          `name:"Command" type:"Repeated"`
	Arg                                        []string                                                          `name:"Arg" type:"Repeated"`
	VolumeMount                                *[]CreateContainerGroupVolumeMount                                `name:"VolumeMount" type:"Repeated"`
	Port                                       *[]CreateContainerGroupPort                                       `name:"Port" type:"Repeated"`
	EnvironmentVar                             *[]CreateContainerGroupEnvironmentVar                             `name:"EnvironmentVar" type:"Repeated"`
	Stdin                                      requests.Boolean                                                  `name:"Stdin"`
	StdinOnce                                  requests.Boolean                                                  `name:"StdinOnce"`
	Tty                                        requests.Boolean                                                  `name:"Tty"`
	Gpu                                        requests.Integer                                                  `name:"Gpu"`
	LifecyclePostStartHandlerHttpGetHost       string                                                            `name:"LifecyclePostStartHandlerHttpGetHost"`
	LifecyclePostStartHandlerHttpGetPort       requests.Integer                                                  `name:"LifecyclePostStartHandlerHttpGetPort"`
	LifecyclePostStartHandlerHttpGetPath       string                                                            `name:"LifecyclePostStartHandlerHttpGetPath"`
	LifecyclePostStartHandlerHttpGetScheme     string                                                            `name:"LifecyclePostStartHandlerHttpGetScheme"`
	LifecyclePostStartHandlerHttpGetHttpHeader *[]CreateContainerGroupLifecyclePostStartHandlerHttpGetHttpHeader `name:"LifecyclePostStartHandlerHttpGetHttpHeader" type:"Repeated"`
	LifecyclePostStartHandlerExec              []string                                                          `name:"LifecyclePostStartHandlerExec" type:"Repeated"`
	LifecyclePostStartHandlerTcpSocketHost     string                                                            `name:"LifecyclePostStartHandlerTcpSocketHost"`
	LifecyclePostStartHandlerTcpSocketPort     requests.Integer                                                  `name:"LifecyclePostStartHandlerTcpSocketPort"`
	LifecyclePreStopHandlerHttpGetHost         string                                                            `name:"LifecyclePreStopHandlerHttpGetHost"`
	LifecyclePreStopHandlerHttpGetPort         requests.Integer                                                  `name:"LifecyclePreStopHandlerHttpGetPort"`
	LifecyclePreStopHandlerHttpGetPath         string                                                            `name:"LifecyclePreStopHandlerHttpGetPath"`
	LifecyclePreStopHandlerHttpGetScheme       string                                                            `name:"LifecyclePreStopHandlerHttpGetScheme"`
	LifecyclePreStopHandlerHttpGetHttpHeader   *[]CreateContainerGroupLifecyclePreStopHandlerHttpGetHttpHeader   `name:"LifecyclePreStopHandlerHttpGetHttpHeader" type:"Repeated"`
	LifecyclePreStopHandlerExec                []string                                                          `name:"LifecyclePreStopHandlerExec" type:"Repeated"`
	LifecyclePreStopHandlerTcpSocketHost       string                                                            `name:"LifecyclePreStopHandlerTcpSocketHost"`
	LifecyclePreStopHandlerTcpSocketPort       requests.Integer                                                  `name:"LifecyclePreStopHandlerTcpSocketPort"`
	ReadinessProbe                             CreateContainerGroupReadinessProbe                                `name:"ReadinessProbe" type:"Struct"`
	LivenessProbe                              CreateContainerGroupLivenessProbe                                 `name:"LivenessProbe" type:"Struct"`
	SecurityContext                            CreateContainerGroupSecurityContext                               `name:"SecurityContext" type:"Struct"`
}

type CreateContainerGroupVolume struct {
	Name             string                               `name:"Name"`
	Type             string                               `name:"Type"`
	NFSVolume        CreateContainerGroupNFSVolume        `name:"NFSVolume" type:"Struct"`
	ConfigFileVolume CreateContainerGroupConfigFileVolume `name:"ConfigFileVolume" type:"Struct"`
	EmptyDirVolume   CreateContainerGroupEmptyDirVolume   `name:"EmptyDirVolume" type:"Struct"`
	DiskVolume       CreateContainerGroupDiskVolume       `name:"DiskVolume" type:"Struct"`
	FlexVolume       CreateContainerGroupFlexVolume       `name:"FlexVolume" type:"Struct"`
	HostPathVolume   CreateContainerGroupHostPathVolume   `name:"HostPathVolume" type:"Struct"`
}

type CreateContainerGroupInitContainer struct {
	Name            string                                `name:"Name"`
	Image           string                                `name:"Image"`
	Cpu             requests.Float                        `name:"Cpu"`
	Memory          requests.Float                        `name:"Memory"`
	WorkingDir      string                                `name:"WorkingDir"`
	ImagePullPolicy string                                `name:"ImagePullPolicy"`
	Command         []string                              `name:"Command" type:"Repeated"`
	Arg             []string                              `name:"Arg" type:"Repeated"`
	VolumeMount     *[]CreateContainerGroupVolumeMount    `name:"VolumeMount" type:"Repeated"`
	Port            *[]CreateContainerGroupPort           `name:"Port" type:"Repeated"`
	EnvironmentVar  *[]CreateContainerGroupEnvironmentVar `name:"EnvironmentVar" type:"Repeated"`
	Gpu             requests.Integer                      `name:"Gpu"`
	SecurityContext CreateContainerGroupSecurityContext   `name:"SecurityContext" type:"Struct"`
}

type CreateContainerGroupHostAliase struct {
	Ip       string   `name:"Ip"`
	Hostname []string `name:"Hostname" type:"Repeated"`
}

type CreateContainerGroupArn struct {
	RoleArn       string `name:"RoleArn"`
	RoleType      string `name:"RoleType"`
	AssumeRoleFor string `name:"AssumeRoleFor"`
}

type CreateContainerGroupDnsConfig struct {
	NameServer []string                      `name:"NameServer"`
	Search     []string                      `name:"Search"`
	Option     *[]CreateContainerGroupOption `name:"Option"`
}

type CreateContainerGroupSecurityContext struct {
	Sysctl *[]CreateContainerGroupSysctl `name:"Sysctl"`
}

type CreateContainerGroupVolumeMount struct {
	MountPath string           `name:"MountPath"`
	ReadOnly  requests.Boolean `name:"ReadOnly"`
	Name      string           `name:"Name"`
	SubPath   string           `name:"SubPath"`
}

type CreateContainerGroupPort struct {
	Protocol string           `name:"Protocol"`
	Port     requests.Integer `name:"Port"`
}

type CreateContainerGroupEnvironmentVar struct {
	Key      string                       `name:"Key"`
	Value    string                       `name:"Value"`
	FieldRef CreateContainerGroupFieldRef `name:"FieldRef" type:"Struct"`
}

type CreateContainerGroupFieldRef struct {
	FieldPath string `name:"FieldPath"`
}

type CreateContainerGroupLifecyclePostStartHandlerHttpGetHttpHeader struct {
	Name  string `name:"Name"`
	Value string `name:"Value"`
}

type CreateContainerGroupLifecyclePreStopHandlerHttpGetHttpHeader struct {
	Name  string `name:"Name"`
	Value string `name:"Value"`
}

type CreateContainerGroupReadinessProbe struct {
	InitialDelaySeconds requests.Integer              `name:"InitialDelaySeconds"`
	PeriodSeconds       requests.Integer              `name:"PeriodSeconds"`
	SuccessThreshold    requests.Integer              `name:"SuccessThreshold"`
	FailureThreshold    requests.Integer              `name:"FailureThreshold"`
	TimeoutSeconds      requests.Integer              `name:"TimeoutSeconds"`
	HttpGet             CreateContainerGroupHttpGet   `name:"HttpGet"`
	Exec                CreateContainerGroupExec      `name:"Exec"`
	TcpSocket           CreateContainerGroupTcpSocket `name:"TcpSocket"`
}

type CreateContainerGroupHttpGet struct {
	Path   string           `name:"Path"`
	Port   requests.Integer `name:"Port"`
	Scheme string           `name:"Scheme"`
}

type CreateContainerGroupExec struct {
	Command []string `name:"Command"`
}

type CreateContainerGroupTcpSocket struct {
	Port requests.Integer `name:"Port"`
}

type CreateContainerGroupLivenessProbe struct {
	InitialDelaySeconds requests.Integer              `name:"InitialDelaySeconds"`
	PeriodSeconds       requests.Integer              `name:"PeriodSeconds"`
	SuccessThreshold    requests.Integer              `name:"SuccessThreshold"`
	FailureThreshold    requests.Integer              `name:"FailureThreshold"`
	TimeoutSeconds      requests.Integer              `name:"TimeoutSeconds"`
	HttpGet             CreateContainerGroupHttpGet   `name:"HttpGet"`
	Exec                CreateContainerGroupExec      `name:"Exec"`
	TcpSocket           CreateContainerGroupTcpSocket `name:"TcpSocket"`
}

type CreateContainerGroupNFSVolume struct {
	Server   string           `name:"Server"`
	Path     string           `name:"Path"`
	ReadOnly requests.Boolean `name:"ReadOnly"`
}

type CreateContainerGroupConfigFileVolume struct {
	ConfigFileToPath *[]CreateContainerGroupConfigFileToPath `name:"ConfigFileToPath"`
	DefaultModel     requests.Integer                        `name:"DefaultModel"`
}

type CreateContainerGroupConfigFileToPath struct {
	Content string           `name:"Content"`
	Path    string           `name:"Path"`
	Mode    requests.Integer `name:"Mode"`
}

type CreateContainerGroupEmptyDirVolume struct {
	Medium string `name:"Medium"`
}

type CreateContainerGroupDiskVolume struct {
	DiskId   string           `name:"DiskId"`
	FsType   string           `name:"FsType"`
	DiskSize requests.Integer `name:"DiskSize"`
}

type CreateContainerGroupFlexVolume struct {
	Driver  string `name:"Driver"`
	FsType  string `name:"FsType"`
	Options string `name:"Options"`
}

type CreateContainerGroupHostPathVolume struct {
	Type string `name:"Type"`
	Path string `name:"Path"`
}

type CreateContainerGroupOption struct {
	Name  string `name:"Name"`
	Value string `name:"Value"`
}

type CreateContainerGroupSysctl struct {
	Name  string `name:"Name"`
	Value string `name:"Value"`
}

// CreateContainerGroupResponse is the response struct for api CreateContainerGroup
type CreateContainerGroupResponse struct {
	*responses.BaseResponse
	RequestId        string `json:"RequestId" xml:"RequestId"`
	ContainerGroupId string `json:"ContainerGroupId" xml:"ContainerGroupId"`
}

// CreateCreateContainerGroupRequest creates a request to invoke CreateContainerGroup API
func CreateCreateContainerGroupRequest() (request *CreateContainerGroupRequest) {
	request = &CreateContainerGroupRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Eci", "2018-08-08", "CreateContainerGroup", "eci", "openAPI")
	return
}

// CreateCreateContainerGroupResponse creates a response to parse from CreateContainerGroup response
func CreateCreateContainerGroupResponse() (response *CreateContainerGroupResponse) {
	response = &CreateContainerGroupResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
