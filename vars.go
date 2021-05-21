package kama

import (
	"fmt"
	"math"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"reflect"

	"github.com/sanity-io/litter"
)

// vari represents a variable
type vari struct {
	Name      string
	Kind      reflect.Kind
	Signature []string
	Fields    []string
	Methods   []string
	Val       string
}

func newVari(i interface{}) vari {
	iType := reflect.TypeOf(i)
	valueOfi := reflect.ValueOf(i)

	if iType == nil {
		// TODO: maybe there is a way in reflect to diffrentiate the various types of nil
		return vari{
			Name:      "nil",
			Kind:      reflect.Ptr,
			Signature: []string{"nil"},
			Val:       dump(valueOfi, false)}
	}

	typeKind := getKind(i)
	typeName := getName(i)
	typeSig := getSignature(i)

	var fields = getAllFields(i)
	var methods = trimMethods(getAllMethods(i))

	return vari{
		Name:      typeName,
		Kind:      typeKind,
		Signature: typeSig,
		Fields:    fields,
		Methods:   methods,
		Val:       dump(valueOfi, false),
	}

}

func (v vari) String() string {
	nLf := func(x []string) []string {
		var fm = []string{}
		if len(x) <= 0 {
			return fm
		}

		for _, c := range x {
			fm = append(fm, "\n\t"+c)
		}
		fm = append(fm, "\n\t")
		return fm
	}

	return fmt.Sprintf(
		`
[
NAME: %v
KIND: %v
SIGNATURE: %v
FIELDS: %v
METHODS: %v
SNIPPET: %s
]
`,
		v.Name,
		v.Kind,
		v.Signature,
		nLf(v.Fields),
		nLf(v.Methods),
		v.Val,
	)
}

func dump(val reflect.Value, compact bool) string {
	/*
		`compact` indicates whether the struct should be laid in one line or not
	*/

	iType := val.Type()
	maxL := 720

	if iType == nil {
		// TODO: handle this better
		return "Nil NotImplemented"
	}

	dumpStruct := func(v reflect.Value, fromPtr bool, compact bool) string {
		/*
			`fromPtr` indicates whether the struct is a value or a pointer; `T{}` vs `&T{}`
			`compact` indicates whether the struct should be laid in one line or not
		*/
		// This logic is only required until similar logic is implemented in sanity-io/litter
		// see:
		// - https://github.com/sanity-io/litter/issues/34
		// - https://github.com/sanity-io/litter/pull/43

		typeName := v.Type().Name()
		if fromPtr {
			typeName = "&" + typeName
		}

		fmt.Println("typeName, compact: ", typeName, compact)
		sep := "\n"
		if compact {
			sep = ""
		}

		vt := v.Type()
		s := fmt.Sprintf("%s{%s", typeName, sep)

		numFields := v.NumField()
		for i := 0; i < numFields; i++ {
			vtf := vt.Field(i)
			fieldd := v.Field(i)
			if unicode.IsUpper(rune(vtf.Name[0])) {
				// `.Interface()` only works for exported fields
				val := dump(fieldd, compact)
				s = s + "  " + vtf.Name + ": " + val + fmt.Sprintf(",%s", sep)
			}
		}
		s = s + "}"
		return s
	}

	dumpSlice := func(v reflect.Value, compact bool) string {
		//dumps slices & arrays

		//TODO: (BUG)
		// inline funcs ``inherit`` variables from their enclosing func.
		// The main func(dump) has a param called `compact` which is `false`
		// If `dumpSlice` is called with `compact==true`; it will have `compact==false`
		//because of the ``inheritance``.
		// We need to kill this inline funcs

		fmt.Println("dunpslie compact: ", compact)
		maxL = 10
		numEntries := val.Len()
		constraint := int(math.Min(float64(numEntries), float64(maxL)))
		typeName := iType.String()

		s := typeName + "{"
		for i := 0; i < constraint; i++ {
			elm := val.Index(i) // todo: call dump on this
			s = s + dump(elm, compact) + ","
		}
		if numEntries > constraint {
			remainder := numEntries - constraint
			s = s + fmt.Sprintf(" ...<%d more redacted>..", remainder)
		}
		s = s + "}"
		return s
	}

	switch iType.Kind() {
	case reflect.String:
		maxL = 50
		// TODO: constraint by maxL
		return fmt.Sprint(val)
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64:
		return fmt.Sprint(val)
	case reflect.Struct:
		// the reason we are doing this is because sanity-io/litter has no way to compact
		// arrays/slices/maps that are inside structs.
		// This logic can be discarded if sanity-io/litter implements similar.
		// see: https://github.com/sanity-io/litter/pull/43
		fromPtr := false
		compact := false
		return dumpStruct(val, fromPtr, compact)
	case reflect.Ptr:
		v := val.Elem()
		if v.IsValid() {
			if v.Type().Kind() == reflect.Struct {
				// the reason we are doing this is because sanity-io/litter has no way to compact
				// arrays/slices/maps that are inside structs.
				// This logic can be discarded if sanity-io/litter implements similar.
				// see: https://github.com/sanity-io/litter/pull/43
				fromPtr := true
				compact := false
				return dumpStruct(v, fromPtr, compact)
			}
		}
	case reflect.Array,
		reflect.Slice:
		// In future we could restrict compaction only to arrays/slices/maps that are of primitive(basic) types
		// see: https://github.com/sanity-io/litter/pull/43
		cpt := true
		return dumpSlice(val, cpt)
	case reflect.Map:
		// In future we could restrict compaction only to arrays/slices/maps that are of primitive(basic) types
		// see: https://github.com/sanity-io/litter/pull/43
		maxL = 50
	}

	x := 9
	if x < 5 {
		sq := litter.Options{
			StripPackageNames: true,
			HidePrivateFields: true,
			HideZeroValues:    false,
			FieldExclusions:   regexp.MustCompile(`^(XXX_.*)$`), // XXX_ is a prefix of fields generated by protoc-gen-go
			Separator:         " "}

		s := sq.Sdump(val)
		if len(s) <= maxL {
			maxL = len(s)
			return s[:maxL]
		}
	}

	_ = maxL
	return "NotImplemented"
}

func getName(i interface{}) string {
	iType := reflect.TypeOf(i)
	typeKind := iType.Kind()

	typeName := ""
	switch typeKind {
	case reflect.Func:
		typeName = runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
		if typeName == "" {
			typeName = "." + iType.Elem().Name()
		}
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		typeName = iType.String()
	case reflect.Ptr:
		valueI := reflect.ValueOf(i).Elem()
		valueType := valueI.Type()
		typeName = valueType.PkgPath() + "." + valueType.Name()
	default:
		typeName = typeKind.String()
	}

	if iType.PkgPath() != "" {
		typeName = iType.PkgPath() + "." + iType.Name()
	}

	return typeName
}

func getKind(i interface{}) reflect.Kind {
	iType := reflect.TypeOf(i)
	typeKind := iType.Kind()

	if typeKind != reflect.Ptr {
		return typeKind
	} else {
		// the passed in type maybe be a `*T` so lets find the kind of `T`
		valueI := reflect.ValueOf(i).Elem()
		valueType := valueI.Type()
		valueTypeKind := valueType.Kind()
		return valueTypeKind
	}
}

func getSignature(i interface{}) []string {
	iType := reflect.TypeOf(i)
	typeKind := iType.Kind()

	var allSignatures = []string{iType.String()}

	if typeKind == reflect.Ptr {
		// the passed in type maybe be a `*T` so lets find the signature of `T`
		valueI := reflect.ValueOf(i).Elem()
		vType := valueI.Type()
		sig := vType.String()
		allSignatures = append(allSignatures, sig)
	} else if typeKind == reflect.Struct {
		// the passed in type is a `T` so lets also find the signature of `*T`
		// NB: We are currently only allowing structs here. But we could expand to other types
		ptrOfT := reflect.PtrTo(iType)
		sig := ptrOfT.String()
		allSignatures = append(allSignatures, sig)
	}

	return allSignatures
}

// getFields finds all the fields(if any) of type `T` and `*T`
func getAllFields(i interface{}) []string {
	iType := reflect.TypeOf(i)
	typeKind := iType.Kind()

	var allFields = []string{}
	if typeKind == reflect.Ptr {
		// the passed in type maybe be a `*T struct{}` so lets also find methods of `T struct{}`
		valueI := reflect.ValueOf(i).Elem()
		allFields = getFields(valueI.Type())
	} else if typeKind == reflect.Struct {
		allFields = getFields(iType)
	}

	return allFields
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
			fields = append(fields, f.Name+" "+f.Type.String())
		}
	}

	return fields
}

// getAllMethods finds all the methods of type `T` and `*T`
func getAllMethods(i interface{}) []string {
	iType := reflect.TypeOf(i)

	var allMethods = []string{}
	var methodsOfPassedInType = []string{}
	var methodsOfT = []string{}
	var methodsOfPointerT = []string{}

	methodsOfPassedInType = getMethods(iType)

	if iType.Kind() == reflect.Ptr {
		// the passed in type is a `*T` so lets also find methods of `T`
		valueI := reflect.ValueOf(i).Elem()
		methodsOfT = getMethods(valueI.Type())
	} else {
		// the passed in type is a `T` so lets also find methods of `*T`
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

// trimMethods removes any duplicated methods.
// if a method is applicable for both type `T` and `*T`, then `trimMethods` will
// just remove the one for `*T`
func trimMethods(methods []string) []string {

	// contains tells whether a contains x.
	contains := func(a []string, x string) bool {
		for _, n := range a {
			if x == n {
				return true
			}
		}
		return false
	}

	var trimmedMethods = []string{}
	var TmethNames []string

	// first add all methods for type `T`
	for _, meth := range methods {
		if !strings.Contains(meth, "*") {
			trimmedMethods = append(trimmedMethods, meth)
			methName := strings.Split(meth, " ")[0]
			TmethNames = append(TmethNames, methName)
		}
	}
	// then add methods for `*T` but only if they do not exist for `T`
	for _, meth := range methods {
		if strings.Contains(meth, "*") {
			methName := strings.Split(meth, " ")[0]
			if !contains(TmethNames, methName) {
				trimmedMethods = append(trimmedMethods, meth)
			}
		}
	}

	return trimmedMethods
}
