package web

import (
	"github.com/codemaestro64/skeleton/web/context"
	"github.com/codemaestro64/skeleton/web/controllers"
	"github.com/labstack/echo/v4"
)

type Handler func(*context.AppContext)

func (s *Server) registerRoutes() {
	homeController := controllers.NewHomeController()

	s.GET("/", homeController.GetIndex)
}

func (s *Server) GET(path string, handler Handler) {
	s.echo.GET(path, s.resolveHandler(handler))
}

func (s *Server) POST(path string, handler Handler) {
	s.echo.POST(path, s.resolveHandler(handler))
}

func (s *Server) PUT(path string, handler Handler) {
	s.echo.PUT(path, s.resolveHandler(handler))
}

func (s *Server) PATCH(path string, handler Handler) {
	s.echo.PATCH(path, s.resolveHandler(handler))
}

func (s *Server) DELETE(path string, handler Handler) {
	s.echo.DELETE(path, s.resolveHandler(handler))
}

func (s *Server) resolveHandler(handler Handler) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		c := &context.AppContext{
			Context: ctx,
			Cache: s.cache,
		}

		handler(c)
		return nil
	}
}
