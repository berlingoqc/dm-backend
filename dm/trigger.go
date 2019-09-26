package dm

import (
	"github.com/berlingoqc/dm-backend/rpcproxy"
	"github.com/berlingoqc/dm-backend/tr/pipeline"
	"github.com/berlingoqc/dm-backend/tr/triggers"
)

// DownloadEvent ...
type DownloadEvent string

const (
	// OnDownloadStart ...
	OnDownloadStart DownloadEvent = "onDownloadStart"
	// OnDownloadOver ...
	OnDownloadOver DownloadEvent = "onDownloadOver"
)

func messageTrapper(msg rpcproxy.WSMessage) {
	rpcCall := msg.Data

	if fileHandler, ok := pipeline.Handlers[msg.Namespace]; ok {
		if event := fileHandler.GetEvent(rpcCall.Method); event != "" {
			if event == string(OnDownloadOver) {
				println("DOWNLOAD IS OVER")
				if filePath, err := fileHandler.GetFilePath(rpcCall.Params); err == nil {
					triggers.TriggerEventChannel <- triggers.TriggerEvent{
						Event: event,
						File:  filePath,
					}
				} else {
					println("ERROR getting filepath ", err.Error())
				}
			}
		} else {
			println("No event for this thing")
		}
	}
}
