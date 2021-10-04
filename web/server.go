package web

import (
	//"context"
	//"os"
	//"os/signal"
	//"time"
	"net/http"
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
	s.logger.Info().Str("port", s.config.App.Port).Msg("Serving...")
	http.ListenAndServe(s.config.App.Port, s.router)
	s.Shutdown()

	/**go func() {
		if err := s.echo.Start(s.config.App.Port); err != nil {
			s.logger.Info().Msg("shutting down the server")
		}
	}()

	s.Shutdown()**/
	return nil
}

func (s *Server) Shutdown() {
	/**quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	// shutdown http server
	s.logger.Info().Msg("Shutting down server...")
	if err := s.echo.Shutdown(ctx); err != nil {
		s.logger.Fatal().Msg(err.Error())
	}**/
}
