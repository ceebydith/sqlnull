package sqlnull

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"
)

// NullValue wraps a target variable to handle SQL null values.
type NullValue struct {
	target any
}

// Scan implements the sql.Scanner interface for NullValue.
func (v *NullValue) Scan(src any) error {
	// Validate the target and create a sql.Scanner.
	null, targetType, err := validate(v.target)
	if err != nil {
		return err
	}

	// Use the sql.Scanner to scan the source value.
	if err := null.Scan(src); err != nil {
		return err
	}

	val := reflect.ValueOf(v.target).Elem()
	if v, err := null.(driver.Valuer).Value(); err != nil {
		return err
	} else if v == nil {
		// Set the target to its zero value if the source is null.
		val.Set(reflect.Zero(targetType.Elem()))
	} else {
		// Convert the value to the target type.
		newval := reflect.ValueOf(v)
		if !val.Elem().CanAddr() {
			val.Set(reflect.New(targetType.Elem().Elem()))
		}
		val.Elem().Set(newval.Convert(targetType.Elem().Elem()))
	}

	return nil
}

// Target returns a NullValue wrapper if the target is valid, otherwise returns the target itself.
func Target(target any) any {
	if target == nil {
		return new(any)
	}
	if _, _, err := validate(target); err == nil {
		return New(target)
	}
	return target
}

// Scanner wraps multiple targets with NullValue.
func Scanner(targets ...any) []any {
	var result []any

	for _, target := range targets {
		result = append(result, Target(target))
	}

	return result
}

// New creates a new NullValue for a given target.
func New(target any) *NullValue {
	return &NullValue{
		target: target,
	}
}

// validate checks if the target type is supported and returns the corresponding sql.Scanner.
func validate(target any) (sql.Scanner, reflect.Type, error) {
	targetType := reflect.TypeOf(target)
	if targetType.Kind() == reflect.Ptr && targetType.Elem().Kind() == reflect.Ptr {
		switch targetType.Elem().Elem().Kind() {
		case reflect.Bool:
			return &sql.NullBool{}, targetType, nil
		case reflect.Uint8:
			return &sql.NullByte{}, targetType, nil
		case reflect.Int8, reflect.Int16, reflect.Uint16:
			return &sql.NullInt16{}, targetType, nil
		case reflect.Int32, reflect.Uint32:
			return &sql.NullInt32{}, targetType, nil
		case reflect.Int64, reflect.Uint64, reflect.Int, reflect.Uint:
			return &sql.NullInt64{}, targetType, nil
		case reflect.String:
			return &sql.NullString{}, targetType, nil
		case reflect.Float32, reflect.Float64:
			return &sql.NullFloat64{}, targetType, nil
		case reflect.Struct:
			if targetType.Elem().Elem() == reflect.TypeOf(time.Time{}) {
				return &sql.NullTime{}, targetType, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("NullValue for %T type is not supported", target)
}