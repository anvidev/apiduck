package apiduck

import (
	"reflect"
	"strings"
)

// structToSlice converts a struct into a slice of body fields. It handles nested structs and slices recursively.
func structToSlice(i interface{}) []BodyField {
	reflectValue := reflect.ValueOf(i)
	if reflectValue.Kind() == reflect.Pointer {
		// if pointer, dereference it
		reflectValue = reflectValue.Elem()
	}
	if reflectValue.Kind() != reflect.Struct {
		return nil
	}

	reflectType := reflectValue.Type()
	fieldCount := reflectValue.NumField()
	var result []BodyField

	for i := range fieldCount {
		field := reflectValue.Field(i)
		fieldType := reflectType.Field(i)

		jsonTag := fieldType.Tag.Get("json")
		description := fieldType.Tag.Get("description")
		validateTag := fieldType.Tag.Get("validate")
		apidocTag := fieldType.Tag.Get("apidoc")
		exampleTag := fieldType.Tag.Get("example")
		enumTag := fieldType.Tag.Get("enum")
		defaultTag := fieldType.Tag.Get("default")

		if !fieldType.IsExported() || jsonTag == "-" || strings.Contains(apidocTag, "ignore") {
			continue
		}

		fieldName := fieldType.Name
		fieldKind := field.Kind()
		jsonName := fieldName

		if jsonTag != "" {
			jsonName = strings.Split(jsonTag, ",")[0]
		}

		validation := parseValidationRules(validateTag)
		bodyFieldType := getTypeString(fieldType.Type)

		bodyField := BodyField{
			Name:        jsonName,
			Type:        bodyFieldType,
			JSONName:    jsonName,
			Required:    strings.Contains(validateTag, "required"),
			Description: description,
			Validation:  validation,
		}

		// Handle examples
		if exampleTag != "" {
			bodyField.Example = parseExampleValue(exampleTag, bodyFieldType)
		}

		// Handle enums
		if enumTag != "" {
			bodyField.Enum = parseEnumValues(enumTag)
		}

		// Handle defaults
		if defaultTag != "" {
			bodyField.Default = parseExampleValue(defaultTag, bodyFieldType)
		}

		// Handle nested structures
		if fieldKind == reflect.Struct || (fieldKind == reflect.Pointer && field.Elem().Kind() == reflect.Struct) {
			bodyField.Fields = structToSlice(field.Interface())
		} else if fieldKind == reflect.Slice {
			sliceType := fieldType.Type.Elem()
			bodyField.Type = "[]" + getTypeString(sliceType)

			if sliceType.Kind() == reflect.Struct || (sliceType.Kind() == reflect.Pointer && sliceType.Elem().Kind() == reflect.Struct) {
				bodyField.Fields = structToSlice(reflect.New(sliceType).Interface())
			}
		}

		result = append(result, bodyField)
	}

	return result
}

// getTypeString returns a clean type string
func getTypeString(t reflect.Type) string {
	typeStr := t.String()
	if strings.Contains(typeStr, ".") {
		parts := strings.Split(typeStr, ".")
		return parts[len(parts)-1]
	}
	return typeStr
}

// parseValidationRules parses validation tags into a map
func parseValidationRules(validateTag string) map[string]string {
	rules := make(map[string]string)
	if validateTag == "" {
		return rules
	}

	for _, rule := range strings.Split(validateTag, ",") {
		rule := strings.TrimSpace(rule)
		parts := strings.SplitN(rule, "=", 2)
		if len(parts) == 2 {
			rules[parts[0]] = parts[1]
		} else if len(parts) == 1 {
			// Handle simple flags like "required"
			rules[parts[0]] = "true"
		}
	}
	return rules
}

// parseExampleValue attempts to parse example values based on type
func parseExampleValue(example, fieldType string) interface{} {
	switch fieldType {
	case "int", "int32", "int64":
		return example
	case "float32", "float64":
		return example
	case "bool":
		return example == "true"
	default:
		return example
	}
}

// parseEnumValues parses comma-separated enum values
func parseEnumValues(enumTag string) []interface{} {
	values := strings.Split(enumTag, ",")
	result := make([]interface{}, len(values))
	for i, v := range values {
		result[i] = strings.TrimSpace(v)
	}
	return result
}
