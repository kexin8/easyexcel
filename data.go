package easyexcel

import (
	"reflect"
	"strconv"
)

// rowToStruct 将行数据转换为结构体
func rowToStruct[T any](row []string, fields Fields) (data T, err error) {
	refValue := getStructReflectValue(&data)

	for i, cell := range row {
		f := fields[i]
		if f == nil {
			continue
		}

		fv := refValue.FieldByName(f.name)

		if !fv.CanSet() {
			continue
		}

		if f.readConverters != nil && len(f.readConverters) > 0 {
			if converter, ok := f.readConverters[cell]; ok {
				cell = converter
			}
		}

		if cell == "" {
			continue
		}

		var (
			v string = cell
		)

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(v)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return data, err
			}
			fv.SetInt(val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return data, err
			}
			fv.SetUint(val)
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return data, err
			}
			fv.SetFloat(val)
		case reflect.Bool:
			val, err := strconv.ParseBool(v)
			if err != nil {
				return data, err
			}
			fv.SetBool(val)

		default:
			//TODO: 暂时不支持其他类型
			return
		}
	}
	return
}

// structToRow 将结构体转换为行数据
func structToRow[T any](data T, fields Fields) (row []string, err error) {
	refValue := getStructReflectValue(data)

	row = make([]string, len(fields))
	for _, f := range fields {
		if f == nil {
			continue
		}

		fv := refValue.FieldByName(f.name)

		var (
			v string
		)

		switch fv.Kind() {
		case reflect.String:
			v = fv.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v = strconv.FormatInt(fv.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v = strconv.FormatUint(fv.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			v = strconv.FormatFloat(fv.Float(), 'f', -1, 64)
		case reflect.Bool:
			v = strconv.FormatBool(fv.Bool())
		default:
			//TODO: 暂时不支持其他类型
		}

		if v == "" {
			row = append(row, f.defaultValue)
			continue
		}

		if f.writeConverters != nil && len(f.writeConverters) > 0 {
			if converter, ok := f.writeConverters[v]; ok {
				v = converter
			}
		}

		row[f.sort] = v
	}

	return
}

// getStructReflectValue 获取结构体的反射值
func getStructReflectValue(v any) reflect.Value {
	refValue := reflect.ValueOf(v)

	// 如果是指针类型，获取指针指向的值
	for refValue.Kind() == reflect.Ptr {
		if refValue.IsNil() && refValue.CanAddr() {
			refValue.Set(reflect.New(refValue.Type().Elem()))
		}
		refValue = refValue.Elem()
	}
	return refValue
}
