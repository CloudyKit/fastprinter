package fastprinter

import (
	"fmt"
	"io"
	"reflect"
)

var (
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
)

func PrintValue(w io.Writer, v reflect.Value) (int, error) {
	t := v.Type()
	k := t.Kind()

	if t.Implements(fmtStringerType) {
		return PrintString(w, v.Interface().(fmt.Stringer).String())
	}

	if t.Implements(errorType) {
		return PrintString(w, v.Interface().(error).Error())
	}

	if k == reflect.String {
		return PrintString(w, v.String())
	}

	if k >= reflect.Int && k <= reflect.Int64 {
		return PrintInt(w, v.Int())
	}

	if k >= reflect.Uint && k <= reflect.Uint64 {
		return PrintUint(w, v.Uint())
	}

	if k == reflect.Float64 || k == reflect.Float64 {
		return PrintFloat(w, v.Float())
	}

	if k == reflect.Bool {
		return PrintBool(w, v.Bool())
	}

	if k == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		return w.Write(v.Bytes())
	}
	return fmt.Fprint(w, v.Interface())
}

func Print(w io.Writer, i interface{}) (int, error) {

	v := reflect.ValueOf(i)
	t := v.Type()
	k := t.Kind()

	if t.Implements(fmtStringerType) {
		return PrintString(w, i.(fmt.Stringer).String())
	}

	if t.Implements(errorType) {
		return PrintString(w, i.(error).Error())
	}

	if k == reflect.String {
		return PrintString(w, v.String())
	}

	if k >= reflect.Int && k <= reflect.Int64 {
		return PrintInt(w, v.Int())
	}

	if k >= reflect.Uint && k <= reflect.Uint64 {
		return PrintUint(w, v.Uint())
	}

	if k == reflect.Bool {
		return PrintBool(w, v.Bool())
	}

	return fmt.Print(i)
}
