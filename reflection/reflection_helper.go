package reflection

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

var (
	types    map[string]reflect.Type
	packages map[string][]reflect.Type
)

// discoverTypes initializes types and packages
func init() {
	types = make(map[string]reflect.Type)
	packages = make(map[string][]reflect.Type)

	discoverTypes()
}

func discoverTypes() {
	typ := reflect.TypeOf(0)
	sections, offset := typelinks2()
	for i, offs := range offset {
		rodata := sections[i]
		for _, off := range offs {
			emptyInterface := (*emptyInterface)(unsafe.Pointer(&typ))
			emptyInterface.data = resolveTypeOff(rodata, off)
			if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
				// just discover pointer types, but we also register this pointer type actual struct type to the registry
				loadedTypePtr := typ
				loadedType := typ.Elem()

				pkgTypes := packages[loadedType.PkgPath()]
				pkgTypesPtr := packages[loadedTypePtr.PkgPath()]

				if pkgTypes == nil {
					pkgTypes = []reflect.Type{}
					packages[loadedType.PkgPath()] = pkgTypes
				}
				if pkgTypesPtr == nil {
					pkgTypesPtr = []reflect.Type{}
					packages[loadedTypePtr.PkgPath()] = pkgTypesPtr
				}

				if strings.Contains(loadedType.String(), "Test") {
					n := GetFullTypeNameByType(loadedType)
					n2 := GetTypeNameByType(loadedType)
					fmt.Println(n)
					fmt.Println(n2)
				}

				types[GetFullTypeNameByType(loadedType)] = loadedType
				types[GetFullTypeNameByType(loadedTypePtr)] = loadedTypePtr
			}
		}
	}
}

func RegisterType(typ reflect.Type) {
	types[GetFullTypeName(typ)] = typ
}

func RegisterTypeWithKey(key string, typ reflect.Type) {
	types[key] = typ
}

func GetAllRegisteredTypes() []reflect.Type {
	var typeSlice []reflect.Type

	for _, typ := range types {
		typeSlice = append(typeSlice, typ)
	}
	return typeSlice
}

// TypeByName returns the type by its exact name, containing its namespace
func TypeByName(typeName string) reflect.Type {
	if typ, ok := types[typeName]; ok {
		return typ
	}
	return nil
}

// TypesByContainingName returns the types if typename is containing this `typeName`
func TypesByContainingName(typeName string) []reflect.Type {
	var containingTypes []reflect.Type
	for name, typ := range types {
		if strings.Contains(name, typeName) {
			containingTypes = append(containingTypes, typ)
		}
	}
	return containingTypes
}

// TypeByContainingName returns the type if typename is containing this `typeName`
func TypeByContainingName(typeName string) reflect.Type {
	for name, typ := range types {
		if strings.Contains(name, typeName) {
			return typ
		}
	}
	return nil
}

func TypeByNameAndImplementedInterface[TInterface interface{}](typeName string) reflect.Type {
	// https://stackoverflow.com/questions/7132848/how-to-get-the-reflect-type-of-an-interface
	implementedInterface := GetTypeFromGeneric[TInterface]()

	if typ, ok := types[typeName]; ok {
		if typ.Implements(implementedInterface) {
			return typ
		}
	}
	return nil
}

func TypesImplementedInterfaceWithFilterTypes[TInterface interface{}](
	types []reflect.Type,
) []reflect.Type {
	// https://stackoverflow.com/questions/7132848/how-to-get-the-reflect-type-of-an-interface
	implementedInterface := GetTypeFromGeneric[TInterface]()

	var res []reflect.Type
	for _, t := range types {
		if t.Implements(implementedInterface) {
			res = append(res, t)
		}
	}

	return res
}

// GetFullTypeName returns the full name of the type by its package name
func GetFullTypeName(input interface{}) string {
	t := reflect.TypeOf(input)
	return t.String()
}

func GetFullTypeNameByType(typ reflect.Type) string {
	return typ.String()
}

// GetTypeName returns the name of the type without its package name
func GetTypeName(input interface{}) string {
	t := reflect.TypeOf(input)
	if t.Kind() != reflect.Ptr {
		return t.Name()
	}

	return fmt.Sprintf("*%s", t.Elem().Name())
}

func GetNonPointerTypeName(input interface{}) string {
	t := reflect.TypeOf(input)
	if t.Kind() != reflect.Ptr {
		return t.Name()
	}

	return t.Elem().Name()
}

func GetTypeNameByType(typ reflect.Type) string {
	if typ.Kind() != reflect.Ptr {
		return typ.Name()
	}

	return fmt.Sprintf("*%s", typ.Elem().Name())
}

// TypeByPackageName return the type by its package and name
func TypeByPackageName(pkgPath string, name string) reflect.Type {
	if pkgTypes, ok := packages[pkgPath]; ok {
		for _, typ := range pkgTypes {
			if typ.Name() == name {
				return typ
			}
		}
	}
	return nil
}

func TypesByPackageName(pkgPath string) []reflect.Type {
	if pkgTypes, ok := packages[pkgPath]; ok {
		return pkgTypes
	}
	return nil
}

func GetTypeFromGeneric[T interface{}]() reflect.Type {
	res := reflect.TypeOf((*T)(nil)).Elem()
	return res
}

func GetBaseType(value interface{}) interface{} {
	if reflect.ValueOf(value).Kind() == reflect.Pointer {
		return reflect.ValueOf(value).Elem().Interface()
	}

	return value
}

func GetReflectType(value interface{}) reflect.Type {
	if reflect.TypeOf(value).Kind() == reflect.Pointer &&
		reflect.TypeOf(value).Elem().Kind() == reflect.Interface {
		return reflect.TypeOf(value).Elem()
	}

	res := reflect.TypeOf(value)
	return res
}

func GetBaseReflectType(value interface{}) reflect.Type {
	if reflect.ValueOf(value).Kind() == reflect.Pointer {
		return reflect.TypeOf(reflect.ValueOf(value).Elem().Interface())
	}

	return reflect.TypeOf(value)
}

func GenericInstanceByT[T any]() T {
	// https://stackoverflow.com/questions/7132848/how-to-get-the-reflect-type-of-an-interface
	typ := GetTypeFromGeneric[T]()
	return getInstanceFromType(typ).(T)
}

func InstanceByType(typ reflect.Type) interface{} {
	return getInstanceFromType(typ)
}

// InstanceByTypeName return an empty instance of the type by its name
// If the type is a pointer type, it will return a pointer instance of the type and
// if the type is a struct type, it will return an empty struct
func InstanceByTypeName(name string) interface{} {
	typ := TypeByName(name)

	return getInstanceFromType(typ)
}

func InstanceByTypeNameAndImplementedInterface[TInterface interface{}](name string) interface{} {
	typ := TypeByNameAndImplementedInterface[TInterface](name)

	return getInstanceFromType(typ)
}

// InstancePointerByTypeName return an empty pointer instance of the type by its name
// If the type is a pointer type, it will return a pointer instance of the type and
// if the type is a struct type, it will return a pointer to the struct
func InstancePointerByTypeName(name string) interface{} {
	typ := TypeByName(name)
	if typ.Kind() == reflect.Ptr {
		res := reflect.New(typ.Elem()).Interface()
		return res
	}

	return reflect.New(typ).Interface()
}

// InstanceByPackageName return an empty instance of the type by its name and package name
// If the type is a pointer type, it will return a pointer instance of the type and
// if the type is a struct type, it will return an empty struct
func InstanceByPackageName(pkgPath string, name string) interface{} {
	typ := TypeByPackageName(pkgPath, name)

	return getInstanceFromType(typ)
}

func getInstanceFromType(typ reflect.Type) interface{} {
	if typ.Kind() == reflect.Ptr {
		res := reflect.New(typ.Elem()).Interface()
		return res
	}

	return reflect.Zero(typ).Interface()
	// return reflect.New(typ).Elem().Interface()
}

func TypesImplementedInterface[TInterface interface{}]() []reflect.Type {
	// https://stackoverflow.com/questions/7132848/how-to-get-the-reflect-type-of-an-interface
	implementedInterface := GetTypeFromGeneric[TInterface]()

	var res []reflect.Type
	for _, t := range types {
		if t.Implements(implementedInterface) {
			res = append(res, t)
		}
	}

	return res
}
