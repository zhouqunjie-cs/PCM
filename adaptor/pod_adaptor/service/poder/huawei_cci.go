package poder

import (
	"context"
	"fmt"
	util "github.com/alibabacloud-go/tea-utils/service"
	"sync"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbpod"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbtenant"

	corev1 "k8s.io/api/core/v1"
	huaweicci "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/JCCE-nudt/PCM/common/tenanter"
	"github.com/pkg/errors"
)

var huaweiClientMutex sync.Mutex

type HuaweiCci struct {
	cli      *huaweicci.Clientset
	region   tenanter.Region
	tenanter tenanter.Tenanter
}

func (cci *HuaweiCci) GetPodRegion(ctx context.Context, req *pbpod.GetPodRegionReq) (*pbpod.GetPodRegionResp, error) {
	//todo
	var (
		regions []*pbtenant.Region
	)
	huaweiRegionName, _ := tenanter.GetHuaweiRegionName(5)
	region := &pbtenant.Region{
		Id:   5,
		Name: huaweiRegionName,
	}
	regions = append(regions, region)
	resp := &pbpod.GetPodRegionResp{
		Regions: regions,
	}

	return resp, nil
}

//CCI auth through iam
const (
	apiVersion  = "client.authentication.k8s.io/v1beta1"
	iamEndpoint = "https://iam.myhuaweicloud.com"
)

func newHuaweiCciClient(region tenanter.Region, tenant tenanter.Tenanter) (Poder, error) {
	var (
		client *huaweicci.Clientset
		err    error
	)
	cciEndpoint := "https://cci." + region.GetName() + ".myhuaweicloud.com"
	cciConfig, err := clientcmd.BuildConfigFromFlags(cciEndpoint, "")
	if err != nil {
		return nil, err
	}

	switch t := tenant.(type) {
	case *tenanter.AccessKeyTenant:
		huaweiClientMutex.Lock()
		var optionArgs []string
		optionArgs = append(optionArgs, fmt.Sprintf("--iam-endpoint=%s", iamEndpoint))
		optionArgs = append(optionArgs, fmt.Sprintf("--project-name=%s", region.GetName()))
		optionArgs = append(optionArgs, fmt.Sprintf("--ak=%s", t.GetId()))
		optionArgs = append(optionArgs, fmt.Sprintf("--sk=%s", t.GetSecret()))
		cciConfig.ExecProvider = &api.ExecConfig{
			Command:    "cci-iam-authenticator",
			APIVersion: apiVersion,
			Args:       append([]string{"token"}, optionArgs...),
			Env:        make([]api.ExecEnvVar, 0),
		}
		client, err = huaweicci.NewForConfig(cciConfig)
		huaweiClientMutex.Unlock()
	default:
	}

	if err != nil {
		return nil, errors.Wrap(err, "init huawei pod client error")
	}

	return &HuaweiCci{
		cli:      client,
		region:   region,
		tenanter: tenant,
	}, nil
}

func (cci *HuaweiCci) CreatePod(ctx context.Context, req *pbpod.CreatePodReq) (*pbpod.CreatePodResp, error) {
	pod := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "core/V1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.PodName,
			Namespace: req.Namespace,
			Labels:    map[string]string{"name": "test_api"},
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyAlways,
			Containers: []corev1.Container{
				{
					Name:  req.ContainerName,
					Image: req.ContainerImage,
					Resources: corev1.ResourceRequirements{
						Limits: map[corev1.ResourceName]resource.Quantity{
							corev1.ResourceCPU:    resource.MustParse(req.CpuPod),
							corev1.ResourceMemory: resource.MustParse(req.MemoryPod),
						},
					},
				},
			},
		},
		Status: corev1.PodStatus{},
	}

	resp, err := cci.cli.CoreV1().Pods(req.Namespace).Create(context.TODO(), &pod, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "Huaweiyun CreatePod error")
	}

	isFinished := false
	if len(resp.UID) > 0 {
		isFinished = true
	}

	glog.Infof("--------------------Huawei CCI Instance created--------------------")

	return &pbpod.CreatePodResp{
		Finished:  isFinished,
		RequestId: "Create huawei pod request ID:" + resp.GenerateName,
		PodId:     string(resp.Generation),
		PodName:   resp.Name,
	}, nil
}

func (cci *HuaweiCci) DeletePod(ctx context.Context, req *pbpod.DeletePodReq) (*pbpod.DeletePodResp, error) {

	err := cci.cli.CoreV1().Pods(req.GetNamespace()).Delete(context.TODO(), req.PodName, metav1.DeleteOptions{})

	isFinished := true
	if err != nil {
		isFinished = false
		return nil, errors.Wrap(err, "Huaweiyun DeletePod error")
	}

	glog.Infof("--------------------Huawei CCI Instance deleted--------------------")
	glog.Infof(*util.ToJSONString(util.ToMap(err)))

	return &pbpod.DeletePodResp{
		Finished:  isFinished,
		RequestId: "Delete huawei pod request ID:" + req.PodName,
		PodId:     req.PodName,
		PodName:   req.PodName,
	}, nil
}

func (cci *HuaweiCci) UpdatePod(ctx context.Context, req *pbpod.UpdatePodReq) (*pbpod.UpdatePodResp, error) {

	qresp, err := cci.cli.CoreV1().Pods(req.GetNamespace()).Get(context.TODO(), req.PodName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "Huaweiyun UpdatePod error")
	}
	pod := corev1.Pod{
		TypeMeta: qresp.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.PodName,
			Namespace: req.Namespace,
			Labels:    map[string]string{"name": req.Labels},
		},
		Spec:   qresp.Spec,
		Status: qresp.Status,
	}
	pod.Spec.Containers[0].Image = req.ContainerImage
	resp, err := cci.cli.CoreV1().Pods(req.Namespace).Update(context.TODO(), &pod, metav1.UpdateOptions{})
	glog.Info("Huawei update pod resp", resp)
	if err != nil {
		return nil, errors.Wrap(err, "Huaweiyun UpdatePod error")
	}

	glog.Infof("--------------------Huawei CCI Instance updated--------------------")

	isFinished := false
	if len(resp.UID) > 0 {
		isFinished = true
	}

	return &pbpod.UpdatePodResp{
		Finished:  isFinished,
		RequestId: "Update huawei pod request ID:" + resp.GenerateName,
		PodId:     string(resp.Generation),
		PodName:   resp.Name,
	}, nil
}

func (cci *HuaweiCci) ListPodDetail(ctx context.Context, req *pbpod.ListPodDetailReq) (*pbpod.ListPodDetailResp, error) {

	resp, err := cci.cli.CoreV1().Pods(req.GetNamespace()).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, errors.Wrap(err, "Huaweiyun ListDetail pod error")
	}
	glog.Info("Huaweiyun ListDetail pod success", resp.Items)
	var pods = make([]*pbpod.PodInstance, len(resp.Items))
	for k, v := range resp.Items {
		pods[k] = &pbpod.PodInstance{
			Provider:       pbtenant.CloudProvider_huawei,
			AccountName:    cci.tenanter.AccountName(),
			PodId:          string(v.GetUID()),
			PodName:        v.Name,
			RegionId:       cci.region.GetId(),
			ContainerImage: v.Spec.Containers[0].Image,
			ContainerName:  v.Spec.Containers[0].Name,
			CpuPod:         v.Spec.Containers[0].Resources.Requests.Cpu().String(),
			MemoryPod:      v.Spec.Containers[0].Resources.Requests.Memory().String(),
			Namespace:      v.Namespace,
			Status:         string(v.Status.Phase),
		}
	}

	glog.Infof("--------------------Huawei CCI Instance updated--------------------")

	isFinished := false
	if len(pods) < int(req.PageSize) {
		isFinished = true
	}

	return &pbpod.ListPodDetailResp{
		Pods:       pods,
		Finished:   isFinished,
		PageNumber: req.PageNumber + 1,
		PageSize:   req.PageSize,
	}, nil
}
