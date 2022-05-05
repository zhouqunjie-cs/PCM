package tenanter

import (
	"strings"

	"github.com/zhouqunjie-cs/PCM/lan_trans/idl/pbtenant"

	"github.com/pkg/errors"
)

var (
	ErrNoExistAliRegionId     = errors.New("region id not exist in ali")
	ErrNoExistTencentRegionId = errors.New("region id not exist in tencent")
	ErrNoExistHuaweiRegionId  = errors.New("region id not exist in huawei")
	ErrNoExistAwsRegionId     = errors.New("region id not exist in aws")
)

type Region interface {
	GetId() int32
	GetName() string
}

type region struct {
	provider   pbtenant.CloudProvider
	regionId   int32
	regionName string
}

func NewRegion(provider pbtenant.CloudProvider, regionId int32) (Region, error) {
	r := &region{
		provider: provider,
		regionId: regionId,
	}
	var err error

	switch provider {
	case pbtenant.CloudProvider_ali:
		r.regionName, err = GetAliRegionName(regionId)
	case pbtenant.CloudProvider_tencent:
		r.regionName, err = GetTencentRegionName(regionId)
	case pbtenant.CloudProvider_huawei:
		r.regionName, err = GetHuaweiRegionName(regionId)
		//case pbtenant.CloudProvider_aws:
		//	r.regionName, err = getAwsRegionName(regionId)
	}

	return r, err
}

func (r *region) GetName() string {
	return r.regionName
}

func (r *region) GetId() int32 {
	return r.regionId
}

func GetAllRegionIds(provider pbtenant.CloudProvider) (regions []Region) {
	switch provider {
	case pbtenant.CloudProvider_ali:
		for rId := range pbtenant.AliRegionId_name {
			if rId != int32(pbtenant.AliRegionId_ali_all) {
				region, _ := NewRegion(provider, rId)
				regions = append(regions, region)
			}
		}
	case pbtenant.CloudProvider_tencent:
		for rId := range pbtenant.TencentRegionId_name {
			if rId != int32(pbtenant.TencentRegionId_tc_all) {
				region, _ := NewRegion(provider, rId)
				regions = append(regions, region)
			}
		}
	case pbtenant.CloudProvider_huawei:
		for rId := range pbtenant.HuaweiRegionId_name {
			if rId != int32(pbtenant.HuaweiRegionId_hw_all) {
				region, _ := NewRegion(provider, rId)
				regions = append(regions, region)
			}
		}
	case pbtenant.CloudProvider_k8s:
		for rId := range pbtenant.K8SRegionId_name {
			region, _ := NewRegion(provider, rId)
			regions = append(regions, region)
		}
	}
	return
}

// GetAliRegionName prefix ali_
func GetAliRegionName(regionId int32) (string, error) {
	name, ok := pbtenant.AliRegionId_name[regionId]
	if !ok || regionId == int32(pbtenant.AliRegionId_ali_all) {
		return "", errors.WithMessagef(ErrNoExistAliRegionId, "input region id is %d", regionId)
	}
	region := strings.ReplaceAll(name, "_", "-")
	return region[4:], nil
}

// GetAliRegionId prefix ali_
func GetAliRegionId(regionName string) (int32, error) {
	regionName = "ali_" + strings.ReplaceAll(regionName, "-", "_")
	id, ok := pbtenant.AliRegionId_value[regionName]
	if !ok || regionName == "" {
		return 0, errors.WithMessagef(ErrNoExistAliRegionId, "input region id is %s", regionName)
	}
	return id, nil
}

// GetTencentRegionName prefix tencent
func GetTencentRegionName(regionId int32) (string, error) {
	name, ok := pbtenant.TencentRegionId_name[regionId]
	if !ok || regionId == int32(pbtenant.TencentRegionId_tc_all) {
		return "", errors.WithMessagef(ErrNoExistTencentRegionId, "input region id is %d", regionId)
	}
	region := strings.ReplaceAll(name, "_", "-")
	return region[3:], nil
}

// GetHuaweiRegionName prefix huawei
func GetHuaweiRegionName(regionId int32) (string, error) {
	name, ok := pbtenant.HuaweiRegionId_name[regionId]
	if !ok || regionId == int32(pbtenant.HuaweiRegionId_hw_all) {
		return "", errors.WithMessagef(ErrNoExistHuaweiRegionId, "input region id is %d", regionId)
	}
	region := strings.ReplaceAll(name, "_", "-")
	return region[3:], nil
}

// GetHuaweiRegionId prefix huawei
func GetHuaweiRegionId(regionName string) (int32, error) {
	regionName = "hw_" + strings.ReplaceAll(regionName, "-", "_")
	id, ok := pbtenant.AliRegionId_value[regionName]
	if !ok || regionName == "" {
		return 0, errors.WithMessagef(ErrNoExistAliRegionId, "input region id is %s", regionName)
	}
	return id, nil
}

// GetK8SRegionName prefix ali_
func GetK8SRegionName(regionId int32) (string, error) {
	name, ok := pbtenant.AliRegionId_name[regionId]
	if !ok || regionId == int32(pbtenant.AliRegionId_ali_all) {
		return "", errors.WithMessagef(ErrNoExistAliRegionId, "input region id is %d", regionId)
	}
	region := strings.ReplaceAll(name, "_", "-")
	return region[4:], nil
}

// GetK8SRegionId prefix ali_
func GetK8SRegionId(regionName string) (int32, error) {
	regionName = "ali_" + strings.ReplaceAll(regionName, "-", "_")
	id, ok := pbtenant.AliRegionId_value[regionName]
	if !ok || regionName == "" {
		return 0, errors.WithMessagef(ErrNoExistAliRegionId, "input region id is %s", regionName)
	}
	return id, nil
}

// GetAwsRegionName prefix aws_
func GetAwsRegionName(regionId int32) (string, error) {
	name, ok := pbtenant.AwsRegionId_name[regionId]
	if !ok || regionId == int32(pbtenant.AwsRegionId_aws_all) {
		return "", errors.WithMessagef(ErrNoExistAwsRegionId, "input region id is %d", regionId)
	}
	region := strings.ReplaceAll(name, "_", "-")
	return region[4:], nil
}
