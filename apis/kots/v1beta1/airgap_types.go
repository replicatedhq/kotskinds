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

type Architecture string

const (
	ArchitectureAmd64 Architecture = "amd64"
	ArchitectureArm64 Architecture = "arm64"
)

// AirgapSpec defines the desired state of AirgapSpec
type AirgapSpec struct {
	AirgapReleaseMeta        `json:",inline"`
	ChannelID                string                    `json:"channelID,omitempty"`
	ChannelName              string                    `json:"channelName,omitempty"`
	Signature                []byte                    `json:"signature,omitempty"`
	AppSlug                  string                    `json:"appSlug,omitempty"`
	IsRequired               bool                      `json:"isRequired,omitempty"`
	RequiredReleases         []AirgapReleaseMeta       `json:"requiredReleases,omitempty"`
	SavedImages              []string                  `json:"savedImages,omitempty"`
	EmbeddedClusterArtifacts *EmbeddedClusterArtifacts `json:"embeddedClusterArtifacts,omitempty"`
	UncompressedSize         int64                     `json:"uncompressedSize,omitempty"`
	Architecture             Architecture              `json:"architecture,omitempty"`
	Format                   string                    `json:"format,omitempty"`
	ReplicatedChartNames     []string                  `json:"replicatedChartNames,omitempty"`
}

// AirgapStatus defines airgap release metadata
type AirgapReleaseMeta struct {
	VersionLabel           string `json:"versionLabel,omitempty"`
	ReleaseNotes           string `json:"releaseNotes,omitempty"`
	UpdateCursor           string `json:"updateCursor,omitempty"`
	EmbeddedClusterVersion string `json:"embeddedClusterVersion,omitempty"`
}

// EmbeddedClusterArtifacts maps embedded cluster artifacts to their path
type EmbeddedClusterArtifacts struct {
	Charts              string                  `json:"charts,omitempty"`
	ImagesAmd64         string                  `json:"imagesAmd64,omitempty"`
	BinaryAmd64         string                  `json:"binaryAmd64,omitempty"`
	Metadata            string                  `json:"metadata,omitempty"`
	Registry            EmbeddedClusterRegistry `json:"registry,omitempty"`
	AdditionalArtifacts map[string]string       `json:"additionalArtifacts,omitempty"`
}

// Total returns the total amount of embedded cluster artifacts contained in
// the airgap bundle. Sums up the amount of charts, images, binaries, metadata
// and additional artifacts.
func (e *EmbeddedClusterArtifacts) Total() int {
	total := 0
	if e == nil {
		return total
	}
	if e.Charts != "" {
		total++
	}
	if e.ImagesAmd64 != "" {
		total++
	}
	if e.BinaryAmd64 != "" {
		total++
	}
	if e.Metadata != "" {
		total++
	}
	total += len(e.AdditionalArtifacts)
	total += len(e.Registry.SavedImages)
	return total
}

// EmbeddedClusterRegistry holds a directory from where a images can be read and later
// pushed to the embedded cluster registry. Format inside the directory is the same as
// the registry storage format.
type EmbeddedClusterRegistry struct {
	Dir         string   `json:"dir,omitempty"`
	SavedImages []string `json:"savedImages,omitempty"`
}

// AirgapStatus defines the observed state of Airgap
type AirgapStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// Airgap is the Schema for the airgap API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Airgap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AirgapSpec   `json:"spec,omitempty"`
	Status AirgapStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AirgapList contains a list of Airgaps
type AirgapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Airgap `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Airgap{}, &AirgapList{})
}
