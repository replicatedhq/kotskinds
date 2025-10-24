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

package crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"

	"github.com/pkg/errors"
)

// Default public keys for signature verification
var PublicKeysRSA = map[string][]byte{
	"1d3f7f6b50714fe7b895554dd65773b0": []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAugyKfZV2gIDaY1Rzkjoo
fbNywGa04sGQIAqYwifMay2e2xzqRwswTRHQnr9SIWypkN86Cfn6QzOB8kkjERC1
DPNdsiKdjBFdcLaxxdyHgrXLgfdzhh6We+Lpq19JT5LCK3PXleZgt/a0aRBpIc1l
xKs57d8MTWUTVh3W3WYi6LbqAPScdmSiG7A145HhKXmmtZFEv4puE5dKmS5lkV2d
VU789XWrNFk74FKKHVwYMdppqAabB6cRBmU8YFiVEULOn+d1FtKRbO/vv/fbA9nX
PUG/1PgEQHogP+3cC4J7b7s9+kBmtHkpSq9x+OUu/5B+nT21dooS6adfQiI8iB/+
NQIDAQAB
-----END PUBLIC KEY-----`), // Dev

	"bdee56560cfb43c9b28bf98eacafa646": []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwdSHE8v64QH/yELBoPBl
GanhS3AD5vMAaqLLFnftwjmDKrxWwqNB9w1GVJWb5gVLvt/UlE/k+HVr5HFdomVI
TMvnvxhD0UvNyGFuUbXBMvQPPW9joR48LcCBLZl+RZTqR5HRhsIbujiExRDnteaq
mU1jG/oVlQkRoyOYrObTeoD0BdcZAr2PdGvgvJvpZduZtrKvjvsSJEBYExoPtko+
8AqhMBAI+qX1/SMix21qpmYSYLNeqN2Pplna0p2MK8yyaHY8KSqTF90ZJF1+P0ZF
MLt6S8/6PIX9WD+vFqmDpW1GCkB+p2OfxsYiAIX1ej98Ck3hoPQnOuiFIovV8aFQ
bQIDAQAB
-----END PUBLIC KEY-----`), // Production

	"de2c275656d04b1bb0f15cf70f0ea2a2": []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2hHg1HER6NYlsqBs+B+B
txibtctT6YB5kxgE1sz7UmVnlcLs+Olc4OZJwD4vLsEU60SVW0HRoTfaGaradv0R
GUIxlFRSOnzjZEMkm/YKL3sdPQigi2m9O0P5tC9LQvzk49dFg5HJxiLODCgWwJ9g
q3pGs8OaAc0dop/tqUE7WqQfHLWJdTPP5pVDLDWybfAO4OmgVmx+oVXdCfMVlOzu
num6SOF+eBuERXQGbEfnd6eSRVokWhfMCfXNPTYtq14DaK9tvX4uzHsub+Asn6UN
OBIAESJntpZfdDDrNqbfOQYql2rqx1lJtU7lVFbTQTkKhj4teInEGO6FvLzy0UE9
swIDAQAB
-----END PUBLIC KEY-----`), // Staging
}

// customPublicKeyRSA allows overriding the default public keys with a single custom key
var customPublicKeyRSA *rsa.PublicKey

// AppSigningKeys contains the public key used to verify license and entitlement signatures
type AppSigningKeys struct {
	PublicKeyRSA *rsa.PublicKey
}

// SetCustomPublicKey sets a custom public key to use instead of the default public keys
func SetCustomPublicKeyRSA(publicKeyPEM string) error {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return errors.New("failed to decode PEM block")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return errors.Wrap(err, "failed to parse public key")
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("public key is not an RSA public key")
	}

	customPublicKeyRSA = rsaPubKey
	return nil
}

// ResetCustomPublicKey clears any custom public key and reverts to default public keys
func ResetCustomPublicKeyRSA() {
	customPublicKeyRSA = nil
}

// OuterSignature represents the outer layer of the license signature
type OuterSignature struct {
	LicenseData    []byte `json:"licenseData"`
	InnerSignature []byte `json:"innerSignature"`
}

// InnerSignature represents the inner layer of the license signature
type InnerSignature struct {
	LicenseSignature   []byte `json:"licenseSignature,omitempty"`
	V2LicenseSignature []byte `json:"v2LicenseSignature,omitempty"`
	PublicKey          string `json:"publicKey"`
	KeySignature       []byte `json:"keySignature,omitempty"`
	V2KeySignature     []byte `json:"v2KeySignature,omitempty"`
}

// KeySignature represents a key signature structure
type KeySignature struct {
	Signature   []byte `json:"signature"`
	GlobalKeyID string `json:"globalKeyId"`
}

// VerifySignature verifies an RSA-PSS signature using the specified hash algorithm
func VerifySignatureRSA(message, signature []byte, publicKeyPEM string, hashAlgo crypto.Hash) error {
	// Parse the public key
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return errors.New("failed to decode public key PEM")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return errors.Wrap(err, "failed to parse public key")
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("public key is not an RSA public key")
	}

	return VerifySignatureWithKeyRSA(message, signature, rsaPubKey, hashAlgo)
}

// VerifySignatureWithKeyRSA verifies an RSA-PSS signature using the specified hash algorithm and RSA public key
func VerifySignatureWithKeyRSA(message, signature []byte, publicKey *rsa.PublicKey, hashAlgo crypto.Hash) error {
	// Hash the message
	hash := hashAlgo.New()
	hash.Write(message)
	hashed := hash.Sum(nil)

	// Verify the signature using RSA-PSS
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto

	err := rsa.VerifyPSS(publicKey, hashAlgo, hashed, signature, &opts)
	if err != nil {
		return errors.Wrap(err, "signature verification failed")
	}

	return nil
}

// parsePublicKey parses a PEM-encoded RSA public key
func parsePublicKeyRSA(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to decode public key PEM")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse public key")
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not an RSA public key")
	}

	return rsaPubKey, nil
}

// FindGlobalPublicKeyRSA finds the global public key for the given key ID
func FindGlobalPublicKeyRSA(keyID string) (string, error) {
	if customPublicKeyRSA != nil {
		// If custom key is set, use it
		pubBytes, err := x509.MarshalPKIXPublicKey(customPublicKeyRSA)
		if err != nil {
			return "", errors.Wrap(err, "failed to marshal custom public key")
		}
		return string(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubBytes,
		})), nil
	}

	globalKey, ok := PublicKeysRSA[keyID]
	if !ok {
		return "", errors.Errorf("global public key not found for key ID: %s", keyID)
	}

	return string(globalKey), nil
}

// DecodeLicenseSignature decodes a base64-encoded signature and returns the outer and inner signature structures
// along with the parsed app signing keys
func DecodeLicenseSignature(signature []byte) (*OuterSignature, *InnerSignature, *AppSigningKeys, error) {
	// Unmarshal outer signature
	var outerSig OuterSignature
	if err := json.Unmarshal(signature, &outerSig); err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to unmarshal outer signature")
	}

	// Unmarshal inner signature
	var innerSig InnerSignature
	if err := json.Unmarshal(outerSig.InnerSignature, &innerSig); err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to unmarshal inner signature")
	}

	// Parse the app public key, which is used to sign entitlement signatures
	appPubKey, err := parsePublicKeyRSA(innerSig.PublicKey)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to parse app public key")
	}

	appKeys := &AppSigningKeys{
		PublicKeyRSA: appPubKey,
	}

	return &outerSig, &innerSig, appKeys, nil
}
