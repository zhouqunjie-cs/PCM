package poder

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/zhouqunjie-cs/PCM/common/tenanter"
	"github.com/zhouqunjie-cs/PCM/lan_trans/idl/pbpod"
	"github.com/zhouqunjie-cs/PCM/lan_trans/idl/pbtenant"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sync"
)

var k8sClientMutex sync.Mutex

type Config struct {
	Host  string
	Token string
	Port  int
}

type K8SPoder struct {
	cli      *k8s.Clientset
	region   tenanter.Region
	tenanter tenanter.Tenanter
}

func (k K8SPoder) GetPodRegion(ctx context.Context, req *pbpod.GetPodRegionReq) (*pbpod.GetPodRegionResp, error) {
	//todo
	var (
		regions []*pbtenant.Region
	)
	huaweiRegionName, _ := tenanter.GetK8SRegionName(0)
	region := &pbtenant.Region{
		Id:   0,
		Name: huaweiRegionName,
	}
	regions = append(regions, region)
	resp := &pbpod.GetPodRegionResp{
		Regions: regions,
	}

	return resp, nil
}

func newK8SClient(tenant tenanter.Tenanter) (Poder, error) {
	var (
		client *k8s.Clientset
		err    error
	)

	switch t := tenant.(type) {
	case *tenanter.AccessKeyTenant:

		kubeConf := &rest.Config{
			Host:        fmt.Sprintf("%s:%d", t.GetUrl(), 6443),
			BearerToken: t.GetToken(),
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true,
			},
		}
		k8sClientMutex.Lock()
		client, err = k8s.NewForConfig(kubeConf)
		k8sClientMutex.Unlock()
	default:

	}

	if err != nil {
		return nil, errors.Wrap(err, "init k8s client error")
	}

	return &K8SPoder{
		cli:      client,
		region:   nil,
		tenanter: tenant,
	}, nil
}

func (k *K8SPoder) CreatePod(ctx context.Context, req *pbpod.CreatePodReq) (*pbpod.CreatePodResp, error) {

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

	resp, err := k.cli.CoreV1().Pods(req.Namespace).Create(context.TODO(), &pod, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "K8S CreatePod error")
	}

	glog.Infof("--------------------K8S Pod Instance created--------------------")

	isFinished := false
	if len(resp.UID) > 0 {
		isFinished = true
	}

	return &pbpod.CreatePodResp{
		Finished:  isFinished,
		RequestId: "K8S pod Name:" + resp.Name,
		PodId:     string(resp.UID),
		PodName:   resp.Name,
	}, nil

}

func (k K8SPoder) DeletePod(ctx context.Context, req *pbpod.DeletePodReq) (*pbpod.DeletePodResp, error) {

	podName := req.PodName
	fmt.Println("K8S ContainerGroup:", podName, " Deleted")
	err := k.cli.CoreV1().Pods(req.Namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})

	glog.Infof("--------------------K8S Pod Instance deleted--------------------")

	isFinished := true
	if err != nil {
		isFinished = false
		return nil, errors.Wrap(err, "K8S DeletePod error")
	}

	return &pbpod.DeletePodResp{
		Finished:  isFinished,
		RequestId: "K8S pod Name:" + req.PodName,
		PodId:     req.PodName,
		PodName:   req.PodName,
	}, nil
}

func (k K8SPoder) UpdatePod(ctx context.Context, req *pbpod.UpdatePodReq) (*pbpod.UpdatePodResp, error) {

	qresp, err := k.cli.CoreV1().Pods(req.GetNamespace()).Get(context.TODO(), req.PodName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "K8S UpdatePod error")
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
	resp, err := k.cli.CoreV1().Pods(req.Namespace).Update(context.TODO(), &pod, metav1.UpdateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "K8S UpdatePod error")
	}

	glog.Infof("--------------------K8S Pod Instance updated--------------------")

	isFinished := false
	if len(resp.UID) > 0 {
		isFinished = true
	}

	return &pbpod.UpdatePodResp{
		Finished:  isFinished,
		RequestId: "K8S pod Name:" + req.PodName,
		PodId:     string(resp.UID),
		PodName:   req.PodName,
	}, nil

}

func (k K8SPoder) ListPodDetail(ctx context.Context, req *pbpod.ListPodDetailReq) (*pbpod.ListPodDetailResp, error) {
	resp, err := k.cli.CoreV1().Pods(req.GetNamespace()).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "K8S ListDetail pod error")
	}

	var pods = make([]*pbpod.PodInstance, len(resp.Items))
	for k, v := range resp.Items {
		pods[k] = &pbpod.PodInstance{
			Provider:       pbtenant.CloudProvider_k8s,
			AccountName:    req.AccountName,
			PodId:          string(v.GetUID()),
			PodName:        v.Name,
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
