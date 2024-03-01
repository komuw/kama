package kama

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// vari represents a variable
type vari struct {
	Name      string
	Kind      reflect.Kind
	Signature []string
	Fields    []string
	Methods   []string
	Val       string
}

func newVari(i interface{}) vari {
	iType := reflect.TypeOf(i)
	valueOfi := reflect.ValueOf(i)
	hideZeroValues := false
	indentLevel := 0

	if iType == nil {
		// TODO: maybe there is a way in reflect to diffrentiate the various types of nil
		return vari{
			Name:      "nil",
			Kind:      reflect.Ptr,
			Signature: []string{"nil"},
			Val:       dump(valueOfi, hideZeroValues, indentLevel),
		}
	}
	if iType.Kind() == reflect.Pointer && valueOfi.IsNil() && valueOfi.IsZero() {
		return vari{
			Name:      "unknown",
			Kind:      reflect.Ptr,
			Signature: []string{fmt.Sprintf("%v", iType)},
			Val:       dump(valueOfi, hideZeroValues, indentLevel),
		}
	}

	typeKind := getKind(i)
	typeName := getName(i)
	typeSig := getSignature(i)

	fields := getAllFields(i)
	methods := trimMethods(getAllMethods(i))

	return vari{
		Name:      typeName,
		Kind:      typeKind,
		Signature: typeSig,
		Fields:    fields,
		Methods:   methods,
		Val:       dump(valueOfi, hideZeroValues, indentLevel),
	}
}

func (v vari) String() string {
	nLf := func(x []string) []string {
		fm := []string{}
		if len(x) <= 0 {
			return fm
		}

		for _, c := range x {
			fm = append(fm, "\n\t"+c)
		}
		fm = append(fm, "\n\t")
		return fm
	}

	w := &bytes.Buffer{}
	stackp(w)

	return fmt.Sprintf(
		`
[
NAME: %v
KIND: %v
SIGNATURE: %v
FIELDS: %v
METHODS: %v
STACK_TRACE: [
%v]
SNIPPET: %s
]
`,
		v.Name,
		v.Kind,
		v.Signature,
		nLf(v.Fields),
		nLf(v.Methods),
		strings.TrimRight(w.String(), "\n"),
		v.Val,
	)
}

func getName(i interface{}) string {
	iType := reflect.TypeOf(i)
	typeKind := iType.Kind()

	typeName := ""
	switch typeKind { //nolint:exhaustive
	case reflect.Func:
		typeName = runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
		if typeName == "" {
			typeName = "." + iType.Elem().Name()
		}
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		typeName = iType.String()
	case reflect.Ptr:
		valueI := reflect.ValueOf(i).Elem()
		valueType := valueI.Type()
		typeName = valueType.PkgPath() + "." + valueType.Name()
	default:
		typeName = typeKind.String()
	}

	if iType.PkgPath() != "" {
		typeName = iType.PkgPath() + "." + iType.Name()
	}

	return typeName
}

func getKind(i interface{}) reflect.Kind {
	iType := reflect.TypeOf(i)
	typeKind := iType.Kind()

	if typeKind != reflect.Ptr {
		return typeKind
	} else {
		// the passed in type maybe be a `*T` so lets find the kind of `T`
		valueI := reflect.ValueOf(i).Elem()
		valueType := valueI.Type()
		valueTypeKind := valueType.Kind()
		return valueTypeKind
	}
}

func getSignature(i interface{}) []string {
	iType := reflect.TypeOf(i)
	typeKind := iType.Kind()

	allSignatures := []string{iType.String()}

	if typeKind == reflect.Ptr {
		// the passed in type maybe be a `*T` so lets find the signature of `T`
		valueI := reflect.ValueOf(i).Elem()
		vType := valueI.Type()
		sig := vType.String()
		allSignatures = append(allSignatures, sig)
	} else if typeKind == reflect.Struct {
		// the passed in type is a `T` so lets also find the signature of `*T`
		// NB: We are currently only allowing structs here. But we could expand to other types
		ptrOfT := reflect.PtrTo(iType)
		sig := ptrOfT.String()
		allSignatures = append(allSignatures, sig)
	}

	return allSignatures
}

// getFields finds all the fields(if any) of type `T` and `*T`
func getAllFields(i interface{}) []string {
	iType := reflect.TypeOf(i)
	typeKind := iType.Kind()

	allFields := []string{}
	if typeKind == reflect.Ptr {
		// the passed in type maybe be a `*T struct{}` so lets also find methods of `T struct{}`
		valueI := reflect.ValueOf(i).Elem()
		allFields = getFields(valueI.Type())
	} else if typeKind == reflect.Struct {
		allFields = getFields(iType)
	}

	return allFields
}

// getFields finds all the fields(if any) of a type
func getFields(iType reflect.Type) []string {
	fields := []string{}
	typeKind := iType.Kind()

	if typeKind == reflect.Struct {
		// TODO: If a structField happens to be a func,
		// We should enhance that signature.
		// Look at `dumpFunc()` for implementation
		// https://github.com/komuw/kama/issues/38

		numFields := iType.NumField()
		for i := 0; i < numFields; i++ {
			f := iType.Field(i)
			if f.PkgPath != "" {
				// private field
				continue
			}
			name := f.Name + " " + f.Type.String()
			fields = append(fields, name)
		}
	}

	return fields
}

// getAllMethods finds all the methods of type `T` and `*T`
func getAllMethods(i interface{}) []string {
	iType := reflect.TypeOf(i)

	allMethods := []string{}
	methodsOfT := []string{}
	methodsOfPointerT := []string{}

	methodsOfPassedInType := getMethods(iType)

	if iType.Kind() == reflect.Ptr {
		// the passed in type is a `*T` so lets also find methods of `T`
		valueI := reflect.ValueOf(i).Elem()
		methodsOfT = getMethods(valueI.Type())
	} else {
		// the passed in type is a `T` so lets also find methods of `*T`
		ptrOfT := reflect.PtrTo(iType)
		methodsOfPointerT = getMethods(ptrOfT)
	}

	allMethods = append(allMethods, methodsOfT...)
	allMethods = append(allMethods, methodsOfPointerT...)
	allMethods = append(allMethods, methodsOfPassedInType...)

	return allMethods
}

// getMethods finds all the methods of type.
func getMethods(iType reflect.Type) []string {
	methods := []string{}
	numMethods := iType.NumMethod()
	for i := 0; i < numMethods; i++ {
		meth := iType.Method(i)
		if meth.PkgPath != "" {
			// private method
			continue
		}
		methName := meth.Name
		methSig := meth.Type.String() // type signature

		// TODO: maybe we should try and also add argument names if any.
		// currently the signature is displayed as;
		//   func(main.Foo, int, int) int
		// it would be cooler to display as;
		//   func(main.Foo, price int, commission int) int
		// https://github.com/komuw/kama/issues/39
		methods = append(methods, methName+" "+methSig)
	}

	return methods
}

// trimMethods removes any duplicated methods.
// if a method is applicable for both type `T` and `*T`, then `trimMethods` will
// just remove the one for `*T`
func trimMethods(methods []string) []string {
	// contains tells whether a contains x.
	contains := func(a []string, x string) bool {
		for _, n := range a {
			if x == n {
				return true
			}
		}
		return false
	}

	trimmedMethods := []string{}
	var TmethNames []string

	// first add all methods for type `T`
	for _, meth := range methods {
		if !strings.Contains(meth, "*") {
			trimmedMethods = append(trimmedMethods, meth)
			methName := strings.Split(meth, " ")[0]
			TmethNames = append(TmethNames, methName)
		}
	}
	// then add methods for `*T` but only if they do not exist for `T`
	for _, meth := range methods {
		if strings.Contains(meth, "*") {
			methName := strings.Split(meth, " ")[0]
			if !contains(TmethNames, methName) {
				trimmedMethods = append(trimmedMethods, meth)
			}
		}
	}

	return trimmedMethods
}
