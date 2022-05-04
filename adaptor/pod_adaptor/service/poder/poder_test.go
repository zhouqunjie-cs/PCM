package poder

import (
	"context"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbpod"
	"testing"

	"github.com/JCCE-nudt/PCM/common/tenanter"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbtenant"
)

func TestEcser_ListDetail(t *testing.T) {
	region, _ := tenanter.NewRegion(pbtenant.CloudProvider_ali, int32(pbtenant.AliRegionId_ali_cn_hangzhou))
	ali, _ := NewPodClient(pbtenant.CloudProvider_ali, region, aliTenant[0])
	aliFailed, _ := NewPodClient(pbtenant.CloudProvider_ali, region, tenanter.NewTenantWithAccessKey("empty", "", "", "", ""))

	region, _ = tenanter.NewRegion(pbtenant.CloudProvider_tencent, int32(pbtenant.TencentRegionId_tc_ap_beijing))
	tc, _ := NewPodClient(pbtenant.CloudProvider_tencent, region, tcTenant[0])
	tcFailed, _ := NewPodClient(pbtenant.CloudProvider_tencent, region, tenanter.NewTenantWithAccessKey("empty", "", "", "", ""))

	region, _ = tenanter.NewRegion(pbtenant.CloudProvider_huawei, int32(pbtenant.HuaweiRegionId_hw_cn_southwest_2))
	hw, _ := NewPodClient(pbtenant.CloudProvider_huawei, region, hwTenant[0])
	// hwFailed, _ := newHuaweiEcsClient(int32(pbtenant.HuaweiRegionId_hw_cn_north_1), tenanter.NewTenantWithAccessKey("empty", "", "", ""))

	region, _ = tenanter.NewRegion(pbtenant.CloudProvider_k8s, int32(pbtenant.AwsRegionId_aws_us_east_2))
	aws, _ := NewPodClient(pbtenant.CloudProvider_k8s, region, awsTenant[0])

	// google, _ := NewGoogleEcsClient(tenanter.NewTenantWithAccessKey("", ""))

	type args struct {
		req *pbpod.ListPodDetailReq
	}
	tests := []struct {
		name    string
		fields  Poder
		args    args
		wantErr bool
	}{
		{name: "ali wrong cli", fields: aliFailed, args: args{&pbpod.ListPodDetailReq{PageNumber: 1, PageSize: 1}}, wantErr: true},
		{name: "ali wrong page number", fields: ali, args: args{&pbpod.ListPodDetailReq{PageNumber: 0, PageSize: 1}}, wantErr: true},
		{name: "ali wrong page size", fields: ali, args: args{&pbpod.ListPodDetailReq{PageNumber: 1, PageSize: 0}}, wantErr: true},
		{name: "ali right cli", fields: ali, args: args{&pbpod.ListPodDetailReq{PageNumber: 1, PageSize: 10}}, wantErr: false},

		{name: "tc wrong cli", fields: tcFailed, args: args{&pbpod.ListPodDetailReq{PageNumber: 1, PageSize: 1}}, wantErr: true},
		{name: "tc wrong page number", fields: tc, args: args{&pbpod.ListPodDetailReq{PageNumber: 0, PageSize: 1}}, wantErr: true},
		{name: "tc wrong page size", fields: tc, args: args{&pbpod.ListPodDetailReq{PageNumber: 1, PageSize: 0}}, wantErr: true},
		{name: "tc right cli", fields: tc, args: args{&pbpod.ListPodDetailReq{PageNumber: 1, PageSize: 10}}, wantErr: false},

		// {name: "hw wrong cli", fields: hwFailed, args: args{pageNumber: 1, pageSize: 1}, wantErr: true},
		{name: "hw right cli", fields: hw, args: args{&pbpod.ListPodDetailReq{PageNumber: 1, PageSize: 10}}, wantErr: false},

		{name: "aws right cli", fields: aws, args: args{&pbpod.ListPodDetailReq{PageNumber: 1, PageSize: 10}}, wantErr: false},

		// {name: "right cli", fields: google, args: args{pageNumber: 1, pageSize: 10}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.fields.ListPodDetail(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListDetail() error = %+v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%+v", err)
			if err == nil {
				t.Log(resp)
			}
		})
	}
}
