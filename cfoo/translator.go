package cfoo

import (
	"reflect"

	"regexp"
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
	case reflect.Slice:
		inputSlice := input.([]interface{})
		var translated []interface{}
		for _, element := range inputSlice {
			translated = append(translated, translate(element))
		}
		return translated
	default:
		panic("Unsupported type " + inputValue.Kind().String())
	}
}

func translateString(input string) interface{} {
	getAttRegex := regexp.MustCompile(`\$\((.+)\[(.+)]\)`)
	getAttValue := getAttRegex.FindStringSubmatch(input)
	if getAttValue != nil {
		getAttMap := make(map[string][]string)
		attrs := []string{getAttValue[1], getAttValue[2]}
		getAttMap["Fn::GetAtt"] = attrs
		return getAttMap
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
