package v1beta1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_OptionalValue_ValuesStringOrMap verifies that spec.optionalValues[].values decodes from
// BOTH a YAML mapping (the form KOTS sees after template rendering) and a scalar string (a
// `repl{{ ... }}` template the Vendor Portal sees before rendering). The scalar form must NOT
// error, otherwise the whole HelmChart CR is dropped and the release is falsely flagged "not
// installable with KOTS". See sc-138613 / itrs-replicated#380.
func Test_OptionalValue_ValuesStringOrMap(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantString bool
	}{
		{
			name:       "map form",
			input:      `{"when":"true","recursiveMerge":true,"values":{"operator":{"enabled":true}}}`,
			wantString: false,
		},
		{
			name:       "templated repl{{ }} scalar is tolerated",
			input:      `{"when":"true","recursiveMerge":false,"values":"repl{{ ConfigOption \"override\" | nindent 8 }}"}`,
			wantString: true,
		},
		{
			name:       "{{repl }} alt-delimiter scalar is tolerated",
			input:      `{"when":"true","values":"{{repl ConfigOption \"override\" | nindent 8 }}"}`,
			wantString: true,
		},
		{
			name:       "non-templated scalar is tolerated, not an error",
			input:      `{"when":"true","values":"this is not a map"}`,
			wantString: true,
		},
		{
			name:       "omitted values",
			input:      `{"when":"false"}`,
			wantString: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ov OptionalValue
			require.NoError(t, json.Unmarshal([]byte(tt.input), &ov))

			if tt.wantString {
				assert.Nil(t, ov.Values, "a scalar values must leave the typed map empty (type-indeterminate)")
				assert.NotEmpty(t, ov.valuesString, "a scalar values must be preserved as a string")
			} else {
				assert.Empty(t, ov.valuesString)
			}

			// A marshal/unmarshal round-trip must be idempotent: a string stays a string, a map
			// stays a map.
			out, err := json.Marshal(ov)
			require.NoError(t, err)

			var rt OptionalValue
			require.NoError(t, json.Unmarshal(out, &rt))
			assert.Equal(t, ov, rt, "OptionalValue must survive a marshal/unmarshal round-trip")
		})
	}
}

// Test_HelmChartSpec_ToleratesTemplatedScalarOptionalValues proves the bug fix at the spec level:
// a HelmChartSpec whose optionalValues[].values is a templated scalar must decode without error so
// the CR is not dropped from the installability check.
func Test_HelmChartSpec_ToleratesTemplatedScalarOptionalValues(t *testing.T) {
	spec := `{"chart":{"name":"issue380","chartVersion":"1.0.0"},"optionalValues":[{"when":"true","recursiveMerge":true,"values":"repl{{ ConfigOption \"override\" | nindent 8 }}"}]}`

	var s HelmChartSpec
	require.NoError(t, json.Unmarshal([]byte(spec), &s))

	assert.Equal(t, "issue380", s.Chart.Name)
	require.Len(t, s.OptionalValues, 1)
	assert.Nil(t, s.OptionalValues[0].Values, "templated scalar leaves the typed map empty")
}

// Test_OptionalValue_DeepCopyPreservesValuesString locks the invariant that a scalar values
// survives DeepCopy/DeepCopyInto. Correctness holds today only because valuesString is a value-type
// string copied by the generated `*out = *in`; this test fails loudly if that ever stops being true
// (e.g. the field becomes a pointer/slice and the generated deepcopy stops covering it).
func Test_OptionalValue_DeepCopyPreservesValuesString(t *testing.T) {
	var ov OptionalValue
	require.NoError(t, json.Unmarshal([]byte(`{"when":"true","values":"repl{{ ConfigOption \"override\" | nindent 8 }}"}`), &ov))
	require.NotEmpty(t, ov.ValuesString())

	cp := ov.DeepCopy()
	assert.Equal(t, ov.ValuesString(), cp.ValuesString(), "scalar must survive DeepCopy")

	var into OptionalValue
	ov.DeepCopyInto(&into)
	assert.Equal(t, ov.ValuesString(), into.ValuesString(), "scalar must survive DeepCopyInto")

	orig, err := json.Marshal(ov)
	require.NoError(t, err)
	copied, err := json.Marshal(cp)
	require.NoError(t, err)
	assert.JSONEq(t, string(orig), string(copied), "a deep copy must re-marshal identically")
}
