package cr_ee

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

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ListRepoBuildRecord invokes the cr.ListRepoBuildRecord API synchronously
// api document: https://help.aliyun.com/api/cr/listrepobuildrecord.html
func (client *Client) ListRepoBuildRecord(request *ListRepoBuildRecordRequest) (response *ListRepoBuildRecordResponse, err error) {
	response = CreateListRepoBuildRecordResponse()
	err = client.DoAction(request, response)
	return
}

// ListRepoBuildRecordWithChan invokes the cr.ListRepoBuildRecord API asynchronously
// api document: https://help.aliyun.com/api/cr/listrepobuildrecord.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListRepoBuildRecordWithChan(request *ListRepoBuildRecordRequest) (<-chan *ListRepoBuildRecordResponse, <-chan error) {
	responseChan := make(chan *ListRepoBuildRecordResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListRepoBuildRecord(request)
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

// ListRepoBuildRecordWithCallback invokes the cr.ListRepoBuildRecord API asynchronously
// api document: https://help.aliyun.com/api/cr/listrepobuildrecord.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListRepoBuildRecordWithCallback(request *ListRepoBuildRecordRequest, callback func(response *ListRepoBuildRecordResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListRepoBuildRecordResponse
		var err error
		defer close(result)
		response, err = client.ListRepoBuildRecord(request)
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

// ListRepoBuildRecordRequest is the request struct for api ListRepoBuildRecord
type ListRepoBuildRecordRequest struct {
	*requests.RpcRequest
	RepoId     string           `position:"Query" name:"RepoId"`
	InstanceId string           `position:"Query" name:"InstanceId"`
	PageNo     requests.Integer `position:"Query" name:"PageNo"`
	PageSize   requests.Integer `position:"Query" name:"PageSize"`
}

// ListRepoBuildRecordResponse is the response struct for api ListRepoBuildRecord
type ListRepoBuildRecordResponse struct {
	*responses.BaseResponse
	ListRepoBuildRecordIsSuccess bool               `json:"IsSuccess" xml:"IsSuccess"`
	Code                         string             `json:"Code" xml:"Code"`
	RequestId                    string             `json:"RequestId" xml:"RequestId"`
	PageNo                       int                `json:"PageNo" xml:"PageNo"`
	PageSize                     int                `json:"PageSize" xml:"PageSize"`
	TotalCount                   string             `json:"TotalCount" xml:"TotalCount"`
	BuildRecords                 []BuildRecordsItem `json:"BuildRecords" xml:"BuildRecords"`
}

// CreateListRepoBuildRecordRequest creates a request to invoke ListRepoBuildRecord API
func CreateListRepoBuildRecordRequest() (request *ListRepoBuildRecordRequest) {
	request = &ListRepoBuildRecordRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("cr", "2018-12-01", "ListRepoBuildRecord", "acr", "openAPI")
	request.Method = requests.POST
	return
}

// CreateListRepoBuildRecordResponse creates a response to parse from ListRepoBuildRecord response
func CreateListRepoBuildRecordResponse() (response *ListRepoBuildRecordResponse) {
	response = &ListRepoBuildRecordResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}