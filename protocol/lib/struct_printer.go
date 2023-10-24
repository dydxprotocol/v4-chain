package lib

import (
	"fmt"
	"reflect"
	"strings"
)

// GetStructFieldsString prints the fields of a struct with their values. This function is useful for
// logging certain proto messages which contain a customtype. The String() defined by the custom type
// is skipped if the top level String() method of the parent proto is called, so here we deconstruct the
// proto into individual fields so that the custom type's String() method is called. This method was created
// specifically with the SerializableInt custom type in mind.
func GetStructFieldsString(i interface{}) string {
	v := reflect.ValueOf(i)

	// If the provided value is a pointer, dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	if t.Kind() != reflect.Struct {
		return "provided value is not a struct"
	}

	var fieldStrings []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s: %v", field.Name, value))
	}

	return strings.Join([]string{"(", strings.Join(fieldStrings, ", "), ")"}, "")
}
