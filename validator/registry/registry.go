// Copyright (c) 2018-2025, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package registry provides a registration mechanism for validators
package registry

import (
	"reflect"
	"strings"
)

// Validator defines the interface that all validators must implement
type Validator interface {
	// Name returns a human-readable name for this validator
	Name() string

	// Matches returns true if this validator should handle the given tag
	Matches(tag string) bool

	// Validate performs the actual validation
	Validate(value reflect.Value, tag string) (bool, error)
}

// validators holds all registered validators
var validators []Validator

// Register adds a validator to the registry
func Register(v Validator) {
	validators = append(validators, v)
}

// FindValidator returns the first validator that matches the given tag
func FindValidator(tag string) Validator {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return nil
	}

	for _, validator := range validators {
		if validator.Matches(tag) {
			return validator
		}
	}

	return nil
}

// AllValidators returns a copy of all registered validators
func AllValidators() []Validator {
	result := make([]Validator, len(validators))
	copy(result, validators)
	return result
}

// ClearValidators removes all registered validators (mainly for testing)
func ClearValidators() {
	validators = nil
}
