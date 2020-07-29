package main

import (
	"fmt"
	"runtime"

	"reflect"
)


// TODO: merge methods/fields of T and *T
// https://play.golang.org/p/aQbEhI8WDP0

// vari represents a variable
type vari struct {
	name      string
	kind      reflect.Kind
	signature string
	fields    []string
	methods   []string
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
		v.name,
		v.kind,
		v.signature,
		v.fields,
		v.methods,
	)
}

func newVari(i interface{}) vari {
	iType := reflect.TypeOf(i)
	if iType == nil {
		// TODO: maybe there is a way in reflect to diffrentiate the various types of nil
		return vari{
			name:      "nil",
			kind:      reflect.Ptr,
			signature: "nil"}
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

	var fields = []string{}
	if typeKind == reflect.Struct {
		numFields := iType.NumField()
		for i := 0; i < numFields; i++ {
			f := iType.Field(i)
			if f.PkgPath != "" {
				// private field
				continue
			}
			fields = append(fields, f.PkgPath+"."+f.Name+",")
		}
	}

	var methods = []string{}
	numMethods := iType.NumMethod()
	for i := 0; i < numMethods; i++ {
		meth := iType.Method(i)
		if meth.PkgPath != "" {
			// private method
			continue
		}
		methName := meth.PkgPath + "." + meth.Name
		methSig := meth.Type.String() // type signature

		// TODO: maybe we should try and also add argument names if any.
		// currently the signature is displayed as;
		//   func(main.Foo, int, int) int
		// it would be cooler to display as;
		//   func(main.Foo, price int, commission int) int
		methods = append(methods, "\n\t"+methName+fmt.Sprintf("\n\t\t%v", methSig))
	}
	methods = append(methods, "\n\t")

	return vari{
		name:      typeName,
		kind:      typeKind,
		signature: typeSig,
		fields:    fields,
		methods:   methods,
	}

}
