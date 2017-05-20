package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jrzimmerman/bestrida-server-go/handlers"
	"github.com/jrzimmerman/bestrida-server-go/models"
	"github.com/jrzimmerman/bestrida-server-go/utils"
)

func main() {
	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	port := utils.GetEnvString("PORT")

	mux := handlers.API()

	// close DB connection
	defer models.Close()

	srv := &http.Server{Addr: ":" + port, Handler: mux}

	go func() {
		log.Printf("Listening on http://0.0.0.0:%s\n", port)

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stopChan

	log.Println("\nShutting down the server...")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	srv.Shutdown(ctx)

	log.Println("Server gracefully stopped")
}
