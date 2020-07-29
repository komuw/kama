package main

import (
	"fmt"
	"go/types"
	"net/http"
	"runtime"
	"strings"

	"reflect"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

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

	// // TODO: dict is an odd name
	// var dict = []string{preamble}
	// fmt.Println(dict)
}

func newVari(i interface{}) vari {
	iType := reflect.TypeOf(i)
	typeKind := iType.Kind()
	if iType == nil {
		// TODO: maybe there is a way in reflect to diffrentiate the various types of nil
		return vari{
			name:      "nil",
			kind:      typeKind,
			signature: "nil"}
	}

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
	var res interface{}
	var err error

	if reflect.TypeOf(i).Kind() == reflect.String {
		i := i.(string)
		res, err = newPaki(i)
		if err != nil {
			res = newVari(i)
		}
	} else {
		res = newVari(i)
	}

	fmt.Println(res)

}

// paki represents a Go package
type paki struct {
	name string
	path string

	constants []string
	variables []string
	functions []string
	types     []string
}

func (p paki) String() string {
	return fmt.Sprintf(
		`
[
NAME: %v
PATH: %v
CONSTANTS: %v
VARIABLES: %v
FUNCTIONS: %v
TYPES: %v
]`,
		p.name,
		p.path,
		p.constants,
		p.variables,
		p.functions,
		p.types,
	)
}

func newPaki(pattern string) (paki, error) {
	// patterns := []string{"pattern=net/http"}
	//patterns := []string{"pattern=os"}

	const (
		// TODO: move this to a global variable
		// TODO: whittle this down to only what we need(I think we only need `NeedTypes`)
		loadAll = packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedTypes |
			packages.NeedTypesSizes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedDeps
	)

	// Although `packages.Load` accepts a slice of multiple items, for `dir` we only accept one.
	patterns := []string{fmt.Sprintf("pattern=%s", pattern)}

	// A higher numbered modes cause Load to return more information,
	cfg := &packages.Config{Mode: loadAll}
	pkgs, err := packages.Load(
		// Load passes most patterns directly to the underlying build tool, but all patterns with the prefix "query=",
		// where query is a non-empty string of letters from [a-z], are reserved and may be interpreted as query operators.
		// Two query operators are currently supported: "file" and "pattern".
		// See: https://pkg.go.dev/golang.org/x/tools/go/packages?tab=doc#pkg-overview
		cfg, patterns...,
	)
	if err != nil {
		return paki{}, err

	}
	// if packages.PrintErrors(pkgs) > 0 {
	// 	// TODO: maybe we do not need this
	// 	log.Fatal("PrintErrors")
	// }

	pkg := pkgs[0]
	constantSlice, varSlice, typeSlice, funcSlice, methodSlice := pkgScope(pkg)

	type2Methods := associateTypeMethods(typeSlice, methodSlice)
	var finalTypeSlice = []string{}
	for typ, methSlice := range type2Methods {
		meths := ""
		for _, v := range methSlice {
			v = strings.TrimPrefix(v, "method ")
			meths = meths + fmt.Sprintf("\n\t%s", v)
		}
		finalTypeSlice = append(finalTypeSlice, fmt.Sprintf("\n%v%v", typ, meths))
	}

	return paki{
		name: pkg.Name,
		path: pkg.PkgPath,

		constants: constantSlice,
		variables: varSlice,
		functions: funcSlice,
		types:     finalTypeSlice,
	}, nil
}

func pkgScope(pkg *packages.Package) ([]string, []string, []string, []string, []string) {
	// package members (TypeCheck or WholeProgram mode)

	constVarTypFunc := []string{} //holds top level constants, variables, types & functions
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

			constVarTypFunc = append(constVarTypFunc, types.ObjectString(obj, qual))

			// lets get methods of types
			if _, ok := obj.(*types.TypeName); ok {
				for _, meth := range typeutil.IntuitiveMethodSet(obj.Type(), nil) {
					// look into: `godex.combinedMethodSet()`
					// https://github.com/golang/tools/blob/e96c4e24768da594adeb5eed27c8ecd547a3d4f1/cmd/godex/print.go#L347-L373
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
	funcSlice := []string{}
	for _, v := range constVarTypFunc {
		if strings.HasPrefix(v, "func") {
			v = strings.TrimPrefix(v, "func ")
			funcSlice = append(funcSlice, fmt.Sprintf("\n\t%v", v))
		}
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

	return constantSlice, varSlice, typeSlice, funcSlice, methodSlice
}
func associateTypeMethods(typeSlice, methodSlice []string) map[string][]string {
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

func main() {
	defer panicHandler()

	dir("archive/tar")
	dir("compress/flate")
	dir(&http.Request{})
	dir(http.Request{})
	dir("github.com/pkg/errors")
}
