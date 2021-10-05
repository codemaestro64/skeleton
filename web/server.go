package web

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/codemaestro64/skeleton/config"
	"github.com/codemaestro64/skeleton/lib/cache"
	"github.com/codemaestro64/skeleton/lib/logger"
	appContext "github.com/codemaestro64/skeleton/web/context"
	"github.com/codemaestro64/skeleton/web/models"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
	config *config.Config
	db     *models.Database
	cache  *cache.Cache
	logger *logger.Logger
	pool   sync.Pool
}

func NewServer(cfg *config.Config, logger *logger.Logger) (*Server, error) {
	s := &Server{
		config: cfg,
		router: chi.NewRouter(),
		logger: logger,
	}
	s.pool.New = func() interface{} {
		return appContext.New()
	}

	// open database connection
	logger.Info().Msg("Attempting database connection...")
	db, err := models.Connect(cfg.Database)
	if err != nil {
		return nil, err
	}
	s.db = db

	logger.Info().Msg("Initializing cache...")
	c, err := cache.New(cfg)
	if err != nil {
		return nil, err
	}
	s.cache = c

	// register routes
	logger.Info().Msg("Registering routes...")
	s.registerRoutes()

	return s, nil
}

func (s *Server) Serve() error {
	ctx, cancel := context.WithCancel(context.Background())

	srv := &http.Server{Addr: s.config.App.Port, Handler: s}
	go func(srv *http.Server) {
		s.logger.Info().Str("port", s.config.App.Port).Msg("Serving...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal().Msg(err.Error())

			// TODO recover from error
		}

	}(srv)

	go s.listenForShutdown(cancel)

	<-ctx.Done()
	s.logger.Info().Msg("Shutting down server...")
	srv.Shutdown(ctx)

	return nil
}

func (s *Server) listenForShutdown(cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit

	s.logger.Info().Msg("Quit signal received...")

	// close database connection
	if s.db != nil {
		s.logger.Info().Msg("Closing database connection...")
		if err := s.db.Disconnect(); err != nil {
			s.logger.Fatal().Msg(err.Error())
		}
	}

	// flush cache
	if s.cache != nil {
		s.logger.Info().Msg("Flushing cache...")
		s.cache.Flush()
	}

	cancel()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}
