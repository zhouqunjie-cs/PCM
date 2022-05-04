package server

import (
	"context"

	"github.com/JCCE-nudt/PCM/adaptor/vm_adaptor/server/ecs"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbecs"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateMultipleEcs return create cloudy ecs
func (s *Server) CreateMultipleEcs(ctx context.Context, reqs *pbecs.CreateEcsMultipleReq) (*pbecs.CreateEcsMultipleResp, error) {
	resp, err := ecs.CreateMultipleEcs(ctx, reqs)
	if err != nil {
		glog.Errorf("ListEcsDetail error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

// CreateEcs return create ecs
func (s *Server) CreateEcs(ctx context.Context, req *pbecs.CreateEcsReq) (*pbecs.CreateEcsResp, error) {
	resp, err := ecs.CreateEcs(ctx, req)
	if err != nil {
		glog.Errorf("ListEcsDetail error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

// DeleteEcs return Delete ecs
func (s *Server) DeleteEcs(ctx context.Context, req *pbecs.DeleteEcsReq) (*pbecs.DeleteEcsResp, error) {
	resp, err := ecs.DeleteEcs(ctx, req)
	if err != nil {
		glog.Errorf("ListEcsDetail error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

// UpdateEcs return Update ecs
func (s *Server) UpdateEcs(ctx context.Context, req *pbecs.UpdateEcsReq) (*pbecs.UpdateEcsResp, error) {
	resp, err := ecs.UpdateEcs(ctx, req)
	if err != nil {
		glog.Errorf("ListEcsDetail error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

// ListEcsDetail return ecs detail
func (s *Server) ListEcsDetail(ctx context.Context, req *pbecs.ListDetailReq) (*pbecs.ListDetailResp, error) {
	resp, err := ecs.ListDetail(ctx, req)
	if err != nil {
		glog.Errorf("ListEcsDetail error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

//ListEcs return ecs list
func (s *Server) ListEcs(ctx context.Context, req *pbecs.ListReq) (*pbecs.ListResp, error) {
	resp, err := ecs.List(ctx, req)
	if err != nil {
		glog.Errorf("ListEcs error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

// ListEcsAll return all ecs
func (s *Server) ListEcsAll(ctx context.Context, req *pbecs.ListAllReq) (*pbecs.ListResp, error) {
	resp, err := ecs.ListAll(ctx)
	if err != nil {
		glog.Errorf("ListEcsAll error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}

// ActionEcs return ecs action
func (s *Server) ActionEcs(ctx context.Context, req *pbecs.ActionReq) (*pbecs.ActionResp, error) {
	resp, err := ecs.ActionEcs(ctx, req)
	if err != nil {
		glog.Errorf("ActionEcs error %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return resp, nil
}
