package sanitize

import (
	"reflect"

	"github.com/microcosm-cc/bluemonday"
)

var strictPolicy = bluemonday.StrictPolicy() //strips all HTML

func SanitizeString(s string) string {
	return strictPolicy.Sanitize(s)
}

// sanitize all string fields in a struct using reflection
func SanitizeStruct(v any) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String && field.CanSet() {
			field.SetString(SanitizeString(field.String()))
		}

		if field.Kind() == reflect.Struct && field.CanSet() {
			SanitizeStruct(field.Addr().Interface())
		}
	}

}
