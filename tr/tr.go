package tr

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/berlingoqc/dm-backend/tr/pipeline"
	"github.com/berlingoqc/dm-backend/tr/triggers"
)

// Settings ...
type Settings struct {
	ConcurrentPipeline int `json:"concurrent_pipeline"`
}

var (
	chTriggerEvent  chan triggers.PipelineTrigger
	chClosingSignal chan interface{}
	settings        Settings
)

var stopingPipelineRunnerCh chan interface{}

// InitPipelineModule ...
func InitPipelineModule(settings Settings) chan interface{} {
	triggers.Triggers["manual"] = &triggers.ManualFileTrigger{}
	triggers.Triggers["file_watch"] = &triggers.FileWatchTrigger{}
	folderPath := pipeline.GetWorkingPath()
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		println("Loading pipeline ", f.Name())
		id := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
		pip, err := pipeline.GetPipelineFile(id)
		if err != nil {
			panic(err)
		}
		pipeline.Pipelines[id] = *pip
	}

	chTriggerEvent, chClosingSignal = triggers.InitTriggers()
	stopingPipelineRunnerCh = make(chan interface{})
	go func() {
		// The stack of pipeline waiting to be execute
		stackPipeline := []triggers.PipelineTrigger{}
		// The channel to get when pipeline will stop
		runningPipeline := make(map[string]chan int)
	Loop:
		for {
			select {
			case _ = <-stopingPipelineRunnerCh:
				break Loop
			case d := <-chTriggerEvent:
				stackPipeline = append(stackPipeline, d)
				break
			default:
				nbr := getPipelineRunning(&runningPipeline)
				if len(stackPipeline) != 0 && nbr < settings.ConcurrentPipeline {
					if te := getNextPipeline(stackPipeline); te != nil {
						if p, ok := pipeline.Pipelines[te.PipelineID]; ok {
							status, err := pipeline.StartPipeline(te.File, &p, te.Data)
							if err != nil {
								println(err)
							}
							runningPipeline[status.Pipeline] = status.ChanPipelineSignal
							if len(stackPipeline) == 0 {
								stackPipeline = stackPipeline[0:0]
							} else {
								stackPipeline = stackPipeline[1:len(stackPipeline)]
							}
						} else {
							// Pipeline doesnt exists
						}
					} else {
						// No next pipeline
					}
				} else {
					// Cant run pipeline now to many are running
				}
				break
			}
		}

		println("END OF THIS")
	}()

	return stopingPipelineRunnerCh
}

// StopPipelineModule ...
func StopPipelineModule() {
	for i := 0; i < 5; i++ {
		chClosingSignal <- i
	}
	stopingPipelineRunnerCh <- 0
}

func getPipelineRunning(m *map[string]chan int) int {
	i := 0
	for id, ch := range *m {
		select {
		case _ = <-ch:
			println("Deletin pipeline with an id because its over baby blue")
			delete(*m, id)
			break
		default:
			i = i + 1
		}
	}
	return i
}

func getNextPipeline(list []triggers.PipelineTrigger) *triggers.PipelineTrigger {
	if len(list) > 0 {
		pt := &list[0]
		list = list[1:len(list)]
		return pt
	}
	return nil
}
