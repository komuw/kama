package main

import (
	"fmt"
	"go/types"
	"log"
	"runtime"
	"strings"

	"reflect"

	"golang.org/x/tools/go/packages"
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
	cool(pkgs[0])
	// litter.Dump(pkgs[0])
	// dir(pkgs[0])
	// for _, pkg := range pkgs {
	// 	fmt.Println(pkg.ID, pkg.GoFiles)
	// }
}

func cool(pkg *packages.Package) {
	// package members (TypeCheck or WholeProgram mode)

	items := []string{}
	if pkg.Types != nil {
		qual := types.RelativeTo(pkg.Types)
		scope := pkg.Types.Scope()

		// fmt.Println("scope: ", scope)
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			if !obj.Exported() {
				// skip unexported names
				continue
			}
			// fmt.Println("obj type: ", obj.Type())
			// fmt.Println("obj name: ", obj.Name())
			// fmt.Println("ok: ", types.ObjectString(obj, qual))
			items = append(items, types.ObjectString(obj, qual))

		}
		// for _, name := range scope.Names() {
		// 	obj := scope.Lookup(name)
		// 	// if _, ok := obj.(*types.TypeName); ok {
		// 	// 	for _, meth := range typeutil.IntuitiveMethodSet(obj.Type(), nil) {
		// 	// 		if !meth.Obj().Exported() {
		// 	// 			continue // skip unexported names
		// 	// 		}
		// 	// 		fmt.Printf("\t%s\n", types.SelectionString(meth, qual))
		// 	// 	}
		// 	// }
		// }
	}

	constantSlice := []string{}
	varSlice := []string{}
	typeSlice := []string{}
	for _, v := range items {
		if strings.HasPrefix(v, "const") {
			constantSlice = append(constantSlice, v)
		} else if strings.HasPrefix(v, "var") {
			varSlice = append(varSlice, v)
		} else if strings.HasPrefix(v, "type") {
			typeSlice = append(typeSlice, v)
		}
	}
	fmt.Println("constantSlice: ", constantSlice)
	fmt.Println("varSlice: ", varSlice)
	fmt.Println("typeSlice: ", typeSlice)
}
func main() {
	defer panicHandler()

	pkgInfo([]string{"pattern=archive/tar"})

}
