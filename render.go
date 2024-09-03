package iss

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const DefaultTagKey string = "json"

//const DefaultTagKey string = "csv"

// Иногда приходят null значение
func parseStringWithDefaultValue(fieldValue interface{}) string {
	if fieldValue == nil {
		return ""
	}
	return fieldValue.(string)
}

// setFloatField
func parseFloatWithDefaultValue(fieldValue interface{}) float64 {
	if fieldValue == nil {
		return 0
	}
	return fieldValue.(float64)
}

func parseIntWithDefaultValue(fieldValue interface{}) int {
	if fieldValue == nil {
		return 0
	}
	return int(fieldValue.(float64))
}
func parseInt64WithDefaultValue(fieldValue interface{}) int64 {
	if fieldValue == nil {
		return 0
	}
	return int64(fieldValue.(float64))
}

// Unmarshal парсинг массивов. По аналогии с csv
func Unmarshal(header []string, data [][]interface{}, destination interface{}) error {
	// получим значения
	sliceValPtr := reflect.ValueOf(destination)
	if sliceValPtr.Kind() != reflect.Ptr {
		return fmt.Errorf("must be a pointer to a slice of structs")
	}

	sliceVal := sliceValPtr.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("must be a pointer to a slice of structs")
	}

	structType := sliceVal.Type().Elem()
	if structType.Kind() != reflect.Struct {
		return fmt.Errorf("must be a pointer to a slice of structs")
	}

	// создадим map с названием колонок
	headerMap := make(map[string]int, len(header))
	for k, name := range header {
		headerMap[name] = k
	}
	// в цикле по данным
	for _, row := range data {
		newValue := reflect.New(structType).Elem()
		//slog.Info("Unmarshal", slog.Any("newVal", newVal))
		// парсим одну строку
		err := unmarshalRow(headerMap, row, newValue, DefaultTagKey)
		if err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, newValue))
	}
	return nil
}

// unmarshalRow распарсим заданную строку
func unmarshalRow(headerMap map[string]int, rowData []interface{}, vv reflect.Value, tagKey string) error {
	vt := vv.Type()
	for i := 0; i < vv.NumField(); i++ {
		typeField := vt.Field(i)
		name := typeField.Tag.Get(tagKey)
		if name == "" {
			name = typeField.Name
		}
		pos, ok := headerMap[name]
		if !ok {
			continue
		}
		val := rowData[pos]
		field := vv.Field(i)
		//slog.Info(" unmarshalOne ", "pos", pos, slog.Any("val", val), slog.Any("field", field), "valRv.Kind()", field.Kind())
		// TODO добавить setField
		switch field.Kind() {
		case reflect.Float64, reflect.Float32:
			field.SetFloat(parseFloatWithDefaultValue(val))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(parseInt64WithDefaultValue(val))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			//field.SetUint(i)
		case reflect.String:
			field.SetString(parseStringWithDefaultValue(val))
		case reflect.Bool:
			//b, err := strconv.ParseBool(val)
			//if err != nil {
			//	return err
			//}
			//field.SetBool(b)
		default:
			return fmt.Errorf("cannot handle field of kind %v", field.Kind())
		}

	}
	return nil
}

// код из проекта go-csv-tag
// https://github.com/artonge/go-csv-tag/blob/master/load.go

// Map the provided content to the destination using the header and the tags.
// @param header: the csv header to match with the struct's tags.
// @param content: the content to put in destination.
// @param destination: the destination where to put the file's content.

func UnmarshalCSV(header []string, content [][]string, destination interface{}, tagKey string) error {
	//func mapToDestination(header []string, content [][]string, destination interface{}, tagKey string) error {
	if destination == nil {
		return fmt.Errorf("destination slice is nil")
	}

	if reflect.TypeOf(destination).Elem().Kind() != reflect.Slice {
		return fmt.Errorf("destination is not a slice")
	}

	// создадим map с названием колонок
	headerMap := make(map[string]int)
	for i, name := range header {
		headerMap[strings.TrimSpace(name)] = i
	}

	// Create the slice to put the values in.
	sliceRv := reflect.MakeSlice(
		reflect.ValueOf(destination).Elem().Type(),
		len(content),
		len(content),
	)

	for i, line := range content {
		emptyStruct := sliceRv.Index(i)

		for j := 0; j < emptyStruct.NumField(); j++ {
			propertyTag := emptyStruct.Type().Field(j).Tag.Get(tagKey)
			if propertyTag == "" {
				continue
			}

			propertyPosition, ok := headerMap[propertyTag]
			if !ok {
				continue
			}

			err := storeValue(line[propertyPosition], emptyStruct.Field(j))
			if err != nil {
				return fmt.Errorf("line: %v to slice: %v:\n	==> %v", line, emptyStruct, err)
			}
		}
	}

	reflect.ValueOf(destination).Elem().Set(sliceRv)

	return nil
}

func storeValue(rawValue string, valRv reflect.Value) error {
	rawValue = strings.TrimSpace(rawValue)
	switch valRv.Kind() {
	case reflect.String:
		valRv.SetString(rawValue)
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uint:
		value, err := strconv.ParseUint(rawValue, 10, 64)
		if err != nil && rawValue != "" {
			return fmt.Errorf("error parsing uint '%v':\n	==> %v", rawValue, err)
		}
		valRv.SetUint(value)
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int:
		i, err := toInt(rawValue)
		if err != nil {
			return fmt.Errorf("error parsing int '%v':\n	==> %v", rawValue, err)
		}
		valRv.SetInt(i)
		//value, err := strconv.ParseInt(rawValue, 10, 64)
		//if err != nil && rawValue != "" {
		//	return fmt.Errorf("error parsing int '%v':\n	==> %v", rawValue, err)
		//}
		//valRv.SetInt(value)
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		f, err := toFloat(rawValue)
		if err != nil {
			return fmt.Errorf("error parsing float '%v':\n	==> %v", rawValue, err)
		}
		valRv.SetFloat(f)
		//value, err := strconv.ParseFloat(rawValue, 64)
		//if err != nil && rawValue != "" {
		//	return fmt.Errorf("error parsing float '%v':\n	==> %v", rawValue, err)
		//}
		//valRv.SetFloat(value)
	case reflect.Bool:
		value, err := strconv.ParseBool(rawValue)
		if err != nil && rawValue != "" {
			return fmt.Errorf("error parsing bool '%v':\n	==> %v", rawValue, err)
		}
		valRv.SetBool(value)
	}

	return nil
}

// из проекта "github.com/gocarina/gocsv"
func toFloat(in interface{}) (float64, error) {
	//if in == nil {
	//	return 0, nil
	//}
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := strings.TrimSpace(inValue.String())
		if s == "" {
			return 0, nil
		}
		s = strings.Replace(s, ",", ".", -1)
		return strconv.ParseFloat(s, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(inValue.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return inValue.Float(), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to float")
}

func toInt(in interface{}) (int64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := strings.TrimSpace(inValue.String())
		if s == "" {
			return 0, nil
		}
		out := strings.SplitN(s, ".", 2)
		return strconv.ParseInt(out[0], 0, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return inValue.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(inValue.Float()), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to int")
}

func toString(in interface{}) (string, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		return inValue.String(), nil
	case reflect.Bool:
		b := inValue.Bool()
		if b {
			return "true", nil
		}
		return "false", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%v", inValue.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%v", inValue.Uint()), nil
	case reflect.Float32:
		return strconv.FormatFloat(inValue.Float(), byte('f'), -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(inValue.Float(), byte('f'), -1, 64), nil
	}
	return "", fmt.Errorf("No known conversion from " + inValue.Type().String() + " to string")
}
