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

package v1beta1

import (
	"crypto"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	kotscrypto "github.com/replicatedhq/kotskinds/pkg/crypto"
	"github.com/replicatedhq/kotskinds/pkg/licensewrapper/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Int Type = iota
	String
	Bool
)

type Type int

type EntitlementValue struct {
	Type    Type   `json:"-"`
	IntVal  int64  `json:"-"`
	StrVal  string `json:"-"`
	BoolVal bool   `json:"-"`
	// notSet tracks whether the value key was missing in YAML/JSON during unmarshaling.
	// This distinguishes between missing value (should be nil) vs explicit zero value.
	// Defaults to false, set to true when value key is absent.
	notSet bool `json:"-"`
}

func (entitlementValue *EntitlementValue) Value() interface{} {
	// When value field was not present in YAML (notSet=true), return nil
	// to match the pre-consolidation behavior. This produces "<nil>" via fmt.Sprint()
	// and allows us to verify the signature
	if entitlementValue.notSet {
		return nil
	}

	if entitlementValue.Type == Int {
		return entitlementValue.IntVal
	} else if entitlementValue.Type == Bool {
		return entitlementValue.BoolVal
	}

	return entitlementValue.StrVal
}

func (entitlementValue EntitlementValue) MarshalJSON() ([]byte, error) {
	switch entitlementValue.Type {
	case Int:
		return json.Marshal(entitlementValue.IntVal)

	case String:
		return json.Marshal(entitlementValue.StrVal)

	case Bool:
		return json.Marshal(entitlementValue.BoolVal)

	default:
		return []byte{}, fmt.Errorf("impossible EntitlementValue.Type")
	}
}

type EntitlementField struct {
	Title       string                     `json:"title,omitempty"`
	Description string                     `json:"description,omitempty"`
	Value       EntitlementValue           `json:"-"`
	ValueRaw    json.RawMessage            `json:"value,omitempty"`
	ValueType   string                     `json:"valueType,omitempty"`
	IsHidden    bool                       `json:"isHidden,omitempty"`
	Signature   *EntitlementFieldSignature `json:"signature,omitempty"`
}

func (ef *EntitlementField) UnmarshalJSON(data []byte) error {
	// Define a type alias to prevent infinite recursion
	type Alias EntitlementField

	// Parse into a map to check for "value" key presence
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Unmarshal all fields using default behavior
	// This works because EntitlementValue has json:"-" on all fields,
	// so the "value" field will be skipped by the default unmarshaler
	if err := json.Unmarshal(data, (*Alias)(ef)); err != nil {
		return err
	}

	// Handle the Value field specially
	if rawValue, hasValue := raw["value"]; hasValue {
		if err := unmarshalEntitlementValue(&ef.Value, rawValue); err != nil {
			return errors.Wrap(err, "failed to unmarshal entitlement value")
		}
		// notSet stays false (default) since value was present
	} else {
		// Value key is missing, mark it as not set
		ef.Value.notSet = true
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for EntitlementField.
// This ensures that entitlement values are correctly marshaled whether the
// license was loaded from YAML/JSON or created programmatically in code.
func (ef EntitlementField) MarshalJSON() ([]byte, error) {
	// Define a type alias to prevent infinite recursion
	type Alias EntitlementField

	// If the entitlement was created programmatically (ValueRaw not populated),
	// marshal the Value field to ensure the "value" key appears in output
	if len(ef.ValueRaw) == 0 && !ef.Value.notSet {
		valueBytes, err := json.Marshal(ef.Value)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal entitlement value")
		}
		ef.ValueRaw = valueBytes
	}

	// Use the alias type to marshal all fields with standard behavior
	// The Signature field is now a pointer, so omitempty will work correctly
	return json.Marshal((Alias)(ef))
}

// unmarshalEntitlementValue manually unmarshals a value into an EntitlementValue
func unmarshalEntitlementValue(ev *EntitlementValue, data []byte) error {
	// Try to detect the type and unmarshal accordingly
	if len(data) > 0 && data[0] == '"' {
		// String value
		ev.Type = String
		return json.Unmarshal(data, &ev.StrVal)
	}

	// Try integer
	intValue, err := strconv.ParseInt(string(data), 10, 64)
	if err == nil {
		ev.Type = Int
		ev.IntVal = intValue
		return nil
	}

	// Try boolean
	boolValue, err := strconv.ParseBool(string(data))
	if err == nil {
		ev.Type = Bool
		ev.BoolVal = boolValue
		return nil
	}

	return errors.New("unknown license value type")
}

type EntitlementFieldSignature struct {
	V1 []byte `json:"v1,omitempty"`
}

type Channel struct {
	ChannelID             string `json:"channelID"`
	ChannelName           string `json:"channelName,omitempty"`
	ChannelSlug           string `json:"channelSlug,omitempty"`
	IsDefault             bool   `json:"isDefault,omitempty"`
	Endpoint              string `json:"endpoint,omitempty"`
	ReplicatedProxyDomain string `json:"replicatedProxyDomain,omitempty"`
	IsSemverRequired      bool   `json:"isSemverRequired,omitempty"`
}

// LicenseSpec defines the desired state of LicenseSpec
type LicenseSpec struct {
	Signature                         []byte                      `json:"signature"`
	AppSlug                           string                      `json:"appSlug"`
	Endpoint                          string                      `json:"endpoint,omitempty"`
	ReplicatedProxyDomain             string                      `json:"replicatedProxyDomain,omitempty"`
	CustomerID                        string                      `json:"customerID,omitempty"`
	CustomerName                      string                      `json:"customerName,omitempty"`
	CustomerEmail                     string                      `json:"customerEmail,omitempty"`
	ChannelID                         string                      `json:"channelID,omitempty"`
	ChannelName                       string                      `json:"channelName,omitempty"`
	Channels                          []Channel                   `json:"channels,omitempty"`
	LicenseSequence                   int64                       `json:"licenseSequence,omitempty"`
	LicenseID                         string                      `json:"licenseID"`
	LicenseType                       string                      `json:"licenseType,omitempty"`
	IsAirgapSupported                 bool                        `json:"isAirgapSupported,omitempty"`
	IsGitOpsSupported                 bool                        `json:"isGitOpsSupported,omitempty"`
	IsIdentityServiceSupported        bool                        `json:"isIdentityServiceSupported,omitempty"`
	IsGeoaxisSupported                bool                        `json:"isGeoaxisSupported,omitempty"`
	IsSnapshotSupported               bool                        `json:"isSnapshotSupported,omitempty"`
	IsDisasterRecoverySupported       bool                        `json:"isDisasterRecoverySupported,omitempty"`
	IsSupportBundleUploadSupported    bool                        `json:"isSupportBundleUploadSupported,omitempty"`
	IsSemverRequired                  bool                        `json:"isSemverRequired,omitempty"`
	IsEmbeddedClusterDownloadEnabled  bool                        `json:"isEmbeddedClusterDownloadEnabled,omitempty"`
	IsEmbeddedClusterMultiNodeEnabled bool                        `json:"isEmbeddedClusterMultiNodeEnabled,omitempty"`
	Entitlements                      map[string]EntitlementField `json:"entitlements,omitempty"`
}

// LicenseStatus defines the observed state of License
type LicenseStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// License is the Schema for the license API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
type License struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LicenseSpec   `json:"spec,omitempty"`
	Status LicenseStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LicenseList contains a list of Licenses
type LicenseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []License `json:"items"`
}

func init() {
	SchemeBuilder.Register(&License{}, &LicenseList{})
}

// ValidateLicense validates the entire v1beta1 license signature and all entitlement signatures
// Returns the app signing keys on success (used for validating entitlement signatures), or an error if validation fails
// if the plaintext license data is different than the signed data, the license object will be modified to ensure all fields are from the signed data, and an error will still be returned.
func (l *License) ValidateLicense() (*kotscrypto.AppSigningKeys, error) {
	// Decode and parse the signature
	outerSig, innerSig, appKeys, err := kotscrypto.DecodeLicenseSignature(l.Spec.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode signature")
	}

	if len(innerSig.KeySignature) == 0 {
		return nil, errors.New("v1 key signature not found")
	}

	var keySig kotscrypto.KeySignature
	if err := json.Unmarshal(innerSig.KeySignature, &keySig); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal key signature")
	}

	globalPubKey, err := kotscrypto.FindGlobalPublicKeyRSA(keySig.GlobalKeyID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find global public key")
	}

	// verify that the application public key is signed by the global key
	if err := kotscrypto.VerifySignatureRSA([]byte(innerSig.PublicKey), keySig.Signature, globalPubKey, crypto.MD5); err != nil {
		return nil, errors.Wrap(err, "v1 key signature verification failed")
	}

	if len(innerSig.LicenseSignature) == 0 {
		return nil, errors.New("v1 license signature not found")
	}

	// verify that the license data is signed by the application key
	if err := kotscrypto.VerifySignatureRSA(outerSig.LicenseData, innerSig.LicenseSignature, innerSig.PublicKey, crypto.MD5); err != nil {
		return nil, errors.Wrap(err, "v1 license signature verification failed")
	}

	// validate that each entitlement value is signed by the application key
	for fieldName, field := range l.Spec.Entitlements {
		// the entitlement values are still covered as part of the entire license body, and some old license files did not include entitlement signatures
		if field.Signature == nil || len(field.Signature.V1) == 0 {
			continue
		}
		if err := field.ValidateSignature(appKeys); err != nil {
			return nil, errors.Wrapf(err, "entitlement validation failed for field: %s", fieldName)
		}
	}

	// Validate license data matches what was signed
	if err := l.compareLicenseData(outerSig.LicenseData); err != nil {
		return nil, errors.Wrap(err, "license data validation failed")
	}

	return appKeys, nil
}

// ValidateSignature validates a single v1beta1 entitlement signature
func (in *EntitlementField) ValidateSignature(appKeys *kotscrypto.AppSigningKeys) error {
	// Check if signature is present
	if in.Signature == nil || len(in.Signature.V1) == 0 {
		return errors.New("v1 signature not found for entitlement")
	}

	// Get the value as a string
	value := in.Value.Value()
	message := []byte(fmt.Sprint(value))

	// Verify the signature using the crypto package
	if err := kotscrypto.VerifySignatureWithKeyRSA(message, in.Signature.V1, appKeys.PublicKeyRSA, crypto.MD5); err != nil {
		return errors.Wrap(err, "v1 entitlement signature verification failed")
	}

	return nil
}

// compareLicenseData decodes the provided license json (from a signature) and compares it to the license data in the calling license
// this function modifies the calling license object to ensure all fields are from the signed data, not the outer unsigned yaml
func (l *License) compareLicenseData(signedJSON []byte) error {
	// Decode the license JSON
	var signedData License
	if err := json.Unmarshal(signedJSON, &signedData); err != nil {
		return errors.Wrap(err, "failed to unmarshal signed license data")
	}

	// ensure that the signature object is updated to the signed data
	originalSignature := l.Spec.Signature
	outerData := l.Spec.DeepCopy()
	l.Spec = signedData.Spec
	l.Spec.Signature = originalSignature

	// Compare each field in outerData against the decoded license data
	if outerData.AppSlug != signedData.Spec.AppSlug {
		return &types.LicenseDataValidationError{
			FieldName:   "appSlug",
			SignedValue: signedData.Spec.AppSlug,
			ActualValue: outerData.AppSlug,
		}
	}
	if outerData.Endpoint != signedData.Spec.Endpoint {
		return &types.LicenseDataValidationError{
			FieldName:   "endpoint",
			SignedValue: signedData.Spec.Endpoint,
			ActualValue: outerData.Endpoint,
		}
	}
	if outerData.ReplicatedProxyDomain != signedData.Spec.ReplicatedProxyDomain {
		return &types.LicenseDataValidationError{
			FieldName:   "replicatedProxyDomain",
			SignedValue: signedData.Spec.ReplicatedProxyDomain,
			ActualValue: outerData.ReplicatedProxyDomain,
		}
	}
	if outerData.CustomerID != signedData.Spec.CustomerID {
		return &types.LicenseDataValidationError{
			FieldName:   "customerID",
			SignedValue: signedData.Spec.CustomerID,
			ActualValue: outerData.CustomerID,
		}
	}
	if outerData.CustomerName != signedData.Spec.CustomerName {
		return &types.LicenseDataValidationError{
			FieldName:   "customerName",
			SignedValue: signedData.Spec.CustomerName,
			ActualValue: outerData.CustomerName,
		}
	}
	if outerData.CustomerEmail != signedData.Spec.CustomerEmail {
		return &types.LicenseDataValidationError{
			FieldName:   "customerEmail",
			SignedValue: signedData.Spec.CustomerEmail,
			ActualValue: outerData.CustomerEmail,
		}
	}
	if outerData.ChannelID != signedData.Spec.ChannelID {
		return &types.LicenseDataValidationError{
			FieldName:   "channelID",
			SignedValue: signedData.Spec.ChannelID,
			ActualValue: outerData.ChannelID,
		}
	}
	if outerData.ChannelName != signedData.Spec.ChannelName {
		return &types.LicenseDataValidationError{
			FieldName:   "channelName",
			SignedValue: signedData.Spec.ChannelName,
			ActualValue: outerData.ChannelName,
		}
	}
	if outerData.LicenseSequence != signedData.Spec.LicenseSequence {
		return &types.LicenseDataValidationError{
			FieldName:   "licenseSequence",
			SignedValue: fmt.Sprintf("%d", signedData.Spec.LicenseSequence),
			ActualValue: fmt.Sprintf("%d", outerData.LicenseSequence),
		}
	}
	if outerData.LicenseID != signedData.Spec.LicenseID {
		return &types.LicenseDataValidationError{
			FieldName:   "licenseID",
			SignedValue: signedData.Spec.LicenseID,
			ActualValue: outerData.LicenseID,
		}
	}
	if outerData.LicenseType != signedData.Spec.LicenseType {
		return &types.LicenseDataValidationError{
			FieldName:   "licenseType",
			SignedValue: signedData.Spec.LicenseType,
			ActualValue: outerData.LicenseType,
		}
	}
	if outerData.IsAirgapSupported != signedData.Spec.IsAirgapSupported {
		return &types.LicenseDataValidationError{
			FieldName:   "isAirgapSupported",
			SignedValue: fmt.Sprintf("%t", signedData.Spec.IsAirgapSupported),
			ActualValue: fmt.Sprintf("%t", outerData.IsAirgapSupported),
		}
	}
	if outerData.IsGitOpsSupported != signedData.Spec.IsGitOpsSupported {
		return &types.LicenseDataValidationError{
			FieldName:   "isGitOpsSupported",
			SignedValue: fmt.Sprintf("%t", signedData.Spec.IsGitOpsSupported),
			ActualValue: fmt.Sprintf("%t", outerData.IsGitOpsSupported),
		}
	}
	if outerData.IsIdentityServiceSupported != signedData.Spec.IsIdentityServiceSupported {
		return &types.LicenseDataValidationError{
			FieldName:   "isIdentityServiceSupported",
			SignedValue: fmt.Sprintf("%t", signedData.Spec.IsIdentityServiceSupported),
			ActualValue: fmt.Sprintf("%t", outerData.IsIdentityServiceSupported),
		}
	}
	if outerData.IsGeoaxisSupported != signedData.Spec.IsGeoaxisSupported {
		return &types.LicenseDataValidationError{
			FieldName:   "isGeoaxisSupported",
			SignedValue: fmt.Sprintf("%t", signedData.Spec.IsGeoaxisSupported),
			ActualValue: fmt.Sprintf("%t", outerData.IsGeoaxisSupported),
		}
	}
	if outerData.IsSnapshotSupported != signedData.Spec.IsSnapshotSupported {
		return &types.LicenseDataValidationError{
			FieldName:   "isSnapshotSupported",
			SignedValue: fmt.Sprintf("%t", signedData.Spec.IsSnapshotSupported),
			ActualValue: fmt.Sprintf("%t", outerData.IsSnapshotSupported),
		}
	}
	if outerData.IsDisasterRecoverySupported != signedData.Spec.IsDisasterRecoverySupported {
		return &types.LicenseDataValidationError{
			FieldName:   "isDisasterRecoverySupported",
			SignedValue: fmt.Sprintf("%t", signedData.Spec.IsDisasterRecoverySupported),
			ActualValue: fmt.Sprintf("%t", outerData.IsDisasterRecoverySupported),
		}
	}
	if outerData.IsSupportBundleUploadSupported != signedData.Spec.IsSupportBundleUploadSupported {
		return &types.LicenseDataValidationError{
			FieldName:   "isSupportBundleUploadSupported",
			SignedValue: fmt.Sprintf("%t", signedData.Spec.IsSupportBundleUploadSupported),
			ActualValue: fmt.Sprintf("%t", outerData.IsSupportBundleUploadSupported),
		}
	}
	if outerData.IsSemverRequired != signedData.Spec.IsSemverRequired {
		return &types.LicenseDataValidationError{
			FieldName:   "isSemverRequired",
			SignedValue: fmt.Sprintf("%t", signedData.Spec.IsSemverRequired),
			ActualValue: fmt.Sprintf("%t", outerData.IsSemverRequired),
		}
	}
	if outerData.IsEmbeddedClusterDownloadEnabled != signedData.Spec.IsEmbeddedClusterDownloadEnabled {
		return &types.LicenseDataValidationError{
			FieldName:   "isEmbeddedClusterDownloadEnabled",
			SignedValue: fmt.Sprintf("%t", signedData.Spec.IsEmbeddedClusterDownloadEnabled),
			ActualValue: fmt.Sprintf("%t", outerData.IsEmbeddedClusterDownloadEnabled),
		}
	}
	if outerData.IsEmbeddedClusterMultiNodeEnabled != signedData.Spec.IsEmbeddedClusterMultiNodeEnabled {
		return &types.LicenseDataValidationError{
			FieldName:   "isEmbeddedClusterMultiNodeEnabled",
			SignedValue: fmt.Sprintf("%t", signedData.Spec.IsEmbeddedClusterMultiNodeEnabled),
			ActualValue: fmt.Sprintf("%t", outerData.IsEmbeddedClusterMultiNodeEnabled),
		}
	}

	// Compare channels (order matters for slices)
	if len(outerData.Channels) != len(signedData.Spec.Channels) {
		return &types.LicenseDataValidationError{
			FieldName:   "channels length",
			SignedValue: fmt.Sprintf("%d", len(signedData.Spec.Channels)),
			ActualValue: fmt.Sprintf("%d", len(outerData.Channels)),
		}
	}
	for i, channel := range outerData.Channels {
		if channel.ChannelID != signedData.Spec.Channels[i].ChannelID {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("channels[%d].channelID", i),
				SignedValue: signedData.Spec.Channels[i].ChannelID,
				ActualValue: channel.ChannelID,
			}
		}
		if channel.ChannelName != signedData.Spec.Channels[i].ChannelName {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("channels[%d].channelName", i),
				SignedValue: signedData.Spec.Channels[i].ChannelName,
				ActualValue: channel.ChannelName,
			}
		}
		if channel.ChannelSlug != signedData.Spec.Channels[i].ChannelSlug {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("channels[%d].channelSlug", i),
				SignedValue: signedData.Spec.Channels[i].ChannelSlug,
				ActualValue: channel.ChannelSlug,
			}
		}
		if channel.IsDefault != signedData.Spec.Channels[i].IsDefault {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("channels[%d].isDefault", i),
				SignedValue: fmt.Sprintf("%t", signedData.Spec.Channels[i].IsDefault),
				ActualValue: fmt.Sprintf("%t", channel.IsDefault),
			}
		}
		if channel.Endpoint != signedData.Spec.Channels[i].Endpoint {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("channels[%d].endpoint", i),
				SignedValue: signedData.Spec.Channels[i].Endpoint,
				ActualValue: channel.Endpoint,
			}
		}
		if channel.ReplicatedProxyDomain != signedData.Spec.Channels[i].ReplicatedProxyDomain {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("channels[%d].replicatedProxyDomain", i),
				SignedValue: signedData.Spec.Channels[i].ReplicatedProxyDomain,
				ActualValue: channel.ReplicatedProxyDomain,
			}
		}
		if channel.IsSemverRequired != signedData.Spec.Channels[i].IsSemverRequired {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("channels[%d].isSemverRequired", i),
				SignedValue: fmt.Sprintf("%t", signedData.Spec.Channels[i].IsSemverRequired),
				ActualValue: fmt.Sprintf("%t", channel.IsSemverRequired),
			}
		}
	}

	// Compare entitlements (order doesn't matter for maps)
	if len(outerData.Entitlements) != len(signedData.Spec.Entitlements) {
		return &types.LicenseDataValidationError{
			FieldName:   "entitlements length",
			SignedValue: fmt.Sprintf("%d", len(signedData.Spec.Entitlements)),
			ActualValue: fmt.Sprintf("%d", len(outerData.Entitlements)),
		}
	}
	for fieldName, field := range l.Spec.Entitlements {
		signedField, exists := signedData.Spec.Entitlements[fieldName]
		if !exists {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("entitlements[%s]", fieldName),
				SignedValue: "",
				ActualValue: "<missing>",
			}
		}
		if field.Title != signedField.Title {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("entitlements[%s].title", fieldName),
				SignedValue: signedField.Title,
				ActualValue: field.Title,
			}
		}
		if field.Description != signedField.Description {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("entitlements[%s].description", fieldName),
				SignedValue: signedField.Description,
				ActualValue: field.Description,
			}
		}
		if field.ValueType != signedField.ValueType {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("entitlements[%s].valueType", fieldName),
				SignedValue: signedField.ValueType,
				ActualValue: field.ValueType,
			}
		}
		if field.IsHidden != signedField.IsHidden {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("entitlements[%s].isHidden", fieldName),
				SignedValue: fmt.Sprintf("%t", signedField.IsHidden),
				ActualValue: fmt.Sprintf("%t", field.IsHidden),
			}
		}
		if field.Value.Type != signedField.Value.Type {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("entitlements[%s].value.type", fieldName),
				SignedValue: fmt.Sprintf("%d", signedField.Value.Type),
				ActualValue: fmt.Sprintf("%d", field.Value.Type),
			}
		}
		if field.Value.Value() != signedField.Value.Value() {
			return &types.LicenseDataValidationError{
				FieldName:   fmt.Sprintf("entitlements[%s].value", fieldName),
				SignedValue: fmt.Sprintf("%v", signedField.Value.Value()),
				ActualValue: fmt.Sprintf("%v", field.Value.Value()),
			}
		}
	}

	return nil
}
