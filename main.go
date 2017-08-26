package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/jrzimmerman/bestrida-server-go/handlers"
	"github.com/jrzimmerman/bestrida-server-go/models"
	"github.com/jrzimmerman/bestrida-server-go/utils"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

func main() {
	port := utils.GetEnvString("PORT")
	mux := handlers.API()
	// close DB connection
	defer models.Close()
	c := cron.New()
	c.AddFunc("0 * * * * *", func() {
		log.Print("Starting cron complete")
		handlers.CronComplete()
		log.Print("Ending cron complete")
	})
	c.Start()
	defer c.Stop()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// We want to report the listener is closed.
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		log.Infof("Listening on http://localhost:%s\n", port)

		srv.ListenAndServe()
		wg.Done()
	}()

	// Listen for an interrupt signal from the OS. Use a buffered
	// channel because of how the signal package is implemented.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt)

	// Wait for a signal to shutdown.
	<-osSignals

	// Create a context to attempt a graceful 5 second shutdown.
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Println("\nShutting down the server...")

	// Attempt the graceful shutdown by closing the listener and
	// completing all inflight requests.
	if err := srv.Shutdown(ctx); err != nil {
		log.Debugf("shutdown : Graceful shutdown did not complete in %v : %v", timeout, err)

		// Looks like we timedout on the graceful shutdown. Kill it hard.
		if err := srv.Close(); err != nil {
			log.WithError(err).Errorf("shutdown : Error killing server : %v", err)
		}
	}

	// Wait for the listener to report it is closed.
	wg.Wait()
	log.Println("Server gracefully stopped")
}
