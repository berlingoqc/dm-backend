package pipeline

import (
	"io/ioutil"
	"path"
	"strings"
)

/*
* Load the pipeline file from the saved location , could be done later on if we want configuration
 */
func init() {
	folderPath := GetWorkingPath()
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		println("Loading pipeline ", f.Name())
		id := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
		pipeline, err := getPipelineFile(id)
		if err != nil {
			panic(err)
		}
		Pipelines[id] = *pipeline
	}
}
