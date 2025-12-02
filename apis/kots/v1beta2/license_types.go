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
}

func (entitlementValue *EntitlementValue) Value() interface{} {
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

func (entitlementValue *EntitlementValue) UnmarshalJSON(value []byte) error {
	if value[0] == '"' {
		entitlementValue.Type = String
		return json.Unmarshal(value, &entitlementValue.StrVal)
	}

	intValue, err := strconv.ParseInt(string(value), 10, 64)
	if err == nil {
		entitlementValue.Type = Int
		entitlementValue.IntVal = intValue
		return nil
	}

	boolValue, err := strconv.ParseBool(string(value))
	if err == nil {
		entitlementValue.Type = Bool
		entitlementValue.BoolVal = boolValue
		return nil
	}

	return errors.New("unknown license value type")
}

type EntitlementField struct {
	Title       string                    `json:"title,omitempty"`
	Description string                    `json:"description,omitempty"`
	Value       EntitlementValue          `json:"value,omitempty"`
	ValueType   string                    `json:"valueType,omitempty"`
	IsHidden    bool                      `json:"isHidden,omitempty"`
	Signature   EntitlementFieldSignature `json:"signature,omitempty"`
}

type EntitlementFieldSignature struct {
	V2 []byte `json:"v2,omitempty"` // SHA-256 signature for v1beta2
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

// LicenseSpec defines the desired state of LicenseSpec for v1beta2
type LicenseSpec struct {
	Signature                         []byte                      `json:"signature"` // SHA-256 signature
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
// License is the Schema for the license API for v1beta2
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// NOTE: No +kubebuilder:storageversion annotation - v1beta1 remains storage version
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

// ValidateLicense validates the entire v1beta2 license signature and all entitlement signatures
// Returns the app signing keys on success (used for validating entitlement signatures), or an error if validation fails
func (l *License) ValidateLicense() (*kotscrypto.AppSigningKeys, error) {
	// Decode and parse the signature
	outerSig, innerSig, appKeys, err := kotscrypto.DecodeLicenseSignature(l.Spec.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode signature")
	}

	if len(innerSig.V2KeySignature) == 0 {
		return nil, errors.New("v2 key signature not found")
	}

	var keySig kotscrypto.KeySignature
	if err := json.Unmarshal(innerSig.V2KeySignature, &keySig); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal v2 key signature")
	}

	globalPubKey, err := kotscrypto.FindGlobalPublicKeyRSA(keySig.GlobalKeyID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find global public key")
	}

	// verify that the application public key is signed by the global key
	if err := kotscrypto.VerifySignatureRSA([]byte(innerSig.PublicKey), keySig.Signature, globalPubKey, crypto.SHA256); err != nil {
		return nil, errors.Wrap(err, "v2 key signature verification failed")
	}

	if len(innerSig.V2LicenseSignature) == 0 {
		return nil, errors.New("v2 license signature not found")
	}

	// verify that the license data is signed by the application key
	if err := kotscrypto.VerifySignatureRSA(outerSig.LicenseData, innerSig.V2LicenseSignature, innerSig.PublicKey, crypto.SHA256); err != nil {
		return nil, errors.Wrap(err, "v2 license signature verification failed")
	}

	// validate that each entitlement value is signed by the application key
	for fieldName, field := range l.Spec.Entitlements {
		// the entitlement values are still covered as part of the entire license body, and some old license files did not include entitlement signatures
		if len(field.Signature.V2) == 0 {
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

// ValidateSignature validates a single v1beta2 entitlement signature
func (e *EntitlementField) ValidateSignature(appKeys *kotscrypto.AppSigningKeys) error {
	// Check if signature is present
	if len(e.Signature.V2) == 0 {
		return errors.New("v2 signature not found for entitlement")
	}

	// Get the value as a string
	value := e.Value.Value()
	message := []byte(fmt.Sprint(value))

	// Verify the signature using the crypto package
	if err := kotscrypto.VerifySignatureWithKeyRSA(message, e.Signature.V2, appKeys.PublicKeyRSA, crypto.SHA256); err != nil {
		return errors.Wrap(err, "v2 entitlement signature verification failed")
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
