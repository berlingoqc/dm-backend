package dm

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/berlingoqc/dm-backend/tr"
	"github.com/berlingoqc/dm-backend/tr/triggers"

	"github.com/berlingoqc/dm-backend/file"
	"github.com/berlingoqc/dm-backend/program"
	"github.com/berlingoqc/dm-backend/rpcproxy"
	"github.com/berlingoqc/dm-backend/webserver"

	"github.com/gorilla/mux"

	"github.com/berlingoqc/dm-backend/aria2"

	// load les tasks de base
	_ "github.com/berlingoqc/dm-backend/tr/task/impl"

	"github.com/berlingoqc/dm-backend/tr/pipeline"
	"github.com/berlingoqc/dm-backend/tr/task"

	"github.com/berlingoqc/find-download-link/api"
	"github.com/berlingoqc/find-download-link/indexer"
)

// FrontEnd ..
type FrontEnd struct {
	Serve    bool   `json:"serve"`
	Location string `json:"location"`
}

// Config ...
type Config struct {
	URL              string                                  `json:"url"`
	FrontEnd         FrontEnd                                `json:"front-end"`
	Handler          map[string]*rpcproxy.RPCHandlerEndpoint `json:"handler"`
	Program          []*program.Settings                     `json:"program"`
	Security         *webserver.SecurityConfig               `json:"security"`
	FindDownloadLink indexer.Settings                        `json:"find-download-link"`
	Pipeline         tr.Settings                             `json:"pipeline"`
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
	rpcproxy.WSMessageChannel = make(chan rpcproxy.WSMessage)

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
	}

	// Configure les triggers utilisés pour le pipeline
	wsTrapper := &triggers.WSTrapper{
		Handler: map[string]triggers.WSEventHandler{
			"aria2": &aria2.FileHandler{
				Config: config.Handler["aria2"],
			},
		},
		Events: map[int64]triggers.WatchInfo{},
	}
	triggers.Triggers["websocket"] = wsTrapper

	// Initiliaze le module de pipeline
	tr.InitPipelineModule(config.Pipeline)

	// Ajoute les modules RPC
	localHandler := &rpcproxy.LocalHandler{
		Handlers: make(map[string]interface{}),
	}
	localHandler.Handlers["task"] = &task.RPCTask{}
	localHandler.Handlers["pipeline"] = &pipeline.RPCPipeline{}
	localHandler.Handlers["trigger"] = &triggers.RPC{}
	localHandler.Handlers["fe"] = &file.RPC{}
	localHandler.Handlers["program"] = &program.RPC{}
	localHandler.Handlers["proxyws"] = &rpcproxy.RPCWS{}

	var err error

	// Section configuration de find-download-link

	indexer.SetSettings(config.FindDownloadLink)
	indexer.FeedBack = func(ns, event string, data interface{}) {
		rpcproxy.SendMessageWS("dm", ns, event, data)
	}

	localHandler.Handlers["findDownload"], err = api.GetFindDownloadAPI()
	if err != nil {
		return nil, err
	}
	localHandler.Handlers["findDownloadDaemon"], err = api.GetDaemonFindDownloadAPI()
	if err != nil {
		return nil, err
	}

	// Ajout le hander local avec la reflexion sur les existants
	rpcproxy.RegisterLocalHandler("dm", localHandler)

	// Set the pipeline feedback to web socket
	pipeline.FeedBack = func(namespace, event string, data interface{}) {
		rpcproxy.SendMessageWS("dm", namespace, event, data)
	}

	r := mux.NewRouter()

	// Enregistre le proxy RPC dans le router
	handler := rpcproxy.Register(r)
	// Servir le front-end si nécessaire
	if config.FrontEnd.Serve {
		fileServer := http.FileServer(http.Dir(config.FrontEnd.Location))
		r.PathPrefix("/").Handler(fileServer)
	}

	// Ajout l'auth si nécessaire
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
