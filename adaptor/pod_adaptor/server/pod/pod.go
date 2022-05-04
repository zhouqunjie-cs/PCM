package pod

import (
	"context"
	"fmt"
	"sync"

	"github.com/JCCE-nudt/PCM/adaptor/pod_adaptor/service/poder"
	"github.com/JCCE-nudt/PCM/common/tenanter"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbpod"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbtenant"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// GetPodRegion get the available region for pod
func GetPodRegion(ctx context.Context, req *pbpod.GetPodRegionReq) (resp *pbpod.GetPodRegionResp, err error) {
	var (
		regionInit tenanter.Region
		regions    []*pbtenant.Region
	)

	switch req.GetProvider() {
	case pbtenant.CloudProvider_ali:
		regionInit, _ = tenanter.NewRegion(req.GetProvider(), 2)
	case pbtenant.CloudProvider_tencent:
		regionInit, _ = tenanter.NewRegion(req.GetProvider(), 5)
	case pbtenant.CloudProvider_huawei:
		regionInit, _ = tenanter.NewRegion(req.GetProvider(), 5)
	}
	tenanters, err := tenanter.GetTenanters(req.GetProvider())
	if err != nil {
		return nil, errors.WithMessage(err, "getTenanters error")
	}

	for _, tenant := range tenanters {

		pod, err := poder.NewPodClient(req.GetProvider(), regionInit, tenant)
		if err != nil {
			return nil, errors.WithMessage(err, "NewPodClient error")
		}
		request := &pbpod.GetPodRegionReq{
			Provider: req.GetProvider(),
		}
		resp, err := pod.GetPodRegion(ctx, request)
		if err != nil {
			return nil, errors.Wrap(err, "GetPodRegion error")
		}
		for _, region := range resp.GetRegions() {
			regions = append(regions, region)
		}
	}

	return &pbpod.GetPodRegionResp{Regions: regions}, nil
}

func CreatePods(ctx context.Context, req *pbpod.CreatePodsReq) (*pbpod.CreatePodsResp, error) {
	var (
		wg         sync.WaitGroup
		requestIds = make([]string, 0)
	)
	wg.Add(len(req.CreatePodReq))
	c := make(chan string)
	for k := range req.CreatePodReq {
		reqPod := req.CreatePodReq[k]
		go func() {
			defer wg.Done()
			resp, err := CreatePod(ctx, reqPod)
			if err != nil || resp == nil {
				fmt.Println(errors.Wrap(err, "Batch pod creation error"))
				return
			}
			c <- resp.RequestId
		}()

	}
	go func() {
		defer close(c)
		wg.Wait()
	}()
	isFinished := false
	if len(requestIds) > 0 {
		isFinished = true
	}

	for v := range c {
		requestIds = append(requestIds, v)
	}

	return &pbpod.CreatePodsResp{
		Finished:  isFinished,
		RequestId: requestIds,
	}, nil
}

func CreatePod(ctx context.Context, req *pbpod.CreatePodReq) (*pbpod.CreatePodResp, error) {
	var (
		pod poder.Poder
	)

	tenanters, err := tenanter.GetTenanters(req.Provider)
	if err != nil {
		return nil, errors.WithMessage(err, "getTenanters error")
	}

	region, err := tenanter.NewRegion(req.Provider, req.RegionId)
	if err != nil {
		return nil, errors.WithMessagef(err, "provider %v regionId %v", req.Provider, req.RegionId)
	}

	for _, tenant := range tenanters {
		if req.AccountName == "" || tenant.AccountName() == req.AccountName {
			if pod, err = poder.NewPodClient(req.Provider, region, tenant); err != nil {
				return nil, errors.WithMessage(err, "NewPodClient error")
			}
			break
		}
	}

	return pod.CreatePod(ctx, req)
}

func DeletePod(ctx context.Context, req *pbpod.DeletePodReq) (*pbpod.DeletePodResp, error) {
	var (
		pod poder.Poder
	)

	tenanters, err := tenanter.GetTenanters(req.Provider)
	if err != nil {
		return nil, errors.WithMessage(err, "getTenanters error")
	}

	region, err := tenanter.NewRegion(req.Provider, req.RegionId)
	if err != nil {
		return nil, errors.WithMessagef(err, "provider %v regionId %v", req.Provider, req.RegionId)
	}

	for _, tenant := range tenanters {
		if req.AccountName == "" || tenant.AccountName() == req.AccountName {
			if pod, err = poder.NewPodClient(req.Provider, region, tenant); err != nil {
				return nil, errors.WithMessage(err, "NewPodClient error")
			}
			break
		}
	}

	return pod.DeletePod(ctx, req)
}

func UpdatePod(ctx context.Context, req *pbpod.UpdatePodReq) (*pbpod.UpdatePodResp, error) {
	var (
		pod poder.Poder
	)

	tenanters, err := tenanter.GetTenanters(req.Provider)
	if err != nil {
		return nil, errors.WithMessage(err, "getTenanters error")
	}

	region, err := tenanter.NewRegion(req.Provider, req.RegionId)
	if err != nil {
		return nil, errors.WithMessagef(err, "provider %v regionId %v", req.Provider, req.RegionId)
	}

	for _, tenant := range tenanters {
		if req.AccountName == "" || tenant.AccountName() == req.AccountName {
			if pod, err = poder.NewPodClient(req.Provider, region, tenant); err != nil {
				return nil, errors.WithMessage(err, "NewPodClient error")
			}
			break
		}
	}

	return pod.UpdatePod(ctx, req)
}

func ListPodDetail(ctx context.Context, req *pbpod.ListPodDetailReq) (*pbpod.ListPodDetailResp, error) {
	var (
		pod poder.Poder
	)

	tenanters, err := tenanter.GetTenanters(req.Provider)
	if err != nil {
		return nil, errors.WithMessage(err, "getTenanters error")
	}

	region, err := tenanter.NewRegion(req.Provider, req.RegionId)
	if err != nil {
		return nil, errors.WithMessagef(err, "provider %v regionId %v", req.Provider, req.RegionId)
	}

	for _, tenant := range tenanters {
		if req.AccountName == "" || tenant.AccountName() == req.AccountName {
			if pod, err = poder.NewPodClient(req.Provider, region, tenant); err != nil {
				return nil, errors.WithMessage(err, "NewPodClient error")
			}
			break
		}
	}

	return pod.ListPodDetail(ctx, req)
}

func ListPod(ctx context.Context, req *pbpod.ListPodReq) (*pbpod.ListPodResp, error) {
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
		pods  []*pbpod.PodInstance
	)

	tenanters, err := tenanter.GetTenanters(req.Provider)
	if err != nil {
		return nil, errors.WithMessage(err, "getTenanters error")
	}
	//get the available region for product
	reqPodRegion := &pbpod.GetPodRegionReq{Provider: req.GetProvider()}
	respPodRegion, err := GetPodRegion(ctx, reqPodRegion)
	if err != nil {
		return nil, errors.WithMessage(err, "getPodRegion error")
	}

	wg.Add(len(tenanters) * len(respPodRegion.Regions))
	for _, t := range tenanters {
		for _, region := range respPodRegion.Regions {
			go func(tenant tenanter.Tenanter, region tenanter.Region) {
				defer wg.Done()
				pod, err := poder.NewPodClient(req.Provider, region, tenant)
				if err != nil {
					glog.Errorf("New Pod Client error %v", err)
					return
				}

				request := &pbpod.ListPodDetailReq{
					Provider:    req.Provider,
					AccountName: tenant.AccountName(),
					RegionId:    region.GetId(),
					Namespace:   req.Namespace,
					PageNumber:  1,
					PageSize:    100,
					NextToken:   "",
				}
				for {
					resp, err := pod.ListPodDetail(ctx, request)
					if err != nil {
						glog.Errorf("ListDetail error %v", err)
						return
					}
					mutex.Lock()
					pods = append(pods, resp.Pods...)
					mutex.Unlock()
					if resp.Finished {
						break
					}
					request.PageNumber, request.PageSize, request.NextToken = resp.PageNumber, resp.PageSize, resp.NextToken
				}
			}(t, region)

		}
	}
	wg.Wait()

	return &pbpod.ListPodResp{Pods: pods}, nil
}

func ListPodAll(ctx context.Context) (*pbpod.ListPodResp, error) {
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
		pods  []*pbpod.PodInstance
	)

	wg.Add(len(pbtenant.CloudProvider_name))
	for k := range pbtenant.CloudProvider_name {
		go func(provider int32) {
			defer wg.Done()

			//针对私有K8S集群，调用listAll时默认只查询ListPodDetailReq namespace下的pod
			if provider == 3 {
				resp, err := ListPod(ctx, &pbpod.ListPodReq{Provider: pbtenant.CloudProvider(provider), Namespace: "pcm"})
				if err != nil {
					glog.Errorf("List error %v", err)
					return
				}
				mutex.Lock()
				pods = append(pods, resp.Pods...)
				mutex.Unlock()
			} else {
				resp, err := ListPod(ctx, &pbpod.ListPodReq{Provider: pbtenant.CloudProvider(provider)})
				if err != nil {
					glog.Errorf("List error %v", err)
					return
				}
				mutex.Lock()
				pods = append(pods, resp.Pods...)
				mutex.Unlock()
			}

		}(k)
	}

	wg.Wait()

	return &pbpod.ListPodResp{Pods: pods}, nil
}
