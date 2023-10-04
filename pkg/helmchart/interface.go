package helmchart

import (
	kotsv1beta1 "github.com/replicatedhq/kotskinds/apis/kots/v1beta1"
	kotsv1beta2 "github.com/replicatedhq/kotskinds/apis/kots/v1beta2"
)

// HelmChartInterface represents any kots.io HelmChart (v1beta1 or v1beta2)
type HelmChartInterface interface {
	GetAPIVersion() string
	GetChartName() string
	GetChartVersion() string
	GetReleaseName() string
	GetDirName() string
	GetNamespace() string
	GetUpgradeFlags() []string
	GetWeight() int64
	GetHelmVersion() string
	GetBuilderValues() (map[string]interface{}, error)
	SetChartNamespace(namespace string)
}

// v1beta1 and v1beta2 HelmChart structs must implement HelmChartInterface
var _ HelmChartInterface = (*kotsv1beta1.HelmChart)(nil)
var _ HelmChartInterface = (*kotsv1beta2.HelmChart)(nil)
