package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New() // Initialize validator
}
func ValidateRequiredFields(data interface{}) string {
	var allErrors []string

	// Gunakan json tag sebagai nama field
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	value := reflect.ValueOf(data)

	switch value.Kind() {
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			err := validate.Struct(value.Index(i).Interface())
			if err != nil {
				for _, e := range err.(validator.ValidationErrors) {
					field := e.Field()
					if e.Tag() == "required" {
						allErrors = append(allErrors, fmt.Sprintf("Index %d - Field '%s' is required", i, field))
					} else {
						allErrors = append(allErrors, fmt.Sprintf("Index %d - Field '%s' failed on '%s'", i, field, e.Tag()))
					}
				}
			}
		}
	case reflect.Struct:
		err := validate.Struct(data)
		if err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				field := e.Field()
				if e.Tag() == "required" {
					allErrors = append(allErrors, fmt.Sprintf("Field '%s' is required", field))
				} else {
					allErrors = append(allErrors, fmt.Sprintf("Field '%s' failed on '%s'", field, e.Tag()))
				}
			}
		}
	default:
		return "Invalid input"
	}

	return strings.Join(allErrors, "\n")
}

type DynamicStruct map[string]interface{}

func RemoveField(original any, fieldName string) DynamicStruct {
	// Konversi struct ke dalam bentuk JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		return nil
	}

	// Unmarshal JSON ke dalam map
	var resultMap map[string]interface{}
	err = json.Unmarshal(jsonData, &resultMap)
	if err != nil {
		return nil
	}

	// Hapus field yang diinginkan
	delete(resultMap, fieldName)

	return resultMap
}
func ValidateRequiredFieldsOld(data interface{}) string {
	var allErrors []string

	value := reflect.ValueOf(data)

	if value.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i).Interface()
			err := validate.Struct(elem)

			if err != nil {
				if _, ok := err.(*validator.InvalidValidationError); ok {
					return "Invalid input"
				}
				var missingFields []string
				for _, err := range err.(validator.ValidationErrors) {
					fieldName := err.StructField()
					missingFields = append(missingFields, "'"+fieldName+"'")
				}
				missingFieldsStr := strings.Join(missingFields, ", ")
				if len(missingFields) > 1 {
					allErrors = append(allErrors, fmt.Sprintf("Index %d: %s fields are required and cannot be empty", i, missingFieldsStr))
				} else {
					allErrors = append(allErrors, fmt.Sprintf("Index %d: %s field is required and cannot be empty", i, missingFieldsStr))
				}
			}
		}
	} else if value.Kind() == reflect.Struct {
		// Jika data bukan array, validasi objek tunggal
		err := validate.Struct(data)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				return "Invalid input"
			}
			var missingFields []string
			for _, err := range err.(validator.ValidationErrors) {
				fieldName := err.StructField()
				missingFields = append(missingFields, "'"+fieldName+"'")
			}
			missingFieldsStr := strings.Join(missingFields, ", ")
			if len(missingFields) > 1 {
				return missingFieldsStr + " fields are required and cannot be empty"
			} else {
				return missingFieldsStr + " field is required and cannot be empty"
			}
		}
	} else {
		return "Invalid input"
	}

	return strings.Join(allErrors, "\n")
}

func GetIdField(entity interface{}) string {
	v := reflect.ValueOf(entity)

	// Kalau pointer → ambil elemennya
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Pastikan struct
	if v.Kind() != reflect.Struct {
		return ""
	}

	t := v.Type()

	// Loop semua field
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		// Cari field dengan nama persis "Id" atau tag json:"_id"
		if field.Name == "Id" {
			val := v.Field(i)
			if val.Kind() == reflect.String {
				return val.String()
			}
		}

		if tag := field.Tag.Get("json"); tag == "_id" {
			val := v.Field(i)
			if val.Kind() == reflect.String {
				return val.String()
			}
		}
	}

	return ""
}
