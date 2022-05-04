package server

import (
	"context"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/demo"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbecs"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbpod"
)

type Server struct {
	// 使用unsafe可以强制让编译器检查是否实现了相关方法
	demo.UnsafeDemoServiceServer
	pbecs.UnsafeEcsServiceServer
	pbpod.UnsafePodServiceServer
}

func (s *Server) Echo(ctx context.Context, req *demo.StringMessage) (*demo.StringMessage, error) {
	return &demo.StringMessage{
		Value: "Welcome to JCCE PCM",
	}, nil
}
