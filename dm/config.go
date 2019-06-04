package dm

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/berlingoqc/dm/file"
	"github.com/berlingoqc/dm/rpcproxy"
	"github.com/berlingoqc/dm/tr"
	"github.com/berlingoqc/dm/webserver"

	"github.com/gorilla/mux"

	// load le module
	_ "github.com/berlingoqc/dm/aria2"
	// load les tasks de base
	_ "github.com/berlingoqc/dm/tr/tasks"
)

// Config ...
type Config struct {
	URL     string                                  `json:"url"`
	Handler map[string]*rpcproxy.RPCHandlerEndpoint `json:"handler"`
}

// RPCMethod ...
type RPCMethod struct{}

// Hello ...
func (r *RPCMethod) Hello(world string) (string, error) {
	return "Hello " + world, errors.New("CCA")
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
		if handler, ok := tr.Handlers[k]; ok {
			handler.SetConfig(i)
		}
	}

	localHandler := &rpcproxy.LocalHandler{
		Handlers: make(map[string]interface{}),
	}
	localHandler.Handlers["test"] = &RPCMethod{}
	localHandler.Handlers["task"] = &tr.RPCTask{}
	localHandler.Handlers["pipeline"] = &tr.RPCPipeline{}

	rpcproxy.RegisterLocalHandler("dm", localHandler)

	//rpcproxy.Handlers["dm"] = localHandler

	r := mux.NewRouter()
	rpcproxy.Register(r)

	tr.Pipelines["cpMovie"] = tr.Pipeline{
		ID:   "cpMovie",
		Name: "cpMovie",
		Node: tr.TaskNode{
			TaskID: "CPP",
			Params: map[string]interface{}{
				"destination": "/home/wq/Project/dm/",
			},
			NextNode: []tr.TaskNode{
				tr.TaskNode{
					TaskID: "ZIP",
					Params: map[string]interface{}{
						"destination": "/home/wq/Project/dm/d",
						"methode":     "unzip",
					},
				},
			},
		},
	}

	tr.RegisterPipeline["/var/share/Download//d.zip"] = "cpMovie"

	return &webserver.WebServer{
		Logger:      log.New(os.Stdout, "", 0),
		Mux:         r,
		ChannelStop: make(chan os.Signal, 1),
		Hs: &http.Server{
			Addr:    config.URL,
			Handler: r,
		},
	}, nil
}
