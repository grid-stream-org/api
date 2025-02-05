package logic

import (
	"reflect"
)

// ExtractBody dynamically extracts non-empty fields from a struct, used for extracting body of requests
func ExtractBody(input any) map[string]any {
	updates := make(map[string]any)

	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem() 
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Get the JSON or DB field name (fallback to struct field name)
		fieldName := field.Tag.Get("json")
		if fieldName == "" {
			fieldName = field.Tag.Get("bigquery")
		}
		if fieldName == "" {
			fieldName = field.Name 
		}

		// Only add non-zero values (non-empty strings, non-zero numbers, etc.)
		if !value.IsZero() {
			updates[fieldName] = value.Interface()
		}
	}

	return updates
}
