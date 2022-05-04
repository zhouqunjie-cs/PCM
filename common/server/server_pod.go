package server

import (
	"context"

	"github.com/JCCE-nudt/PCM/adaptor/pod_adaptor/server/pod"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbpod"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetProdRegions get available region for product
func (s *Server) GetProdRegions(ctx context.Context, req *pbpod.GetPodRegionReq) (*pbpod.GetPodRegionResp, error) {
	resp, err := pod.GetPodRegion(ctx, req)
	if err != nil {
		glog.Errorf("CreatePods error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

// CreatePods create multiple pod on multiple clouds
func (s *Server) CreatePods(ctx context.Context, req *pbpod.CreatePodsReq) (*pbpod.CreatePodsResp, error) {
	resp, err := pod.CreatePods(ctx, req)
	if err != nil {
		glog.Errorf("CreatePods error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

// CreatePod create pod on one cloud
func (s *Server) CreatePod(ctx context.Context, req *pbpod.CreatePodReq) (*pbpod.CreatePodResp, error) {
	resp, err := pod.CreatePod(ctx, req)
	if err != nil {
		glog.Errorf("CreatePod error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

// DeletePod delete specified pod
func (s *Server) DeletePod(ctx context.Context, req *pbpod.DeletePodReq) (*pbpod.DeletePodResp, error) {
	resp, err := pod.DeletePod(ctx, req)
	if err != nil {
		glog.Errorf("DeletePod error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

// UpdatePod update specified pod
func (s *Server) UpdatePod(ctx context.Context, req *pbpod.UpdatePodReq) (*pbpod.UpdatePodResp, error) {
	resp, err := pod.UpdatePod(ctx, req)
	if err != nil {
		glog.Errorf("UpdatePod error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) ListPodDetail(ctx context.Context, req *pbpod.ListPodDetailReq) (*pbpod.ListPodDetailResp, error) {
	resp, err := pod.ListPodDetail(ctx, req)
	if err != nil {
		glog.Errorf("ListPodDetail error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) ListPod(ctx context.Context, req *pbpod.ListPodReq) (*pbpod.ListPodResp, error) {
	resp, err := pod.ListPod(ctx, req)
	if err != nil {
		glog.Errorf("ListPod error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) ListPodAll(ctx context.Context, req *pbpod.ListPodAllReq) (*pbpod.ListPodResp, error) {
	resp, err := pod.ListPodAll(ctx)
	if err != nil {
		glog.Errorf("ListPodAll error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}
