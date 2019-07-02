package webserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// SecurityConfig ...
type SecurityConfig struct {
	Secure  bool   `json:"secure"`
	Cert    string `json:"cert"`
	Key     string `json:"key"`
	AuthKey string `json:"authkey"`
}

// WebServer my webserver that work with my modules implements of IWebServer
type WebServer struct {
	// Logger of the web site and its module
	Logger *log.Logger
	// Mux is the base router of the website
	Mux *mux.Router
	// Hs is the http server that run
	Hs *http.Server
	// ChannelStop is the channel to stop the webserver
	ChannelStop chan os.Signal
	// Settings de la securité
	Security *SecurityConfig
}

func (w *WebServer) start() {
	var err error
	if w.Security.Secure {
		println("Starting mode secure")
		err = w.Hs.ListenAndServeTLS(w.Security.Cert, w.Security.Key)
	} else {
		println("Starting mode unsecure")
		err = w.Hs.ListenAndServe()
	}
	if err != http.ErrServerClosed {
		w.Logger.Fatal(err)
	}
}

// StartAsync demarre le serveur web
func (w *WebServer) StartAsync() {
	// Crée mon channel pour le signal d'arret
	w.ChannelStop = make(chan os.Signal, 1)

	signal.Notify(w.ChannelStop, os.Interrupt, syscall.SIGTERM)

	go func() {
		w.start()
	}()
}

// Start demarrer le server web de facon synchrone
func (w *WebServer) Start() {
	w.start()
}

// Stop arrête le serveur web
func (w *WebServer) Stop() {
	w.Logger.Println("Fermeture du serveur")
	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c() // release les ressources du context

	w.Hs.Shutdown(ctx)

	w.Logger.Println("Serveur eteint ...")
}
