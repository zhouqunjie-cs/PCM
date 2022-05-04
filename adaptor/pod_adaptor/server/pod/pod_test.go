package pod

import (
	"context"
	"testing"

	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbpod"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbtenant"
)

func TestListDetail(t *testing.T) {
	type args struct {
		req *pbpod.ListPodDetailReq
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "ali", args: args{req: &pbpod.ListPodDetailReq{Provider: pbtenant.CloudProvider_ali, RegionId: int32(pbtenant.AliRegionId_ali_cn_hangzhou), PageNumber: 1, PageSize: 10}}, wantErr: false},
		{name: "tencent", args: args{req: &pbpod.ListPodDetailReq{Provider: pbtenant.CloudProvider_tencent, RegionId: int32(pbtenant.TencentRegionId_tc_ap_beijing), PageNumber: 1, PageSize: 10}}, wantErr: false},
		{name: "aws", args: args{req: &pbpod.ListPodDetailReq{Provider: pbtenant.CloudProvider_aws, RegionId: int32(pbtenant.AwsRegionId_aws_us_east_2), PageNumber: 1, PageSize: 10}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListPodDetail(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListDetail() error = %+v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}

func TestList(t *testing.T) {
	type args struct {
		req *pbpod.ListPodReq
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "ali", args: args{req: &pbpod.ListPodReq{Provider: pbtenant.CloudProvider_ali}}, wantErr: false},
		{name: "tencent", args: args{req: &pbpod.ListPodReq{Provider: pbtenant.CloudProvider_tencent}}, wantErr: false},
		{name: "huawei", args: args{req: &pbpod.ListPodReq{Provider: pbtenant.CloudProvider_huawei}}, wantErr: false},
		{name: "aws", args: args{req: &pbpod.ListPodReq{Provider: pbtenant.CloudProvider_aws}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListPod(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %+v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}

func TestListAll(t *testing.T) {
	type args struct {
		req *pbpod.ListPodReq
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "all", args: args{req: &pbpod.ListPodReq{}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListPodAll(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAll() error = %+v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}
