package dm

import (
	"encoding/json"

	"github.com/berlingoqc/dm/rpcproxy"
	"github.com/berlingoqc/dm/tr"
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
	rpcCall := rpcproxy.RPCCall{}
	if err := json.Unmarshal(msg.Data, &rpcCall); err != nil {
		println("FAILED to unmarshall RPCCALL event ", err.Error())
	}

	if fileHandler, ok := tr.Handlers[msg.Namespace]; ok {
		if event := fileHandler.GetEvent(rpcCall.Method); event != "" {
			if event == string(OnDownloadOver) {
				println("DOWNLOAD IS OVER")
				if filePath, err := fileHandler.GetFilePath(rpcCall.Params); err == nil {
					tr.TriggerEventChannel <- tr.TriggerEvent{
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

	} else {
		println("Cant handle this namespace in filehandler ", msg.Namespace)
	}
}
