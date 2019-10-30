package tr

import (
	"github.com/berlingoqc/dm-backend/tr/pipeline"
	"github.com/berlingoqc/dm-backend/tr/triggers"
)

// InitPipelineModule ...
func InitPipelineModule() {
	ch := triggers.InitTriggers()
	go func() {
		for {
			d := <-ch
			println(d.Data["PATH"])
			println("Starting pipeline " + d.PipelineID)
			m := map[string]interface{}{}
			for k, v := range d.Data {
				m[k] = v
			}
			if p, ok := pipeline.Pipelines[d.PipelineID]; ok {
				pipeline.StartPipeline(d.File, &p, m)
			}
		}

	}()
}
