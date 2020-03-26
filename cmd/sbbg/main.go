package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"networkgaming.co.uk/techtest/pkg/game"
	"networkgaming.co.uk/techtest/pkg/socket"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Starting server")

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("role", "gozer").
		Str("host", "localhost").
		Logger()

	ctx, cancel := context.WithCancel(context.Background())

	simple := game.NewGame(game.NewRNG(game.MaxNum))
	engineConfig := &game.EngineConfig{
		GameSpeed:    1 * time.Second,
		WaitingCount: 10,
		ManualRun:    false,
	}
	gameEngine := game.NewEngine(simple, engineConfig)
	gameBroadcaster := game.NewBroadcaster(gameEngine.Event)
	gameBroadcaster.Start(ctx)
	gameEngine.Start()

	webSocketHandler := socket.New(gameBroadcaster)
	joinGameHandler := game.NewJoinGameHandler(gameEngine.Action)

	var allowedOrigins []string
	allowedOrigins = append(allowedOrigins, "http://localhost:8091")

	crossOrigin := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300,
	})
	router := chi.NewRouter()
	router.Use(chiMiddleware.Timeout(60 * time.Second))
	router.Use(hlog.NewHandler(logger))
	router.Use(crossOrigin.Handler)

	router.Route("/", func(r chi.Router) {
		r.Get("/", healthCheckHandler)
	})

	router.Route("/subscribe", func(r chi.Router) {
		r.Get("/", webSocketHandler.Subscribe)
	})

	router.Route("/join", func(r chi.Router) {
		r.Post("/", joinGameHandler.JoinGame)
	})

	srv := &http.Server{
		Handler:      router,
		Addr:         "localhost:8089",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		startMsg := fmt.Sprintln("Starting Server on port 8089")
		logger.Info().Msg(startMsg)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv, cancel)
}

func waitForShutdown(srv *http.Server, cancel context.CancelFunc) {

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan
	cancel()
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Info().Msg("Shutting down")
	os.Exit(0)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintln("It's Working, Simple Browser Based Game Server is up and running boyee!")))
}
