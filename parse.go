package apiduck

import (
	"reflect"
	"strings"
)

/*
	type field struct {
		fields []field

		name         string
		description  string
		typ          string
		defaultValue any

		example any

		required     bool
		enum      []any
		minimum   *float64
		maximum   *float64
	}
*/
func parseStruct(v any) []Field {
	fields := []Field{}

	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Pointer {
		// dereference pointer
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fields
	}

	reflectType := val.Type()
	numField := val.NumField()

	for i := range numField {
		reflectField := val.Field(i)
		reflectFieldType := reflectType.Field(i)

		jsonTag := reflectFieldType.Tag.Get("json")
		validateTag := reflectFieldType.Tag.Get("validate")
		apiduckTag := reflectFieldType.Tag.Get("apiduck")

		if !reflectFieldType.IsExported() || jsonTag == "-" || jsonTag == "" {
			continue
		}

		fieldType := getTypeString(reflectFieldType.Type)

		field := Field{
			Name: jsonTag,
			Type: fieldType,
		}

		if reflectField.Kind() == reflect.Struct {
			field.Fields = parseStruct(reflectField.Interface())

		}

		parseValidationRules(&field, validateTag)
		parseApiduckTag(&field, apiduckTag)

		fields = append(fields, field)
	}

	return fields
}

func getTypeString(t reflect.Type) string {
	typeStr := t.String()
	if strings.Contains(typeStr, ".") {
		parts := strings.Split(typeStr, ".")
		return parts[len(parts)-1]
	}
	return typeStr
}

func parseValidationRules(field *Field, validateTag string) {
	if validateTag == "" {
		return
	}

	for rule := range strings.SplitSeq(validateTag, ",") {
		rule := strings.TrimSpace(rule)
		parts := strings.SplitN(rule, "=", 2)

		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]

			switch key {
			case "min":
				field.Minimum = &value
			case "max":
				field.Maximum = &value
			case "oneof":
				values := strings.Split(value, " ")
				field.Enums = append(field.Enums, values)
			}
		} else if len(parts) == 1 {
			if parts[0] == "required" {
				field.Req = true
			}
		}
	}
}

func parseApiduckTag(field *Field, tag string) {
	if tag == "" {
		return
	}

	for rule := range strings.SplitSeq(tag, ",") {
		rule := strings.TrimSpace(rule)
		parts := strings.SplitN(rule, "=", 2)

		key := parts[0]
		value := parts[1]

		switch key {
		case "desc":
			field.Desc = value
		case "default":
			field.DefaultValue = value
		case "example":
			field.Ex = value
		}
	}
}

// func parseExampleValue(example, fieldType string) any {
// 	switch fieldType {
// 	case "int", "int32", "int64":
// 		return example
// 	case "float32", "float64":
// 		return example
// 	case "bool":
// 		return example == "true"
// 	default:
// 		return example
// 	}
// }
