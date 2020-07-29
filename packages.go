package main

import (
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

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
