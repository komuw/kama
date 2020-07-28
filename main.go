package main

import (
	"fmt"
	"go/types"
	"log"
	"net/http"
	"runtime"
	"strings"

	"reflect"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

// TODO: If someone passes in, say a struct;
// we should show them its type, methods etc
// but also print it out and its contents
// basically, do what `litter.Dump` would have done

// TODO: maybe add syntax highlighting, maybe make it optional??

// TODO: clean up

// TODO: add of `dir` documentation

// TODO: maybe we should show docs when someone requests for something specific.
// eg if they do `dir(http)` we do not show docs, but if they do `dir(&http.Request{})` we show docs.
// An alternative is only show docs, if someone requests. `dir(i interface{}, config ...dir.Config)`; config is `...` so that it is optional
// where config is a `type Config struct {}`

// TODO: add a command line api.
//   eg; `dir http.Request` or `dir http`
// have a look at `golang.org/x/tools/cmd/godex`

// TODO: this will stutter; `dir.dir(23)`
// maybe it is okay??
// TODO: surface all info for both the type and its pointer.
// currently `dir(&http.Client{})` & `dir(http.Client{})` produces different output; they should NOT
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
NAME: %v
KIND: %v
SIGNATURE: %v
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

	// TODO: make this a constant template?
	// TODO: is using `iType.String()` as the value of `SIGNATURE` correct?
	preamble := fmt.Sprintf(
		`
NAME: %v
KIND: %v
SIGNATURE: %v
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

func pkgInfo(pattern string) {
	// patterns := []string{"pattern=net/http"}
	//patterns := []string{"pattern=os"}

	// Although `packages.Load` accepts a slice of multiple items, for `dir` we only accept one.
	patterns := []string{fmt.Sprintf("pattern=%s", pattern)}

	// A higher numbered modes cause Load to return more information,
	cfg := &packages.Config{Mode: packages.LoadAllSyntax}
	pkgs, err := packages.Load(
		// Load passes most patterns directly to the underlying build tool, but all patterns with the prefix "query=",
		// where query is a non-empty string of letters from [a-z], are reserved and may be interpreted as query operators.
		// Two query operators are currently supported: "file" and "pattern".
		// See: https://pkg.go.dev/golang.org/x/tools/go/packages?tab=doc#pkg-overview
		cfg, patterns...,
	)
	if err != nil {
		log.Fatal(err)

	}
	// if packages.PrintErrors(pkgs) > 0 {
	// 	// TODO: maybe we do not need this
	// 	log.Fatal("PrintErrors")
	// }

	pkg1 := pkgs[0]
	constantSlice, varSlice, typeSlice, methodSlice := cool(pkg1)

	type2Methods := okay(typeSlice, methodSlice)
	var finalTypeSlice = []string{}
	for typ, methSlice := range type2Methods {
		meths := ""
		for _, v := range methSlice {
			v = strings.TrimPrefix(v, "method ")
			meths = meths + fmt.Sprintf("\n\t%s", v)
		}
		finalTypeSlice = append(finalTypeSlice, fmt.Sprintf("\n%v%v", typ, meths))
	}

	preamble := fmt.Sprintf(
		`
NAME: %v
PATH: %v
CONSTANTS: %v
VARIABLES: %v
TYPES: %v
`,
		pkg1.Name,
		pkg1.PkgPath,
		constantSlice,
		varSlice,
		finalTypeSlice,
	)
	// TODO: dict is an odd name
	var dict = []string{preamble}
	fmt.Println(dict)
}

func okay(typeSlice, methodSlice []string) map[string][]string {
	type2Methods := map[string][]string{}
	for _, typ := range typeSlice {
		typName := strings.Split(typ, " ")[1]
		for _, meth := range methodSlice {
			methReceiverName := strings.Split(meth, " ")[1]
			methReceiverName = strings.ReplaceAll(methReceiverName, ")", "")
			methReceiverName = strings.ReplaceAll(methReceiverName, "(", "")
			methReceiverName = strings.ReplaceAll(methReceiverName, "*", "")

			typSaveName := strings.Split(typ, " ")[1] + " " + strings.Split(typ, " ")[2]
			typSaveName = strings.TrimSpace(strings.Split(typSaveName, "{")[0])
			if methReceiverName == typName {
				_, exists := type2Methods[typSaveName]
				if exists {
					methds := type2Methods[typSaveName]
					methds = append(methds, meth)
					type2Methods[typSaveName] = methds
				} else {
					type2Methods[typSaveName] = []string{meth}
				}
			}
		}

	}

	return type2Methods
}

func cool(pkg *packages.Package) ([]string, []string, []string, []string) {
	// package members (TypeCheck or WholeProgram mode)

	constVarTyp := []string{}
	methodSlice := []string{}
	if pkg.Types != nil {
		qual := types.RelativeTo(pkg.Types)
		scope := pkg.Types.Scope()

		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			if !obj.Exported() {
				// skip unexported names
				continue
			}

			// TODO: add top level, Exported functions.
			// currently we do not. see `pkgInfo("github.com/pkg/errors")` <- this is missing the exported Funcs
			// fmt.Println("name: ", name)
			// fmt.Println("obj: ", obj)

			constVarTyp = append(constVarTyp, types.ObjectString(obj, qual))

			// lets get methods of types
			if _, ok := obj.(*types.TypeName); ok {
				for _, meth := range typeutil.IntuitiveMethodSet(obj.Type(), nil) {
					if !meth.Obj().Exported() {
						// skip unexported methods
						continue
					}
					methodSlice = append(methodSlice, types.SelectionString(meth, qual))
				}
			}

		}
	}

	constantSlice := []string{}
	varSlice := []string{}
	typeSlice := []string{}
	for _, v := range constVarTyp {
		if strings.HasPrefix(v, "const") {
			v = strings.TrimPrefix(v, "const ")
			constantSlice = append(constantSlice, fmt.Sprintf("\n\t%v", v))
		} else if strings.HasPrefix(v, "var") {
			v = strings.TrimPrefix(v, "var ")
			varSlice = append(varSlice, fmt.Sprintf("\n\t%v", v))
		} else if strings.HasPrefix(v, "type") {
			typeSlice = append(typeSlice, v)
		}
	}
	// TODO: associate methods with their types from `typeSlice`
	return constantSlice, varSlice, typeSlice, methodSlice
}
func main() {
	defer panicHandler()

	pkgInfo("archive/tar")
	dir(&http.Request{})
	dir(http.Request{})
	pkgInfo("pkg/errors")
}
