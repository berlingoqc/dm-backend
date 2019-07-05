package pipeline

import (
	"errors"
	"reflect"
	"regexp"

	"github.com/berlingoqc/dm-backend/tr/task"
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

func replaceParamPipelineTask(pipeline *Pipeline, data map[string]interface{}) error {
	nodes := make(chan *task.TaskNode, 5)
	nodes <- pipeline.Node
	for {
		select {
		case d := <-nodes:
			for k, v := range d.Params {
				if str, ok := v.(string); ok {
					match := re.FindStringSubmatch(str)
					if len(match) > 1 {
						if replacement, ok := data[match[1]]; ok {
							println("REPLACING ", match[0], "with ", replacement)
							d.Params[k] = re.ReplaceAllString(str, replacement.(string))
						} else {
							return errors.New("Key not found " + match[0])
						}
					}
				}
			}

			for _, v := range d.NextNode {
				if v == nil {
					continue
				}
				nodes <- v
			}
			break
		default:
			return nil
		}
	}
}
