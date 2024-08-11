package iss

import (
	"fmt"
	"reflect"
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
	namePos := make(map[string]int, len(header))
	for k, name := range header {
		namePos[name] = k
	}
	// в цикле по данным
	for _, row := range data {
		newVal := reflect.New(structType).Elem()
		//slog.Info("Unmarshal", slog.Any("newVal", newVal))
		err := unmarshalOne(row, namePos, newVal, DefaultTagKey)
		if err != nil {
			return err
		}
		//slog.Debug("Unmarshal", slog.Any("row", row))
		sliceVal.Set(reflect.Append(sliceVal, newVal))
	}
	return nil
}

// unmarshalOne распарсим заданную строку
func unmarshalOne(row []interface{}, namePos map[string]int, vv reflect.Value, tagKey string) error {
	vt := vv.Type()
	for i := 0; i < vv.NumField(); i++ {
		typeField := vt.Field(i)
		name := typeField.Tag.Get(tagKey)
		if name == "" {
			name = typeField.Name
		}
		pos, ok := namePos[name]
		if !ok {
			continue
		}
		val := row[pos]
		field := vv.Field(i)

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

/*
// из проекта https://github.com/gin-gonic/gin/
// https://github.com/gin-gonic/gin/blob/master/binding/form_mapping.go
func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		tv, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}

		d := time.Duration(1)
		if tf == "unixnano" {
			d = time.Second
		}

		t := time.Unix(tv/int64(d), tv%int64(d))
		value.Set(reflect.ValueOf(t))
		return nil
	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}
*/
