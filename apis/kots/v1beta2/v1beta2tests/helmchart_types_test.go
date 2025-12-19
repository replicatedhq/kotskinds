package v1beta2tests

import (
	"testing"

	kotsv1beta2 "github.com/replicatedhq/kotskinds/apis/kots/v1beta2"
	kotsscheme "github.com/replicatedhq/kotskinds/client/kotsclientset/scheme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/scheme"
)

func Test_HelmChart(t *testing.T) {
	data := `apiVersion: kots.io/v1beta2
kind: HelmChart
metadata:
  name: test
spec:
  # chart identifies a matching chart from a .tgz
  chart:
    name: test
    chartVersion: 1.5.0-beta.2

  # values are used in the customer environment, as a pre-render step
  # these values will be supplied to helm template
  values:
    secretKey: '{{repl ConfigOption "secret_key"}}'
    components:
      worker:
        replicaCount: repl{{ ConfigOption "worker_replica_count"}}

    externalCacheDSN: repl{{ if ConfigOptionEquals "redis_location" "redis_location_external"}}{{repl ConfigOption "external_redis_dsn"}}{{repl end}}

    ingress:
      enabled: repl{{ ConfigOptionEquals "ingress_enabled" "ingress_enabled_yes"}}

  # builder values provide a way to render the chart with all images
  # and manifests. this is used in replicated to create airgap packages
  builder:
    ingress:
      enabled: true
      test: ~
      1: 4
      a: 4.0
      b: ["a", "b"]
`

	kotsscheme.AddToScheme(scheme.Scheme)

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, gvk, err := decode([]byte(data), nil, nil)
	require.NoError(t, err)

	assert.Equal(t, "kots.io", gvk.Group)
	assert.Equal(t, "v1beta2", gvk.Version)
	assert.Equal(t, "HelmChart", gvk.Kind)

	helmChart := obj.(*kotsv1beta2.HelmChart)

	assert.Equal(t, "test", helmChart.Spec.Chart.Name)
	assert.Equal(t, "1.5.0-beta.2", helmChart.Spec.Chart.ChartVersion)

}

func Test_HelmChart_Docs(t *testing.T) {
	data := `apiVersion: kots.io/v1beta2
kind: HelmChart
metadata:
  name: myapp
spec:
  chart:
    name: myapp
    chartVersion: 3.1.7
  releaseName: myapp-release
  docs:
    install: |
      # Installation Guide
      
      This Helm chart deploys MyApp with PostgreSQL support.
      
      ## Prerequisites
      - Kubernetes 1.24+
      - Persistent volume support
      - 4GB RAM minimum
      
      ## Installation
      
      The chart will automatically provision storage and create
      the necessary services.
    
    prerequisites: |
      # Requirements
      
      Before installing, ensure your cluster has:
      - A storage class with dynamic provisioning
      - LoadBalancer support (or NodePort for on-prem)
      - Network access to container registries
    
    troubleshooting: |
      # Common Issues
      
      **Pods stuck in Pending:**
      Check storage class availability with kubectl get sc
      
      **Database connection errors:**
      Verify the postgresql.host value in your configuration.
`

	kotsscheme.AddToScheme(scheme.Scheme)

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, gvk, err := decode([]byte(data), nil, nil)
	require.NoError(t, err)

	assert.Equal(t, "kots.io", gvk.Group)
	assert.Equal(t, "v1beta2", gvk.Version)
	assert.Equal(t, "HelmChart", gvk.Kind)

	helmChart := obj.(*kotsv1beta2.HelmChart)

	assert.Equal(t, "myapp", helmChart.Spec.Chart.Name)
	assert.Equal(t, "3.1.7", helmChart.Spec.Chart.ChartVersion)
	assert.Equal(t, "myapp-release", helmChart.Spec.ReleaseName)

	require.NotNil(t, helmChart.Spec.Docs)
	assert.Equal(t, 3, len(helmChart.Spec.Docs))

	installDoc, ok := helmChart.Spec.Docs["install"]
	require.True(t, ok)
	assert.Contains(t, installDoc, "Installation Guide")
	assert.Contains(t, installDoc, "Kubernetes 1.24+")

	prereqDoc, ok := helmChart.Spec.Docs["prerequisites"]
	require.True(t, ok)
	assert.Contains(t, prereqDoc, "Requirements")
	assert.Contains(t, prereqDoc, "storage class")

	troubleshootingDoc, ok := helmChart.Spec.Docs["troubleshooting"]
	require.True(t, ok)
	assert.Contains(t, troubleshootingDoc, "Common Issues")
	assert.Contains(t, troubleshootingDoc, "Pods stuck in Pending")
}

func Test_HelmChart_Docs_Optional(t *testing.T) {
	data := `apiVersion: kots.io/v1beta2
kind: HelmChart
metadata:
  name: myapp
spec:
  chart:
    name: myapp
    chartVersion: 3.1.7
  releaseName: myapp-release
`

	kotsscheme.AddToScheme(scheme.Scheme)

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, gvk, err := decode([]byte(data), nil, nil)
	require.NoError(t, err)

	assert.Equal(t, "kots.io", gvk.Group)
	assert.Equal(t, "v1beta2", gvk.Version)
	assert.Equal(t, "HelmChart", gvk.Kind)

	helmChart := obj.(*kotsv1beta2.HelmChart)

	// docs field should be optional (nil when not provided)
	assert.Nil(t, helmChart.Spec.Docs)
}
