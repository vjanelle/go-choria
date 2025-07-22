// Copyright (c) 2018-2025, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"reflect"

	"github.com/choria-io/go-choria/validator/registry"
)

// ValidateValue validates any value using the specified validation tag
func ValidateValue[T any](value T, validationTag string) (bool, error) {
	if validationTag == "" {
		return true, nil
	}

	// Find a validator that matches this validation tag
	validator := registry.FindValidator(validationTag)
	if validator == nil {
		// No validator found for this tag
		return false, nil
	}

	// Convert the value to a reflect.Value for the validator
	reflectValue := reflect.ValueOf(value)

	// Run the validation
	return validator.Validate(reflectValue, validationTag)
}
