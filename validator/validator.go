// Copyright (c) 2018-2022, R.I. Pienaar and the Choria Project contributors
//
// SPDX-License-Identifier: Apache-2.0

/*
Package validator provides common validation helpers commonly used
in operations tools.  Additionally structures can be marked up with
tags indicating the validation of individual keys and the entire struct
can be validated in one go
*/
package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/choria-io/go-choria/validator/registry"

	// Import all validator subpackages to ensure their init() functions run
	_ "github.com/choria-io/go-choria/validator/duration"
	_ "github.com/choria-io/go-choria/validator/enum"
	_ "github.com/choria-io/go-choria/validator/ipaddress"
	_ "github.com/choria-io/go-choria/validator/ipv4"
	_ "github.com/choria-io/go-choria/validator/ipv6"
	_ "github.com/choria-io/go-choria/validator/maxlength"
	_ "github.com/choria-io/go-choria/validator/regex"
	_ "github.com/choria-io/go-choria/validator/shellsafe"
)

// ValidateStruct validates all keys in a struct using their validate tag
func ValidateStruct(target any) (bool, error) {
	val := reflect.ValueOf(target)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return validateStructValue(val)
}

// ValidateStructField validates one field in a struct
func ValidateStructField(target any, field string) (bool, error) {
	val := reflect.ValueOf(target)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	valueField := val.FieldByName(field)
	typeField, ok := val.Type().FieldByName(field)
	if !ok {
		return false, fmt.Errorf("unknown field %s", field)
	}

	validation := strings.TrimSpace(typeField.Tag.Get("validate"))

	err := validateStructField(valueField, typeField, validation)
	if err != nil {
		return false, err
	}

	return true, nil
}

func validateStructValue(val reflect.Value) (bool, error) {
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		validation := strings.TrimSpace(typeField.Tag.Get("validate"))

		err := validateStructField(valueField, typeField, validation)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func validateStructField(valueField reflect.Value, typeField reflect.StructField, validation string) error {
	if valueField.Kind() == reflect.Struct {
		ok, err := validateStructValue(valueField)
		if !ok {
			return err
		}
	}

	if validation == "" {
		return nil
	}

	// Find a validator that matches this validation tag
	validator := registry.FindValidator(validation)
	if validator == nil {
		// No validator found for this tag - silently ignore for backwards compatibility
		return nil
	}

	// Run the validation
	ok, err := validator.Validate(valueField, validation)
	if !ok {
		return fmt.Errorf("%s %s validation failed: %s", typeField.Name, validator.Name(), err)
	}

	return nil
}
