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

func replaceParams(m map[string]string, data map[string]string) (map[string]string, error) {
	rtr := make(map[string]string)
	for k, v := range m {
		match := re.FindStringSubmatch(v)
		if len(match) > 1 {
			if replacement, ok := data[match[1]]; ok {
				println("REPLACING ", match[0], "with ", replacement)
				rtr[k] = re.ReplaceAllString(k, replacement)
				continue
			} else {
				return nil, errors.New("Key not found " + match[0])
			}
		}
		rtr[k] = v
	}
	return rtr, nil
}
