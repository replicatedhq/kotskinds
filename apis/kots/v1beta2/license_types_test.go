package v1beta2_test

import (
	"encoding/json"
	"testing"

	kotsv1beta2 "github.com/replicatedhq/kotskinds/apis/kots/v1beta2"
	kotsscheme "github.com/replicatedhq/kotskinds/client/kotsclientset/scheme"
	"github.com/replicatedhq/kotskinds/pkg/crypto"
	"github.com/replicatedhq/kotskinds/pkg/licensewrapper/types"
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

func TestLicenseSignature(t *testing.T) {
	// set the global signing key
	globalKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxHh2OXzDqlQ7kZJ1d4zr
wbpXsSFHcYzr+k6pe+QXLUelAMvlik9NXauIt+YFtEAxNypV+xPCr8ClH5L2qPPb
QBeG0ExxzvRshDMGxm7TXVHzTXQCrD7azS8Va6RsAB4tJMlvymn2uHsQDbShQiOY
RKaRY/KKBmaIcYmysaSvfU8E5Ve9f4478X3u1cPzKUG6dk5j1Nt3nSv3BWINM5ec
IXJQCB+gQVkOjzvA9aRVtLJtFqAoX7A6BfTNqrx35eyBEmzQOo0Mx1JkZDDW4+qC
bhC0kq14IRpwKFIALBhSojfbJelM+gCv3wjF4hrWxAZQzWSPexP1Msof2KbrniEe
LQIDAQAB
-----END PUBLIC KEY-----
`
	err := crypto.SetCustomPublicKeyRSA(globalKey)
	require.NoError(t, err)

	// test license signed with the above global key
	data := `
apiVersion: kots.io/v1beta2
kind: License
metadata:
  name: testcustomer
spec:
  licenseID: test-license-id
  licenseType: trial
  customerName: Test Customer
  appSlug: test-app
  channelID: "1"
  channelName: Stable
  customerEmail: test@example.com
  channels:
  - channelID: "1"
    channelSlug: stable
    channelName: Stable
    isDefault: true
    isSemverRequired: true
  entitlements:
    bool_default_false_set_true:
      title: Boolean Default False Set True
      value: true
      valueType: Boolean
      signature:
        v2: UjAdcioep7RLLOVEm4A0Ig6KOudvMFFE0IWC9wW1it3CNKyH4q2gE5KVUEE5mD1kpQqg0B0arZQo1N1uUjlRfTksprS0eoRzF4qz+vrlRs4sSHNG3bMOWk0q7qGCmfJAsSXRi3mOJ83QxSeLexlmrhZammiTeFy+XmPs+oT/odIuGJ9HeYyuuXG+UQycy0eifI0Mm+6kYAZXnzyMl6dNOrD+ZSfbo8o/iOdKoeo6UoycZVkcrDh0NfXu3Zxw9yPG6xzMN64uN2ujeg0mG6ejRDdAdKUexxuci878XwRsToWfHwNxWS/vF+EOxpPzDPdW0LUr/sN0J94icedx0PCFNA==
    expires_at:
      title: Expiration
      description: License Expiration
      value: "2025-12-31T23:59:59Z"
      valueType: String
      signature:
        v2: e/KrJiLwB45UKkSiRCnfkiXS5yCJGLRvSlLZKh814P0flaGxXQCJ8r56hZvkAN2eG4KrwT5jsFg7LRZmjMvyY49RrjUYaYNFpC5Dp3zS/L07dw/SLZVtHDg8wZyNej9MlvRXgyHHqjl8HxFsUZzsQC9+cDrGY03AjWZCqOPzpHqmn+VRi3r6XDQE0UEOhmS/y66aAWrgzOH8iTJRp9WVtHyGUKCu+mwag5thAUNC/JNRUeUIlxvF7zCtpBvhO9LERBIewqHL3nENi+XZxe7xxNe+Q7zb31u4uvD887sgNsArtdorMZWPJSxqp9T3HwAQ3p0NpYh2OD2IWfOcNPVBYA==
    int_default_0_set_587:
      title: Integer Default 0 Set 587
      value: 587
      valueType: Integer
      signature:
        v2: O4828HBAJWNBEUqFGDpZ1pKwjUC+jVHMerZtCQ8wBny8MoL0BNqp9C7VVe4tfnfKqPknneDoZL02Tpxp5OvnkTG9NwO6C/+t5rV0KsJHrVCr242/J8HW4y8S0QTYfDKKxVW1z2jX+1IE/QBvrkV65EXWwpB8vJayF5vBzEER520+Tbs8QdSEIyPR9BZpnTxsMjgH496RwE4rp1KY/dcRe3mR7iLKFsIPic6EjCrloOR5RPDqRGRsVzFw4p0Hx3dKMeSF3xpvAr6/BPFnD4TnElMXaU5xRnXoR2LatcveiJLoRDZUJSohvnz027t28hI0ECJ5/ixUebNqy3RaLAOHUA==
    string_default_empty_set_teststring:
      title: String Default Empty Set Teststring
      value: teststring
      valueType: String
      signature:
        v2: PxR4yjwV2UL19YrNgQmJFnVdb1cHihrByoCXJokItAPA9s64WtFWLpq1BZ0diku+lRzNQoGcxIWjO8UlHqIvyawmR6NYOetgmO1ujbQf7rioohROVQOn3vNWgzapIDhLAI1EP/FIPzsBfVg0T/1ZNkwxWxMwPHFV3FgJL62cHwXi+UBMNibIIRlojUULOyOSU8jeFuTqxtWSQgGM2qyGVQyl3/lX/lGZSA+98OVhEHUc3HciM0zpAHcko6vgVsh3akC0Kv0P4dvd1TOckiDqyzpuf7lDKDLOkppp0o0gI9lK9mzGDz57tCXr4HkqIgMzF2a11O96637u/prEDcWgqA==
  isNewKotsUiEnabled: true
  isSemverRequired: true
  signature: eyJsaWNlbnNlRGF0YSI6ImV5SmhjR2xXWlhKemFXOXVJam9pYTI5MGN5NXBieTkyTVdKbGRHRXlJaXdpYTJsdVpDSTZJa3hwWTJWdWMyVWlMQ0p0WlhSaFpHRjBZU0k2ZXlKdVlXMWxJam9pZEdWemRHTjFjM1J2YldWeUluMHNJbk53WldNaU9uc2liR2xqWlc1elpVbEVJam9pZEdWemRDMXNhV05sYm5ObExXbGtJaXdpYkdsalpXNXpaVlI1Y0dVaU9pSjBjbWxoYkNJc0ltTjFjM1J2YldWeVRtRnRaU0k2SWxSbGMzUWdRM1Z6ZEc5dFpYSWlMQ0poY0hCVGJIVm5Jam9pZEdWemRDMWhjSEFpTENKamFHRnVibVZzU1VRaU9pSXhJaXdpWTJoaGJtNWxiRTVoYldVaU9pSlRkR0ZpYkdVaUxDSmpkWE4wYjIxbGNrVnRZV2xzSWpvaWRHVnpkRUJsZUdGdGNHeGxMbU52YlNJc0ltTm9ZVzV1Wld4eklqcGJleUpqYUdGdWJtVnNTVVFpT2lJeElpd2lZMmhoYm01bGJGTnNkV2NpT2lKemRHRmliR1VpTENKamFHRnVibVZzVG1GdFpTSTZJbE4wWVdKc1pTSXNJbWx6UkdWbVlYVnNkQ0k2ZEhKMVpTd2lhWE5UWlcxMlpYSlNaWEYxYVhKbFpDSTZkSEoxWlgxZExDSmxiblJwZEd4bGJXVnVkSE1pT25zaVltOXZiRjlrWldaaGRXeDBYMlpoYkhObFgzTmxkRjkwY25WbElqcDdJblJwZEd4bElqb2lRbTl2YkdWaGJpQkVaV1poZFd4MElFWmhiSE5sSUZObGRDQlVjblZsSWl3aWRtRnNkV1VpT25SeWRXVXNJblpoYkhWbFZIbHdaU0k2SWtKdmIyeGxZVzRpTENKemFXZHVZWFIxY21VaU9uc2lkaklpT2lKVmFrRmtZMmx2WlhBM1VreE1UMVpGYlRSQk1FbG5Oa3RQZFdSMlRVWkdSVEJKVjBNNWQxY3hhWFF6UTA1TGVVZzBjVEpuUlRWTFZsVkZSVFZ0UkRGcmNGRnhaekJDTUdGeVdsRnZNVTR4ZFZWcWJGSm1WR3R6Y0hKVE1HVnZVbnBHTkhGNkszWnliRkp6TkhOVFNFNUhNMkpOVDFkck1IRTNjVWREYldaS1FYTlRXRkpwTTIxUFNqZ3pVWGhUWlV4bGVHeHRjbWhhWVcxdGFWUmxSbmtyV0cxUWN5dHZWQzl2WkVsMVIwbzVTR1ZaZVhWMVdFY3JWVkY1WTNrd1pXbG1TVEJOYlNzMmExbEJXbGh1ZW5sTmJEWmtUazl5UkN0YVUyWmliemh2TDJsUFpFdHZaVzgyVlc5NVkxcFdhMk55Ukdnd1RtWllkVE5hZUhjNWVWQkhObmg2VFU0Mk5IVk9NblZxWldjd2JVYzJaV3BTUkdSQlpFdFZaWGg0ZFdOcE9EYzRXSGRTYzFSdlYyWklkMDU0VjFNdmRrWXJSVTk0Y0ZCNlJGQmtWekJNVlhJdmMwNHdTamswYVdObFpIZ3dVRU5HVGtFOVBTSjlmU3dpWlhod2FYSmxjMTloZENJNmV5SjBhWFJzWlNJNklrVjRjR2x5WVhScGIyNGlMQ0prWlhOamNtbHdkR2x2YmlJNklreHBZMlZ1YzJVZ1JYaHdhWEpoZEdsdmJpSXNJblpoYkhWbElqb2lNakF5TlMweE1pMHpNVlF5TXpvMU9UbzFPVm9pTENKMllXeDFaVlI1Y0dVaU9pSlRkSEpwYm1jaUxDSnphV2R1WVhSMWNtVWlPbnNpZGpJaU9pSmxMMHR5U21sTWQwSTBOVlZMYTFOcFVrTnVabXRwV0ZNMWVVTktSMHhTZGxOc1RGcExhRGd4TkZBd1pteGhSM2hZVVVOS09ISTFObWhhZG10QlRqSmxSelJMY25kVU5XcHpSbWMzVEZKYWJXcE5kbmxaTkRsU2NtcFZXV0ZaVGtad1F6VkVjRE42VXk5TU1EZGtkeTlUVEZwV2RFaEVaemgzV25sT1pXbzVUV3gyVWxobmVVaEljV3BzT0VoNFJuTlZXbnB6VVVNNUsyTkVja2RaTUROQmFsZGFRM0ZQVUhwd1NIRnRiaXRXVW1remNqWllSRkZGTUZWRlQyaHRVeTk1TmpaaFFWZHlaM3BQU0RocFZFcFNjRGxYVm5SSWVVZFZTME4xSzIxM1lXYzFkR2hCVlU1REwwcE9VbFZsVlVsc2VIWkdOM3BEZEhCQ2RtaFBPVXhGVWtKSlpYZHhTRXd6YmtWT2FTdFlXbmhsTjNoNFRtVXJVVGQ2WWpNeGRUUjFka1E0T0RkelowNXpRWEowWkc5eVRWcFhVRXBUZUhGd09WUXpTSGRCVVROd01FNXdXV2d5VDBReVNWZG1UMk5PVUZaQ1dVRTlQU0o5ZlN3aWFXNTBYMlJsWm1GMWJIUmZNRjl6WlhSZk5UZzNJanA3SW5ScGRHeGxJam9pU1c1MFpXZGxjaUJFWldaaGRXeDBJREFnVTJWMElEVTROeUlzSW5aaGJIVmxJam8xT0Rjc0luWmhiSFZsVkhsd1pTSTZJa2x1ZEdWblpYSWlMQ0p6YVdkdVlYUjFjbVVpT25zaWRqSWlPaUpQTkRneU9FaENRVXBYVGtKRlZYRkdSMFJ3V2pGd1MzZHFWVU1yYWxaSVRXVnlXblJEVVRoM1FtNTVPRTF2VERCQ1RuRndPVU0zVmxabE5IUm1ibVpMY1ZCcmJtNWxSRzlhVERBeVZIQjRjRFZQZG01clZFYzVUbmRQTmtNdkszUTFjbFl3UzNOS1NISldRM0l5TkRJdlNqaElWelI1T0ZNd1VWUlpaa1JMUzNoV1Z6RjZNbXBZS3pGSlJTOVJRblp5YTFZMk5VVllWM2R3UWpoMlNtRjVSalYyUW5wRlJWSTFNakFyVkdKek9GRmtVMFZKZVZCU09VSmFjRzVVZUhOTmFtZElORGsyVW5kRk5ISndNVXRaTDJSalVtVXpiVkkzYVV4TFJuTkpVR2xqTmtWcVEzSnNiMDlTTlZKUVJIRlNSMUp6Vm5wR2R6UndNRWg0TTJSTFRXVlRSak40Y0haQmNqWXZRbEJHYmtRMFZHNUZiRTFZWVZVMWVGSnVXRzlTTWt4aGRHTjJaV2xLVEc5U1JGcFZTbE52YUhadWVqQXlOM1F5T0doSk1FVkRTalV2YVhoVlpXSk9jWGt6VW1GTVFVOUlWVUU5UFNKOWZTd2ljM1J5YVc1blgyUmxabUYxYkhSZlpXMXdkSGxmYzJWMFgzUmxjM1J6ZEhKcGJtY2lPbnNpZEdsMGJHVWlPaUpUZEhKcGJtY2dSR1ZtWVhWc2RDQkZiWEIwZVNCVFpYUWdWR1Z6ZEhOMGNtbHVaeUlzSW5aaGJIVmxJam9pZEdWemRITjBjbWx1WnlJc0luWmhiSFZsVkhsd1pTSTZJbE4wY21sdVp5SXNJbk5wWjI1aGRIVnlaU0k2ZXlKMk1pSTZJbEI0VWpSNWFuZFdNbFZNTVRsWmNrNW5VVzFLUm01V1pHSXhZMGhwYUhKQ2VXOURXRXB2YTBsMFFWQkJPWE0yTkZkMFJsZE1jSEV4UWxvd1pHbHJkU3RzVW5wT1VXOUhZM2hKVjJwUE9GVnNTSEZKZG5saGQyMVNOazVaVDJWMFoyMVBNWFZxWWxGbU4zSnBiMjlvVWs5V1VVOXVNM1pPVjJkNllYQkpSR2hNUVVreFJWQXZSa2xRZW5OQ1psWm5NRlF2TVZwT2EzZDRWM2hOZDFCSVJsWXpSbWRLVERZeVkwaDNXR2tyVlVKTlRtbGlTVWxTYkc5cVZWVk1UM2xQVTFVNGFtVkdkVlJ4ZUhSWFUxRm5SMDB5Y1hsSFZsRjViRE12YkZndmJFZGFVMEVyT1RoUFZtaEZTRlZqTTBoamFVMHdlbkJCU0dOcmJ6WjJaMVp6YUROaGEwTXdTM1l3VURSa2RtUXhWRTlqYTJsRWNYbDZjSFZtTjJ4RVMwUk1UMnR3Y0hBd2J6Qm5TVGxzU3psdGVrZEVlalUzZEVOWWNqUklhM0ZKWjAxNlJqSmhNVEZQT1RZMk16ZDFMM0J5UlVSalYyZHhRVDA5SW4xOWZTd2lhWE5PWlhkTGIzUnpWV2xGYm1GaWJHVmtJanAwY25WbExDSnBjMU5sYlhabGNsSmxjWFZwY21Wa0lqcDBjblZsZlgwPSIsImlubmVyU2lnbmF0dXJlIjoiZXlKMk1reHBZMlZ1YzJWVGFXZHVZWFIxY21VaU9pSnZMM013V1V4b2NsRXhUMlJPTTFkM2NsY3dRMU42T0VWRWQyVjVaWGRvTkVseGIycFhWeTlwTnprNVUyTlRWVU5FVUVWVGNUTkxhMnRoYjBWclFXbDFlVGg1V2pZdmEzQXZNSHBNWVhWSk5uSmxSbUVyZEVKUWEyVmpZbVJuV2sxcldVVllOM1JFSzA1a1dqVnhka1JLUld0c1lUWjBWRkJRVld4Uk0zVkROVXgwU0ZwS1kyYzJOSFIzT1ZSNWRsWjFMMjR2UW5oa1F6QlRWWG8yZGpReVVIaHlTbFJsWnl0MmNHaFhhVkZsYkRCR2RWZFBOemM0TDI1ck9GbFphRGhFZVRZeGFtZFpLMVEwUkVkUFFTc3lVMU5NZGxkeGNERm9VbXhSVFhCWE5WUmpiR05sUTB4NlNUSm5aR3RqZHl0bVVDOVVXQ3RwYnpFeFFtOVBWMDVtYW0xV05HcFdSWEF4V0Zvd1QyTlJlR012UkV3elltdE1XalJLZFdoSWNWQlZObFptTmxWU1dYZFFVbUpZS3pobmJsZFpibGxuUTJaTFprWnBWa3RvY1dsVFpVWnBXbXc0ZURsU09XaGhVbmQ0T1c5MVIzYzlQU0lzSW5CMVlteHBZMHRsZVNJNklpMHRMUzB0UWtWSFNVNGdVRlZDVEVsRElFdEZXUzB0TFMwdFhHNU5TVWxDU1dwQlRrSm5hM0ZvYTJsSE9YY3dRa0ZSUlVaQlFVOURRVkU0UVUxSlNVSkRaMHREUVZGRlFYb3pNRTB6WXpRM2IyVjVUakZUYVdkeFIwTlFYRzVhUjI0elJYZHJUSEpwVjBnMFNrZzBlWGx1UjBJMWVFNXZLMlJwWTBJdlVHOVNaSFZQVTFNeldVdHVla1J1WTI1NWEyVTVkMHhsVnl0SlEySTNiVmhtWEc0NVdtaDZVM00yV25CS1lrSlhjVWsxTjJGT1duQXJUVEF3YkdST2RrZFVVWEZOZDFSaWNESm1UVVZNUkVKa01qY3hZa1l5UWxadVUwNUliQ3QxVVhSWFhHNU5iRU1yVDBoUmRXaHZZbk00ZEVSUE5sUlRTWGRGWWsxYWFYWk1RVGxMVG5WdGVqRkNiMDVUWTJaMGNucDFVVXQwTTBaVEx6aGFVR2x2YjBKSFdtSk1YRzVyTUVoUmVXWTJWVVpUWVRnM0syeFRTMjlzVFdkRFlsVTBWMHR5V2xKTlRFdE5lbm96ZFZORmJIRnBia2hsUmxoc1RVMXhUM05FZEVFMmJUZFVWbVJYWEc1SmRqSm5iVE53Vld0dU1sSXhiSEJ0T0RKSFZWSnZUVzB5VUhCblUyVTJTV0p0VWxOcGRsSlliWGhyU1ZRd1NUVTJPVnBGYzFSemFYWllXVEYwT0VFMVhHNXRkMGxFUVZGQlFseHVMUzB0TFMxRlRrUWdVRlZDVEVsRElFdEZXUzB0TFMwdFhHNGlMQ0oyTWt0bGVWTnBaMjVoZEhWeVpTSTZJbVY1U25waFYyUjFXVmhTTVdOdFZXbFBhVXBDVm10T1NVc3pWbWxoUlRWUVkwVldSVkpWTVZkT1IzQlRZakF4ZG1ScWFIUk9SRnB3VGpJMU1HSXlkR0ZOYkhCQ1l6TlJlVmxUZEVwUldFNVBUVlU1UmxkdVdreE9lWFJyV2xSa01XVllhekZhVlRWNVRqQkpNMW95V2tsa1YwVXdWakpvUTFScGRFeE5NR2h3U3pOWk5XTXhUVFJUYTFweVdXeGtkMDFUT0RCWmJsSlNUV3hGTW1GdE5UQlhWRm96VWtWV01WcHViRWRYVm1NMVN6SmFNRmt4UW1GTmEzUjNZM3ByZVZsc2FFWlZiRXBwWTJrNVJWcFZhRFJOYTNCSVlWZHNXVnBGUmxOYVJWVTFWbXhvTWxORVVuUldSbXhXVWxoYVVHTnRPV3hXYkVKMldWaFNRMk13YUhGaU0yUkpaVlU1TmxSclVrOWxiVW95VjFWV1IxTllTbTVrU0ZaWVRqTmFTRk5JU2pWalZFbzFUMGRLYmxKVmVIaGxiRW8wVjFoV2RsWjZiSFZpVkUxNFltMTRSR0pJWkRGV00wcFdUMGRzV1ZGV1VYbE1lbHB3VWtaRk1GTkhXbHBTYkVWMlV6QkthVTlGY0hWUmJWcDNWa2RKZUUxWFZuVlZiRXAxV1RKbk5XSnRVVE5OTW1NelpXdGFOVkpJVW5CVmJUUXhXa1ZzVjA1cVZrdFNNV2hyVmpOdmQyVklaSFJTYkVwS1pHcFdkbUZGT1ZSalZXTjJVMGhDV1U5WGFHOU5NV1IyWTNwb01XUnRNVTVOTTBKUVpHdHpNRTVWVmt4Wk1tTTVVRk5KYzBsdFpITmlNa3BvWWtWMGJHVlZiR3RKYW05cFpFZFdlbVJETVc1aVJ6bHBXVmQzZEdFeVZqVk1WMnhyU1c0d1BTSjkifQ==

`

	kotsscheme.AddToScheme(scheme.Scheme)

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, gvk, err := decode([]byte(data), nil, nil)
	require.NoError(t, err)

	assert.Equal(t, "kots.io", gvk.Group)
	assert.Equal(t, "v1beta2", gvk.Version)
	assert.Equal(t, "License", gvk.Kind)

	license := obj.(*kotsv1beta2.License)
	assert.Equal(t, "test-license-id", license.Spec.LicenseID)

	appKeys, err := license.ValidateLicense()
	require.NoError(t, err, "the unchanged license should be valid")

	// change a non-spec field and validate that things do not break
	license.ObjectMeta.Name = "changed"
	_, err = license.ValidateLicense()
	require.NoError(t, err, "changing non-spec fields should not invalidate the license")

	// change a spec field and validate that the signature is no longer valid
	license.Spec.AppSlug = "changed"
	_, err = license.ValidateLicense()
	require.Error(t, err, "changing spec fields should invalidate the license")
	var ldve *types.LicenseDataValidationError
	require.ErrorAs(t, err, &ldve, "error should be a LicenseDataValidationError")
	require.True(t, types.IsLicenseDataValidationError(err), "error should be a LicenseDataValidationError")
	require.Equal(t, "test-app", license.Spec.AppSlug, "app slug should be returned to the signed value after validation")

	// validate entitlement signatures
	for k, val := range license.Spec.Entitlements {
		err = val.ValidateSignature(appKeys)
		require.NoError(t, err, "entitlement %s signature should be valid", k)
	}

	// change entitlement value and revalidate
	ent, ok := license.Spec.Entitlements["int_default_0_set_587"]
	require.True(t, ok, "int_default_0_set_587 should be an entitlement")
	ent.Value.IntVal = 33
	err = ent.ValidateSignature(appKeys)
	require.Error(t, err, "changing entitlement value should break signatures")
}
