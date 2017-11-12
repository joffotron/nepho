package cfoo

import (
	"fmt"
	"reflect"

	"regexp"

	"github.com/k0kubun/pp"
)

func Translate(input map[string]interface{}) interface{} {
	translated := make(map[string]interface{})
	for key, value := range input {
		translated[key] = translate(value)
	}

	return translated
}

func translate(input interface{}) interface{} {
	inputValue := reflect.ValueOf(input)
	switch inputValue.Kind() {
	case reflect.String:
		return translateString(input.(string))
	case reflect.Map:
		translated := make(map[string]interface{})
		for key, value := range input.(map[interface{}]interface{}) {
			switch key.(type) {
			case string:
				translated[key.(string)] = translate(value)
			}
		}
		return translated
	default:
		panic("Unsupported type " + inputValue.Kind().String())
	}

}

func translateRecursiveOld(cpy, original reflect.Value) {
	switch original.Kind() {
	case reflect.Ptr:
		// If it is a pointer we need to unwrap and call once again
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		cpy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		translateRecursiveOld(cpy.Elem(), originalValue)
	case reflect.Interface:
		// If it is an interface (which is very similar to a pointer), do basically the
		// same as for the pointer. Though a pointer is not the same as an interface so
		// note that we have to call Elem() after creating a new object because otherwise
		// we would end up with an actual pointer
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		translateRecursiveOld(copyValue, originalValue)
		cpy.Set(copyValue)

		// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			translateRecursiveOld(cpy.Field(i), original.Field(i))
		}

		// If it is a slice we create a new slice and translate each element
	case reflect.Slice:
		cpy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			translateRecursiveOld(cpy.Index(i), original.Index(i))
		}

		// If it is a map we create a new map and translate each value
	case reflect.Map:
		cpy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			translateRecursiveOld(copyValue, originalValue)
			cpy.SetMapIndex(key, copyValue)
		}

		// Otherwise we cannot traverse anywhere so this finishes the the recursion
	case reflect.String:
		pp.Println(original.String())
		// If it is a string translate it (yay finally we're doing what we came for)
		translated := translateValue(original.Interface().(string))

		if translated.Kind() == reflect.Map {
			cpy = reflect.MakeMap(translated.Type())
			translateRecursiveOld(cpy, translated)
		} else {
			cpy.Set(translated)
		}
	default:
		// And everything else will simply be taken from the original
		cpy.Set(original)
	}
}

func translateValue(i string) reflect.Value {
	return reflect.ValueOf(nil)
}

func translateString(input string) interface{} {
	getAttRegex := regexp.MustCompile(`\$\((.+)\[(.+)]\)`)
	getAttValue := getAttRegex.FindStringSubmatch(input)
	if getAttValue != nil {
		return fmt.Sprintf(`{ "Fn::GetAtt" : [ "%s", "%s" ] }`, getAttValue[1], getAttValue[2])
	}

	refRegex := regexp.MustCompile(`\$\((.+)\)`)
	refValue := refRegex.FindStringSubmatch(input)
	if refValue != nil {
		newValue := make(map[string]string)
		newValue["Ref"] = refValue[1]
		return newValue
	}

	return input
}
