package v1beta2_test

import (
	"encoding/json"
	"testing"

	v1beta2 "github.com/replicatedhq/kotskinds/apis/kots/v1beta2"
	kotsscheme "github.com/replicatedhq/kotskinds/client/kotsclientset/scheme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

// Test_OptionalValue_YAMLDecode_StringOrMap proves the string-or-map tolerance on
// optionalValues[].values works for YAML input, not just JSON — and that YAML and JSON decode to
// the same thing.
//
// There is deliberately NO UnmarshalYAML on OptionalValue. Every real decode path converts YAML to
// JSON and then calls json.Unmarshal: sigs.k8s.io/yaml does it directly, and the k8s
// UniversalDeserializer (the path kots and vandoor actually use) routes YAML through its JSON/YAML
// serializer. So the custom UnmarshalJSON fires either way. This matches the repo's existing
// dual-type convention (multitype.BoolOrString, MappedChartValue implement JSON hooks only). Adding
// UnmarshalYAML would be redundant. See sc-138613 / itrs-replicated#380.
//
// This is an external test package (v1beta2_test) on purpose: it decodes via
// client/kotsclientset/scheme, which imports this API package, so an in-package test would be an
// import cycle. It only needs the exported ValuesString()/Values, so external is fine.
func Test_OptionalValue_YAMLDecode_StringOrMap(t *testing.T) {
	tests := []struct {
		name       string
		yamlDoc    string
		jsonDoc    string
		wantString bool
	}{
		{
			name: "templated repl{{ }} scalar",
			yamlDoc: `apiVersion: kots.io/v1beta2
kind: HelmChart
metadata:
  name: issue380
spec:
  chart:
    name: issue380
    chartVersion: 1.0.0
  optionalValues:
    - when: "true"
      values: 'repl{{ ConfigOption "override" | nindent 8 }}'`,
			jsonDoc:    `{"apiVersion":"kots.io/v1beta2","kind":"HelmChart","metadata":{"name":"issue380"},"spec":{"chart":{"name":"issue380","chartVersion":"1.0.0"},"optionalValues":[{"when":"true","values":"repl{{ ConfigOption \"override\" | nindent 8 }}"}]}}`,
			wantString: true,
		},
		{
			name: "{{repl }} alt-delimiter scalar",
			yamlDoc: `apiVersion: kots.io/v1beta2
kind: HelmChart
metadata:
  name: issue380
spec:
  chart:
    name: issue380
    chartVersion: 1.0.0
  optionalValues:
    - when: "true"
      values: '{{repl ConfigOption "override" | nindent 8 }}'`,
			jsonDoc:    `{"apiVersion":"kots.io/v1beta2","kind":"HelmChart","metadata":{"name":"issue380"},"spec":{"chart":{"name":"issue380","chartVersion":"1.0.0"},"optionalValues":[{"when":"true","values":"{{repl ConfigOption \"override\" | nindent 8 }}"}]}}`,
			wantString: true,
		},
		{
			name: "map values",
			yamlDoc: `apiVersion: kots.io/v1beta2
kind: HelmChart
metadata:
  name: issue380
spec:
  chart:
    name: issue380
    chartVersion: 1.0.0
  optionalValues:
    - when: "true"
      values:
        operator:
          enabled: true`,
			jsonDoc:    `{"apiVersion":"kots.io/v1beta2","kind":"HelmChart","metadata":{"name":"issue380"},"spec":{"chart":{"name":"issue380","chartVersion":"1.0.0"},"optionalValues":[{"when":"true","values":{"operator":{"enabled":true}}}]}}`,
			wantString: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. sigs.k8s.io/yaml: YAML -> JSON -> json.Unmarshal (fires UnmarshalJSON).
			var viaYAML v1beta2.HelmChart
			require.NoError(t, yaml.Unmarshal([]byte(tt.yamlDoc), &viaYAML))

			// 2. The k8s UniversalDeserializer — how kots/vandoor actually decode manifests.
			obj, _, err := kotsscheme.Codecs.UniversalDeserializer().Decode([]byte(tt.yamlDoc), nil, nil)
			require.NoError(t, err)
			viaScheme, ok := obj.(*v1beta2.HelmChart)
			require.True(t, ok, "decoded object must be a *v1beta2.HelmChart")

			// 3. Direct JSON decode, for the equivalence assertion.
			var viaJSON v1beta2.HelmChart
			require.NoError(t, json.Unmarshal([]byte(tt.jsonDoc), &viaJSON))

			for _, hc := range []*v1beta2.HelmChart{&viaYAML, viaScheme, &viaJSON} {
				require.Len(t, hc.Spec.OptionalValues, 1)
				ov := hc.Spec.OptionalValues[0]
				if tt.wantString {
					assert.Nil(t, ov.Values, "a scalar values must leave the typed map empty")
					assert.NotEmpty(t, ov.ValuesString(), "a scalar values must be preserved as a string")
				} else {
					assert.NotNil(t, ov.Values, "a map values must populate the typed map")
					assert.Empty(t, ov.ValuesString())
				}
			}

			// YAML (both routes) must decode equivalently to JSON.
			assert.Equal(t, viaJSON.Spec.OptionalValues, viaYAML.Spec.OptionalValues, "sigs.k8s.io/yaml decode must match JSON decode")
			assert.Equal(t, viaJSON.Spec.OptionalValues, viaScheme.Spec.OptionalValues, "UniversalDeserializer decode must match JSON decode")
		})
	}
}
