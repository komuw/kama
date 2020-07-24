package main

import (
	"fmt"
	"go/types"
	"log"
	"runtime"
	"strings"

	"reflect"

	"github.com/sanity-io/litter"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

// TODO: clean up

// TODO: add documentation

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

func pkgInfo(patterns []string) {
	// patterns := []string{"pattern=net/http"}
	//patterns := []string{"pattern=os"}

	// A higher numbered modes cause Load to return more information,
	cfg := &packages.Config{Mode: packages.LoadAllSyntax}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)

	}
	// if packages.PrintErrors(pkgs) > 0 {
	// 	// TODO: maybe we do not need this
	// 	log.Fatal("PrintErrors")
	// }

	pkg1 := pkgs[0]
	// pkg1.Imports: %v
	fmt.Printf(`
Name: %v
PkgPath: %v
ExportFile: %v
`,
		pkg1.Name,
		pkg1.PkgPath,
		pkg1.ExportFile,
		// pkg1.TypesInfo,

		// pkg1.TypesInfo,
	)
	constantSlice, varSlice, typeSlice, methodSlice := cool(pkgs[0])
	okay(constantSlice, varSlice, typeSlice, methodSlice)
	// litter.Dump(pkgs[0])
	// dir(pkgs[0])
	// for _, pkg := range pkgs {
	// 	fmt.Println(pkg.ID, pkg.GoFiles)
	// }
}

func okay(constantSlice, varSlice, typeSlice, methodSlice []string) {
	// fmt.Println("constantSlice: ", constantSlice)
	// fmt.Println("varSlice: ", varSlice)
	// fmt.Println("typeSlice: ", typeSlice)
	// TODO: associate methods with their types from `typeSlice`
	// fmt.Println("methodSlice: ", methodSlice)

	type2Methods := map[string][]string{}
	for _, typ := range typeSlice {
		typName := strings.Split(typ, " ")[1]
		for _, meth := range methodSlice {
			methReceiverName := strings.Split(meth, " ")[1]
			methReceiverName = strings.ReplaceAll(methReceiverName, ")", "")
			methReceiverName = strings.ReplaceAll(methReceiverName, "(", "")
			methReceiverName = strings.ReplaceAll(methReceiverName, "*", "")

			if methReceiverName == typName {
				_, exists := type2Methods[typ]
				if exists {
					methds := type2Methods[typ]
					methds = append(methds, meth)
					type2Methods[typ] = methds
				} else {
					type2Methods[typ] = []string{meth}
				}
			}
		}

	}
	fmt.Println("type2Methods:")
	litter.Dump(type2Methods)
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
			constantSlice = append(constantSlice, v)
		} else if strings.HasPrefix(v, "var") {
			varSlice = append(varSlice, v)
		} else if strings.HasPrefix(v, "type") {
			typeSlice = append(typeSlice, v)
		}
	}
	// TODO: associate methods with their types from `typeSlice`
	return constantSlice, varSlice, typeSlice, methodSlice
}
func main() {
	defer panicHandler()

	pkgInfo([]string{"pattern=archive/tar"})

}
