package dm

import (
	"log"
	"net/http"
	"os"

	"github.com/berlingoqc/dm-backend/file"
	"github.com/berlingoqc/dm-backend/program"
	"github.com/berlingoqc/dm-backend/rpcproxy"
	"github.com/berlingoqc/dm-backend/tr"
	"github.com/berlingoqc/dm-backend/webserver"

	"github.com/gorilla/mux"

	// load le module
	_ "github.com/berlingoqc/dm-backend/aria2"
	"github.com/berlingoqc/dm-backend/tr/pipeline"
	"github.com/berlingoqc/dm-backend/tr/task"

	// load les tasks de base
	_ "github.com/berlingoqc/dm-backend/tr/tasks"
)

// Config ...
type Config struct {
	URL     string                                  `json:"url"`
	Handler map[string]*rpcproxy.RPCHandlerEndpoint `json:"handler"`
	Program []*program.Settings                     `json:"program"`
}

// Load ...
func Load(filepath string) (*webserver.WebServer, error) {
	var config = &Config{}

	if err := file.LoadJSON(filepath, config); err != nil {
		return nil, err
	}

	// Demarre ou pas les programmes
	program.Start(config.Program)

	// Configure la fonction d'handle des messages websocket
	rpcproxy.WSMessageTrapper = messageTrapper

	// Configure les handler pour le proxy rpc
	for k, i := range config.Handler {
		if handler, ok := rpcproxy.Handlers[k]; ok {
			println("Setting config for ", k)
			handler.SetConfig(i)
			// Demarre le client websocket pour le handler
			rpcproxy.StartWebSocketClient(*handler.GetConfig())
		}
		if handler, ok := pipeline.Handlers[k]; ok {
			handler.SetConfig(i)
		}
	}

	localHandler := &rpcproxy.LocalHandler{
		Handlers: make(map[string]interface{}),
	}
	localHandler.Handlers["task"] = &task.RPCTask{}
	localHandler.Handlers["pipeline"] = &pipeline.RPCPipeline{}
	localHandler.Handlers["tr"] = &tr.RPC{}

	// Ajout le hander local avec la reflexion sur les existants
	rpcproxy.RegisterLocalHandler("dm", localHandler)

	// Set the pipeline feedback to web socket
	pipeline.FeedBack = func(namespace, event string, data interface{}) {
		rpcproxy.SendMessageWS("dm", namespace, event, data)
	}

	r := mux.NewRouter()
	handler := rpcproxy.Register(r)

	return &webserver.WebServer{
		Logger:      log.New(os.Stdout, "", 0),
		Mux:         r,
		ChannelStop: make(chan os.Signal, 1),
		Hs: &http.Server{
			Addr:    config.URL,
			Handler: handler,
		},
	}, nil
}
