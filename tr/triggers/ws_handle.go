package triggers

import (
	"strings"

	"github.com/berlingoqc/dm-backend/rpcproxy"
)

// WSEventHandler ...
type WSEventHandler interface {
	GetFile(receive interface{}, expected interface{}) (string, bool, error)
}

// WSTrapper ...
type WSTrapper struct {
	Events  map[int64]WatchInfo
	Handler map[string]WSEventHandler
}

// AddWatch ...
func (w *WSTrapper) AddWatch(event string, param interface{}, settings *Settings) (int64, error) {
	var i int64 = 11
	w.Events[i] = WatchInfo{
		Event:    event,
		Param:    param,
		Settings: settings,
	}

	return i, nil
}

// GetWatchInfo ...
func (w *WSTrapper) GetWatchInfo() *map[int64]WatchInfo {
	return &w.Events
}

// Init ...
func (w *WSTrapper) Init(triggerChannel chan PipelineTrigger) {
	go func() {
		for {
			ws := <-rpcproxy.WSMessageChannel
			for id, watchInfo := range w.Events {
				if strings.Contains(ws.Data.Method, watchInfo.Event) {
					if handler, ok := w.Handler[ws.Namespace]; ok {
						file, b, _ := handler.GetFile(ws.Data.Params, watchInfo.Param)
						if b {
							triggerChannel <- PipelineTrigger{
								File:       file,
								PipelineID: watchInfo.Settings.PipelineID,
								Data:       watchInfo.Settings.Data,
							}
							if watchInfo.Settings.RemoveAfterRun {
								delete(w.Events, id)
							}
						}
					}
				}
			}
		}
	}()
}
