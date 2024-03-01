package kama

import (
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/exp/slices"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

const (
	packageNeeds = packages.NeedName |
		packages.NeedImports |
		packages.NeedTypes
)

// pak represents a Go package
type pak struct {
	Name      string
	Constants []string
	Variables []string
	Functions []string
	Types     map[string][]string
}

func newPak(pattern string) (pak, error) {
	// Although `packages.Load` accepts a slice of multiple items, for `kama` we only accept one.
	patterns := []string{fmt.Sprintf("pattern=%s", pattern)}

	// A higher numbered modes cause Load to return more information,
	cfg := &packages.Config{Mode: packageNeeds}
	pkgs, err := packages.Load(
		// Load passes most patterns directly to the underlying build tool, but all patterns with the prefix "query=",
		// where query is a non-empty string of letters from [a-z], are reserved and may be interpreted as query operators.
		// Two query operators are currently supported: "file" and "pattern".
		// See: https://pkg.go.dev/golang.org/x/tools/go/packages?tab=doc#pkg-overview
		cfg, patterns...,
	)
	if err != nil {
		return pak{}, err
	}

	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		return pak{}, pkg.Errors[0]
	}
	constantSlice, varSlice, typeSlice, funcSlice, methodSlice := pkgScope(pkg)

	type2Methods := associateTypeMethods(typeSlice, methodSlice)

	return pak{
		Name:      pkg.PkgPath,
		Constants: constantSlice,
		Variables: varSlice,
		Functions: funcSlice,
		Types:     type2Methods,
	}, nil
}

func (p pak) String() string {
	nLf := func(x []string) []string {
		fm := []string{}
		if len(x) <= 1 {
			return fm
		}

		for _, c := range x {
			fm = append(fm, "\n\t"+c)
		}
		fm = append(fm, "\n\t")
		return fm
	}

	var sliceTypeMeths []string
	for typ, meths := range p.Types {
		typ := typ
		meths := meths

		t := "\n\t" + typ
		for _, met := range meths {
			t = t + "\n\t\t" + met
		}
		sliceTypeMeths = append(sliceTypeMeths, t)
	}
	slices.Sort(sliceTypeMeths)

	return fmt.Sprintf(
		`
[
NAME: %v
CONSTANTS: %v
VARIABLES: %v
FUNCTIONS: %v
TYPES: %v
]`,
		p.Name,
		nLf(p.Constants),
		nLf(p.Variables),
		nLf(p.Functions),
		sliceTypeMeths,
	)
}

func pkgScope(pkg *packages.Package) ([]string, []string, []string, []string, []string) {
	// package members (TypeCheck or WholeProgram mode)

	constVarTypFunc := []string{} // holds top level constants, variables, types & functions
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
			funcSlice = append(funcSlice, fmt.Sprintf("%v", v))
		}
		if strings.HasPrefix(v, "const") {
			v = strings.TrimPrefix(v, "const ")
			constantSlice = append(constantSlice, fmt.Sprintf("%v", v))
		} else if strings.HasPrefix(v, "var") {
			v = strings.TrimPrefix(v, "var ")
			varSlice = append(varSlice, fmt.Sprintf("%v", v))
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
				methSaveName := strings.TrimPrefix(meth, "method ")

				_, exists := type2Methods[typSaveName]
				if exists {
					methds := type2Methods[typSaveName]
					methds = append(methds, methSaveName)
					type2Methods[typSaveName] = methds
				} else {
					type2Methods[typSaveName] = []string{methSaveName}
				}
			}
		}

	}

	return type2Methods
}
