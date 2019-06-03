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
}

// StartAsync demarre le serveur web
func (w *WebServer) StartAsync() {
	// Crée mon channel pour le signal d'arret
	w.ChannelStop = make(chan os.Signal, 1)

	signal.Notify(w.ChannelStop, os.Interrupt, syscall.SIGTERM)

	go func() {
		w.Logger.Printf("Starting yawf.ca at %v\n", w.Hs.Addr)
		if err := w.Hs.ListenAndServe(); err != http.ErrServerClosed {
			w.Logger.Fatal(err)
		}
	}()
}

// Start demarrer le server web de facon synchrone
func (w *WebServer) Start() {
	if err := w.Hs.ListenAndServe(); err != nil {
		w.Logger.Fatal(err)
	}
}

// Stop arrête le serveur web
func (w *WebServer) Stop() {
	w.Logger.Println("Fermeture du serveur")
	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c() // release les ressources du context

	w.Hs.Shutdown(ctx)

	w.Logger.Println("Serveur eteint ...")
}
