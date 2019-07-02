package dm

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/berlingoqc/dm-backend/file"
	"github.com/berlingoqc/dm-backend/program"
	"github.com/berlingoqc/dm-backend/rpcproxy"
	"github.com/berlingoqc/dm-backend/tr"
	"github.com/berlingoqc/dm-backend/webserver"

	"github.com/gorilla/mux"

	// load le module
	_ "github.com/berlingoqc/dm-backend/aria2"
	_ "github.com/berlingoqc/dm-backend/ydl"

	"github.com/berlingoqc/dm-backend/tr/pipeline"
	"github.com/berlingoqc/dm-backend/tr/task"

	// load les tasks de base
	_ "github.com/berlingoqc/dm-backend/tr/tasks"
)

// Config ...
type Config struct {
	URL      string                                  `json:"url"`
	Handler  map[string]*rpcproxy.RPCHandlerEndpoint `json:"handler"`
	Program  []*program.Settings                     `json:"program"`
	Security *webserver.SecurityConfig               `json:"security"`
}

// Load ...
func Load(filepath string) (*webserver.WebServer, error) {
	var config = &Config{}

	if err := file.LoadJSON(filepath, config); err != nil {
		return nil, err
	}

	// Demarre ou pas les programmes
	program.Start(config.Program)

	// WAIT TEMPORAIRE
	time.Sleep(1 * time.Second)

	// Configure la fonction d'handle des messages websocket
	rpcproxy.WSMessageTrapper = messageTrapper

	// Configure les handler pour le proxy rpc
	for k, i := range config.Handler {
		if handler, ok := rpcproxy.Handlers[k]; ok {
			println("Setting config for ", k)
			handler.SetConfig(i)
			// Demarre le client websocket pour le handler
			rpcproxy.StartWebSocketClient(*handler.GetConfig())
		} else {
			println("Could not found handler ", k)
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
	localHandler.Handlers["program"] = &program.RPC{}
	localHandler.Handlers["proxyws"] = &rpcproxy.RPCWS{}

	// Ajout le hander local avec la reflexion sur les existants
	rpcproxy.RegisterLocalHandler("dm", localHandler)

	// Set the pipeline feedback to web socket
	pipeline.FeedBack = func(namespace, event string, data interface{}) {
		rpcproxy.SendMessageWS("dm", namespace, event, data)
	}

	r := mux.NewRouter()
	handler := rpcproxy.Register(r)

	r.Use(mux.CORSMethodMiddleware(r))

	if config.Security.AuthKey != "" {
		rpcproxy.ValidToken = func(token string, r *http.Request) error {
			if token != config.Security.AuthKey {
				return errors.New("invalid authkey")
			}
			return nil
		}
		r.Use(rpcproxy.AuthMiddleware)
	}

	return &webserver.WebServer{
		Security:    config.Security,
		Logger:      log.New(os.Stdout, "", 0),
		Mux:         r,
		ChannelStop: make(chan os.Signal, 1),
		Hs: &http.Server{
			Addr:    config.URL,
			Handler: handler,
		},
	}, nil
}
