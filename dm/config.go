package dm

import (
	"log"
	"net/http"
	"os"

	"github.com/berlingoqc/dm-backend/file"
	"github.com/berlingoqc/dm-backend/program"
	"github.com/berlingoqc/dm-backend/rpcproxy"
	"github.com/berlingoqc/dm-backend/webserver"

	"github.com/gorilla/mux"

	// load le module
	_ "github.com/berlingoqc/dm-backend/aria2"
	// load les tasks de base
	"github.com/berlingoqc/dm-backend/tr/pipeline"
	"github.com/berlingoqc/dm-backend/tr/task"
	_ "github.com/berlingoqc/dm-backend/tr/tasks"
)

// Config ...
type Config struct {
	URL     string                                  `json:"url"`
	Handler map[string]*rpcproxy.RPCHandlerEndpoint `json:"handler"`
	Program map[string]*program.Settings            `json:"program"`
}

// Load ...
func Load(filepath string) (*webserver.WebServer, error) {
	var config = &Config{}

	if err := file.LoadJSON(filepath, config); err != nil {
		return nil, err
	}

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

	// Ajout le hander local avec la reflexion sur les existants
	rpcproxy.RegisterLocalHandler("dm", localHandler)

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
