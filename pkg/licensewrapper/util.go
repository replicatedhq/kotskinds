package licensewrapper

import (
	"os"

	"github.com/pkg/errors"
	kotsv1beta1 "github.com/replicatedhq/kotskinds/apis/kots/v1beta1"
	kotsv1beta2 "github.com/replicatedhq/kotskinds/apis/kots/v1beta2"
	kotsscheme "github.com/replicatedhq/kotskinds/client/kotsclientset/scheme"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	kotsscheme.AddToScheme(scheme.Scheme)
}

// LoadLicenseFromPath loads a license from a file path and returns a LicenseWrapper
func LoadLicenseFromPath(licenseFilePath string) (LicenseWrapper, error) {
	licenseData, err := os.ReadFile(licenseFilePath)
	if err != nil {
		return LicenseWrapper{}, errors.Wrap(err, "failed to read license file")
	}

	return LoadLicenseFromBytes(licenseData)
}

// LoadLicenseFromBytes deserializes license YAML/JSON bytes into a LicenseWrapper.
// It automatically detects whether the data contains a v1beta1 or v1beta2 License.
func LoadLicenseFromBytes(data []byte) (LicenseWrapper, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, gvk, err := decode([]byte(data), nil, nil)
	if err != nil {
		return LicenseWrapper{}, errors.Wrap(err, "failed to decode license data")
	}

	if gvk.Group != "kots.io" || gvk.Kind != "License" {
		return LicenseWrapper{}, errors.Errorf("unexpected GVK: %s", gvk.String())
	}

	// Return wrapper with appropriate version populated
	switch gvk.Version {
	case "v1beta1":
		v1License, ok := obj.(*kotsv1beta1.License)
		if !ok {
			return LicenseWrapper{}, errors.New("failed to cast to v1beta1.License")
		}
		return LicenseWrapper{V1: v1License}, nil

	case "v1beta2":
		v2License, ok := obj.(*kotsv1beta2.License)
		if !ok {
			return LicenseWrapper{}, errors.New("failed to cast to v1beta2.License")
		}
		return LicenseWrapper{V2: v2License}, nil

	default:
		return LicenseWrapper{}, errors.Errorf("unsupported license version: %s", gvk.Version)
	}
}
