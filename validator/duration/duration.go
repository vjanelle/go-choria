// Copyright (c) 2020-2021, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

package duration

import (
	"fmt"
	"reflect"
	"time"

	"github.com/choria-io/go-choria/validator/registry"
)

// ValidateString validates that input is a valid duration
func ValidateString(input string) (bool, error) {
	_, err := time.ParseDuration(input)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ValidateStructField validates a struct field holds a valid duration
func ValidateStructField(value reflect.Value, tag string) (bool, error) {
	if value.Kind() != reflect.String {
		return false, fmt.Errorf("only strings can be Duration validated")
	}

	return ValidateString(value.String())
}

// durationValidator implements the registry.Validator interface
type durationValidator struct{}

func (v *durationValidator) Name() string {
	return "duration"
}

func (v *durationValidator) Matches(tag string) bool {
	return tag == "duration"
}

func (v *durationValidator) Validate(value reflect.Value, tag string) (bool, error) {
	return ValidateStructField(value, tag)
}

func init() {
	registry.Register(&durationValidator{})
}
