// Copyright (c) 2020-2021, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

package ipv4

import (
	"fmt"
	"net"
	"reflect"

	"github.com/choria-io/go-choria/validator/registry"
)

// ValidateString validates that the given string is an IPv4 address
func ValidateString(input string) (bool, error) {
	ip := net.ParseIP(input).To4()

	if ip == nil {
		return false, fmt.Errorf("%s is not an IPv4 address", input)
	}

	return true, nil
}

// ValidateStructField validates a struct field holds an IPv4 address
func ValidateStructField(value reflect.Value, tag string) (bool, error) {
	if value.Kind() != reflect.String {
		return false, fmt.Errorf("only strings can be IPv4 validated")
	}

	return ValidateString(value.String())
}

// ipv4Validator implements the registry.Validator interface
type ipv4Validator struct{}

func (v *ipv4Validator) Name() string {
	return "IPv4"
}

func (v *ipv4Validator) Matches(tag string) bool {
	return tag == "ipv4"
}

func (v *ipv4Validator) Validate(value reflect.Value, tag string) (bool, error) {
	return ValidateStructField(value, tag)
}

func init() {
	registry.Register(&ipv4Validator{})
}
