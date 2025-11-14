package structs

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/magaluCloud/mgccli/cmd/common/validator"
	"gopkg.in/yaml.v3"
)

func convertStringToType(str string, fieldType reflect.Type) (reflect.Value, error) {
	switch fieldType.Kind() {
	case reflect.String:
		return reflect.ValueOf(str), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("cannot convert %q to %s: %w", str, fieldType.Name(), err)
		}
		return reflect.ValueOf(val).Convert(fieldType), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("cannot convert %q to %s: %w", str, fieldType.Name(), err)
		}
		return reflect.ValueOf(val).Convert(fieldType), nil
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("cannot convert %q to %s: %w", str, fieldType.Name(), err)
		}
		return reflect.ValueOf(val).Convert(fieldType), nil
	case reflect.Bool:
		val, err := strconv.ParseBool(str)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("cannot convert %q to bool: %w", str, err)
		}
		return reflect.ValueOf(val), nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported type for default value: %s", fieldType.Name())
	}
}

func StructToMap[T any](config T) (map[string]any, error) {
	mapConfig := make(map[string]any)
	fields := reflect.ValueOf(config)
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)

		fieldName := fields.Type().Field(i).Tag.Get("yaml")
		fieldName = strings.Split(fieldName, ",")[0]
		if fieldName == "" {
			fieldName = fields.Type().Field(i).Name
		}
		mapConfig[fieldName] = field.Interface()

		if field.IsZero() {
			defaultValue := fields.Type().Field(i).Tag.Get("default")
			if defaultValue != "" {
				mapConfig[fieldName] = defaultValue
			}
		}

		if mapConfig[fieldName] == nil && field.IsZero() {
			continue
		}

		validateTag := fields.Type().Field(i).Tag.Get("validate")
		if validateTag != "" {
			err := validator.NewValidator(mapConfig[fieldName], validateTag).Validate()
			if err != nil {
				return nil, err
			}
		}
	}
	return mapConfig, nil
}

func InitConfig[T any]() T {
	newObject := new(T)
	fields := reflect.ValueOf(newObject).Elem()
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		fieldName := fields.Type().Field(i).Tag.Get("yaml")
		fieldName = strings.Split(fieldName, ",")[0]
		if fieldName == "" {
			fieldName = fields.Type().Field(i).Name
		}
		defaultValue := fields.Type().Field(i).Tag.Get("default")
		if defaultValue != "" {
			convertedValue, err := convertStringToType(defaultValue, field.Type())
			if err != nil {
				// Log error but continue - default value will remain zero value
				// In production, you might want to handle this differently
				continue
			}
			field.Set(convertedValue)
		}
	}
	return *newObject
}

func Set[T any](structPtr *T, name string, value any) error {
	fields := reflect.ValueOf(structPtr).Elem()
	sucess := false
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		fieldName := fields.Type().Field(i).Tag.Get("yaml")
		fieldName = strings.ReplaceAll(fieldName, ",omitempty", "")
		if fieldName == "" {
			fieldName = fields.Type().Field(i).Name
		}
		if fieldName == name {
			if value == nil {
				field.Set(reflect.Zero(field.Type()))
				return nil
			}

			validateTag := fields.Type().Field(i).Tag.Get("validate")
			if validateTag != "" {
				err := validator.NewValidator(value, validateTag).Validate()
				if err != nil {
					return fmt.Errorf("invalid value for config %s: %w", name, err)
				}
			}

			typeOfField := field.Type()
			switch typeOfField.Kind() {
			case reflect.String:
				field.Set(reflect.ValueOf(value.(string)).Convert(typeOfField))
				sucess = true
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				val, err := strconv.ParseInt(value.(string), 10, 64)
				if err != nil {
					return fmt.Errorf("invalid value for config %s: %w", name, err)
				}
				field.Set(reflect.ValueOf(val).Convert(typeOfField))
				sucess = true
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				val, err := strconv.ParseUint(value.(string), 10, 64)
				if err != nil {
					return fmt.Errorf("invalid value for config %s: %w", name, err)
				}
				field.Set(reflect.ValueOf(val).Convert(typeOfField))
				sucess = true
			case reflect.Float32, reflect.Float64:
				field.Set(reflect.ValueOf(value.(float64)).Convert(typeOfField))
				sucess = true
			case reflect.Bool:
				val, err := strconv.ParseBool(value.(string))
				if err != nil {
					return fmt.Errorf("invalid value for config %s: %w", name, err)
				}
				field.Set(reflect.ValueOf(val).Convert(typeOfField))
				sucess = true
			default:
				return fmt.Errorf("unsupported type for value: %s", typeOfField.Name())
			}
		}
	}

	if !sucess {
		return fmt.Errorf("config %s not found", name)
	}
	return nil
}

func LoadFileToStruct[T any](filePath string) (T, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return InitConfig[T](), nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return *new(T), err
	}

	fileContent := new(T)
	err = yaml.Unmarshal(data, fileContent)
	if err != nil {
		return *new(T), err
	}
	return *fileContent, nil
}
