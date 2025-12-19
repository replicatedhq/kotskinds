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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TerraformSpec defines the desired state of Terraform
type TerraformSpec struct {
	// Filename references a .tgz file containing Terraform modules
	Filename string `json:"filename"`

	// MinTerraformVersion specifies the minimum Terraform version required
	MinTerraformVersion string `json:"minTerraformVersion,omitempty"`

	// Docs provides documentation for various user-defined steps
	Docs map[string]string `json:"docs,omitempty"`
}

// TerraformStatus defines the observed state of Terraform
type TerraformStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// Terraform is the Schema for the terraform API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Terraform struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TerraformSpec   `json:"spec,omitempty"`
	Status TerraformStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TerraformList contains a list of Terraform resources
type TerraformList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Terraform `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Terraform{}, &TerraformList{})
}
