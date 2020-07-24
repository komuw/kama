package main

import (
	"bufio"
	"fmt"
	"os"

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
	iType := reflect.TypeOf(i)
	typeName := iType.PkgPath() + "." + iType.Name()
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
		methods = append(methods, "\n\t"+methName+fmt.Sprintf(" <<%v>> ", methSig))
	}
	methods = append(methods, "\n\t")

	preamble := fmt.Sprintf(
		`
NAME: %v,
KIND: %v,
STRING: %v,
FIELDS: %v
METHODS: %v
`,
		typeName,
		typeKind,
		typeString,
		fields,
		methods,
	)
	var dict = []string{preamble}
	fmt.Println(dict)
}

func main() {
	defer panicHandler()

	dir(34)
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

	// dir(io.Reader{})
}
