package ecser

import (
	"context"
	"fmt"
	"github.com/JCCE-nudt/PCM/common/tenanter"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbecs"
	"github.com/JCCE-nudt/PCM/lan_trans/idl/pbtenant"
	"github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
	harvClient "github.com/harvester/harvester/pkg/generated/clientset/versioned"
	"github.com/longhorn/longhorn-manager/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8s "k8s.io/client-go/kubernetes"
	kubirtv1 "kubevirt.io/client-go/api/v1"
	"strconv"
	"strings"
	"time"
)

const (
	prefix                   = "harvesterhci.io"
	vmAnnotationPVC          = prefix + "/volumeClaimTemplates"
	vmAnnotationNetworkIps   = "network.harvesterhci.io/ips"
	defaultCloudInitUserData = "#cloud-config\npackage_update: true\npackages:\n  - qemu-guest-agent\nruncmd:\n  - - systemctl\n    - enable\n    - '--now'\n    - qemu-guest-agent\n"
)

type Config struct {
	Host  string
	Token string
	Port  int
}

type HarVMer struct {
	k8sCli   *k8s.Clientset
	harvCli  *harvClient.Clientset
	region   tenanter.Region
	tenanter tenanter.Tenanter
}

func newHarvesterClient(tenant tenanter.Tenanter) (Ecser, error) {
	var (
		k8sclient       *k8s.Clientset
		harvesterClient *harvClient.Clientset
		err             error
	)
	switch t := tenant.(type) {
	case *tenanter.AccessKeyTenant:
		k8sclient, err = GetKubernetesClient(t.GetUrl(), t.GetToken())
		if err != nil {
			return nil, err
		}
		harvesterClient, err = GetHarvesterClient(t.GetUrl(), t.GetToken())
		if err != nil {
			return nil, err
		}
	default:

	}
	if err != nil {
		return nil, errors.Wrap(err, "init harvester client error")
	}
	return &HarVMer{
		k8sCli:   k8sclient,
		harvCli:  harvesterClient,
		region:   nil,
		tenanter: tenant,
	}, nil
}

func (h *HarVMer) CreateEcs(ctx context.Context, req *pbecs.CreateEcsReq) (resp *pbecs.CreateEcsResp, err error) {
	var (
		vmTemplate *kubirtv1.VirtualMachineInstanceTemplateSpec
		vmImage    *v1beta1.VirtualMachineImage
	)
	if req.ImageId != "" {
		vmImage, err = h.harvCli.HarvesterhciV1beta1().VirtualMachineImages(req.GetNamespace()).Get(context.TODO(), req.ImageId, k8smetav1.GetOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "get vm image error")
		}
	} else {
		return nil, errors.Wrap(err, "Image ID given does not exist!")
	}
	storageClassName := vmImage.Status.StorageClassName
	vmNameBase := req.InstanceName

	vmLabels := map[string]string{
		prefix + "/creator": "harvester",
	}
	vmiLabels := vmLabels
	_amount := req.Amount
	if _amount == 0 {
		return nil, fmt.Errorf("VM count provided is 0, no VM will be created")
	}

	repAmount := 0
	InstanceIds := make([]string, 0)
	for i := 1; i <= int(_amount); i++ {
		var (
			vmName         string
			secretRandomID string
		)
		randomID := util.RandomID()
		if _amount > 1 {
			vmName = vmNameBase + "-" + fmt.Sprint(i)
			secretRandomID = vmNameBase + "-" + randomID
		} else {
			vmName = vmNameBase
			secretRandomID = vmNameBase + "-" + randomID
		}

		vmiLabels[prefix+"/vmName"] = vmName
		vmiLabels[prefix+"/vmNamePrefix"] = vmNameBase
		diskRandomID := util.RandomID()
		pvcName := vmName + "-disk-0-" + diskRandomID
		pvcAnnotation := "[{\"metadata\":{\"name\":\"" + pvcName + "\",\"annotations\":{\"harvesterhci.io/imageId\":\"" + req.GetNamespace() + "/" + req.GetImageId() + "\"}},\"spec\":{\"accessModes\":[\"ReadWriteMany\"],\"resources\":{\"requests\":{\"storage\":\"" + req.GetDiskSize() + "\"}},\"volumeMode\":\"Block\",\"storageClassName\":\"" + storageClassName + "\"}}]"

		vmTemplate, err = buildVMTemplate(int(req.GetCpu()), req.GetMemory(), req.GetSshKey(), h.harvCli, pvcName, vmiLabels, vmNameBase, secretRandomID)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		vm := &kubirtv1.VirtualMachine{
			ObjectMeta: k8smetav1.ObjectMeta{
				Name:      vmName,
				Namespace: req.GetNamespace(),
				Annotations: map[string]string{
					vmAnnotationPVC:        pvcAnnotation,
					vmAnnotationNetworkIps: "[]",
				},
				Labels: vmLabels,
			},
			Spec: kubirtv1.VirtualMachineSpec{
				Running:  NewTrue(),
				Template: vmTemplate,
			},
		}
		resp, err1 := h.harvCli.KubevirtV1().VirtualMachines(req.GetNamespace()).Create(context.TODO(), vm, k8smetav1.CreateOptions{})
		if err1 != nil {
			return nil, errors.Wrap(err, "VM create failed")
		}

		var sshKey *v1beta1.KeyPair
		cloudInitSSHSection := ""
		if req.GetSshKey() != "" {
			sshArr := strings.Split(req.GetSshKey(), "/")
			if len(sshArr) != 2 {
				return nil, errors.New("sshKeyName should be in format namespace/name")
			}
			sshKey, err = h.harvCli.HarvesterhciV1beta1().KeyPairs(sshArr[0]).Get(context.TODO(), sshArr[1], k8smetav1.GetOptions{})
			if err != nil {
				return nil, errors.Wrap(err, "error during getting keypair from Harvester")
			}
			cloudInitSSHSection = "\nssh_authorized_keys:\n  - 	>-\n" + sshKey.Spec.PublicKey + "\n"
			logrus.Debugf("SSH Key Name %s given does exist!", req.GetSshKey())
		}
		if req.UserDataTemplate == "" {
			req.UserDataTemplate = defaultCloudInitUserData
		}
		// Create the secret for the VM
		if _, secreterr := createCloudInitDataFromSecret(h.k8sCli, vmName, resp.ObjectMeta.UID, secretRandomID, req.Namespace, req.UserDataTemplate+cloudInitSSHSection, req.NetworkDataTemplate); secreterr != nil {
			logrus.Errorf("Create secret failed, %s", secreterr)
			return nil, errors.Wrap(secreterr, "Create cloud init data from secret failed")
		}
		InstanceIds = append(InstanceIds, string(resp.UID))
		repAmount++
	}
	isFinished := false
	if int32(repAmount) == req.Amount {
		isFinished = true
	}
	return &pbecs.CreateEcsResp{
		Provider:       pbtenant.CloudProvider_harvester,
		AccountName:    h.tenanter.AccountName(),
		InstanceIdSets: InstanceIds,
		Finished:       isFinished,
	}, nil
}

//buildVMTemplate creates a *kubirtv1.VirtualMachineInstanceTemplateSpec from the CLI Flags and some computed values
func buildVMTemplate(vCpu int, memory, sshKeyName string, c *harvClient.Clientset,
	pvcName string, vmiLabels map[string]string, vmName string, secretName string) (vmTemplate *kubirtv1.VirtualMachineInstanceTemplateSpec, err error) {
	vmTemplate = nil
	_memory := resource.MustParse(memory)
	if sshKeyName != "" {
		sshArr := strings.Split(sshKeyName, "/")
		if len(sshArr) != 2 {
			return nil, errors.New("sshKeyName should be in format namespace/name")
		}
		_, keyerr := c.HarvesterhciV1beta1().KeyPairs(sshArr[0]).Get(context.TODO(), sshArr[1], k8smetav1.GetOptions{})
		if keyerr != nil {
			return nil, errors.Wrap(keyerr, "error during getting keypair from Harvester")
		}
		logrus.Debugf("SSH Key Name %s given does exist!", sshKeyName)
	}
	logrus.Debug("CloudInit: ")
	vmTemplate = &kubirtv1.VirtualMachineInstanceTemplateSpec{
		ObjectMeta: k8smetav1.ObjectMeta{
			Annotations: vmiAnnotations(pvcName, sshKeyName),
			Labels:      vmiLabels,
		},
		Spec: kubirtv1.VirtualMachineInstanceSpec{
			Hostname: vmName,
			Networks: []kubirtv1.Network{
				{
					Name: "default",
					NetworkSource: kubirtv1.NetworkSource{
						Multus: &kubirtv1.MultusNetwork{
							NetworkName: "default/service-network",
						},
					},
				},
			},
			Volumes: []kubirtv1.Volume{
				{
					Name: "disk-0",
					VolumeSource: kubirtv1.VolumeSource{
						PersistentVolumeClaim: &kubirtv1.PersistentVolumeClaimVolumeSource{
							PersistentVolumeClaimVolumeSource: v1.PersistentVolumeClaimVolumeSource{
								ClaimName: pvcName,
							},
						},
					},
				},
				{
					Name: "cloudinitdisk",
					VolumeSource: kubirtv1.VolumeSource{
						CloudInitNoCloud: &kubirtv1.CloudInitNoCloudSource{
							UserDataSecretRef:    &v1.LocalObjectReference{Name: secretName},
							NetworkDataSecretRef: &v1.LocalObjectReference{Name: secretName},
						},
					},
				},
			},
			Domain: kubirtv1.DomainSpec{
				CPU: &kubirtv1.CPU{
					Cores:   uint32(vCpu),
					Sockets: uint32(1),
					Threads: uint32(1),
				},
				Memory: &kubirtv1.Memory{
					Guest: &_memory,
				},
				Devices: kubirtv1.Devices{
					Inputs: []kubirtv1.Input{
						{
							Bus:  "usb",
							Type: "tablet",
							Name: "tablet",
						},
					},
					Interfaces: []kubirtv1.Interface{
						{
							Name:                   "default",
							Model:                  "virtio",
							InterfaceBindingMethod: kubirtv1.DefaultBridgeNetworkInterface().InterfaceBindingMethod,
						},
					},
					Disks: []kubirtv1.Disk{
						{
							BootOrder: PointerToUint(1),
							Name:      "disk-0",
							DiskDevice: kubirtv1.DiskDevice{
								Disk: &kubirtv1.DiskTarget{
									Bus: "virtio",
								},
							},
						},
						{
							Name: "cloudinitdisk",
							DiskDevice: kubirtv1.DiskDevice{
								Disk: &kubirtv1.DiskTarget{
									Bus: "virtio",
								},
							},
						},
					},
				},
				Resources: kubirtv1.ResourceRequirements{
					Limits: v1.ResourceList{
						"cpu":    resource.MustParse(strconv.Itoa(vCpu)),
						"memory": resource.MustParse(memory),
					},
					Requests: v1.ResourceList{
						"memory": resource.MustParse(memory),
					},
				},
			},
			Affinity: &v1.Affinity{
				PodAntiAffinity: &v1.PodAntiAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
						{
							Weight: int32(1),
							PodAffinityTerm: v1.PodAffinityTerm{
								TopologyKey: "kubernetes.io/hostname",
								LabelSelector: &k8smetav1.LabelSelector{
									MatchLabels: map[string]string{
										prefix + "/vmNamePrefix": vmName,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return
}

// vmiAnnotations generates a map of strings to be injected as annotations from a PVC name and an SSK Keyname
func vmiAnnotations(pvcName string, sshKeyName string) map[string]string {
	sshKey := "[]"
	if sshKeyName != "" {
		sshKey = "[\"" + sshKeyName + "\"]"
	}
	return map[string]string{
		prefix + "/diskNames": "[\"" + pvcName + "\"]",
		prefix + "/sshNames":  sshKey,
	}
}

//CreateCloudInitDataFromSecret creates a cloud-init configmap from a secret
func createCloudInitDataFromSecret(c *k8s.Clientset, vmName string, uid types.UID, secretName, namespace, userData, networkData string) (secret *v1.Secret, err error) {
	toCreate := &v1.Secret{
		TypeMeta: k8smetav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Labels: map[string]string{
				prefix + "/cloud-init-template": "harvester",
			},
			OwnerReferences: []k8smetav1.OwnerReference{
				{
					APIVersion: "kubevirt.io/v1",
					Kind:       "VirtualMachine",
					Name:       vmName,
					UID:        uid,
				},
			},
		},
		Type: "secret",
		Data: map[string][]byte{
			"userdata":    []byte(userData),
			"networkdata": []byte(networkData),
		},
	}
	resp, err := c.CoreV1().Secrets(namespace).Create(context.TODO(), toCreate, k8smetav1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error during getting cloud-init secret")
	}
	return resp, nil
}

func (h *HarVMer) DeleteEcs(ctx context.Context, req *pbecs.DeleteEcsReq) (resp *pbecs.DeleteEcsResp, err error) {
	if req.Namespace == "" {
		return nil, errors.New("namespace is required")
	}
	vm, err := h.harvCli.KubevirtV1().VirtualMachines(req.GetNamespace()).Get(context.TODO(), req.GetInstanceName(), k8smetav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "VM does not exist")
	}
	err = h.harvCli.KubevirtV1().VirtualMachines(req.GetNamespace()).Delete(context.TODO(), req.GetInstanceName(), k8smetav1.DeleteOptions{})
	if err != nil {
		logrus.Errorf("delete vm error: %v", err)
		return nil, errors.Wrap(err, "VM  could not be deleted successfully: %w")
	}
	//delete vm disk
	if req.DiskName != "" {
		for _, delName := range strings.Split(req.DiskName, ",") {
			for _, disk := range vm.Spec.Template.Spec.Volumes {
				if disk.Name == delName {
					ClaimName := disk.VolumeSource.PersistentVolumeClaim.ClaimName
					err1 := h.k8sCli.CoreV1().PersistentVolumeClaims(req.GetNamespace()).Delete(context.TODO(), ClaimName, k8smetav1.DeleteOptions{})
					if err1 != nil {
						logrus.Errorf("delete pvc failed,err:%v", err1)
						return nil, errors.Wrap(err, "VM disk not be deleted successfully")
					}
				}
			}
		}
	}
	return &pbecs.DeleteEcsResp{
		Provider:    pbtenant.CloudProvider_harvester,
		AccountName: h.tenanter.AccountName(),
	}, nil
}

func (h *HarVMer) UpdateEcs(ctx context.Context, req *pbecs.UpdateEcsReq) (resp *pbecs.UpdateEcsResp, err error) {
	//查询删除的vm
	vm, err := h.harvCli.KubevirtV1().VirtualMachines(req.GetNamespace()).Get(context.TODO(), req.GetInstanceName(), k8smetav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "VM does not exist")
	}
	if req.Cpu != "" && req.Memory != "" {
		vm.Spec.Template.Spec.Domain.Resources = kubirtv1.ResourceRequirements{
			Limits: v1.ResourceList{
				"cpu":    resource.MustParse(req.Cpu),
				"memory": resource.MustParse(req.Memory),
			},
		}
	}
	vm.ObjectMeta.Annotations["field.cattle.io/description"] = req.Description
	if req.Cpu != "" {
		j, err := strconv.ParseUint(req.Cpu, 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "cpu is not a number")
		}
		vm.Spec.Template.Spec.Domain.CPU = &kubirtv1.CPU{
			Cores:   uint32(j),
			Sockets: uint32(1),
			Threads: uint32(1),
		}
	}
	if req.Memory != "" {
		_memory := resource.MustParse(req.Memory)
		vm.Spec.Template.Spec.Domain.Memory = &kubirtv1.Memory{
			Guest: &_memory,
		}
	}
	if err != nil {
		return nil, errors.Wrap(err, "Harvester client connection failed")
	}
	//update
	_, err = h.harvCli.KubevirtV1().VirtualMachines(req.GetNamespace()).Update(context.TODO(), vm, k8smetav1.UpdateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "VM update failed")
	}
	if req.IsRestart {
		//重启
		err = restartVmByName(h.harvCli, req.GetNamespace(), req.GetInstanceName())
		if err != nil {
			return nil, errors.Wrap(err, "VM restart failed")
		}
	}
	return &pbecs.UpdateEcsResp{
		Provider:    pbtenant.CloudProvider_harvester,
		AccountName: h.tenanter.AccountName(),
	}, nil
}

func (h *HarVMer) ListDetail(ctx context.Context, req *pbecs.ListDetailReq) (resp *pbecs.ListDetailResp, err error) {
	vmList, err := h.harvCli.KubevirtV1().VirtualMachines(req.GetNamespace()).List(context.TODO(), k8smetav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "VM list failed")
	}
	vmiList, err := h.harvCli.KubevirtV1().VirtualMachineInstances(req.GetNamespace()).List(context.TODO(), k8smetav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "VMI list failed")
	}
	vmiMap := map[string]kubirtv1.VirtualMachineInstance{}
	for _, vmi := range vmiList.Items {
		vmiMap[vmi.Name] = vmi
	}
	var ecses = make([]*pbecs.EcsInstance, len(vmList.Items))
	for k, vm := range vmList.Items {
		running := *vm.Spec.Running
		var state string
		if running {
			state = "Running"
		} else {
			state = "Not Running"
		}
		IP := make([]string, 0)
		if vmiMap[vm.Name].Status.Interfaces == nil {
			IP = append(IP, "")
		} else {
			IP = append(IP, vmiMap[vm.Name].Status.Interfaces[0].IP)
		}
		ecses[k] = &pbecs.EcsInstance{
			Provider:     pbtenant.CloudProvider_harvester,
			AccountName:  h.tenanter.AccountName(),
			Status:       state,
			InstanceName: vm.Name,
			Node:         vmiMap[vm.Name].Status.NodeName,
			Cpu:          vm.Spec.Template.Spec.Domain.Resources.Limits.Cpu().String(),
			Memory:       vm.Spec.Template.Spec.Domain.Resources.Limits.Memory().String(),
			PublicIps:    IP,
			CreationTime: vm.CreationTimestamp.String(),
			Description:  vm.ObjectMeta.Annotations["field.cattle.io/description"],
			Namespace:    vm.Namespace,
		}
	}
	isFinished := false
	if len(ecses) > 0 {
		isFinished = true
	}
	return &pbecs.ListDetailResp{
		Ecses:    ecses,
		Finished: isFinished,
	}, nil
}

func (h *HarVMer) ActionEcs(ctx context.Context, req *pbecs.ActionReq) (resp *pbecs.ActionResp, err error) {
	status := ""
	switch req.GetActionType() {
	case pbecs.ActionType_start:
		err := startVmByName(h.harvCli, req.GetNamespace(), req.GetVmName())
		if err != nil {
			return nil, err
		}
		status = "Running"
	case pbecs.ActionType_stop:
		err := stopVmByName(h.harvCli, req.GetNamespace(), req.GetVmName())
		if err != nil {
			return nil, err
		}
		status = "Off"
	case pbecs.ActionType_restart:
		err := restartVmByName(h.harvCli, req.GetNamespace(), req.GetVmName())
		if err != nil {
			return nil, err
		}
		status = "Running"
	}
	return &pbecs.ActionResp{
		Provider:    pbtenant.CloudProvider_harvester,
		AccountName: h.tenanter.AccountName(),
		Status:      status,
	}, nil
}

//startVmByName starts a VM by first issuing a GET using the VM name, then updating the resulting VM object
func startVmByName(c *harvClient.Clientset, namespace, vmName string) error {
	vm, err := c.KubevirtV1().VirtualMachines(namespace).Get(context.TODO(), vmName, k8smetav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "VM not found")
	}
	*vm.Spec.Running = true
	_, err = c.KubevirtV1().VirtualMachines(namespace).Update(context.TODO(), vm, k8smetav1.UpdateOptions{})
	if err != nil {
		return errors.Wrap(err, "VM start failed")
	}
	return nil
}

//stopVmByName will stop a VM by first finding it by its name and then call stopBMbyRef function
func stopVmByName(c *harvClient.Clientset, namespace, vmName string) error {
	vm, err := c.KubevirtV1().VirtualMachines(namespace).Get(context.TODO(), vmName, k8smetav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "VM not found")
	}
	*vm.Spec.Running = false
	_, err = c.KubevirtV1().VirtualMachines(namespace).Update(context.TODO(), vm, k8smetav1.UpdateOptions{})
	if err != nil {
		return errors.Wrap(err, "VM stop failed")
	}
	return nil
}

//restartVMbyName will restart a VM by first finding it by its name and then call restartVMbyRef function
func restartVmByName(c *harvClient.Clientset, namespace, vmName string) error {
	vm, err := c.KubevirtV1().VirtualMachines(namespace).Get(context.TODO(), vmName, k8smetav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "VM not found")
	}
	err = stopVmByName(c, namespace, vm.Name)
	if err != nil {
		return errors.Wrap(err, "VM stop failed")
	}
	select {
	case <-time.Tick(1 * time.Second):
		return startVmByName(c, namespace, vm.Name)
	}
}
