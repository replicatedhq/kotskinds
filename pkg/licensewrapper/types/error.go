package types

import (
	"errors"
	"fmt"
)

type LicenseDataValidationError struct {
	FieldName   string
	SignedValue string
	ActualValue string
}

func (e *LicenseDataValidationError) Error() string {
	return fmt.Sprintf("license data validation error: %s field has changed to %q (license) from %q (within signature)", e.FieldName, e.ActualValue, e.SignedValue)
}

// return true if the error is a LicenseDataValidationError
func (e *LicenseDataValidationError) Is(target error) bool {
	_, ok := target.(*LicenseDataValidationError)
	return ok
}

func IsLicenseDataValidationError(err error) bool {
	var ldve *LicenseDataValidationError
	return errors.As(err, &ldve)
}
