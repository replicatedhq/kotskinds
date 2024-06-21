package v1beta1tests

import (
	"testing"

	"github.com/replicatedhq/kotskinds/apis/kots/v1beta1"
	"github.com/stretchr/testify/require"
)

func TestEmbeddedClusterArtifactsTotal(t *testing.T) {
	for _, tt := range []struct {
		name   string
		airgap *v1beta1.Airgap
		want   int
	}{
		{
			name:   "undefined embedded cluster artifacts",
			airgap: &v1beta1.Airgap{},
			want:   0,
		},
		{
			name: "no images and not additional",
			airgap: &v1beta1.Airgap{
				Spec: v1beta1.AirgapSpec{
					EmbeddedClusterArtifacts: &v1beta1.EmbeddedClusterArtifacts{
						Charts:      "charts",
						ImagesAmd64: "images",
						BinaryAmd64: "binary",
						Metadata:    "metadata",
					},
				},
			},
			want: 4,
		},
		{
			name: "only additional artifacts",
			airgap: &v1beta1.Airgap{
				Spec: v1beta1.AirgapSpec{
					EmbeddedClusterArtifacts: &v1beta1.EmbeddedClusterArtifacts{
						AdditionalArtifacts: map[string]string{
							"foo": "bar",
							"baz": "qux",
						},
					},
				},
			},
			want: 2,
		},
		{
			name: "only embedded cluster images",
			airgap: &v1beta1.Airgap{
				Spec: v1beta1.AirgapSpec{
					EmbeddedClusterArtifacts: &v1beta1.EmbeddedClusterArtifacts{
						Registry: v1beta1.EmbeddedClusterRegistry{
							SavedImages: []string{"foo", "bar", "baz"},
						},
					},
				},
			},
			want: 3,
		},
		{
			name: "various embedded cluster artifacts",
			airgap: &v1beta1.Airgap{
				Spec: v1beta1.AirgapSpec{
					EmbeddedClusterArtifacts: &v1beta1.EmbeddedClusterArtifacts{
						Charts:      "charts",
						ImagesAmd64: "images",
						BinaryAmd64: "binary",
						Metadata:    "metadata",
						AdditionalArtifacts: map[string]string{
							"foo": "bar",
							"baz": "qux",
						},
						Registry: v1beta1.EmbeddedClusterRegistry{
							SavedImages: []string{"foo", "bar", "baz"},
						},
					},
				},
			},
			want: 9,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.airgap.Spec.EmbeddedClusterArtifacts.Total(), tt.want)
		})
	}
}
