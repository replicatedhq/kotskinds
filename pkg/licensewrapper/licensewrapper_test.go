package licensewrapper

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	kotsv1beta1 "github.com/replicatedhq/kotskinds/apis/kots/v1beta1"
	kotsv1beta2 "github.com/replicatedhq/kotskinds/apis/kots/v1beta2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/v1beta1.yaml
var testdataV1Beta1 []byte

//go:embed testdata/v1beta2.yaml
var testdataV1Beta2 []byte

//go:embed testdata/missing-values.yaml
var testdataMissingValues []byte

//go:embed testdata/blank-values.yaml
var testdataBlankValues []byte

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
	customerID: test-customer-id
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
	customerID: test-customer-id
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
	assert.Equal(t, "test-customer-id", wrapper.GetCustomerID())
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
	assert.Equal(t, "test-customer-id", wrapper.GetCustomerID())
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
	assert.Equal(t, "", wrapper.GetCustomerID())
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
	// Comprehensive test ensuring all methods work for v1beta1
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
		{"GetCustomerID", func() interface{} { return wrapper.GetCustomerID() }, "test-customer-id"},
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
		// New methods added for additional field access
		{"GetReplicatedProxyDomain", func() interface{} { return wrapper.GetReplicatedProxyDomain() }, ""},
		{"IsDisasterRecoverySupported", func() interface{} { return wrapper.IsDisasterRecoverySupported() }, false},
		{"IsEmbeddedClusterDownloadEnabled", func() interface{} { return wrapper.IsEmbeddedClusterDownloadEnabled() }, false},
		{"IsEmbeddedClusterMultiNodeEnabled", func() interface{} { return wrapper.IsEmbeddedClusterMultiNodeEnabled() }, false},
	}

	for _, m := range methods {
		t.Run(m.name, func(t *testing.T) {
			result := m.testFunc()
			assert.Equal(t, m.expected, result)
		})
	}

	// Typed nil values require assert.Nil instead of assert.Equal
	t.Run("GetSignature", func(t *testing.T) {
		assert.Nil(t, wrapper.GetSignature())
	})

	t.Run("GetChannels", func(t *testing.T) {
		assert.Nil(t, wrapper.GetChannels())
	})

	t.Run("GetEntitlements", func(t *testing.T) {
		assert.Nil(t, wrapper.GetEntitlements())
	})
}

func TestLicenseWrapper_AllMethods_V1Beta2(t *testing.T) {
	// Comprehensive test ensuring all methods work for v1beta2
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
		{"GetCustomerID", func() interface{} { return wrapper.GetCustomerID() }, "test-customer-id"},
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
		// New methods added for additional field access
		{"GetReplicatedProxyDomain", func() interface{} { return wrapper.GetReplicatedProxyDomain() }, ""},
		{"IsDisasterRecoverySupported", func() interface{} { return wrapper.IsDisasterRecoverySupported() }, false},
		{"IsEmbeddedClusterDownloadEnabled", func() interface{} { return wrapper.IsEmbeddedClusterDownloadEnabled() }, false},
		{"IsEmbeddedClusterMultiNodeEnabled", func() interface{} { return wrapper.IsEmbeddedClusterMultiNodeEnabled() }, false},
	}

	for _, m := range methods {
		t.Run(m.name, func(t *testing.T) {
			result := m.testFunc()
			assert.Equal(t, m.expected, result)
		})
	}

	// Special cases that need custom assertions
	t.Run("GetSignature", func(t *testing.T) {
		signature := wrapper.GetSignature()
		assert.NotNil(t, signature)
		assert.True(t, len(signature) > 0)
	})

	t.Run("GetChannels", func(t *testing.T) {
		channels := wrapper.GetChannels()
		assert.NotNil(t, channels)
		assert.Empty(t, channels)
	})

	t.Run("GetEntitlements", func(t *testing.T) {
		assert.Nil(t, wrapper.GetEntitlements())
	})
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

func TestLicenseWrapper_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		wrapper  *LicenseWrapper
		expected bool
	}{
		{
			name:     "nil wrapper",
			wrapper:  nil,
			expected: true,
		},
		{
			name:     "empty wrapper with both V1 and V2 nil",
			wrapper:  &LicenseWrapper{},
			expected: true,
		},
		{
			name: "wrapper with V1 license",
			wrapper: &LicenseWrapper{
				V1: &kotsv1beta1.License{},
			},
			expected: false,
		},
		{
			name: "wrapper with V2 license",
			wrapper: &LicenseWrapper{
				V2: &kotsv1beta2.License{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.wrapper.IsEmpty()
			assert.Equal(t, tt.expected, result, "IsEmpty() returned unexpected value for %s", tt.name)
		})
	}
}

func TestLicenseWrapper_IsEmpty_WithLoadedLicenses(t *testing.T) {
	// Test with actual loaded v1beta1 license
	t.Run("loaded v1beta1 license is not empty", func(t *testing.T) {
		wrapper, err := LoadLicenseFromBytes([]byte(v1beta1LicenseYAML))
		require.NoError(t, err)
		assert.False(t, wrapper.IsEmpty())
	})

	// Test with actual loaded v1beta2 license
	t.Run("loaded v1beta2 license is not empty", func(t *testing.T) {
		wrapper, err := LoadLicenseFromBytes([]byte(v1beta2LicenseYAML))
		require.NoError(t, err)
		assert.False(t, wrapper.IsEmpty())
	})
}

func TestLicenseWrapper_VerifySignature(t *testing.T) {
	tests := []struct {
		name        string
		wrapper     *LicenseWrapper
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil wrapper returns error",
			wrapper:     nil,
			expectError: true,
			errorMsg:    "license wrapper is empty",
		},
		{
			name:        "empty wrapper returns error",
			wrapper:     &LicenseWrapper{},
			expectError: true,
			errorMsg:    "license wrapper is empty",
		},
		{
			name: "wrapper with V1 license but no signature returns error",
			wrapper: &LicenseWrapper{
				V1: &kotsv1beta1.License{},
			},
			expectError: true,
			// Error will come from ValidateLicense, not our wrapper
		},
		{
			name: "wrapper with V2 license but no signature returns error",
			wrapper: &LicenseWrapper{
				V2: &kotsv1beta2.License{},
			},
			expectError: true,
			// Error will come from ValidateLicense, not our wrapper
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.wrapper.VerifySignature()
			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLicenseWrapper_VerifySignature_WithTestData(t *testing.T) {
	tests := []struct {
		name        string
		licenseData []byte
		expectError bool
		errorMsg    string
	}{
		{
			name:        "v1beta1 license with valid signatures",
			licenseData: testdataV1Beta1,
			expectError: false,
		},
		{
			name:        "v1beta2 license with valid signatures",
			licenseData: testdataV1Beta2,
			expectError: false,
		},
		{
			name:        "license with missing entitlement field values",
			licenseData: testdataMissingValues,
			expectError: false,
		},
		{
			name:        "license with blank entitlement field values",
			licenseData: testdataBlankValues,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load the license from embedded testdata
			wrapper, err := LoadLicenseFromBytes(tt.licenseData)
			require.NoError(t, err, "failed to load license from embedded data")

			// Verify the signature
			err = wrapper.VerifySignature()

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err, "signature verification failed for embedded data")
			}
		})
	}
}
