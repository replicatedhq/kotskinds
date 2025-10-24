package v1beta2_test

import (
	"encoding/json"
	"testing"

	kotsv1beta2 "github.com/replicatedhq/kotskinds/apis/kots/v1beta2"
	kotsscheme "github.com/replicatedhq/kotskinds/client/kotsclientset/scheme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/scheme"
)

func Test_License(t *testing.T) {
	data := `apiVersion: kots.io/v1beta2
kind: License
metadata:
  name: local
spec:
  licenseID: abcdef
  appSlug: my-app
  endpoint: 'http://localhost:30016'
  entitlements:
    expires_at:
      title: Expiration
      description: License Expiration
      value: ""
      signature:
        v2: SUE9PQ==
    has-product-2:
      title: Has Product 2
      value: "test"
    is_vip:
      title: Is VIP
      value: false
    num_seats:
      title: Number Of Seats
      value: 10
    sdzf:
      title: sdf
      value: 1
    test:
      title: test
      value: "123asd"
  signature: SUE9PQ==`

	kotsscheme.AddToScheme(scheme.Scheme)

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, gvk, err := decode([]byte(data), nil, nil)
	require.NoError(t, err)

	assert.Equal(t, "kots.io", gvk.Group)
	assert.Equal(t, "v1beta2", gvk.Version)
	assert.Equal(t, "License", gvk.Kind)

	license := obj.(*kotsv1beta2.License)

	assert.Equal(t, "abcdef", license.Spec.LicenseID)
	assert.NotEmpty(t, license.Spec.Signature)

	entitlements := license.Spec.Entitlements
	assert.Len(t, entitlements, 6)

	expiresAt := entitlements["expires_at"]
	assert.NotNil(t, expiresAt)
	assert.Equal(t, "Expiration", expiresAt.Title)
	assert.Equal(t, "License Expiration", expiresAt.Description)
	assert.NotEmpty(t, expiresAt.Signature.V2)

	numSeats := entitlements["num_seats"]
	assert.NotNil(t, numSeats)
	assert.Equal(t, "Number Of Seats", numSeats.Title)
	assert.Equal(t, "", numSeats.Description)

	testField := entitlements["test"]
	assert.NotNil(t, testField)
	assert.Equal(t, "test", testField.Title)
	assert.Equal(t, "123asd", testField.Value.Value())
}

func Test_SignatureField(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected []byte
		hasError bool
	}{
		{
			name: "signature field is present",
			jsonData: `{
				"licenseID": "test-id",
				"appSlug": "test-app",
				"signature": "SUE9PQ=="
			}`,
			expected: []byte("IA=="),
			hasError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var spec kotsv1beta2.LicenseSpec
			err := json.Unmarshal([]byte(test.jsonData), &spec)
			if test.hasError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, spec.Signature)
			assert.NotEmpty(t, spec.LicenseID)
			assert.NotEmpty(t, spec.AppSlug)
		})
	}
}

func Test_EntitlementSignatureV2(t *testing.T) {
	tests := []struct {
		name           string
		jsonData       string
		expectedV2     []byte
		entitlementKey string
	}{
		{
			name: "v2 signature present",
			jsonData: `{
				"licenseID": "test-id",
				"appSlug": "test-app",
				"entitlements": {
					"test-entitlement": {
						"title": "Test",
						"value": "test-value",
						"signature": {
							"v2": "djJzaWc="
						}
					}
				}
			}`,
			expectedV2:     []byte("v2sig"),
			entitlementKey: "test-entitlement",
		},
		{
			name: "no signature",
			jsonData: `{
				"licenseID": "test-id",
				"appSlug": "test-app",
				"entitlements": {
					"test-entitlement": {
						"title": "Test",
						"value": "test-value"
					}
				}
			}`,
			expectedV2:     nil,
			entitlementKey: "test-entitlement",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var spec kotsv1beta2.LicenseSpec
			err := json.Unmarshal([]byte(test.jsonData), &spec)
			require.NoError(t, err)

			entitlement := spec.Entitlements[test.entitlementKey]
			assert.NotNil(t, entitlement)
			assert.Equal(t, test.expectedV2, entitlement.Signature.V2)
		})
	}
}

func Test_IsEmbeddedClusterMultiNodeEnabled(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected bool
	}{
		{
			name: "field is missing, should default to false",
			jsonData: `{
				"licenseID": "test-id",
				"appSlug": "test-app"
			}`,
			expected: false,
		},
		{
			name: "field is explicitly set to false",
			jsonData: `{
				"licenseID": "test-id",
				"appSlug": "test-app",
				"isEmbeddedClusterMultiNodeEnabled": false
			}`,
			expected: false,
		},
		{
			name: "field is explicitly set to true",
			jsonData: `{
				"licenseID": "test-id",
				"appSlug": "test-app",
				"isEmbeddedClusterMultiNodeEnabled": true
			}`,
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var spec kotsv1beta2.LicenseSpec
			err := json.Unmarshal([]byte(test.jsonData), &spec)
			require.NoError(t, err)
			assert.Equal(t, test.expected, spec.IsEmbeddedClusterMultiNodeEnabled)
			assert.NotEmpty(t, spec.LicenseID)
			assert.NotEmpty(t, spec.AppSlug)
		})
	}
}
