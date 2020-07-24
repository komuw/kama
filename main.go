package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"

	"reflect"

	"github.com/bradfitz/iter"
	"github.com/komuw/meli"
)

type Foo struct {
	Prop string
}

func (f Foo) String() string {
	return fmt.Sprintf("Foo(%v)", f.Prop)
}

func (f Foo) Bar() string {
	return f.Prop
}
func (f Foo) Add(a, b int) int {
	return a + b
}
func (f Foo) private(s string) string {
	return s
}

func dir(i interface{}) {
	// TODO: from the documentation of reflect.Type interface:
	// Not all methods apply to all kinds of types. Restrictions,
	// if any, are noted in the documentation for each method.
	// Use the Kind method to find out the kind of type before
	// calling kind-specific methods. Calling a method
	// inappropriate to the kind of type causes a run-time panic.
	//
	// TODO: we should check the kinds before calling any methods on the `Type`
	// to make sure they are allowed.

	iType := reflect.TypeOf(i)
	if iType == nil {
		// TODO: make this template a constant
		// TODO: stop repeating myself
		// TODO: maybe there is a way in reflect to diffrentiate the various types of nil
		preamble := fmt.Sprintf(
			`
NAME: %v,
KIND: %v,
SIGNATURE: %v,
FIELDS: %v
METHODS: %v
`,
			iType,
			iType,
			iType,
			iType,
			iType,
		)
		// TODO: dict is an odd name
		var dict = []string{preamble}
		fmt.Println(dict)
		return
	}

	typeName := iType.PkgPath() + "." + iType.Name()
	if typeName == "." {
		typeName = runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
		if typeName == "" {
			typeName = "." + iType.Elem().Name()
		}
	}

	typeKind := iType.Kind()
	typeString := iType.String()

	var fields = []string{}
	if typeKind == reflect.Struct {
		numFields := iType.NumField()
		for i := 0; i < numFields; i++ {
			f := iType.Field(i)
			fields = append(fields, f.PkgPath+"."+f.Name+",")
		}
	}

	var methods = []string{}
	numMethods := iType.NumMethod()
	for i := 0; i < numMethods; i++ {
		meth := iType.Method(i)
		methName := meth.PkgPath + "." + meth.Name
		methSig := meth.Type.String() // type signature
		methods = append(methods, "\n\t"+methName+fmt.Sprintf("\n\t\t%v", methSig))
	}
	methods = append(methods, "\n\t")

	// TODO: make this a constant template?
	// TODO: is using `iType.String()` as the value of `SIGNATURE` correct?
	preamble := fmt.Sprintf(
		`
NAME: %v,
KIND: %v,
SIGNATURE: %v,
FIELDS: %v
METHODS: %v
`,
		typeName,
		typeKind,
		typeString,
		fields,
		methods,
	)
	// TODO: dict is an odd name
	var dict = []string{preamble}
	fmt.Println(dict)
}

func myFunc(arg1 string, arg2 int) {

}
func main() {
	defer panicHandler()

	foo := Foo{}
	dir(foo)
	dir(bufio.Scanner{})
	dir(iter.N)
	dir(iter.N(89))

	dc := &meli.DockerContainer{
		ComposeService: meli.ComposeService{Image: "busybox"},
		LogMedium:      os.Stdout,
		FollowLogs:     true}

	dir(dc)
	dir(myFunc)

	// dir(io.Reader{})
}
