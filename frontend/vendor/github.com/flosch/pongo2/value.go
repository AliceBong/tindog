package pongo2

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Value struct {
	val  reflect.Value
	safe bool // used to indicate whether a Value needs explicit escaping in the template
}

// Converts any given value to a pongo2.Value
// Usually being used within own functions passed to a template
// through a Context or within filter functions.
//
// Example:
//     AsValue("my string")
func AsValue(i interface{}) *Value {
	return &Value{
		val: reflect.ValueOf(i),
	}
}

// Like AsValue, but does not apply the 'escape' filter.
func AsSafeValue(i interface{}) *Value {
	return &Value{
		val:  reflect.ValueOf(i),
		safe: true,
	}
}

func (v *Value) getResolvedValue() reflect.Value {
	if v.val.IsValid() && v.val.Kind() == reflect.Ptr {
		return v.val.Elem()
	}
	return v.val
}

// Checks whether the underlying value is a string
func (v *Value) IsString() bool {
	return v.getResolvedValue().Kind() == reflect.String
}

// Checks whether the underlying value is a bool
func (v *Value) IsBool() bool {
	return v.getResolvedValue().Kind() == reflect.Bool
}

// Checks whether the underlying value is a float
func (v *Value) IsFloat() bool {
	return v.getResolvedValue().Kind() == reflect.Float32 ||
		v.getResolvedValue().Kind() == reflect.Float64
}

// Checks whether the underlying value is an integer
func (v *Value) IsInteger() bool {
	return v.getResolvedValue().Kind() == reflect.Int ||
		v.getResolvedValue().Kind() == reflect.Int8 ||
		v.getResolvedValue().Kind() == reflect.Int16 ||
		v.getResolvedValue().Kind() == reflect.Int32 ||
		v.getResolvedValue().Kind() == reflect.Int64 ||
		v.getResolvedValue().Kind() == reflect.Uint ||
		v.getResolvedValue().Kind() == reflect.Uint8 ||
		v.getResolvedValue().Kind() == reflect.Uint16 ||
		v.getResolvedValue().Kind() == reflect.Uint32 ||
		v.getResolvedValue().Kind() == reflect.Uint64
}

// Checks whether the underlying value is either an integer
// or a float.
func (v *Value) IsNumber() bool {
	return v.IsInteger() || v.IsFloat()
}

// Checks whether the underlying value is NIL
func (v *Value) IsNil() bool {
	//fmt.Printf("%+v\n", v.getResolvedValue().Type().String())
	return !v.getResolvedValue().IsValid()
}

// Returns a string for the underlying value. If this value is not
// of type string, pongo2 tries to convert it. Currently the following
// types for underlying values are supported:
//
//     1. string
//     2. int/uint (any size)
//     3. float (any precision)
//     4. bool
//     5. time.Time
//     6. String() will be called on the underlying value if provided
//
// NIL values will lead to an empty string. Unsupported types are leading
// to their respective type name.
func (v *Value) String() string {
	if v.IsNil() {
		return ""
	}

	switch v.getResolvedValue().Kind() {
	case reflect.String:
		return v.getResolvedValue().String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.getResolvedValue().Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.getResolvedValue().Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v.getResolvedValue().Float())
	case reflect.Bool:
		if v.Bool() {
			return "True"
		} else {
			return "False"
		}
	case reflect.Struct:
		if t, ok := v.Interface().(fmt.Stringer); ok {
			return t.String()
		}
	}

	logf("Value.String() n