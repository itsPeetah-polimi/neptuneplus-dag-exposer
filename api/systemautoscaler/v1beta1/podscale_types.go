/*
Copyright 2025.

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

package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PodScaleSpec defines the desired state of PodScale
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodScaleList is a list of PodScale resources
type PodScaleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []PodScale `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodScale defines the mapping between a `ServiceLevelAgreement` and a
// `Pod` matching the selector. It also keeps track of the resource values
// computed by `Recommender` and adjusted by `Contention Manager`.
type PodScale struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodScaleSpec   `json:"spec"`
	Status PodScaleStatus `json:"status"`
}

// PodScaleSpec is the spec for a PodScale resource
type PodScaleSpec struct {
	Namespace        string          `json:"namespace"`
	SLA              string          `json:"serviceLevelAgreement"`
	Pod              string          `json:"pod"`
	Service          string          `json:"service"`
	Container        string          `json:"container"`
	DesiredResources v1.ResourceList `json:"desired,omitempty" protobuf:"bytes,3,rep,name=desired,casttype=ResourceList,castkey=ResourceName"`
}

// PodScaleStatus contains the resources patched by the
// `Contention Manager` according to the available node resources
// and other pods' SLA
type PodScaleStatus struct {
	CappedResources v1.ResourceList `json:"capped,omitempty" protobuf:"bytes,3,rep,name=actual,casttype=ResourceList,castkey=ResourceName"`
	ActualResources v1.ResourceList `json:"actual,omitempty" protobuf:"bytes,3,rep,name=actual,casttype=ResourceList,castkey=ResourceName"`
}

func init() {
	SchemeBuilder.Register(&PodScale{}, &PodScaleList{})
}
