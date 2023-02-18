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

	logf("Value.String() not implemented for type: %s\n", v.getResolvedValue().Kind().String())
	return v.getResolvedValue().String()
}

// Returns the underlying value as an integer (converts the underlying
// value, if necessary). If it's not possible to convert the underlying value,
// it will return 0.
func (v *Value) Integer() int {
	switch v.getResolvedValue().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(v.getResolvedValue().Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(v.getResolvedValue().Uint())
	case reflect.Float32, reflect.Float64:
		return int(v.getResolvedValue().Float())
	case reflect.String:
		// Try to convert from string to int (base 10)
		f, err := strconv.ParseFloat(v.getResolvedValue().String(), 64)
		if err != nil {
			return 0
		}
		return int(f)
	default:
		logf("Value.Integer() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return 0
	}
}

// Returns the underlying value as a float (converts the underlying
// value, if necessary). If it's not possible to convert the underlying value,
// it will return 0.0.
func (v *Value) Float() float64 {
	switch v.getResolvedValue().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.getResolvedValue().Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.getResolvedValue().Uint())
	case reflect.Float32, reflect.Float64:
		return v.getResolvedValue().Float()
	case reflect.String:
		// Try to convert from string to float64 (base 10)
		f, err := strconv.ParseFloat(v.getResolvedValue().String(), 64)
		if err != nil {
			return 0.0
		}
		return f
	default:
		logf("Value.Float() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return 0.0
	}
}

// Returns the underlying value as bool. If the value is not bool, false
// will always be returned. If you're looking for true/false-evaluation of the
// underlying value, have a look on the IsTrue()-function.
func (v *Value) Bool() bool {
	switch v.getResolvedValue().Kind() {
	case reflect.Bool:
		return v.getResolvedValue().Bool()
	default:
		logf("Value.Bool() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return false
	}
}

// Tries to evaluate the underlying value the Pythonic-way:
//
// Returns TRUE in one the following cases:
//
//     * int != 0
//     * uint != 0
//     * float != 0.0
//     * len(array/chan/map/slice/string) > 0
//     * bool == true
//     * underlying value is a struct
//
// Otherwise returns always FALSE.
func (v *Value) IsTrue() bool {
	switch v.getResolvedValue().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.getResolvedValue().Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.getResolvedValue().Uint() != 0
	case reflect.Float32, reflect.Float64:
		return v.getResolvedValue().Float() != 0
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.getResolvedValue().Len() > 0
	case reflect.Bool:
		return v.getResolvedValue().Bool()
	case reflect.Struct:
		return true // struct instance is always true
	default:
		logf("Value.IsTrue() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return false
	}
}

// Tries to negate the underlying value. It's mainly used for
// the NOT-operator and in conjunction with a call to
// return_value.IsTrue() afterwards.
//
// Example:
//     AsValue(1).Negate().IsTrue() == false
func (v *Value) Negate() *Value {
	switch v.getResolvedValue().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Integer() != 0 {
			return AsValue(0)
		} else {
			return AsValue(1)
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() != 0.0 {
			return AsValue(float64(0.0))
		} else {
			return AsValue(float64(1.1))
		}
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return AsValue(v.getResolvedValue().Len() == 0)
	case reflect.Bool:
		return AsValue(!v.getResolvedValue().Bool())
	default:
		logf("Value.IsTrue() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return AsValue(true)
	}
}

// Returns the length for an array, chan, map, slice or string.
// Otherwise it will return 0.
func (v *Value) Len() int {
	switch v.getResolvedValue().Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return v.getResolvedValue().Len()
	case reflect.String:
		runes := []rune(v.getResolvedValue().String())
		return len(runes)
	default:
		logf("Value.Len() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return 0
	}
}

// Slices an array, slice 