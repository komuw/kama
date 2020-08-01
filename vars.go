package main

import (
	"fmt"
	"runtime"

	"reflect"
)

// TODO: merge methods/fields of T and *T
// 1. https://play.golang.org/p/aQbEhI8WDP0
// 2. https://play.golang.org/p/EBhZW6hjb7O
// 3. https://play.golang.org/p/Olb2az0L2iI

// vari represents a variable
type vari struct {
	Name      string
	Kind      reflect.Kind
	Signature string
	Fields    []string
	Methods   []string
}

func (v vari) String() string {
	return fmt.Sprintf(
		`
[
NAME: %v
KIND: %v
SIGNATURE: %v
FIELDS: %v
METHODS: %v
]
`,
		v.Name,
		v.Kind,
		v.Signature,
		v.Fields,
		v.Methods,
	)
}

func newVari(i interface{}) vari {
	iType := reflect.TypeOf(i)
	if iType == nil {
		// TODO: maybe there is a way in reflect to diffrentiate the various types of nil
		return vari{
			Name:      "nil",
			Kind:      reflect.Ptr,
			Signature: "nil"}
	}

	typeKind := iType.Kind()
	typeName := iType.PkgPath() + "." + iType.Name()
	typeSig := iType.String()
	if typeName == "." {
		typeName = runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
		if typeName == "" {
			typeName = "." + iType.Elem().Name()
		}
	}

	var fields = getFields(iType)

	// TODO: If there is a method whose name is reported for both type `T` and `*T` we should only chose the method for `T`
	var methods = getAllMethods(i)

	return vari{
		Name:      typeName,
		Kind:      typeKind,
		Signature: typeSig,
		Fields:    fields,
		Methods:   methods,
	}

}

// getFields finds all the fields(if any) of a type
func getFields(iType reflect.Type) []string {
	var fields = []string{}
	typeKind := iType.Kind()
	if typeKind == reflect.Struct {
		numFields := iType.NumField()
		for i := 0; i < numFields; i++ {
			f := iType.Field(i)
			if f.PkgPath != "" {
				// private field
				continue
			}
			fields = append(fields, f.Name)
		}
	}

	return fields
}

// getAllMethods finds all the methods of type `T` and `*T`
func getAllMethods(i interface{}) []string {
	iType := reflect.TypeOf(i)

	var allMethods = []string{}
	var methodsOfT = []string{}
	var methodsOfPointerT = []string{}
	var methodsOfPassedInType = []string{}

	methodsOfPassedInType = getMethods(iType)

	if iType.Kind() == reflect.Ptr {
		// the passed in is a `*T` so lets also find methods of `T`
		valueI := reflect.ValueOf(i).Elem()
		methodsOfT = getMethods(valueI.Type())
	} else {
		// the passed in is a `T` so lets also find methods of `*T`
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
	var methods = []string{}
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
		methods = append(methods, methName+" "+methSig)
	}

	return methods
}
