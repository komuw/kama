package kama

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
)

func (va vari) dump(val reflect.Value, hideZeroValues bool, indentLevel int) string {
	/*
		`compact` indicates whether the struct should be laid in one line or not
		`hideZeroValues` indicates whether to show zeroValued vars
		`indentLevel` is the number of spaces from the left-most side of the termninal for struct names

		In future(if we ever add compation) we could restrict compaction only to arrays/slices/maps that are of primitive(basic) types
		see:
			1. https://github.com/sanity-io/litter/pull/43
			2. https://github.com/komuw/kama/pull/28
	*/

	deVal := deInterface(val)
	if !deVal.IsValid() {
		return "nil"
	}
	if deVal.Type() == nil {
		return "nil"
	}

	iType := val.Type()
	indentLevel = indentLevel + 1
	iTypeKind := iType.Kind()
	iTypeStr := strings.ReplaceAll(fmt.Sprint(iType), "*", "") // remove the pointer symbol.
	deValStr := fmt.Sprint(deVal)

	{
		if slices.Contains(
			// todo: Ideally this should be handled inside the `dumpInterface` func.
			[]string{
				// This are taken from; https://github.com/golang/go/blob/master/src/context/context.go
				"context.emptyCtx",
				"context.valueCtx",
				"context.backgroundCtx",
				"context.todoCtx",
				"context.withoutCancelCtx",
				"context.timerCtx",
				"context.cancelCtx",
			},
			iTypeStr,
		) {
			// This could be a context.Context type.
			// Let's use the formatting that is provided by the stdlib;
			// https://github.com/golang/go/blob/39effbc105f5c54117a6011af3c48e3c8f14eca9/src/context/context.go#L197-L206
			//
			// This will not handle custom types that implement context.Context
			return deValStr
		}

		if slices.Contains(
			// todo: Ideally this should be handled inside the `dumpInterface` func.
			[]string{
				"errors.errorString",
			},
			iTypeStr,
		) {
			// This will not handle custom types that implement Error interface.
			return "error(" + deValStr + ")"
		}
	}

	switch iTypeKind { //nolint:exhaustive
	case reflect.Invalid:
		return "<invalid>"
	case reflect.String:
		return va.dumpString(val)
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Uintptr:
		return va.dumpNumbers(val)
	case reflect.Struct:
		// We used to use `sanity-io/litter` to do dumping.
		// We however, decided to implement our own dump functionality.
		//
		// The main reason is because sanity-io/litter has no way to compact
		// arrays/slices/maps that are inside structs.
		// This logic can be discarded if sanity-io/litter implements similar.
		// see: https://github.com/sanity-io/litter/pull/43
		fromPtr := false
		return va.dumpStruct(val, fromPtr, hideZeroValues, indentLevel)
	case reflect.Ptr:
		v := val.Elem()
		if v.IsValid() {
			if v.Type().Kind() == reflect.Struct {
				fromPtr := true
				// TODO: should we pass in `val`, itself instead of `v`
				// that way `val.Pointer()` would happen inside `dumpStruct`
				return va.dumpStruct(v, fromPtr, hideZeroValues, indentLevel)
			} else {
				return va.dumpNonStructPointer(v, hideZeroValues, indentLevel)
			}
		} else {
			// `v.IsValid()` returns false if v is the zero Value.
			// If `IsValid` returns false, all other methods except String panic.
			return val.Type().String() + "(nil)"
		}
	case reflect.Array, reflect.Slice:
		return va.dumpSlice(val, hideZeroValues, indentLevel)
	case reflect.Chan:
		return va.dumpChan(val)
	case reflect.Map:
		return va.dumpMap(val, hideZeroValues, indentLevel)
	case reflect.Bool:
		return fmt.Sprint(val)
	case reflect.Func:
		return va.dumpFunc(val)
	case reflect.Complex64, reflect.Complex128:
		return va.dumpComplexNum(val)
	case reflect.UnsafePointer:
		// It is not generally safe to do anything with an unsafe.Pointer
		// see: https://golang.org/pkg/unsafe/#Pointer
		// so we probably want to leave it as is.
		// do note that if we wanted we could get a uintptr via `val.Pointer()`
		return "unsafe.Pointer"
	case reflect.Interface:
		return va.dumpInterface(val)
	default:
		return fmt.Sprintf("%v NotImplemented", iTypeKind)
	}
}

func (va vari) dumpString(v reflect.Value) string {
	// dumps strings

	adder := 1 // this is a custom string type.
	if strings.HasPrefix(fmt.Sprintf("%#v", v), `"`) {
		adder = 2 // the `+2` is important so that the final quote `"` at end of string is not cut off
	}
	numEntries := v.Len()
	if numEntries > 0 && (fmt.Sprintf("%#v", v)[0] == 34) { // 34 is the char("")
		adder = 2
	}
	newLineCount := strings.Count(fmt.Sprintf("%s", v), "\n")

	constraint := int(math.Min(float64(numEntries), float64(va.cfg.MaxLength+50))) + adder + newLineCount

	s := fmt.Sprintf("%#v", v)[:constraint]

	if numEntries > constraint {
		remainder := numEntries - constraint
		s = s + fmt.Sprintf(" ...<%d more redacted>..", remainder)
	}
	if s == "" {
		s = `""`
	}

	return s
}

func (va vari) dumpStruct(v reflect.Value, fromPtr, hideZeroValues bool, indentLevel int) string {
	/*
		`fromPtr` indicates whether the struct is a value or a pointer; `T{}` vs `&T{}`
		`compact` indicates whether the struct should be laid in one line or not
		`hideZeroValues` indicates whether to show zeroValued fields
		`indentLevel` is the number of spaces from the left-most side of the termninal for struct names
	*/
	typeName := v.Type().Name()
	if fromPtr {
		typeName = "&" + typeName
	}

	if indentLevel > va.cfg.MaxIndentLevel {
		return fmt.Sprintf("%v: kama warning(indentation `%d` exceeds max of `%d`. Possible circular reference)", typeName, indentLevel, va.cfg.MaxIndentLevel)
	}

	sep := "\n"
	fieldNameSep := strings.Repeat("  ", indentLevel)
	lastBracketSep := strings.Repeat("  ", indentLevel-1)

	vt := v.Type()
	s := fmt.Sprintf("%s{%s", typeName, sep)

	numFields := v.NumField()
	for i := 0; i < numFields; i++ {
		vtf := vt.Field(i)
		fieldd := v.Field(i)
		if unicode.IsUpper(rune(vtf.Name[0])) || va.cfg.ShowPrivateFields {
			// Only dump public fields, unless the config option for private is turned on.

			if hideZeroValues && isZeroValue(fieldd) {
				continue
			} else {
				// when something is inside a struct, that's when we use compact & hideZeroValues
				cpt := true
				_ = cpt
				hzv := true
				val := va.dump(fieldd, hzv, indentLevel)
				s = s + fieldNameSep + vtf.Name + ": " + val + fmt.Sprintf(",%s", sep)
			}
		}
	}

	if v.IsZero() && indentLevel > 1 {
		s = strings.TrimRight(s, ",\n")
		s = s + "}"
	} else {
		s = s + lastBracketSep + "}"
	}

	return s
}

func (va vari) dumpSlice(v reflect.Value, hideZeroValues bool, indentLevel int) string {
	// dumps slices & arrays

	// In future(if we ever add compation) we could restrict compaction only to arrays/slices/maps that are of primitive(basic) types
	// see:
	//     1. https://github.com/sanity-io/litter/pull/43
	//     2. https://github.com/komuw/kama/pull/28

	numEntries := v.Len()
	constraint := int(math.Min(float64(numEntries), float64(va.cfg.MaxLength)))
	typeName := v.Type().String()

	newline := "\n"
	if numEntries <= 0 {
		// do not use newline.
		newline = ""
	}
	leftSep := "   "

	s := typeName + "{" + newline
	for i := 0; i < constraint; i++ {
		elm := v.Index(i)
		s = s + leftSep + va.dump(elm, hideZeroValues, indentLevel) + "," + newline
	}
	if numEntries > constraint {
		remainder := numEntries - constraint
		ident := strings.Repeat("  ", indentLevel)
		s = s + ident + fmt.Sprintf(" ...<%d more redacted>..", remainder) + newline
	}
	if v.IsZero() {
		s = s + "(nil)}"
	} else if numEntries <= 0 {
		s = s + "}"
	} else {
		ident := strings.Repeat("  ", indentLevel)
		s = strings.TrimRight(s, ",") // maybe use `strings.TrimSuffix`
		s = s + ident + "}"
	}
	return s
}

func (va vari) dumpMap(v reflect.Value, hideZeroValues bool, indentLevel int) string {
	// dumps maps

	// In future(if we ever add compation) we could restrict compaction only to arrays/slices/maps that are of primitive(basic) types
	// see:
	//     1. https://github.com/sanity-io/litter/pull/43
	//     2. https://github.com/komuw/kama/pull/28

	numEntries := v.Len()
	constraint := int(math.Min(float64(numEntries), float64(va.cfg.MaxLength)))
	typeName := v.Type().String()

	newline := "\n"
	leftSep := "   "
	colonSep := " "
	s := typeName + "{" + newline

	// Lets sort the map based on keys. This is done to introduce stability of the output.
	// This is not an important part of the design of kama, however, it makes testing much easier.
	keys := v.MapKeys()
	sort.Slice(keys,
		func(i, j int) bool {
			// it's unfortunate that we have to dump twice. In this func and in the `for range` below.
			return va.dump(keys[i], hideZeroValues, indentLevel) < va.dump(keys[j], hideZeroValues, indentLevel)
		},
	)
	for count, key := range keys {
		mapKey := key
		mapVal := v.MapIndex(key)
		s = s + leftSep + va.dump(mapKey, hideZeroValues, indentLevel) + ":" + colonSep + va.dump(mapVal, hideZeroValues, indentLevel) + ", " + newline
		if count > constraint {
			remainder := numEntries - constraint
			s = s + fmt.Sprintf("%s...<%d more redacted>..", leftSep, remainder)
			break
		}
	}

	s = strings.TrimRight(s, ",\n") // maybe use `strings.TrimSuffix`
	if v.IsZero() {
		s = s + "(nil)}"
	} else if numEntries <= 0 {
		s = s + "}"
	} else {
		ident := strings.Repeat("  ", indentLevel)
		s = s + "\n" + ident + "}"
	}
	return s
}

func (va vari) dumpChan(v reflect.Value) string {
	// dumps channels
	cap := v.Cap()
	len := v.Len()
	direction := v.Type().ChanDir()
	element := v.Type().Elem()
	return fmt.Sprintf("%v %v (len=%d, cap=%d)", direction, element, len, cap)
}

func (va vari) dumpFunc(v reflect.Value) string {
	// dumps functions

	vType := v.Type()
	typeName := vType.String()

	if !strings.Contains(typeName, "(") {
		// ie, typeName is like `http.HandlerFunc` instead of like `func() (io.ReadCloser, error)`
		// we thus need to bring out the actual signature
		numIn := vType.NumIn()
		numOut := vType.NumOut()

		if numIn > 0 {
			typeName = typeName + "("
			for i := 0; i < numIn; i++ {
				arg := vType.In(i)
				typeName = typeName + arg.String() + ", "
			}
			typeName = strings.TrimRight(typeName, ", ") // maybe use `strings.TrimSuffix`
			typeName = typeName + ")"
		}
		if numOut > 0 {
			typeName = typeName + " ("
			for i := 0; i < numOut; i++ {
				arg := vType.Out(i)
				typeName = typeName + arg.String() + ", "
			}
			typeName = strings.TrimRight(typeName, ", ") // maybe use `strings.TrimSuffix`
			typeName = typeName + ")"
		}
	}

	return typeName
}

func (va vari) dumpComplexNum(v reflect.Value) string {
	// dumps complex64 and complex128 numbers
	bits := v.Type().Bits()
	cmp := v.Complex() // returns complex128 even for `reflect.Complex64`
	if bits == 64 {
		return fmt.Sprintf("complex64%v", cmp)
	}
	return fmt.Sprintf("complex128%v", cmp)
}

func (va vari) dumpNonStructPointer(v reflect.Value, hideZeroValues bool, indentLevel int) string {
	// dumps pointer types other than struct.
	// ie; someIntEight := int8(14); kama.Dirp(&someIntEight)
	// dumping for struct pointers is handled in `dumpStruct()`

	pref := "&"
	s := va.dump(v, hideZeroValues, indentLevel)

	if strings.HasPrefix(s, pref) {
		return s
	}
	return pref + s
}

func (va vari) dumpNumbers(v reflect.Value) string {
	// dumps numbers.

	iType := v.Type()
	iTypeKind := iType.Kind()

	switch iTypeKind { //nolint:exhaustive
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,

		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,

		reflect.Float32,
		reflect.Float64,

		reflect.Uintptr:

		name := v.Type().String()
		return fmt.Sprintf("%s(%v)", name, v)
	default:
		return fmt.Sprintf("%v NotImplemented", iTypeKind)
	}
}

func (va vari) dumpInterface(v reflect.Value) string {
	// dump interface

	name := v.Type().String() // eg; `io.Reader` or `error`
	concrete := ""
	actualVal := ""

	if !v.IsNil() {
		elm := v.Elem()
		concrete = elm.Type().String() // eg; `*strings.Reader`

		// The fmt package treats `reflect.Value` specially.
		// It does not call their String method implicitly but instead prints the concrete values they hold.
		switch name {
		// TODO: add more cases here as we recognise how to handle them
		case "error":
			actualVal = fmt.Sprint(elm) // this will be the string content of the error
		case "context.Context":
			actualVal = fmt.Sprint(elm)
		default:
			actualVal = fmt.Sprint(elm)

			// default:
			// 	panic(fmt.Sprintf("dumpInterface unable to handle: `%v`. please open a github issue.", name))
		}
	} else {
		actualVal = "nil"
	}

	vv := name
	if concrete != "" {
		vv = vv + "(" + concrete + ")"
	}
	if actualVal != "" {
		vv = vv + " " + actualVal
	}

	if name == "error" {
		return "error(" + actualVal + ")"
	}
	return vv // s
}

func isPointerValue(v reflect.Value) bool {
	// Taken from https://github.com/sanity-io/litter/blob/v1.5.1/util.go
	// under the MIT license;
	// https://github.com/sanity-io/litter/blob/v1.5.1/LICENSE
	switch v.Kind() { //nolint:exhaustive
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return true
	}
	return false
}

func isZeroValue(v reflect.Value) bool {
	// Taken from https://github.com/sanity-io/litter/blob/v1.5.1/util.go
	// under the MIT license;
	// https://github.com/sanity-io/litter/blob/v1.5.1/LICENSE
	return (isPointerValue(v) && v.IsNil()) ||
		(v.IsValid() && v.CanInterface() && reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface()))
}

// deInterface returns values inside of non-nil interfaces when possible.
// This is useful for data types like structs, arrays, slices, and maps which
// can contain varying types packed inside an interface.
func deInterface(v reflect.Value) reflect.Value {
	// Taken from https://github.com/sanity-io/litter/blob/v1.5.1/util.go
	// under the MIT license;
	// https://github.com/sanity-io/litter/blob/v1.5.1/LICENSE
	if v.Kind() == reflect.Interface && !v.IsNil() {
		v = v.Elem()
	}
	return v
}
