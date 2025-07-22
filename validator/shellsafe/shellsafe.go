// Copyright (c) 2018-2021, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

package shellsafe

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/choria-io/go-choria/validator/registry"
)

// Validate checks if a string is safe to use in a shell without any escapes or redirects
func Validate(input string) (bool, error) {
	badchars := []string{"`", "$", ";", "|", "&&", ">", "<"}

	for _, c := range badchars {
		if strings.Contains(input, c) {
			return false, fmt.Errorf("may not contain '%s'", c)
		}
	}

	return true, nil
}

// ValidateStructField validates a reflect.Value is shellsafe
func ValidateStructField(value reflect.Value, tag string) (bool, error) {
	if value.Kind() != reflect.String {
		return false, errors.New("should be a string")
	}

	return Validate(value.String())
}

// shellsafeValidator implements the registry.Validator interface
type shellsafeValidator struct{}

func (v *shellsafeValidator) Name() string {
	return "shellsafe"
}

func (v *shellsafeValidator) Matches(tag string) bool {
	return tag == "shellsafe"
}

func (v *shellsafeValidator) Validate(value reflect.Value, tag string) (bool, error) {
	return ValidateStructField(value, tag)
}

func init() {
	registry.Register(&shellsafeValidator{})
}
