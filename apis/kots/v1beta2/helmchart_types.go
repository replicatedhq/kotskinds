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

package v1beta2

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"github.com/replicatedhq/kotskinds/multitype"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:validation:Type=""
type MappedChartValue struct {
	Value string `json:"-"`

	valueType string `json:"-"`

	strValue   string  `json:"-"`
	boolValue  bool    `json:"-"`
	floatValue float64 `json:"-"`

	children map[string]*MappedChartValue `json:"-"`
	array    []*MappedChartValue          `json:"-"`
}

func (m MappedChartValue) MarshalJSON() ([]byte, error) {
	val, err := m.getBuiltValue()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get built value")
	}

	return json.Marshal(val)
}

func (m *MappedChartValue) getBuiltValue() (interface{}, error) {
	if m.valueType == "string" {
		return m.strValue, nil
	}
	if m.valueType == "bool" {
		return m.boolValue, nil
	}
	if m.valueType == "float" {
		return m.floatValue, nil
	}
	if m.valueType == "nil" {
		return nil, nil
	}

	if m.valueType == "children" {
		children := map[string]interface{}{}
		for k, v := range m.children {
			childValue, err := v.getBuiltValue()
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get value of child %s", k)
			}
			children[k] = childValue
		}
		return children, nil
	}
	if m.valueType == "array" {
		var elements []interface{}
		for i, v := range m.array {
			elValue, err := v.getBuiltValue()
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get value of child %d", i)
			}
			elements = append(elements, elValue)
		}
		return elements, nil
	}

	return nil, errors.New("unknown value type")
}

func (m *MappedChartValue) UnmarshalJSON(value []byte) error {
	var b interface{}
	err := json.Unmarshal(value, &b)
	if err != nil {
		return err
	}

	if b == nil {
		m.valueType = "nil"
		return nil
	}

	if b, ok := b.(string); ok {
		m.strValue = b
		m.valueType = "string"
		return nil
	}

	if b, ok := b.(bool); ok {
		m.boolValue = b
		m.valueType = "bool"
		return nil
	}

	if b, ok := b.(float64); ok {
		m.floatValue = b
		m.valueType = "float"
		return nil
	}

	if b, ok := b.(map[string]interface{}); ok {
		m.children = make(map[string]*MappedChartValue)
		for k, v := range b {
			vv, err := json.Marshal(v)
			if err != nil {
				return err
			}

			m2 := &MappedChartValue{}
			if err := m2.UnmarshalJSON(vv); err != nil {
				return err
			}

			m.children[k] = m2
		}

		m.valueType = "children"

		return nil
	}

	if b, ok := b.([]interface{}); ok {
		m.array = []*MappedChartValue{}
		for _, v := range b {
			vv, err := json.Marshal(v)
			if err != nil {
				return err
			}

			m2 := &MappedChartValue{}
			if err := m2.UnmarshalJSON(vv); err != nil {
				return err
			}

			m.array = append(m.array, m2)
		}

		m.valueType = "array"

		return nil
	}

	return errors.Errorf("unknown mapped chart value type: %T", b)
}

type ChartIdentifier struct {
	Name         string `json:"name"`
	ChartVersion string `json:"chartVersion"`
}

func (h *HelmChartSpec) GetHelmValues(values map[string]MappedChartValue) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	for k, v := range values {
		value, err := h.renderValue(&v)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to render value at %s", k)
		}

		result[k] = value
	}

	return result, nil
}

func (h *HelmChartSpec) GetReplTmplValues(values map[string]MappedChartValue) (map[string]interface{}, error) {
	newValues := make(map[string]interface{})

	for k, v := range values {
		value, err := h.getReplTmplValue(&v)
		if err != nil || value == nil {
			continue
		}
		newValues[k] = value
	}

	return newValues, nil
}

func (h *HelmChartSpec) getReplTmplValue(value *MappedChartValue) (interface{}, error) {
	if value.valueType == "children" {
		result := map[string]interface{}{}
		for k, v := range value.children {
			built, err := h.getReplTmplValue(v)
			if err != nil || built == nil {
				continue
			}
			result[k] = built
		}
		if len(result) == 0 {
			return nil, nil
		}
		return result, nil
	} else if value.valueType == "array" {
		result := []interface{}{}
		for _, v := range value.array {
			built, err := h.getReplTmplValue(v)
			if err != nil {
				return nil, errors.Wrap(err, "failed to render array value")
			}
			result = append(result, built)
		}
		return result, nil
	} else {
		built, err := value.getBuiltValue()
		if err != nil {
			return nil, errors.Wrap(err, "failed to build value")
		}
		str, ok := built.(string)
		if ok && (strings.Contains(str, "repl{{") || strings.Contains(str, "{{repl")) {
			return built, nil
		}
		return nil, errors.New("value is not string or not repl tmpl function")
	}
}

func GetMapIntersect(m1, m2 map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			continue
		}

		v1map, v1mapOK := v1.(map[string]interface{})
		v2map, v2mapOK := v2.(map[string]interface{})
		if v1mapOK && v2mapOK {
			res[k] = GetMapIntersect(v1map, v2map)
			continue
		}

		res[k] = v2
	}

	return res
}

func MergeHelmChartValues(baseValues map[string]MappedChartValue,
	overlayValues map[string]MappedChartValue) map[string]MappedChartValue {

	result := map[string]MappedChartValue{}
	for k, v := range baseValues {
		if _, exists := overlayValues[k]; !exists {
			result[k] = baseValues[k]
			continue
		}
		if v.valueType != "children" {
			result[k] = overlayValues[k]
		} else {
			result[k] = MappedChartValue{
				valueType: "children",
				children:  mergeValueChildren(v.children, overlayValues[k].children),
			}
		}
	}
	for k, v := range overlayValues {
		if _, exists := baseValues[k]; !exists {
			result[k] = v
		}
	}
	return result
}

func mergeValueChildren(baseValues map[string]*MappedChartValue, overlayValues map[string]*MappedChartValue) map[string]*MappedChartValue {
	result := map[string]*MappedChartValue{}
	for k, v := range baseValues {
		if _, exists := overlayValues[k]; !exists {
			result[k] = baseValues[k]
			continue
		}
		if v.valueType != "children" {
			result[k] = overlayValues[k]
		} else {
			result[k] = &MappedChartValue{
				valueType: "children",
				children:  mergeValueChildren(v.children, overlayValues[k].children),
			}
		}
	}

	for k, v := range overlayValues {
		if _, exists := baseValues[k]; !exists {
			result[k] = v
		}
	}
	return result
}

func (h *HelmChartSpec) renderValue(value *MappedChartValue) (interface{}, error) {
	if value.valueType == "children" {
		result := map[string]interface{}{}
		for k, v := range value.children {
			built, err := h.renderValue(v)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to render child value at key %s", k)
			}
			result[k] = built
		}
		return result, nil
	} else if value.valueType == "array" {
		result := []interface{}{}
		for _, v := range value.array {
			built, err := h.renderValue(v)
			if err != nil {
				return nil, errors.Wrap(err, "failed to render array value")
			}
			result = append(result, built)
		}
		return result, nil
	} else {
		built, err := value.getBuiltValue()
		if err != nil {
			return nil, errors.Wrap(err, "failed to build value")
		}
		return built, nil
	}
}

func (h *HelmChart) GetDirName() string {
	return h.GetReleaseName()
}

func (h *HelmChart) GetReleaseName() string {
	if h.Spec.ReleaseName != "" {
		return h.Spec.ReleaseName
	}
	return h.Spec.Chart.Name
}

func (h *HelmChart) GetNamespace() string {
	return h.Spec.Namespace
}

func (h *HelmChart) GetChartName() string {
	return h.Spec.Chart.Name
}

func (h *HelmChart) GetChartVersion() string {
	return h.Spec.Chart.ChartVersion
}

func (h *HelmChart) GetAPIVersion() string {
	return h.APIVersion
}

func (h *HelmChart) GetUpgradeFlags() []string {
	return h.Spec.HelmUpgradeFlags
}

func (h *HelmChart) GetWeight() int64 {
	return h.Spec.Weight
}

func (h *HelmChart) GetHelmVersion() string {
	return "v3" // v3 is the only supported version for v1beta2
}

func (h *HelmChart) GetBuilderValues() (map[string]interface{}, error) {
	return h.Spec.GetHelmValues(h.Spec.Builder)
}

func (h *HelmChart) SetChartNamespace(namespace string) {
	h.Spec.Namespace = namespace
}

type OptionalValue struct {
	When           string `json:"when"`
	RecursiveMerge bool   `json:"recursiveMerge"`

	// Values holds the map form of the optional values. This is the form KOTS decodes at install
	// time, because templates are rendered before the manifest is unmarshalled.
	Values map[string]MappedChartValue `json:"values,omitempty"`

	// valuesString holds the raw scalar form of values when it is provided as a string, e.g. a
	// `repl{{ ConfigOption "x" | nindent 8 }}` template that renders to a map at install time. The
	// Vendor Portal validates specs before template rendering, so it must tolerate this form
	// instead of dropping the whole HelmChart CR. Exactly one of Values / valuesString is ever set.
	valuesString string `json:"-"`
}

// optionalValueJSON mirrors OptionalValue's serialized shape but with values left as a raw message
// so UnmarshalJSON can decode either a map or a scalar string without recursing into OptionalValue.
type optionalValueJSON struct {
	When           string          `json:"when"`
	RecursiveMerge bool            `json:"recursiveMerge"`
	Values         json.RawMessage `json:"values,omitempty"`
}

// UnmarshalJSON tolerates spec.optionalValues[].values being either a map (the rendered form) or a
// scalar string (a repl template the Vendor Portal sees before rendering). A scalar is stored as a
// string and leaves Values nil (type-indeterminate) so the CR still decodes and is not dropped.
func (o *OptionalValue) UnmarshalJSON(data []byte) error {
	var raw optionalValueJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	o.When = raw.When
	o.RecursiveMerge = raw.RecursiveMerge
	o.Values = nil
	o.valuesString = ""

	trimmed := strings.TrimSpace(string(raw.Values))
	if trimmed == "" || trimmed == "null" {
		return nil
	}

	// Prefer the map form. Fall back to a scalar string (e.g. a repl template) rather than failing
	// the decode and dropping the chart.
	var values map[string]MappedChartValue
	if err := json.Unmarshal(raw.Values, &values); err == nil {
		o.Values = values
		return nil
	}

	var str string
	if err := json.Unmarshal(raw.Values, &str); err == nil {
		o.valuesString = str
		return nil
	}

	return errors.Errorf("optionalValues[].values must be a map or a string, got %q", trimmed)
}

// MarshalJSON re-emits values in whichever form it was decoded from, so a round-trip is lossless.
func (o OptionalValue) MarshalJSON() ([]byte, error) {
	raw := optionalValueJSON{
		When:           o.When,
		RecursiveMerge: o.RecursiveMerge,
	}

	switch {
	case o.valuesString != "":
		encoded, err := json.Marshal(o.valuesString)
		if err != nil {
			return nil, err
		}
		raw.Values = encoded
	case o.Values != nil:
		encoded, err := json.Marshal(o.Values)
		if err != nil {
			return nil, err
		}
		raw.Values = encoded
	}

	return json.Marshal(raw)
}

// ValuesString returns the raw scalar form of optionalValues[].values (a repl template) when values
// was provided as a string rather than a map; it returns "" for the map form.
func (o *OptionalValue) ValuesString() string {
	return o.valuesString
}

// HelmChartSpec defines the desired state of HelmChartSpec
type HelmChartSpec struct {
	Chart            ChartIdentifier             `json:"chart"`
	ReleaseName      string                      `json:"releaseName,omitempty"`
	Exclude          multitype.BoolOrString      `json:"exclude,omitempty"`
	Namespace        string                      `json:"namespace,omitempty"`
	Values           map[string]MappedChartValue `json:"values,omitempty"`
	OptionalValues   []*OptionalValue            `json:"optionalValues,omitempty"`
	Builder          map[string]MappedChartValue `json:"builder,omitempty"`
	Weight           int64                       `json:"weight,omitempty"`
	HelmUpgradeFlags []string                    `json:"helmUpgradeFlags,omitempty"`
	Docs             map[string]string           `json:"docs,omitempty"`
}

// HelmChartStatus defines the observed state of HelmChart
type HelmChartStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// HelmChart is the Schema for the helmchart API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
type HelmChart struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelmChartSpec   `json:"spec,omitempty"`
	Status HelmChartStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HelmChartList contains a list of HelmCharts
type HelmChartList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelmChart `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HelmChart{}, &HelmChartList{})
}
