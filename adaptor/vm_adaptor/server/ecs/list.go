package ecs

import (
	"context"
	"sync"

	"github.com/JCCE-nudt/PCM/adaptor/vm_adaptor/service/ecser"
	"github.com/JCCE-nudt/PCM/common/tenanter"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbecs"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbtenant"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

//CreateMultipleEcs 创建多云ECS
func CreateMultipleEcs(ctx context.Context, reqs *pbecs.CreateEcsMultipleReq) (*pbecs.CreateEcsMultipleResp, error) {
	var (
		wg         sync.WaitGroup
		requestIds = make([]string, 0)
	)
	wg.Add(len(reqs.GetCreateEcsReqs()))
	c := make(chan string, len(reqs.GetCreateEcsReqs()))
	for _, k := range reqs.GetCreateEcsReqs() {
		k := k
		go func() {
			defer wg.Done()
			resp, err := CreateEcs(ctx, k)
			if err != nil {
				glog.Errorf(k.Provider.String()+"CreateEcs error: %v", err)
				c <- k.Provider.String()
				return
			}
			c <- resp.GetRequestId()
		}()
	}
	go func() {
		defer close(c)
		wg.Wait()
	}()
	for v := range c {
		requestIds = append(requestIds, v)
	}
	isFinished := false
	if len(requestIds) > 0 {
		isFinished = true
	}
	return &pbecs.CreateEcsMultipleResp{
		RequestId: requestIds,
		Finished:  isFinished,
	}, nil
}

func CreateEcs(ctx context.Context, req *pbecs.CreateEcsReq) (*pbecs.CreateEcsResp, error) {
	var (
		ecs ecser.Ecser
	)
	tenanters, err := tenanter.GetTenanters(req.Provider)
	region, err := tenanter.NewRegion(req.Provider, req.RegionId)
	if err != nil {
		return nil, errors.Wrap(err, "get tenanters failed")
	}
	for _, tenanter := range tenanters {
		if req.AccountName == "" || tenanter.AccountName() == req.AccountName {
			if ecs, err = ecser.NewEcsClient(req.Provider, region, tenanter); err != nil {
				return nil, errors.WithMessage(err, "NewEcsClient error")
			}
			break
		}
	}
	return ecs.CreateEcs(ctx, req)
}

func DeleteEcs(ctx context.Context, req *pbecs.DeleteEcsReq) (*pbecs.DeleteEcsResp, error) {
	var (
		ecs ecser.Ecser
	)
	tenanters, err := tenanter.GetTenanters(req.Provider)
	region, err := tenanter.NewRegion(req.Provider, req.RegionId)
	if err != nil {
		return nil, errors.Wrap(err, "get tenanters failed")
	}
	for _, tenanter := range tenanters {
		if req.AccountName == "" || tenanter.AccountName() == req.AccountName {
			if ecs, err = ecser.NewEcsClient(req.Provider, region, tenanter); err != nil {
				return nil, errors.WithMessage(err, "NewEcsClient error")
			}
			break
		}
	}
	return ecs.DeleteEcs(ctx, req)
}

func UpdateEcs(ctx context.Context, req *pbecs.UpdateEcsReq) (*pbecs.UpdateEcsResp, error) {
	var (
		ecs ecser.Ecser
	)
	tenanters, err := tenanter.GetTenanters(req.Provider)
	region, err := tenanter.NewRegion(req.Provider, req.RegionId)
	if err != nil {
		return nil, errors.Wrap(err, "get tenanters failed")
	}
	for _, tenanter := range tenanters {
		if req.AccountName == "" || tenanter.AccountName() == req.AccountName {
			if ecs, err = ecser.NewEcsClient(req.Provider, region, tenanter); err != nil {
				return nil, errors.WithMessage(err, "NewEcsClient error")
			}
			break
		}
	}
	return ecs.UpdateEcs(ctx, req)
}

//ListDetail returns the detail of ecs instances
func ListDetail(ctx context.Context, req *pbecs.ListDetailReq) (*pbecs.ListDetailResp, error) {
	var (
		ecs ecser.Ecser
	)

	tenanters, err := tenanter.GetTenanters(req.Provider)
	if err != nil {
		return nil, errors.WithMessage(err, "getTenanters error")
	}

	region, err := tenanter.NewRegion(req.Provider, req.RegionId)
	if err != nil {
		return nil, errors.WithMessagef(err, "provider %v regionId %v", req.Provider, req.RegionId)
	}

	for _, tenanter := range tenanters {
		if req.AccountName == "" || tenanter.AccountName() == req.AccountName {
			if ecs, err = ecser.NewEcsClient(req.Provider, region, tenanter); err != nil {
				return nil, errors.WithMessage(err, "NewEcsClient error")
			}
			break
		}
	}

	return ecs.ListDetail(ctx, req)
}

//List returns the list of ecs instances
func List(ctx context.Context, req *pbecs.ListReq) (*pbecs.ListResp, error) {
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
		ecses []*pbecs.EcsInstance
	)

	tenanters, err := tenanter.GetTenanters(req.Provider)
	if err != nil {
		return nil, errors.WithMessage(err, "getTenanters error")
	}

	regions := tenanter.GetAllRegionIds(req.Provider)

	wg.Add(len(tenanters) * len(regions))
	for _, t := range tenanters {
		for _, region := range regions {
			go func(tenant tenanter.Tenanter, region tenanter.Region) {
				defer wg.Done()
				ecs, err := ecser.NewEcsClient(req.Provider, region, tenant)
				if err != nil {
					glog.Errorf("New Ecs Client error %v", err)
					return
				}

				request := &pbecs.ListDetailReq{
					Provider:    req.Provider,
					AccountName: tenant.AccountName(),
					RegionId:    region.GetId(),
					PageNumber:  1,
					PageSize:    100,
					NextToken:   "",
				}
				for {
					resp, err := ecs.ListDetail(ctx, request)
					if err != nil {
						glog.Errorf("ListDetail error %v", err)
						return
					}
					mutex.Lock()
					ecses = append(ecses, resp.Ecses...)
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

	return &pbecs.ListResp{Ecses: ecses}, nil
}

// ListAll returns all ecs instances
func ListAll(ctx context.Context) (*pbecs.ListResp, error) {
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
		ecses []*pbecs.EcsInstance
	)

	wg.Add(len(pbtenant.CloudProvider_name))
	for k := range pbtenant.CloudProvider_name {
		go func(provider int32) {
			defer wg.Done()

			resp, err := List(ctx, &pbecs.ListReq{Provider: pbtenant.CloudProvider(provider)})
			if err != nil {
				glog.Errorf("List error %v", err)
				return
			}
			mutex.Lock()
			ecses = append(ecses, resp.Ecses...)
			mutex.Unlock()
		}(k)
	}

	wg.Wait()

	return &pbecs.ListResp{Ecses: ecses}, nil
}

func ActionEcs(ctx context.Context, req *pbecs.ActionReq) (*pbecs.ActionResp, error) {
	var (
		ecs ecser.Ecser
	)
	tenanters, err := tenanter.GetTenanters(req.Provider)
	region, err := tenanter.NewRegion(req.Provider, req.RegionId)
	if err != nil {
		return nil, errors.Wrap(err, "get tenanters failed")
	}
	for _, tenanter := range tenanters {
		if req.AccountName == "" || tenanter.AccountName() == req.AccountName {
			if ecs, err = ecser.NewEcsClient(req.Provider, region, tenanter); err != nil {
				return nil, errors.WithMessage(err, "NewEcsClient error")
			}
			break
		}
	}
	return ecs.ActionEcs(ctx, req)
}
