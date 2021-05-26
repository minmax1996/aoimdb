package datatypes

import (
	"errors"
	"reflect"
	"strconv"
)

var Int = reflect.TypeOf(int(1))
var Int32 = reflect.TypeOf(int32(1))
var Int64 = reflect.TypeOf(int64(1))
var String = reflect.TypeOf(string(""))
var Float32 = reflect.TypeOf(float32(1.0))
var Float64 = reflect.TypeOf(float64(1.0))

func ToString(t reflect.Type) string {
	return t.Kind().String()
}

func FromString(s string) reflect.Type {
	switch s {
	case reflect.String.String():
		return String
	case reflect.Int.String():
		return Int
	case reflect.Int32.String():
		return Int32
	case reflect.Int64.String():
		return Int64
	case reflect.Float32.String():
		return Float32
	case reflect.Float64.String():
		return Float64
	default:
		return nil
	}
}

func ConvertToColumnType(value interface{}, resultType reflect.Kind) (interface{}, error) {
	switch v := value.(type) {
	case string:
		switch resultType {
		case reflect.String:
			return v, nil
		case reflect.Int:
			if c, err := strconv.ParseInt(v, 10, 32); err != nil {
				return nil, err
			} else {
				return int(c), err
			}
		case reflect.Int32:
			if c, err := strconv.ParseInt(v, 10, 32); err != nil {
				return nil, err
			} else {
				return int32(c), err
			}
		case reflect.Int64:
			if c, err := strconv.ParseInt(v, 10, 64); err != nil {
				return nil, err
			} else {
				return int64(c), err
			}
		case reflect.Float32:
			if c, err := strconv.ParseFloat(v, 32); err != nil {
				return nil, err
			} else {
				return float32(c), err
			}
		case reflect.Float64:
			if c, err := strconv.ParseFloat(v, 64); err != nil {
				return nil, err
			} else {
				return float64(c), err
			}
		default:
			return nil, errors.New("unknown reflect type")
		}
	default:
		return nil, errors.New("unknown type")
	}
}
