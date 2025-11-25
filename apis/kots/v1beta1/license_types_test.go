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

package v1beta1_test

import (
	"bytes"
	"testing"

	kotsv1beta1 "github.com/replicatedhq/kotskinds/apis/kots/v1beta1"
	kotsscheme "github.com/replicatedhq/kotskinds/client/kotsclientset/scheme"
	"github.com/replicatedhq/kotskinds/pkg/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestLicenseSignature(t *testing.T) {
	// set the global signing key
	globalKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAy5036gItehm9afUCntLs
1FjoWo6Pbp1RmrbaJQX6xG4qE/Lq45C23t+cfFBFyzjUbgOqzOYGa1V8oTUvdZ+C
wWWY12LrIVZtWg3ynbxo6SCJUS1pRA1S/8sQOGMWqRkxSkPnA2bE0je7dttyyd+A
65RuNi4tUjxqEAt0pwogTVAZYCvplHyvOsGqIAjMwxuM86XBbTJiQAaU3d+qwIFg
Y48EHCqMf+cyDFkulbP+cbXafv9GtfTODj9cz6Rz/AIodjivWQ7luD3w3K2JkfJu
ul7aWXO7BKThyLPQ+VvzqR+mHUKOK/c4XXqFh6a351anSZEm/x03brbs9Kwgvdy3
gwIDAQAB
-----END PUBLIC KEY-----
`
	err := crypto.SetCustomPublicKeyRSA(globalKey)
	require.NoError(t, err)

	// test license signed with the above global key
	data := `
apiVersion: kots.io/v1beta1
kind: License
metadata:
  name: testcustomer
spec:
  licenseID: test-license-id
  licenseType: trial
  customerID: test-customer-id
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
        v1: f5y+Fu3Oa8z68N1LlkJxN+dF/3gY2+4DLyptpHe6JVNA04Bi/3H44Zdp5HOVB8sbqvOjXp/wdrXk1Jcev3wuMJXP8Tz48YSOCdzIB7vev7CYbtjzeR0YaVNxxIBlgWwNwiPAthrTKsIVJJQ+DqFNbjRdwdPyZldu+D25EnJ/nNxrS0piYKCokCADf1HRPWTmQEe99WCk2FTbGbQOUtOUoWyFwZO0DHA6Ua1VNEBVgE6XItZSiuQydfI5E+Nnxh7r28vxx/+wUgGlKIJ4b2V4DGgA8NRjPEnca6RpEus9Ne+EG4UtTRt0pKdZcteRXIl/w5778M4eSuKNSVtp826/9Q==
    expires_at:
      title: Expiration
      description: License Expiration
      value: "2025-12-31T23:59:59Z"
      valueType: String
      signature:
        v1: Z0XutI8rb0xUhuKFNhJhgopnicidzSK4Y6zwn+G+I2+kEnMRQw2z8w7qhLZCuIYT6ngSFj6OMUCPa/BKwfnqH2ykKxexCzHicuDDuT/PB+06lTaCj2MSWyjGNxuX/XefP4c/pOjR4wqoKMddJQdcqWn4pjm07Gt1WRVWDGvulMEPzoRLvXS0GxevwlylAWE0mFojtdloaGjXOmnyozztBeNxaEPR12bdEgZLqHlT33T8uVwUPYxBy1RUqXtv8Jfs2dsJGruMTVwliNejlz9haMumvmugL1rX2uSMpStuo7BIY2Vx/ykr/QmOsqhKjz3JvHRDyLTjYugWyClItijc7w==
    int_default_0_set_587:
      title: Integer Default 0 Set 587
      value: 587
      valueType: Integer
      signature:
        v1: lzx/mW+dsOApsgsVoLnw8nVNW5X+aZEshuxwv/HUrw1m/oTsQY0aw+iVGrrEf4OCYsdPQpCzRL6ukyojZhDN1nJBQ5CiV5mlyuk/2BsAc+JczaaMTohE4nDXUKvrlIwnUfiZarUwXvlccRPeI3Lx2xGVuVi3o9uS0QLl6nn6OXRVaaH+oxUT6LgyrZykC9R63TkTUWCh55hWQMUfDfUjY9Pkzgr2Nay4N/7kpU0fvWxSOIrxZNpxQCNQgWRrMyqLF/QzWi0xA2z/QH58PV07b9CaDuS8RCp8mHWZiAZvuaoUPAKG8jtNSCMPjK2G452MFlkSneq6f52Mo3WW3E8v8A==
    string_default_empty_set_teststring:
      title: String Default Empty Set Teststring
      value: teststring
      valueType: String
      signature:
        v1: lnJNxGhgHxSrA2ImL+tdMKFl+vbKcZtq9eKPRhTFHEQntwD8AXi/6JpPPmb29w1kIPzf9FmbqofQbufuJjBBSr0n9IoViyduEgV7atBc9UxyWLObCyhltaxGWvyOjQpRQJEgpKMYEUV0Lsqjp0vdKN3+NSSoDJ+MXJY9xxsrNhMaeiXKGaU7zUYCrfwNfZrXk4pGvqjTOvRZGA1d0TitQXYYxDzpjryy02aWOIZDWBvDAOLXVPgJ5Gir67vBiAX9VEs49hyar6j0zYcRpVWexeDCSKTqlLZXzlthcNoW7NSEx9zPWhaJ4NziIDQCFg2G3OqJB3VJhWOWFwIJQ7YoZA==
  isNewKotsUiEnabled: true
  isSemverRequired: true
  signature: eyJsaWNlbnNlRGF0YSI6ImV5SmhjR2xXWlhKemFXOXVJam9pYTI5MGN5NXBieTkyTVdKbGRHRXhJaXdpYTJsdVpDSTZJa3hwWTJWdWMyVWlMQ0p0WlhSaFpHRjBZU0k2ZXlKdVlXMWxJam9pZEdWemRHTjFjM1J2YldWeUluMHNJbk53WldNaU9uc2liR2xqWlc1elpVbEVJam9pZEdWemRDMXNhV05sYm5ObExXbGtJaXdpYkdsalpXNXpaVlI1Y0dVaU9pSjBjbWxoYkNJc0ltTjFjM1J2YldWeVRtRnRaU0k2SWxSbGMzUWdRM1Z6ZEc5dFpYSWlMQ0poY0hCVGJIVm5Jam9pZEdWemRDMWhjSEFpTENKamFHRnVibVZzU1VRaU9pSXhJaXdpWTJoaGJtNWxiRTVoYldVaU9pSlRkR0ZpYkdVaUxDSmpkWE4wYjIxbGNrVnRZV2xzSWpvaWRHVnpkRUJsZUdGdGNHeGxMbU52YlNJc0ltTm9ZVzV1Wld4eklqcGJleUpqYUdGdWJtVnNTVVFpT2lJeElpd2lZMmhoYm01bGJGTnNkV2NpT2lKemRHRmliR1VpTENKamFHRnVibVZzVG1GdFpTSTZJbE4wWVdKc1pTSXNJbWx6UkdWbVlYVnNkQ0k2ZEhKMVpTd2lhWE5UWlcxMlpYSlNaWEYxYVhKbFpDSTZkSEoxWlgxZExDSmxiblJwZEd4bGJXVnVkSE1pT25zaVltOXZiRjlrWldaaGRXeDBYMlpoYkhObFgzTmxkRjkwY25WbElqcDdJblJwZEd4bElqb2lRbTl2YkdWaGJpQkVaV1poZFd4MElFWmhiSE5sSUZObGRDQlVjblZsSWl3aWRtRnNkV1VpT25SeWRXVXNJblpoYkhWbFZIbHdaU0k2SWtKdmIyeGxZVzRpTENKemFXZHVZWFIxY21VaU9uc2lkakVpT2lKbU5Ya3JSblV6VDJFNGVqWTRUakZNYkd0S2VFNHJaRVl2TTJkWk1pczBSRXg1Y0hSd1NHVTJTbFpPUVRBMFFta3ZNMGcwTkZwa2NEVklUMVpDT0hOaWNYWlBhbGh3TDNka2NsaHJNVXBqWlhZemQzVk5TbGhRT0ZSNk5EaFpVMDlEWkhwSlFqZDJaWFkzUTFsaWRHcDZaVkl3V1dGV1RuaDRTVUpzWjFkM1RuZHBVRUYwYUhKVVMzTkpWa3BLVVN0RWNVWk9ZbXBTWkhka1VIbGFiR1IxSzBReU5VVnVTaTl1VG5oeVV6QndhVmxMUTI5clEwRkVaakZJVWxCWFZHMVJSV1U1T1ZkRGF6SkdWR0pIWWxGUFZYUlBWVzlYZVVaM1drOHdSRWhCTmxWaE1WWk9SVUpXWjBVMldFbDBXbE5wZFZGNVpHWkpOVVVyVG01NGFEZHlNamgyZUhndkszZFZaMGRzUzBsS05HSXlWalJFUjJkQk9FNVNhbEJGYm1OaE5sSndSWFZ6T1U1bEswVkhORlYwVkZKME1IQkxaRnBqZEdWU1dFbHNMM2MxTnpjNFRUUmxVM1ZMVGxOV2RIQTRNall2T1ZFOVBTSjlmU3dpWlhod2FYSmxjMTloZENJNmV5SjBhWFJzWlNJNklrVjRjR2x5WVhScGIyNGlMQ0prWlhOamNtbHdkR2x2YmlJNklreHBZMlZ1YzJVZ1JYaHdhWEpoZEdsdmJpSXNJblpoYkhWbElqb2lNakF5TlMweE1pMHpNVlF5TXpvMU9UbzFPVm9pTENKMllXeDFaVlI1Y0dVaU9pSlRkSEpwYm1jaUxDSnphV2R1WVhSMWNtVWlPbnNpZGpFaU9pSmFNRmgxZEVrNGNtSXdlRlZvZFV0R1RtaEthR2R2Y0c1cFkybGtlbE5MTkZrMmVuZHVLMGNyU1RJcmEwVnVUVkpSZHpKNk9IYzNjV2hNV2tOMVNWbFVObTVuVTBacU5rOU5WVU5RWVM5Q1MzZG1ibkZJTW5sclMzaGxlRU42U0dsamRVUkVkVlF2VUVJck1EWnNWR0ZEYWpKTlUxZDVha2RPZUhWWUwxaGxabEEwWXk5d1QycFNOSGR4YjB0TlpHUktVV1JqY1ZkdU5IQnFiVEEzUjNReFYxSldWMFJIZG5Wc1RVVlFlbTlTVEhaWVV6QkhlR1YyZDJ4NWJFRlhSVEJ0Um05cWRHUnNiMkZIYWxoUGJXNTViM3A2ZEVKbFRuaGhSVkJTTVRKaVpFVm5Xa3h4U0d4VU16TlVPSFZXZDFWUVdYaENlVEZTVlhGWWRIWTRTbVp6TW1SelNrZHlkVTFVVm5kc2FVNWxhbXg2T1doaFRYVnRkbTExWjB3eGNsZ3lkVk5OY0ZOMGRXODNRa2xaTWxaNEwzbHJjaTlSYlU5emNXaExhbm96U25aSVVrUjVURlJxV1hWblYzbERiRWwwYVdwak4zYzlQU0o5ZlN3aWFXNTBYMlJsWm1GMWJIUmZNRjl6WlhSZk5UZzNJanA3SW5ScGRHeGxJam9pU1c1MFpXZGxjaUJFWldaaGRXeDBJREFnVTJWMElEVTROeUlzSW5aaGJIVmxJam8xT0Rjc0luWmhiSFZsVkhsd1pTSTZJa2x1ZEdWblpYSWlMQ0p6YVdkdVlYUjFjbVVpT25zaWRqRWlPaUpzZW5ndmJWY3JaSE5QUVhCelozTldiMHh1ZHpodVZrNVhOVmdyWVZwRmMyaDFlSGQyTDBoVmNuY3hiUzl2VkhOUldUQmhkeXRwVmtkeWNrVm1ORTlEV1hOa1VGRndRM3BTVERaMWEzbHZhbHBvUkU0eGJrcENVVFZEYVZZMWJXeDVkV3N2TWtKelFXTXJTbU42WVdGTlZHOW9SVFJ1UkZoVlMzWnliRWwzYmxWbWFWcGhjbFYzV0hac1kyTlNVR1ZKTTB4NE1uaEhWblZXYVROdk9YVlRNRkZNYkRadWJqWlBXRkpXWVdGSUsyOTRWVlEyVEdkNWNscDVhME01VWpZelZHdFVWVmREYURVMWFGZFJUVlZtUkdaVmFsazVVR3Q2WjNJeVRtRjVORTR2TjJ0d1ZUQm1kbGQ0VTA5SmNuaGFUbkI0VVVOT1VXZFhVbkpOZVhGTVJpOVJlbGRwTUhoQk1ub3ZVVWcxT0ZCV01EZGlPVU5oUkhWVE9GSkRjRGh0U0ZkYWFVRmFkblZoYjFWUVFVdEhPR3AwVGxORFRWQnFTekpITkRVeVRVWnNhMU51WlhFMlpqVXlUVzh6VjFjelJUaDJPRUU5UFNKOWZTd2ljM1J5YVc1blgyUmxabUYxYkhSZlpXMXdkSGxmYzJWMFgzUmxjM1J6ZEhKcGJtY2lPbnNpZEdsMGJHVWlPaUpUZEhKcGJtY2dSR1ZtWVhWc2RDQkZiWEIwZVNCVFpYUWdWR1Z6ZEhOMGNtbHVaeUlzSW5aaGJIVmxJam9pZEdWemRITjBjbWx1WnlJc0luWmhiSFZsVkhsd1pTSTZJbE4wY21sdVp5SXNJbk5wWjI1aGRIVnlaU0k2ZXlKMk1TSTZJbXh1U2s1NFIyaG5TSGhUY2tFeVNXMU1LM1JrVFV0R2JDdDJZa3RqV25SeE9XVkxVRkpvVkVaSVJWRnVkSGRFT0VGWWFTODJTbkJRVUcxaU1qbDNNV3RKVUhwbU9VWnRZbkZ2WmxGaWRXWjFTbXBDUWxOeU1HNDVTVzlXYVhsa2RVVm5WamRoZEVKak9WVjRlVmRNVDJKRGVXaHNkR0Y0UjFkMmVVOXFVWEJTVVVwRlozQkxUVmxGVlZZd1RITnhhbkF3ZG1STFRqTXJUbE5UYjBSS0swMVlTbGs1ZUhoemNrNW9UV0ZsYVZoTFIyRlZOM3BWV1VOeVpuZE9abHB5V0dzMGNFZDJjV3BVVDNaU1drZEJNV1F3VkdsMFVWaFpXWGhFZW5CcWNubDVNREpoVjA5SldrUlhRblpFUVU5TVdGWlFaMG8xUjJseU5qZDJRbWxCV0RsV1JYTTBPV2g1WVhJMmFqQjZXV05TY0ZaWFpYaGxSRU5UUzFSeGJFeGFXSHBzZEdoalRtOVhOMDVUUlhnNWVsQlhhR0ZLTkU1NmFVbEVVVU5HWnpKSE0wOXhTa0l6Vmtwb1YwOVhSbmRKU2xFM1dXOWFRVDA5SW4xOWZTd2lhWE5PWlhkTGIzUnpWV2xGYm1GaWJHVmtJanAwY25WbExDSnBjMU5sYlhabGNsSmxjWFZwY21Wa0lqcDBjblZsZlgwPSIsImlubmVyU2lnbmF0dXJlIjoiZXlKc2FXTmxibk5sVTJsbmJtRjBkWEpsSWpvaVNWRmxialZTWVZaYWNDOWlkSFp5WjBad1FVVTNVMGMxUVZWTFNuQnBUbWRTYkVkc09EUTNMeXN3VFU5Tll6UTVSa1ZpVXpsWVZHTlZMemhZZVN0Q2VsaHRkRkkyTjJ0NlREbFhMMDlxVnpWdU4zRTNaSFE1TTBSNVl5dGxVR2RhTTBkdGNtSTJiRWxLUjJoRVNGWkRRa0pGUlVsa1lXVnRNMU5UY1VaS05HRkpTVTF3VjBSUGNXNVpURUZvTUcwclJFbE5jSGxTYzNwWVVGaGxORloxVlVGWFNtSXpTMlZ0WTI5RVJ6QlRaVFZyVWtNMFNUQXliR00wUW5ob2FVSnpaVEE0ZVVGbGRHZFBWRlZ0WTBjcldsUTRLMWhLYm1OME9YUmllV1ZvYURKYVNrWndURlZoWm5GQk56WnhSelJNVlZkYVlsQXdkamhxV1c5aGJDOTJZbnBMV21WMlp6UnVSR1ZPVVc1Q2FrdEZWSGRpU2pOWWNHeG1iMVJWVGsxQmFGbzBWbkJYWWtka1VWVjRZVXBxTldOak0wSjRhbmhKUzFoVE1XSlhjekpPYzNKckwzSm1lV2QwVlhSeE5WUkpSVUpEUlVaQlBUMGlMQ0p3ZFdKc2FXTkxaWGtpT2lJdExTMHRMVUpGUjBsT0lGQlZRa3hKUXlCTFJWa3RMUzB0TFZ4dVRVbEpRa2xxUVU1Q1oydHhhR3RwUnpsM01FSkJVVVZHUVVGUFEwRlJPRUZOU1VsQ1EyZExRMEZSUlVGNVlraHpNVkEwTDFscVMxaElNbEZxWkVzdmRWeHVabE14VkZsUU5UZDFRVWRxWW1nM1JISTRWbmw1VTFveGVEQm1aRm94YURkamFHRnFlRWRyUnlzNFFXVldSMUl3ZHpKUFpuZFpkbloyTlRScVZsTnBTRnh1VlVKaGNXdGpUR0pGWjJGV05VdEZUVTV5U25WT2RtcG1kMEp6TUVGeWVYaE1XVE5ZTXpCWE4xcGtSU3RYUWpGWGNrWkxWbTV4UmxsYVRGRkJjMUpXUTF4dVIxWkVZM3BDT0hneFZYWlVjMHBUZWxBNFdrMWtWVFU1WjJwSmJXVkNaaXROZWtsWVVVZzVPSHBuZVZCbEwwcHRTbFJ6TjNCcFNVUkdWSFZ4UlZVNFpseHVSalJTUTJSWEt6UlViRkpKYmxWNE4wcDBWMVV5YlROa00yMTZjakZ2WTJkbVRFOUpOSG94ZUZKTU9XUk5jVlJVZW0xUlJHVlNSMU5CVmtobFNtRmpTbHh1UzB4Wk0yUmhRVXB2ZFdwT01XcDZSa2R5ZVRjd1NEWTVXRzVsV0ZCc0x6RTFkU3QzYUZOcE16VlROM05VVFVSUmEwWnhjazB6Y3psemFtOTNka3hSY2x4dU9WRkpSRUZSUVVKY2JpMHRMUzB0UlU1RUlGQlZRa3hKUXlCTFJWa3RMUzB0TFZ4dUlpd2lhMlY1VTJsbmJtRjBkWEpsSWpvaVpYbEtlbUZYWkhWWldGSXhZMjFWYVU5cFNqTmpSbWgyVVc1YU5sVlZUbHBUVjJSeldtNVpjbVJyU25kU2VYUkpXa1paZVZaRGRHNVdha3BwVW1rNVRGVnVWWGhNTWxseVkxZFpNRk50VFROTlJUVndZMnhPY0ZkcVozZGhWbU4yWTBaS1JXSXpWa2xhYkZZMVZHMXdVRk5HUmtsU1ZGcHFVVmQwTTAxSFJrZFVhazVTVkVaV1VtTkdRak5OYldRelMzcE9TbU5GWkZoTk1sSnVWbGhPY0U1R1VYWlBSMnhyUzNwV05XTkhaRzVXZWxKNlZGUmtSbUZXV2xSYWVscHlZVmRHTTA5R2JFVk5WbVI1VVcxU1dsRXlSbFpoVms1UVkyMU5lR1ZzU2xaUFZUbE5ZVlU1ZFdOVmJFOVRSM2hGV2pJd01sVlRkRVJWTWxacldsUk9WVTB5TVdwaGFrcFJWVVJOTTFacVVUUlRNbFpKWlZjeFIyRklTbEpPTTJ4RVVXMVdjV1F3VGpaU2FtUjJVV3RvUTFGVmVFdE9Nbmh3V1RKb1MyTllaRFpXVlhCWFZsaG9VR0l3V25wYVZGVnlWRWQwVGxkVlNuaFZSRVUwVWpOV1JXUlhUblZoTTFwRVpXNUtVVTVGYUU5a00yUkRWRVZLYldKWE5XaE9SMngzVTJrNWRWTnFUbmRpU0ZFeldsZE9iVlF4VWxoWGFsbzFVVmRLVEZGcVJrdGpWMUpxWW5wak1FNUVUbmxsUkVwMFZHMTRXbU5XWTNaV1YyZDJZa1pXWVdWdVVrbGlNR1I1WkRKMFNGbFlUa3hVUkVZd1RtcGpOVTR3UlRsUVUwbHpTVzFrYzJJeVNtaGlSWFJzWlZWc2EwbHFiMmxrUjFaNlpFTXhibUpIT1dsWlYzZDBZVEpXTlV4WGJHdEpiakE5SW4wPSJ9

`

	kotsscheme.AddToScheme(scheme.Scheme)

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, gvk, err := decode([]byte(data), nil, nil)
	require.NoError(t, err)

	assert.Equal(t, "kots.io", gvk.Group)
	assert.Equal(t, "v1beta1", gvk.Version)
	assert.Equal(t, "License", gvk.Kind)

	license := obj.(*kotsv1beta1.License)
	assert.Equal(t, "test-license-id", license.Spec.LicenseID)

	appKeys, err := license.ValidateLicense()
	require.NoError(t, err)

	// change a non-spec field and validate that things do not break
	license.ObjectMeta.Name = "changed"
	_, err = license.ValidateLicense()
	require.NoError(t, err, "changing non-spec fields should not invalidate the license")

	// change a spec field and validate that the signature is no longer valid
	license.Spec.AppSlug = "changed"
	_, err = license.ValidateLicense()
	require.Error(t, err, "changing spec fields should invalidate the license")

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

// TestProgrammaticEntitlementMarshaling tests that entitlement values are correctly
// preserved when marshaling and unmarshaling licenses created programmatically.
// This is a regression test ensuring programmatic license creation works as expected.
func TestProgrammaticEntitlementMarshaling(t *testing.T) {
	tests := []struct {
		name          string
		entitlement   kotsv1beta1.EntitlementField
		expectedValue interface{}
		valueType     string
	}{
		{
			name: "string entitlement",
			entitlement: kotsv1beta1.EntitlementField{
				Title:       "Test String Field",
				Description: "A test string entitlement",
				Value: kotsv1beta1.EntitlementValue{
					Type:   kotsv1beta1.String,
					StrVal: "test-value",
				},
				ValueType: "String",
			},
			expectedValue: "test-value",
			valueType:     "String",
		},
		{
			name: "integer entitlement",
			entitlement: kotsv1beta1.EntitlementField{
				Title:       "Test Integer Field",
				Description: "A test integer entitlement",
				Value: kotsv1beta1.EntitlementValue{
					Type:   kotsv1beta1.Int,
					IntVal: 42,
				},
				ValueType: "Integer",
			},
			expectedValue: int64(42),
			valueType:     "Integer",
		},
		{
			name: "boolean entitlement",
			entitlement: kotsv1beta1.EntitlementField{
				Title:       "Test Boolean Field",
				Description: "A test boolean entitlement",
				Value: kotsv1beta1.EntitlementValue{
					Type:    kotsv1beta1.Bool,
					BoolVal: true,
				},
				ValueType: "Boolean",
			},
			expectedValue: true,
			valueType:     "Boolean",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a license with a programmatically-created entitlement
			license := &kotsv1beta1.License{
				Spec: kotsv1beta1.LicenseSpec{
					LicenseID: "test-license",
					Entitlements: map[string]kotsv1beta1.EntitlementField{
						"test-field": tt.entitlement,
					},
				},
			}

			// Marshal to YAML using the Kubernetes serializer
			kotsscheme.AddToScheme(scheme.Scheme)
			s := scheme.Codecs.LegacyCodec(kotsv1beta1.SchemeGroupVersion)
			var buf bytes.Buffer
			err := s.Encode(license, &buf)
			require.NoError(t, err, "should be able to marshal programmatically-created license")
			yamlBytes := buf.Bytes()

			// Unmarshal back
			decode := scheme.Codecs.UniversalDeserializer().Decode
			obj, _, err := decode(yamlBytes, nil, nil)
			require.NoError(t, err, "should be able to unmarshal back")

			unmarshaled := obj.(*kotsv1beta1.License)
			field, ok := unmarshaled.Spec.Entitlements["test-field"]
			require.True(t, ok, "entitlement should exist after round-trip")

			// Verify the value was preserved
			assert.Equal(t, tt.entitlement.Title, field.Title, "title should be preserved")
			assert.Equal(t, tt.entitlement.Description, field.Description, "description should be preserved")
			assert.Equal(t, tt.entitlement.ValueType, field.ValueType, "value type should be preserved")

			// Verify the actual value matches
			actualValue := field.Value.Value()
			assert.Equal(t, tt.expectedValue, actualValue, "value should be preserved through marshal/unmarshal")

			// Verify the type is correct
			assert.Equal(t, tt.entitlement.Value.Type, field.Value.Type, "entitlement type should be preserved")
		})
	}
}
