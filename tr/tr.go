package tr

import (
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
func InitPipelineModule(settings Settings) {
	chTriggerEvent, chClosingSignal = triggers.InitTriggers()
	stopingPipelineRunnerCh = make(chan interface{})
	go func() {
		// The stack of pipeline waiting to be execute
		stackPipeline := []triggers.PipelineTrigger{}
		// The channel to get when pipeline will stop
		runningPipeline := make(map[string]chan int)
		for {
			select {
			case _ = <-stopingPipelineRunnerCh:
				return
			case d := <-chTriggerEvent:
				stackPipeline = append(stackPipeline, d)
				break
			default:
				nbr := getPipelineRunning(&runningPipeline)
				if nbr < settings.ConcurrentPipeline {
					if te := getNextPipeline(stackPipeline); te != nil {
						if p, ok := pipeline.Pipelines[te.PipelineID]; ok {
							status, err := pipeline.StartPipeline(te.File, &p, te.Data)
							if err != nil {
								println(err)
							}
							runningPipeline[status.Pipeline] = status.ChanPipelineSignal
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
	}()
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
