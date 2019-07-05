package pipeline

import (
	"errors"
	"reflect"
	"regexp"
)

func cloneValue(source interface{}, destin interface{}) {
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(destin).Elem().Set(y.Elem())
	} else {
		destin = x.Interface()
	}
}

var re = regexp.MustCompile("\\$\\{(.*?)\\}")

func replaceParams(m map[string]interface{}, data map[string]interface{}) (map[string]interface{}, error) {
	rtr := make(map[string]interface{})
	for k, v := range m {
		if str, ok := v.(string); ok {
			match := re.FindStringSubmatch(str)
			if len(match) > 1 {
				if replacement, ok := data[match[1]]; ok {
					println("REPLACING ", match[0], "with ", replacement)
					rtr[k] = re.ReplaceAllString(str, replacement.(string))
					continue
				} else {
					return nil, errors.New("Key not found " + match[0])
				}
			}
		}
		rtr[k] = v
	}
	return rtr, nil
}
