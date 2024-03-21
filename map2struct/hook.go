package map2struct

import (
	"github.com/jackc/pgtype"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"gopkg.in/guregu/null.v4"
	"reflect"
	"time"
)

var timeLocation *time.Location

func ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return cast.ToTimeInDefaultLocationE(data, timeLocation)
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
		// Convert it by parsing
	}
}

func ToTimePtrHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(&time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			if data == "" {
				return nil, nil
			}
			v, err := cast.ToTimeInDefaultLocationE(data, timeLocation)
			return &v, err
		case reflect.Float64:
			v := time.Unix(0, int64(data.(float64))*int64(time.Millisecond))
			return &v, nil
		case reflect.Int64:
			v := time.Unix(0, data.(int64)*int64(time.Millisecond))
			return &v, nil
		default:
			return data, nil
		}
		// Convert it by parsing
	}
}

func ToDecimalHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {

		if t == reflect.TypeOf(decimal.Decimal{}) {
			if f.Kind() == reflect.Float64 {
				return decimal.NewFromFloat(data.(float64)), nil
			}

			if input := data.(string); input != "" {
				return decimal.NewFromString(input)
			}
			return decimal.Decimal{}, nil
		}

		return data, nil
	}
}

func ToPgDateHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if t == reflect.TypeOf(pgtype.Date{}) {
			date := pgtype.Date{}
			_ = date.Set(data)
			return date, nil
		}

		return data, nil
	}
}

func ToPgDatePtrHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if t == reflect.TypeOf(&pgtype.Date{}) {
			date := &pgtype.Date{}
			_ = date.Set(data)
			return date, nil
		}

		return data, nil
	}
}

func ToNullableStringHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if t == reflect.TypeOf(null.String{}) {
			return null.StringFrom(data.(string)), nil
		}

		return data, nil
	}
}
