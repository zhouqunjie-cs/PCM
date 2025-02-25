//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package ali

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/zhouqunjie-cs/PCM/adaptor/pod_adaptor/server/pod"
	"github.com/zhouqunjie-cs/PCM/common/tenanter"
	"github.com/zhouqunjie-cs/PCM/lan_trans/idl/pbpod"
	"github.com/zhouqunjie-cs/PCM/lan_trans/idl/pbtenant"
)

// DeleteContainerGroup invokes the eci.DeleteContainerGroup API synchronously
// api document: https://help.aliyun.com/api/eci/deletecontainergroup.html
func (client *Client) DeleteContainerGroup(request *DeleteContainerGroupRequest) (response *DeleteContainerGroupResponse, err error) {

	provider := pbtenant.CloudProvider(request.ProviderId)
	regionId, err := tenanter.GetAliRegionId(request.RegionId)
	podId := request.ContainerGroupId
	podName := request.ContainerGroupName

	requestPCM := &pbpod.DeletePodReq{
		Provider:    provider,
		AccountName: request.AccountName,
		PodId:       podId,
		PodName:     podName,
		Namespace:   request.Namespace,
		RegionId:    regionId,
	}

	resp, err := pod.DeletePod(nil, requestPCM)

	response = &DeleteContainerGroupResponse{
		BaseResponse: nil,
		RequestId:    resp.RequestId,
	}

	return response, err

}

// DeleteContainerGroupWithChan invokes the eci.DeleteContainerGroup API asynchronously
// api document: https://help.aliyun.com/api/eci/deletecontainergroup.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteContainerGroupWithChan(request *DeleteContainerGroupRequest) (<-chan *DeleteContainerGroupResponse, <-chan error) {
	responseChan := make(chan *DeleteContainerGroupResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteContainerGroup(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// DeleteContainerGroupWithCallback invokes the eci.DeleteContainerGroup API asynchronously
// api document: https://help.aliyun.com/api/eci/deletecontainergroup.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteContainerGroupWithCallback(request *DeleteContainerGroupRequest, callback func(response *DeleteContainerGroupResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteContainerGroupResponse
		var err error
		defer close(result)
		response, err = client.DeleteContainerGroup(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// DeleteContainerGroupRequest is the request struct for api DeleteContainerGroup
type DeleteContainerGroupRequest struct {
	*requests.RpcRequest
	/*********PCM param************/
	ProviderId         int32  `position:"Query" name:"ProviderId"`
	AccountName        string `position:"Query" name:"AccountName"`
	Namespace          string `position:"Query" name:"Namespace"`
	ContainerGroupName string `position:"Query" name:"ContainerGroupName"`
	/*********PCM param************/
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	RegionId             string           `position:"Query" name:"RegionId"`
	ContainerGroupId     string           `position:"Query" name:"ContainerGroupId"`
	ClientToken          string           `position:"Query" name:"ClientToken"`
	VkClientVersion      string           `position:"Query" name:"VkClientVersion"`
}

// DeleteContainerGroupResponse is the response struct for api DeleteContainerGroup
type DeleteContainerGroupResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteContainerGroupRequest creates a request to invoke DeleteContainerGroup API
func CreateDeleteContainerGroupRequest() (request *DeleteContainerGroupRequest) {
	request = &DeleteContainerGroupRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Eci", "2018-08-08", "DeleteContainerGroup", "eci", "openAPI")
	return
}

// CreateDeleteContainerGroupResponse creates a response to parse from DeleteContainerGroup response
func CreateDeleteContainerGroupResponse() (response *DeleteContainerGroupResponse) {
	response = &DeleteContainerGroupResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
