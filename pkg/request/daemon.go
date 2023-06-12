/*
 Copyright 2021 The Hybridnet Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package request

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/parnurzeal/gorequest"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	networkingv1 "github.com/alibaba/hybridnet/pkg/apis/networking/v1"
)

// CniDaemonClient is the client to visit cnidaemon
type CniDaemonClient struct {
	*gorequest.SuperAgent
}

// PodRequest is the cnidaemon request format
type PodRequest struct {
	PodName      string `json:"pod_name"`
	PodNamespace string `json:"pod_namespace"`
	ContainerID  string `json:"container_id"`
	NetNs        string `json:"net_ns"`
}

type IPAddress struct {
	// ip with mask
	IP string `json:"ip"`

	Mac      string                 `json:"mac"`
	Gateway  string                 `json:"gateway"`
	Protocol networkingv1.IPVersion `json:"protocol"`
}

// PodResponse is the cnidaemon response format
type PodResponse struct {
	IPAddress     []IPAddress `json:"address"`
	HostInterface string      `json:"host_interface"`
	Err           string      `json:"error"`
}

// PodIPAMRequest is the formatted request body for IPAM
type PodIPAMRequest struct {
	PodName       string `json:"pod_name"`
	PodNamespace  string `json:"pod_namespace"`
	InterfaceName string `json:"interface_name"`
	ContainerID   string `json:"container_id"`
}

// PodIPAMResponse is the formatted response body for IPAM
type PodIPAMResponse struct {
	Addresses []IPAddress `json:"addresses"`
	Err       string      `json:"error"`
}

// NewCniDaemonClient return a new cnidaemonclient
func NewCniDaemonClient(socketAddress string) CniDaemonClient {
	request := gorequest.New()
	request.Transport = &http.Transport{DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
		return net.Dial("unix", socketAddress)
	}}
	return CniDaemonClient{request}
}

// Add pod request
func (cdc CniDaemonClient) Add(podRequest PodRequest) (*PodResponse, error) {
	resp := PodResponse{}
	res, _, errors := cdc.Post("http://dummy/api/v1/add").Send(podRequest).EndStruct(&resp)
	if len(errors) != 0 {
		return nil, errors[0]
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("request ip return %d %s", res.StatusCode, resp.Err)
	}
	return &resp, nil
}

// Del pod request
func (cdc CniDaemonClient) Del(podRequest PodRequest) error {
	res, body, errors := cdc.Post("http://dummy/api/v1/del").Send(podRequest).End()
	if len(errors) != 0 {
		return errors[0]
	}
	if res.StatusCode != 204 {
		return fmt.Errorf("delete ip return %d %s", res.StatusCode, body)
	}
	return nil
}

func (cdc CniDaemonClient) IPAMAdd(request PodIPAMRequest) (*PodIPAMResponse, error) {
	resp := PodIPAMResponse{}

	res, _, errors := cdc.Post("http://dummy/api/v1/ipam/add").Send(request).EndStruct(&resp)
	if len(errors) != 0 {
		return nil, utilerrors.NewAggregate(errors)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response %d, %s", res.StatusCode, resp.Err)
	}

	return &resp, nil
}

func (cdc CniDaemonClient) IPAMDel(request PodIPAMRequest) error {
	res, body, errors := cdc.Post("http://dummy/api/v1/ipam/del").Send(request).End()
	if len(errors) != 0 {
		return utilerrors.NewAggregate(errors)
	}

	if res.StatusCode != 204 {
		return fmt.Errorf("unexpected response %d, %s", res.StatusCode, body)
	}

	return nil
}
