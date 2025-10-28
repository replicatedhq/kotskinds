// Package licensewrapper provides a version-agnostic wrapper for KOTS License resources.
//
// The LicenseWrapper type can hold either a v1beta1.License or v1beta2.License, providing
// a unified interface for accessing license properties regardless of the underlying version.
//
// Usage:
//
//	// Load from file
//	wrapper, err := licensewrapper.LoadLicenseFromPath("/path/to/license.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Load from bytes
//	wrapper, err := licensewrapper.LoadLicenseFromBytes(yamlData)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Check version
//	if wrapper.IsV1() {
//	    // Working with v1beta1
//	}
//	if wrapper.IsV2() {
//	    // Working with v1beta2
//	}
//
//	// Access properties (works with both versions)
//	appSlug := wrapper.GetAppSlug()
//	licenseID := wrapper.GetLicenseID()
//	isAirgap := wrapper.IsAirgapSupported()
package licensewrapper
