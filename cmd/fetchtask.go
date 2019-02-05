package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paveloborin/fetch-task/pkg/cache"

	"github.com/paveloborin/fetch-task/pkg/handlers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"context"

	"github.com/julienschmidt/httprouter"
)

var shutdownChan chan bool
var osSignal chan os.Signal

func init() {
	shutdownChan = make(chan bool, 1)
	osSignal = make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

	zerolog.MessageFieldName = "MESSAGE"
	zerolog.LevelFieldName = "LEVEL"
	log.Logger = log.Output(os.Stderr).With().Str("PROGRAM", "fetchtask").Logger()
}

func main() {
	r := router()
	srv := http.Server{Addr: ":8000", Handler: r}
	log.Info().Msg("starting server on :8000")

	go func() {
		defer close(osSignal)
		defer close(shutdownChan)

		sig := <-osSignal
		log.Info().Msgf("get signal: %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Warn().Err(err).Msg("error shutdown server")
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Err(err)
	}

	<-shutdownChan
	log.Info().Msg("server shutdown completed, exiting")
}

func router() *httprouter.Router {
	r := httprouter.New()
	h := handlers.NewHandler(cache.NewStorage())

	r.POST("/task", h.PostHandler)
	r.GET("/task/:id", h.GetOneByIdHandler)
	r.GET("/task", h.GetAllHandler)
	r.DELETE("/task/:id", h.DeleteHandler)

	return r
}
