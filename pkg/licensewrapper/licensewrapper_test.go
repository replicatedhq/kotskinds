package licensewrapper

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	v1beta1LicenseYAML = `apiVersion: kots.io/v1beta1
kind: License
metadata:
  name: test-license
spec:
  licenseID: test-license-id
  appSlug: test-app-slug
  channelID: test-channel-id
  channelName: Stable
  customerName: Test Customer
  customerEmail: test@example.com
  endpoint: https://replicated.app
  licenseType: trial
  licenseSequence: 2
  isAirgapSupported: true
  isGitOpsSupported: true
  isIdentityServiceSupported: true
  isGeoaxisSupported: true
  isSnapshotSupported: true
  isSupportBundleUploadSupported: true
  isSemverRequired: true
`

	v1beta2LicenseYAML = `apiVersion: kots.io/v1beta2
kind: License
metadata:
  name: test-license
spec:
  licenseID: test-license-id
  appSlug: test-app-slug
  channelID: test-channel-id
  channelName: Stable
  customerName: Test Customer
  customerEmail: test@example.com
  endpoint: https://replicated.app
  licenseType: trial
  licenseSequence: 2
  isAirgapSupported: true
  isGitOpsSupported: true
  isIdentityServiceSupported: true
  isGeoaxisSupported: true
  isSnapshotSupported: true
  isSupportBundleUploadSupported: true
  isSemverRequired: true
  signature: dGVzdC1zaWduYXR1cmU=
`
)

func TestLoadLicenseFromBytes_V1Beta1(t *testing.T) {
	wrapper, err := LoadLicenseFromBytes([]byte(v1beta1LicenseYAML))
	require.NoError(t, err)

	// Version detection
	assert.True(t, wrapper.IsV1())
	assert.False(t, wrapper.IsV2())

	// Basic properties
	assert.Equal(t, "test-license-id", wrapper.GetLicenseID())
	assert.Equal(t, "test-app-slug", wrapper.GetAppSlug())
	assert.Equal(t, "test-channel-id", wrapper.GetChannelID())
	assert.Equal(t, "Stable", wrapper.GetChannelName())
	assert.Equal(t, "Test Customer", wrapper.GetCustomerName())
	assert.Equal(t, "test@example.com", wrapper.GetCustomerEmail())
	assert.Equal(t, "https://replicated.app", wrapper.GetEndpoint())
	assert.Equal(t, "trial", wrapper.GetLicenseType())
	assert.Equal(t, int64(2), wrapper.GetLicenseSequence())

	// Feature flags
	assert.True(t, wrapper.IsAirgapSupported())
	assert.True(t, wrapper.IsGitOpsSupported())
	assert.True(t, wrapper.IsIdentityServiceSupported())
	assert.True(t, wrapper.IsGeoaxisSupported())
	assert.True(t, wrapper.IsSnapshotSupported())
	assert.True(t, wrapper.IsSupportBundleUploadSupported())
	assert.True(t, wrapper.IsSemverRequired())
}

func TestLoadLicenseFromBytes_V1Beta2(t *testing.T) {
	wrapper, err := LoadLicenseFromBytes([]byte(v1beta2LicenseYAML))
	require.NoError(t, err)

	// Version detection
	assert.False(t, wrapper.IsV1())
	assert.True(t, wrapper.IsV2())

	// Basic properties
	assert.Equal(t, "test-license-id", wrapper.GetLicenseID())
	assert.Equal(t, "test-app-slug", wrapper.GetAppSlug())
	assert.Equal(t, "test-channel-id", wrapper.GetChannelID())
	assert.Equal(t, "Stable", wrapper.GetChannelName())
	assert.Equal(t, "Test Customer", wrapper.GetCustomerName())
	assert.Equal(t, "test@example.com", wrapper.GetCustomerEmail())
	assert.Equal(t, "https://replicated.app", wrapper.GetEndpoint())
	assert.Equal(t, "trial", wrapper.GetLicenseType())
	assert.Equal(t, int64(2), wrapper.GetLicenseSequence())

	// Feature flags
	assert.True(t, wrapper.IsAirgapSupported())
	assert.True(t, wrapper.IsGitOpsSupported())
	assert.True(t, wrapper.IsIdentityServiceSupported())
	assert.True(t, wrapper.IsGeoaxisSupported())
	assert.True(t, wrapper.IsSnapshotSupported())
	assert.True(t, wrapper.IsSupportBundleUploadSupported())
	assert.True(t, wrapper.IsSemverRequired())
}

func TestLoadLicenseFromBytes_InvalidYAML(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Invalid YAML syntax",
			input: "not: valid: yaml: [",
		},
		{
			name: "Wrong GVK Group",
			input: `apiVersion: apps/v1
kind: License
metadata:
  name: test
spec:
  licenseID: test-id`,
		},
		{
			name: "Wrong GVK Kind",
			input: `apiVersion: kots.io/v1beta1
kind: Application
metadata:
  name: test
spec:
  licenseID: test-id`,
		},
		{
			name: "Unsupported Version",
			input: `apiVersion: kots.io/v1beta3
kind: License
metadata:
  name: test
spec:
  licenseID: test-id`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadLicenseFromBytes([]byte(tt.input))
			assert.Error(t, err)
		})
	}
}

func TestLoadLicenseFromPath_V1Beta1(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	licensePath := filepath.Join(tmpDir, "license.yaml")

	err := os.WriteFile(licensePath, []byte(v1beta1LicenseYAML), 0644)
	require.NoError(t, err)

	// Load from path
	wrapper, err := LoadLicenseFromPath(licensePath)
	require.NoError(t, err)

	assert.True(t, wrapper.IsV1())
	assert.Equal(t, "test-license-id", wrapper.GetLicenseID())
}

func TestLoadLicenseFromPath_V1Beta2(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	licensePath := filepath.Join(tmpDir, "license.yaml")

	err := os.WriteFile(licensePath, []byte(v1beta2LicenseYAML), 0644)
	require.NoError(t, err)

	// Load from path
	wrapper, err := LoadLicenseFromPath(licensePath)
	require.NoError(t, err)

	assert.True(t, wrapper.IsV2())
	assert.Equal(t, "test-license-id", wrapper.GetLicenseID())
}

func TestLoadLicenseFromPath_FileNotFound(t *testing.T) {
	_, err := LoadLicenseFromPath("/nonexistent/path/license.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read license file")
}

func TestLicenseWrapper_EmptyWrapper(t *testing.T) {
	// Test behavior when wrapper is empty (both fields nil)
	wrapper := LicenseWrapper{}

	// Version checks
	assert.False(t, wrapper.IsV1())
	assert.False(t, wrapper.IsV2())

	// All getters should return zero values
	assert.Equal(t, "", wrapper.GetAppSlug())
	assert.Equal(t, "", wrapper.GetLicenseID())
	assert.Equal(t, "", wrapper.GetLicenseType())
	assert.Equal(t, "", wrapper.GetEndpoint())
	assert.Equal(t, "", wrapper.GetChannelID())
	assert.Equal(t, "", wrapper.GetChannelName())
	assert.Equal(t, "", wrapper.GetCustomerName())
	assert.Equal(t, "", wrapper.GetCustomerEmail())
	assert.Equal(t, int64(0), wrapper.GetLicenseSequence())

	// All boolean getters should return false
	assert.False(t, wrapper.IsAirgapSupported())
	assert.False(t, wrapper.IsGitOpsSupported())
	assert.False(t, wrapper.IsIdentityServiceSupported())
	assert.False(t, wrapper.IsGeoaxisSupported())
	assert.False(t, wrapper.IsSnapshotSupported())
	assert.False(t, wrapper.IsSupportBundleUploadSupported())
	assert.False(t, wrapper.IsSemverRequired())
}

func TestLicenseWrapper_FeatureFlags_False(t *testing.T) {
	// Test v1beta1 with all feature flags disabled
	v1DisabledYAML := `apiVersion: kots.io/v1beta1
kind: License
metadata:
  name: test-license
spec:
  licenseID: test-license-id
  appSlug: test-app-slug
  isAirgapSupported: false
  isGitOpsSupported: false
  isIdentityServiceSupported: false
  isGeoaxisSupported: false
  isSnapshotSupported: false
  isSupportBundleUploadSupported: false
  isSemverRequired: false
`

	wrapper, err := LoadLicenseFromBytes([]byte(v1DisabledYAML))
	require.NoError(t, err)

	assert.False(t, wrapper.IsAirgapSupported())
	assert.False(t, wrapper.IsGitOpsSupported())
	assert.False(t, wrapper.IsIdentityServiceSupported())
	assert.False(t, wrapper.IsGeoaxisSupported())
	assert.False(t, wrapper.IsSnapshotSupported())
	assert.False(t, wrapper.IsSupportBundleUploadSupported())
	assert.False(t, wrapper.IsSemverRequired())
}

func TestLicenseWrapper_AllMethods_V1Beta1(t *testing.T) {
	// Comprehensive test ensuring all 18 methods work for v1beta1
	wrapper, err := LoadLicenseFromBytes([]byte(v1beta1LicenseYAML))
	require.NoError(t, err)

	methods := []struct {
		name     string
		testFunc func() interface{}
		expected interface{}
	}{
		{"IsV1", func() interface{} { return wrapper.IsV1() }, true},
		{"IsV2", func() interface{} { return wrapper.IsV2() }, false},
		{"GetAppSlug", func() interface{} { return wrapper.GetAppSlug() }, "test-app-slug"},
		{"GetLicenseID", func() interface{} { return wrapper.GetLicenseID() }, "test-license-id"},
		{"GetLicenseType", func() interface{} { return wrapper.GetLicenseType() }, "trial"},
		{"GetEndpoint", func() interface{} { return wrapper.GetEndpoint() }, "https://replicated.app"},
		{"GetChannelID", func() interface{} { return wrapper.GetChannelID() }, "test-channel-id"},
		{"GetChannelName", func() interface{} { return wrapper.GetChannelName() }, "Stable"},
		{"GetCustomerName", func() interface{} { return wrapper.GetCustomerName() }, "Test Customer"},
		{"GetCustomerEmail", func() interface{} { return wrapper.GetCustomerEmail() }, "test@example.com"},
		{"GetLicenseSequence", func() interface{} { return wrapper.GetLicenseSequence() }, int64(2)},
		{"IsAirgapSupported", func() interface{} { return wrapper.IsAirgapSupported() }, true},
		{"IsGitOpsSupported", func() interface{} { return wrapper.IsGitOpsSupported() }, true},
		{"IsIdentityServiceSupported", func() interface{} { return wrapper.IsIdentityServiceSupported() }, true},
		{"IsGeoaxisSupported", func() interface{} { return wrapper.IsGeoaxisSupported() }, true},
		{"IsSnapshotSupported", func() interface{} { return wrapper.IsSnapshotSupported() }, true},
		{"IsSupportBundleUploadSupported", func() interface{} { return wrapper.IsSupportBundleUploadSupported() }, true},
		{"IsSemverRequired", func() interface{} { return wrapper.IsSemverRequired() }, true},
	}

	for _, m := range methods {
		t.Run(m.name, func(t *testing.T) {
			result := m.testFunc()
			assert.Equal(t, m.expected, result)
		})
	}
}

func TestLicenseWrapper_AllMethods_V1Beta2(t *testing.T) {
	// Comprehensive test ensuring all 18 methods work for v1beta2
	wrapper, err := LoadLicenseFromBytes([]byte(v1beta2LicenseYAML))
	require.NoError(t, err)

	methods := []struct {
		name     string
		testFunc func() interface{}
		expected interface{}
	}{
		{"IsV1", func() interface{} { return wrapper.IsV1() }, false},
		{"IsV2", func() interface{} { return wrapper.IsV2() }, true},
		{"GetAppSlug", func() interface{} { return wrapper.GetAppSlug() }, "test-app-slug"},
		{"GetLicenseID", func() interface{} { return wrapper.GetLicenseID() }, "test-license-id"},
		{"GetLicenseType", func() interface{} { return wrapper.GetLicenseType() }, "trial"},
		{"GetEndpoint", func() interface{} { return wrapper.GetEndpoint() }, "https://replicated.app"},
		{"GetChannelID", func() interface{} { return wrapper.GetChannelID() }, "test-channel-id"},
		{"GetChannelName", func() interface{} { return wrapper.GetChannelName() }, "Stable"},
		{"GetCustomerName", func() interface{} { return wrapper.GetCustomerName() }, "Test Customer"},
		{"GetCustomerEmail", func() interface{} { return wrapper.GetCustomerEmail() }, "test@example.com"},
		{"GetLicenseSequence", func() interface{} { return wrapper.GetLicenseSequence() }, int64(2)},
		{"IsAirgapSupported", func() interface{} { return wrapper.IsAirgapSupported() }, true},
		{"IsGitOpsSupported", func() interface{} { return wrapper.IsGitOpsSupported() }, true},
		{"IsIdentityServiceSupported", func() interface{} { return wrapper.IsIdentityServiceSupported() }, true},
		{"IsGeoaxisSupported", func() interface{} { return wrapper.IsGeoaxisSupported() }, true},
		{"IsSnapshotSupported", func() interface{} { return wrapper.IsSnapshotSupported() }, true},
		{"IsSupportBundleUploadSupported", func() interface{} { return wrapper.IsSupportBundleUploadSupported() }, true},
		{"IsSemverRequired", func() interface{} { return wrapper.IsSemverRequired() }, true},
	}

	for _, m := range methods {
		t.Run(m.name, func(t *testing.T) {
			result := m.testFunc()
			assert.Equal(t, m.expected, result)
		})
	}
}

func TestLicenseWrapper_NewMethods_V1Beta1(t *testing.T) {
	// Test the 10 new methods added for additional field access
	wrapper, err := LoadLicenseFromBytes([]byte(v1beta1LicenseYAML))
	require.NoError(t, err)

	// GetSignature - v1beta1 test fixture doesn't have signature field
	assert.Nil(t, wrapper.GetSignature())

	// GetReplicatedProxyDomain
	assert.Equal(t, "", wrapper.GetReplicatedProxyDomain()) // Not in test fixture

	// GetChannels
	assert.Nil(t, wrapper.GetChannels()) // Not in test fixture

	// IsDisasterRecoverySupported
	assert.False(t, wrapper.IsDisasterRecoverySupported()) // Not in test fixture

	// IsEmbeddedClusterDownloadEnabled
	assert.False(t, wrapper.IsEmbeddedClusterDownloadEnabled()) // Not in test fixture

	// IsEmbeddedClusterMultiNodeEnabled
	assert.False(t, wrapper.IsEmbeddedClusterMultiNodeEnabled()) // Not in test fixture

	// GetEntitlements - not in test fixture
	assert.Nil(t, wrapper.GetEntitlements())
}

func TestLicenseWrapper_NewMethods_V1Beta2(t *testing.T) {
	// Test the 10 new methods added for additional field access
	wrapper, err := LoadLicenseFromBytes([]byte(v1beta2LicenseYAML))
	require.NoError(t, err)

	// GetSignature - signature is stored as []byte in the spec
	signature := wrapper.GetSignature()
	assert.NotNil(t, signature)
	// The YAML parser interprets the base64 string as bytes
	assert.True(t, len(signature) > 0)

	// GetReplicatedProxyDomain
	assert.Equal(t, "", wrapper.GetReplicatedProxyDomain()) // Not in test fixture

	// GetChannels - v1beta2 has empty channels array
	channels := wrapper.GetChannels()
	assert.NotNil(t, channels)
	assert.Empty(t, channels)

	// IsDisasterRecoverySupported
	assert.False(t, wrapper.IsDisasterRecoverySupported()) // Not in test fixture

	// IsEmbeddedClusterDownloadEnabled
	assert.False(t, wrapper.IsEmbeddedClusterDownloadEnabled()) // Not in test fixture

	// IsEmbeddedClusterMultiNodeEnabled
	assert.False(t, wrapper.IsEmbeddedClusterMultiNodeEnabled()) // Not in test fixture

	// GetEntitlements - not in test fixture
	assert.Nil(t, wrapper.GetEntitlements())
}

func TestLicenseWrapper_EmptyWrapper_NewMethods(t *testing.T) {
	// Test behavior when wrapper is empty (both fields nil) for new methods
	wrapper := LicenseWrapper{}

	// GetSignature
	assert.Nil(t, wrapper.GetSignature())

	// GetReplicatedProxyDomain
	assert.Equal(t, "", wrapper.GetReplicatedProxyDomain())

	// GetChannels
	assert.Nil(t, wrapper.GetChannels())

	// Boolean getters should return false
	assert.False(t, wrapper.IsDisasterRecoverySupported())
	assert.False(t, wrapper.IsEmbeddedClusterDownloadEnabled())
	assert.False(t, wrapper.IsEmbeddedClusterMultiNodeEnabled())

	// GetEntitlements
	assert.Nil(t, wrapper.GetEntitlements())
}
