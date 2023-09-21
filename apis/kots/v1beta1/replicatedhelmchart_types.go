/*
Copyright 2019 Replicated, Inc..

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
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReplicatedHelmChartSpec defines the desired state of ReplicatedHelmChartSpec
type ReplicatedHelmChartSpec struct {
	ChartNames   []string              `json:"chartNames,omitempty"`
	ChartValues  *apiextensionsv1.JSON `json:"chartValues,omitempty"`
	GlobalValues *apiextensionsv1.JSON `json:"globalValues,omitempty"`
}

// ReplicatedHelmChartStatus defines the observed state of ReplicatedHelmChart
type ReplicatedHelmChartStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ReplicatedHelmChart is the Schema for the installation API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type ReplicatedHelmChart struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReplicatedHelmChartSpec   `json:"spec,omitempty"`
	Status ReplicatedHelmChartStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReplicatedHelmChartList contains a list of ReplicatedHelmCharts
type ReplicatedHelmChartList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReplicatedHelmChart `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ReplicatedHelmChart{}, &ReplicatedHelmChartList{})
}
