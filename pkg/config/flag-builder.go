package config

import (
	"errors"
	"reflect"
	"strconv"
)

// ErrInvalidSpecification indicates that a specification is of the wrong type.
var ErrInvalidSpecification = errors.New("specification must be a struct pointer")

// HydrateFromFlags hydrate object form flags
func HydrateFromFlags(conf interface{}) error {
	rv := reflect.ValueOf(conf)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	typ := rv.Type()

	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		field := rv.Field(i)

		var key, desc, val string
		if tag, ok := fieldType.Tag.Lookup("flag"); ok {
			key = tag
		} else {
			continue
		}
		if tag, ok := fieldType.Tag.Lookup("desc"); ok {
			desc = tag
		}
		if tag, ok := fieldType.Tag.Lookup("default"); ok {
			val = tag
		}

		switch fieldType.Type.Kind() {
		case reflect.String:
			field.SetString(val)
			ptr, _ := field.Addr().Interface().(*string)
			setFlagEnvString(ptr, key, desc, val)
		case reflect.Bool:
			bVal, err := strconv.ParseBool(val)
			if err != nil {
				return err
			}
			field.SetBool(bVal)
			ptr, _ := field.Addr().Interface().(*bool)
			setFlagEnvBool(ptr, key, desc, bVal)
		case reflect.Int:
			i64Val, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				return err
			}
			iVal := int(i64Val)
			field.SetInt(i64Val)
			ptr, _ := field.Addr().Interface().(*int)
			setFlagEnvInt(ptr, key, desc, iVal)
		}
	}
	return nil
}
