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

	// Validate license data matches
	if err := l.compareLicenseData(outerSig.LicenseData); err != nil {
		return nil, errors.Wrap(err, "license data validation failed")
	}

	// v1beta2 should use SHA-256 signatures
	if len(innerSig.V2LicenseSignature) == 0 {
		return nil, errors.New("v2 license signature not found")
	}

	// Verify that the license data is signed by the application key
	if err := kotscrypto.VerifySignatureRSA(outerSig.LicenseData, innerSig.V2LicenseSignature, innerSig.PublicKey, crypto.SHA256); err != nil {
		return nil, errors.Wrap(err, "v2 license signature verification failed")
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

	// validate that each entitlement value is signed by the application key
	for fieldName, field := range l.Spec.Entitlements {
		if err := field.ValidateSignature(appKeys); err != nil {
			return nil, errors.Wrapf(err, "entitlement validation failed for field: %s", fieldName)
		}
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
func (l *License) compareLicenseData(signedJSON []byte) error {
	// Decode the license JSON
	var signedData License
	if err := json.Unmarshal(signedJSON, &signedData); err != nil {
		return errors.Wrap(err, "failed to unmarshal signed license data")
	}

	// Compare each field in l.Spec against the decoded license
	if l.Spec.AppSlug != signedData.Spec.AppSlug {
		return errors.Errorf(`"appSlug" field has changed to %q (license) from %q (within signature)`, l.Spec.AppSlug, signedData.Spec.AppSlug)
	}
	if l.Spec.Endpoint != signedData.Spec.Endpoint {
		return errors.Errorf(`"endpoint" field has changed to %q (license) from %q (within signature)`, l.Spec.Endpoint, signedData.Spec.Endpoint)
	}
	if l.Spec.ReplicatedProxyDomain != signedData.Spec.ReplicatedProxyDomain {
		return errors.Errorf(`"replicatedProxyDomain" field has changed to %q (license) from %q (within signature)`, l.Spec.ReplicatedProxyDomain, signedData.Spec.ReplicatedProxyDomain)
	}
	if l.Spec.CustomerName != signedData.Spec.CustomerName {
		return errors.Errorf(`"customerName" field has changed to %q (license) from %q (within signature)`, l.Spec.CustomerName, signedData.Spec.CustomerName)
	}
	if l.Spec.CustomerEmail != signedData.Spec.CustomerEmail {
		return errors.Errorf(`"customerEmail" field has changed to %q (license) from %q (within signature)`, l.Spec.CustomerEmail, signedData.Spec.CustomerEmail)
	}
	if l.Spec.ChannelID != signedData.Spec.ChannelID {
		return errors.Errorf(`"channelID" field has changed to %q (license) from %q (within signature)`, l.Spec.ChannelID, signedData.Spec.ChannelID)
	}
	if l.Spec.ChannelName != signedData.Spec.ChannelName {
		return errors.Errorf(`"channelName" field has changed to %q (license) from %q (within signature)`, l.Spec.ChannelName, signedData.Spec.ChannelName)
	}
	if l.Spec.LicenseSequence != signedData.Spec.LicenseSequence {
		return errors.Errorf(`"licenseSequence" field has changed to %d (license) from %d (within signature)`, l.Spec.LicenseSequence, signedData.Spec.LicenseSequence)
	}
	if l.Spec.LicenseID != signedData.Spec.LicenseID {
		return errors.Errorf(`"licenseID" field has changed to %q (license) from %q (within signature)`, l.Spec.LicenseID, signedData.Spec.LicenseID)
	}
	if l.Spec.LicenseType != signedData.Spec.LicenseType {
		return errors.Errorf(`"licenseType" field has changed to %q (license) from %q (within signature)`, l.Spec.LicenseType, signedData.Spec.LicenseType)
	}
	if l.Spec.IsAirgapSupported != signedData.Spec.IsAirgapSupported {
		return errors.Errorf(`"isAirgapSupported" field has changed to %t (license) from %t (within signature)`, l.Spec.IsAirgapSupported, signedData.Spec.IsAirgapSupported)
	}
	if l.Spec.IsGitOpsSupported != signedData.Spec.IsGitOpsSupported {
		return errors.Errorf(`"isGitOpsSupported" field has changed to %t (license) from %t (within signature)`, l.Spec.IsGitOpsSupported, signedData.Spec.IsGitOpsSupported)
	}
	if l.Spec.IsIdentityServiceSupported != signedData.Spec.IsIdentityServiceSupported {
		return errors.Errorf(`"isIdentityServiceSupported" field has changed to %t (license) from %t (within signature)`, l.Spec.IsIdentityServiceSupported, signedData.Spec.IsIdentityServiceSupported)
	}
	if l.Spec.IsGeoaxisSupported != signedData.Spec.IsGeoaxisSupported {
		return errors.Errorf(`"isGeoaxisSupported" field has changed to %t (license) from %t (within signature)`, l.Spec.IsGeoaxisSupported, signedData.Spec.IsGeoaxisSupported)
	}
	if l.Spec.IsSnapshotSupported != signedData.Spec.IsSnapshotSupported {
		return errors.Errorf(`"isSnapshotSupported" field has changed to %t (license) from %t (within signature)`, l.Spec.IsSnapshotSupported, signedData.Spec.IsSnapshotSupported)
	}
	if l.Spec.IsDisasterRecoverySupported != signedData.Spec.IsDisasterRecoverySupported {
		return errors.Errorf(`"isDisasterRecoverySupported" field has changed to %t (license) from %t (within signature)`, l.Spec.IsDisasterRecoverySupported, signedData.Spec.IsDisasterRecoverySupported)
	}
	if l.Spec.IsSupportBundleUploadSupported != signedData.Spec.IsSupportBundleUploadSupported {
		return errors.Errorf(`"isSupportBundleUploadSupported" field has changed to %t (license) from %t (within signature)`, l.Spec.IsSupportBundleUploadSupported, signedData.Spec.IsSupportBundleUploadSupported)
	}
	if l.Spec.IsSemverRequired != signedData.Spec.IsSemverRequired {
		return errors.Errorf(`"isSemverRequired" field has changed to %t (license) from %t (within signature)`, l.Spec.IsSemverRequired, signedData.Spec.IsSemverRequired)
	}
	if l.Spec.IsEmbeddedClusterDownloadEnabled != signedData.Spec.IsEmbeddedClusterDownloadEnabled {
		return errors.Errorf(`"isEmbeddedClusterDownloadEnabled" field has changed to %t (license) from %t (within signature)`, l.Spec.IsEmbeddedClusterDownloadEnabled, signedData.Spec.IsEmbeddedClusterDownloadEnabled)
	}
	if l.Spec.IsEmbeddedClusterMultiNodeEnabled != signedData.Spec.IsEmbeddedClusterMultiNodeEnabled {
		return errors.Errorf(`"isEmbeddedClusterMultiNodeEnabled" field has changed to %t (license) from %t (within signature)`, l.Spec.IsEmbeddedClusterMultiNodeEnabled, signedData.Spec.IsEmbeddedClusterMultiNodeEnabled)
	}

	// Compare channels (order matters for slices)
	if len(l.Spec.Channels) != len(signedData.Spec.Channels) {
		return errors.Errorf(`"channels" length has changed to %d (license) from %d (within signature)`, len(l.Spec.Channels), len(signedData.Spec.Channels))
	}
	for i, channel := range l.Spec.Channels {
		if channel.ChannelID != signedData.Spec.Channels[i].ChannelID {
			return errors.Errorf(`"channels[%d].channelID" field has changed to %q (license) from %q (within signature)`, i, channel.ChannelID, signedData.Spec.Channels[i].ChannelID)
		}
		if channel.ChannelName != signedData.Spec.Channels[i].ChannelName {
			return errors.Errorf(`"channels[%d].channelName" field has changed to %q (license) from %q (within signature)`, i, channel.ChannelName, signedData.Spec.Channels[i].ChannelName)
		}
		if channel.ChannelSlug != signedData.Spec.Channels[i].ChannelSlug {
			return errors.Errorf(`"channels[%d].channelSlug" field has changed to %q (license) from %q (within signature)`, i, channel.ChannelSlug, signedData.Spec.Channels[i].ChannelSlug)
		}
		if channel.IsDefault != signedData.Spec.Channels[i].IsDefault {
			return errors.Errorf(`"channels[%d].isDefault" field has changed to %t (license) from %t (within signature)`, i, channel.IsDefault, signedData.Spec.Channels[i].IsDefault)
		}
		if channel.Endpoint != signedData.Spec.Channels[i].Endpoint {
			return errors.Errorf(`"channels[%d].endpoint" field has changed to %q (license) from %q (within signature)`, i, channel.Endpoint, signedData.Spec.Channels[i].Endpoint)
		}
		if channel.ReplicatedProxyDomain != signedData.Spec.Channels[i].ReplicatedProxyDomain {
			return errors.Errorf(`"channels[%d].replicatedProxyDomain" field has changed to %q (license) from %q (within signature)`, i, channel.ReplicatedProxyDomain, signedData.Spec.Channels[i].ReplicatedProxyDomain)
		}
		if channel.IsSemverRequired != signedData.Spec.Channels[i].IsSemverRequired {
			return errors.Errorf(`"channels[%d].isSemverRequired" field has changed to %t (license) from %t (within signature)`, i, channel.IsSemverRequired, signedData.Spec.Channels[i].IsSemverRequired)
		}
	}

	// Compare entitlements (order doesn't matter for maps)
	if len(l.Spec.Entitlements) != len(signedData.Spec.Entitlements) {
		return errors.Errorf(`"entitlements" length has changed to %d (license) from %d (within signature)`, len(l.Spec.Entitlements), len(signedData.Spec.Entitlements))
	}
	for fieldName, field := range l.Spec.Entitlements {
		signedField, exists := signedData.Spec.Entitlements[fieldName]
		if !exists {
			return errors.Errorf(`"entitlements[%s]" field not found within signature`, fieldName)
		}
		if field.Title != signedField.Title {
			return errors.Errorf(`"entitlements[%s].title" field has changed to %q (license) from %q (within signature)`, fieldName, field.Title, signedField.Title)
		}
		if field.Description != signedField.Description {
			return errors.Errorf(`"entitlements[%s].description" field has changed to %q (license) from %q (within signature)`, fieldName, field.Description, signedField.Description)
		}
		if field.ValueType != signedField.ValueType {
			return errors.Errorf(`"entitlements[%s].valueType" field has changed to %q (license) from %q (within signature)`, fieldName, field.ValueType, signedField.ValueType)
		}
		if field.IsHidden != signedField.IsHidden {
			return errors.Errorf(`"entitlements[%s].isHidden" field has changed to %t (license) from %t (within signature)`, fieldName, field.IsHidden, signedField.IsHidden)
		}
		if field.Value.Type != signedField.Value.Type {
			return errors.Errorf(`"entitlements[%s].value.type" field has changed to %d (license) from %d (within signature)`, fieldName, field.Value.Type, signedField.Value.Type)
		}
		if field.Value.Value() != signedField.Value.Value() {
			return errors.Errorf(`"entitlements[%s].value" field has changed to %v (license) from %v (within signature)`, fieldName, field.Value.Value(), signedField.Value.Value())
		}
	}

	return nil
}
