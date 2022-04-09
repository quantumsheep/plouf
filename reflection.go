package plouf

import "reflect"

func reflectTypeName(i interface{}) string {
	typ := reflect.TypeOf(i)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ.Name()
}

func reflectValue(i interface{}) reflect.Value {
	value := reflect.ValueOf(i)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return value
}
