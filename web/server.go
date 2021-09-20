package web

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/codemaestro64/skeleton/config"
	"github.com/codemaestro64/skeleton/lib/cache"
	"github.com/codemaestro64/skeleton/lib/logger"
	"github.com/codemaestro64/skeleton/web/models"
	"github.com/labstack/echo/v4"
)

type Server struct {
	echo   *echo.Echo
	config *config.Config
	db     *models.Database
	cache  *cache.Cache
	logger *logger.Logger
}

func NewServer(cfg *config.Config, logger *logger.Logger) (*Server, error) {
	s := &Server{
		config: cfg,
		echo:   echo.New(),
		logger: logger,
	}
	s.echo.HideBanner = true

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
	go func() {
		if err := s.echo.Start(s.config.App.Port); err != nil {
			s.logger.Info().Msg("shutting down the server")
		}
	}()

	s.Shutdown()
	return nil
}

func (s *Server) Shutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// close database connection
	if s.db != nil {
		s.logger.Info().Msg("Closing database connection...")
		if err := s.db.Disconnect(); err != nil {
			s.logger.Fatal().Msg(err.Error())
		}
	}

	// shutdown http server
	s.logger.Info().Msg("Shutting down server...")
	if err := s.echo.Shutdown(ctx); err != nil {
		s.logger.Fatal().Msg(err.Error())
	}
}