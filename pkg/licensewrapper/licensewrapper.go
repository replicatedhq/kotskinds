package licensewrapper

import (
	kotsv1beta1 "github.com/replicatedhq/kotskinds/apis/kots/v1beta1"
	kotsv1beta2 "github.com/replicatedhq/kotskinds/apis/kots/v1beta2"
)

// LicenseWrapper holds either a v1beta1 or v1beta2 license (never both).
// Exactly one field will be non-nil.
type LicenseWrapper struct {
	V1 *kotsv1beta1.License
	V2 *kotsv1beta2.License
}

// EntitlementFieldWrapper holds either a v1beta1 or v1beta2 EntitlementField (never both).
// Exactly one field will be non-nil.
type EntitlementFieldWrapper struct {
	V1 *kotsv1beta1.EntitlementField
	V2 *kotsv1beta2.EntitlementField
}

// IsV1 returns true if this wrapper contains a v1beta1 license
func (w LicenseWrapper) IsV1() bool {
	return w.V1 != nil
}

// IsV2 returns true if this wrapper contains a v1beta2 license
func (w LicenseWrapper) IsV2() bool {
	return w.V2 != nil
}

// GetAppSlug returns the app slug from whichever version is present
func (w LicenseWrapper) GetAppSlug() string {
	if w.V1 != nil {
		return w.V1.Spec.AppSlug
	}
	if w.V2 != nil {
		return w.V2.Spec.AppSlug
	}
	return ""
}

// GetLicenseID returns the license ID from whichever version is present
func (w LicenseWrapper) GetLicenseID() string {
	if w.V1 != nil {
		return w.V1.Spec.LicenseID
	}
	if w.V2 != nil {
		return w.V2.Spec.LicenseID
	}
	return ""
}

// GetLicenseType returns the license type from whichever version is present
func (w LicenseWrapper) GetLicenseType() string {
	if w.V1 != nil {
		return w.V1.Spec.LicenseType
	}
	if w.V2 != nil {
		return w.V2.Spec.LicenseType
	}
	return ""
}

// GetEndpoint returns the endpoint from whichever version is present
func (w LicenseWrapper) GetEndpoint() string {
	if w.V1 != nil {
		return w.V1.Spec.Endpoint
	}
	if w.V2 != nil {
		return w.V2.Spec.Endpoint
	}
	return ""
}

// GetChannelID returns the channel ID from whichever version is present
func (w LicenseWrapper) GetChannelID() string {
	if w.V1 != nil {
		return w.V1.Spec.ChannelID
	}
	if w.V2 != nil {
		return w.V2.Spec.ChannelID
	}
	return ""
}

// GetChannelName returns the channel name from whichever version is present
func (w LicenseWrapper) GetChannelName() string {
	if w.V1 != nil {
		return w.V1.Spec.ChannelName
	}
	if w.V2 != nil {
		return w.V2.Spec.ChannelName
	}
	return ""
}

// GetCustomerName returns the customer name from whichever version is present
func (w LicenseWrapper) GetCustomerName() string {
	if w.V1 != nil {
		return w.V1.Spec.CustomerName
	}
	if w.V2 != nil {
		return w.V2.Spec.CustomerName
	}
	return ""
}

// GetCustomerEmail returns the customer email from whichever version is present
func (w LicenseWrapper) GetCustomerEmail() string {
	if w.V1 != nil {
		return w.V1.Spec.CustomerEmail
	}
	if w.V2 != nil {
		return w.V2.Spec.CustomerEmail
	}
	return ""
}

// GetLicenseSequence returns the license sequence from whichever version is present
func (w LicenseWrapper) GetLicenseSequence() int64 {
	if w.V1 != nil {
		return w.V1.Spec.LicenseSequence
	}
	if w.V2 != nil {
		return w.V2.Spec.LicenseSequence
	}
	return 0
}

// IsAirgapSupported returns whether airgap is supported from whichever version is present
func (w LicenseWrapper) IsAirgapSupported() bool {
	if w.V1 != nil {
		return w.V1.Spec.IsAirgapSupported
	}
	if w.V2 != nil {
		return w.V2.Spec.IsAirgapSupported
	}
	return false
}

// IsGitOpsSupported returns whether GitOps is supported from whichever version is present
func (w LicenseWrapper) IsGitOpsSupported() bool {
	if w.V1 != nil {
		return w.V1.Spec.IsGitOpsSupported
	}
	if w.V2 != nil {
		return w.V2.Spec.IsGitOpsSupported
	}
	return false
}

// IsIdentityServiceSupported returns whether identity service is supported from whichever version is present
func (w LicenseWrapper) IsIdentityServiceSupported() bool {
	if w.V1 != nil {
		return w.V1.Spec.IsIdentityServiceSupported
	}
	if w.V2 != nil {
		return w.V2.Spec.IsIdentityServiceSupported
	}
	return false
}

// IsGeoaxisSupported returns whether geoaxis is supported from whichever version is present
func (w LicenseWrapper) IsGeoaxisSupported() bool {
	if w.V1 != nil {
		return w.V1.Spec.IsGeoaxisSupported
	}
	if w.V2 != nil {
		return w.V2.Spec.IsGeoaxisSupported
	}
	return false
}

// IsSnapshotSupported returns whether snapshots are supported from whichever version is present
func (w LicenseWrapper) IsSnapshotSupported() bool {
	if w.V1 != nil {
		return w.V1.Spec.IsSnapshotSupported
	}
	if w.V2 != nil {
		return w.V2.Spec.IsSnapshotSupported
	}
	return false
}

// IsSupportBundleUploadSupported returns whether support bundle upload is supported from whichever version is present
func (w LicenseWrapper) IsSupportBundleUploadSupported() bool {
	if w.V1 != nil {
		return w.V1.Spec.IsSupportBundleUploadSupported
	}
	if w.V2 != nil {
		return w.V2.Spec.IsSupportBundleUploadSupported
	}
	return false
}

// IsSemverRequired returns whether semver is required from whichever version is present
func (w LicenseWrapper) IsSemverRequired() bool {
	if w.V1 != nil {
		return w.V1.Spec.IsSemverRequired
	}
	if w.V2 != nil {
		return w.V2.Spec.IsSemverRequired
	}
	return false
}

// GetSignature returns the license signature from whichever version is present
func (w LicenseWrapper) GetSignature() []byte {
	if w.V1 != nil {
		return w.V1.Spec.Signature
	}
	if w.V2 != nil {
		return w.V2.Spec.Signature
	}
	return nil
}

// GetReplicatedProxyDomain returns the replicated proxy domain from whichever version is present
func (w LicenseWrapper) GetReplicatedProxyDomain() string {
	if w.V1 != nil {
		return w.V1.Spec.ReplicatedProxyDomain
	}
	if w.V2 != nil {
		return w.V2.Spec.ReplicatedProxyDomain
	}
	return ""
}

// GetChannels returns the list of channels from whichever version is present
// Channel type is identical in both v1beta1 and v1beta2, so we return v1beta1.Channel
func (w LicenseWrapper) GetChannels() []kotsv1beta1.Channel {
	if w.V1 != nil {
		return w.V1.Spec.Channels
	}
	if w.V2 != nil {
		// Channel types are identical, safe to convert
		channels := make([]kotsv1beta1.Channel, len(w.V2.Spec.Channels))
		for i, ch := range w.V2.Spec.Channels {
			channels[i] = kotsv1beta1.Channel{
				ChannelID:             ch.ChannelID,
				ChannelName:           ch.ChannelName,
				ChannelSlug:           ch.ChannelSlug,
				IsDefault:             ch.IsDefault,
				Endpoint:              ch.Endpoint,
				ReplicatedProxyDomain: ch.ReplicatedProxyDomain,
				IsSemverRequired:      ch.IsSemverRequired,
			}
		}
		return channels
	}
	return nil
}

// IsDisasterRecoverySupported returns whether disaster recovery is supported from whichever version is present
func (w LicenseWrapper) IsDisasterRecoverySupported() bool {
	if w.V1 != nil {
		return w.V1.Spec.IsDisasterRecoverySupported
	}
	if w.V2 != nil {
		return w.V2.Spec.IsDisasterRecoverySupported
	}
	return false
}

// IsEmbeddedClusterDownloadEnabled returns whether embedded cluster download is enabled from whichever version is present
func (w LicenseWrapper) IsEmbeddedClusterDownloadEnabled() bool {
	if w.V1 != nil {
		return w.V1.Spec.IsEmbeddedClusterDownloadEnabled
	}
	if w.V2 != nil {
		return w.V2.Spec.IsEmbeddedClusterDownloadEnabled
	}
	return false
}

// IsEmbeddedClusterMultiNodeEnabled returns whether embedded cluster multi-node is enabled from whichever version is present
func (w LicenseWrapper) IsEmbeddedClusterMultiNodeEnabled() bool {
	if w.V1 != nil {
		return w.V1.Spec.IsEmbeddedClusterMultiNodeEnabled
	}
	if w.V2 != nil {
		return w.V2.Spec.IsEmbeddedClusterMultiNodeEnabled
	}
	return false
}

// GetEntitlements returns the entitlements map from whichever version is present
// Returns wrapped entitlements for version-agnostic access
func (w LicenseWrapper) GetEntitlements() map[string]EntitlementFieldWrapper {
	if w.V1 != nil {
		if w.V1.Spec.Entitlements == nil {
			return nil
		}
		wrapped := make(map[string]EntitlementFieldWrapper, len(w.V1.Spec.Entitlements))
		for key, ent := range w.V1.Spec.Entitlements {
			entCopy := ent
			wrapped[key] = EntitlementFieldWrapper{V1: &entCopy}
		}
		return wrapped
	}
	if w.V2 != nil {
		if w.V2.Spec.Entitlements == nil {
			return nil
		}
		wrapped := make(map[string]EntitlementFieldWrapper, len(w.V2.Spec.Entitlements))
		for key, ent := range w.V2.Spec.Entitlements {
			entCopy := ent
			wrapped[key] = EntitlementFieldWrapper{V2: &entCopy}
		}
		return wrapped
	}
	return nil
}

// EntitlementFieldWrapper accessor methods

// GetTitle returns the entitlement title from whichever version is present
func (w EntitlementFieldWrapper) GetTitle() string {
	if w.V1 != nil {
		return w.V1.Title
	}
	if w.V2 != nil {
		return w.V2.Title
	}
	return ""
}

// GetDescription returns the entitlement description from whichever version is present
func (w EntitlementFieldWrapper) GetDescription() string {
	if w.V1 != nil {
		return w.V1.Description
	}
	if w.V2 != nil {
		return w.V2.Description
	}
	return ""
}

// GetValue returns the entitlement value from whichever version is present
// EntitlementValue type is identical in both versions
func (w EntitlementFieldWrapper) GetValue() kotsv1beta1.EntitlementValue {
	if w.V1 != nil {
		return w.V1.Value
	}
	if w.V2 != nil {
		// Safe to cast since EntitlementValue is identical in both versions
		return kotsv1beta1.EntitlementValue{
			Type:    kotsv1beta1.Type(w.V2.Value.Type),
			IntVal:  w.V2.Value.IntVal,
			StrVal:  w.V2.Value.StrVal,
			BoolVal: w.V2.Value.BoolVal,
		}
	}
	return kotsv1beta1.EntitlementValue{}
}

// GetValueType returns the entitlement value type from whichever version is present
func (w EntitlementFieldWrapper) GetValueType() string {
	if w.V1 != nil {
		return w.V1.ValueType
	}
	if w.V2 != nil {
		return w.V2.ValueType
	}
	return ""
}

// IsHidden returns whether the entitlement is hidden from whichever version is present
func (w EntitlementFieldWrapper) IsHidden() bool {
	if w.V1 != nil {
		return w.V1.IsHidden
	}
	if w.V2 != nil {
		return w.V2.IsHidden
	}
	return false
}

// GetSignature returns the entitlement signature from whichever version is present
// Abstracts the V1/V2 signature field difference
func (w EntitlementFieldWrapper) GetSignature() []byte {
	if w.V1 != nil {
		return w.V1.Signature.V1
	}
	if w.V2 != nil {
		return w.V2.Signature.V2
	}
	return nil
}
