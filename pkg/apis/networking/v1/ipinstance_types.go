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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IPInstanceSpec defines the desired state of IPInstance
type IPInstanceSpec struct {
	// +kubebuilder:validation:Required
	Network string `json:"network"`
	// +kubebuilder:validation:Required
	Subnet string `json:"subnet"`
	// +kubebuilder:validation:Required
	Address Address `json:"address"`
}

// IPInstanceStatus defines the observed state of IPInstance
type IPInstanceStatus struct {
	// +kubebuilder:validation:Optional
	NodeName string `json:"nodeName"`
	// +kubebuilder:validation:Optional
	Phase IPPhase `json:"phase"`
	// +kubebuilder:validation:Optional
	PodName string `json:"podName"`
	// +kubebuilder:validation:Optional
	PodNamespace string `json:"podNamespace"`
	// +kubebuilder:validation:Optional
	SandboxID string `json:"sandboxID"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="IP",type=string,JSONPath=`.spec.address.ip`
// +kubebuilder:printcolumn:name="Gateway",type=string,JSONPath=`.spec.address.gateway`
// +kubebuilder:printcolumn:name="PodName",type=string,JSONPath=`.status.podName`
// +kubebuilder:printcolumn:name="Node",type=string,JSONPath=`.status.nodeName`
// +kubebuilder:printcolumn:name="Subnet",type=string,JSONPath=`.spec.subnet`
// +kubebuilder:printcolumn:name="Network",type=string,JSONPath=`.spec.network`

// IPInstance is the Schema for the ipinstances API
type IPInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IPInstanceSpec   `json:"spec,omitempty"`
	Status IPInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IPInstanceList contains a list of IPInstance
type IPInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IPInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IPInstance{}, &IPInstanceList{})
}
