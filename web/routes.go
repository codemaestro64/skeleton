package web

import (
	"github.com/codemaestro64/skeleton/web/context"
	"github.com/codemaestro64/skeleton/web/controllers"
)

type HandlerFunc func(*context.AppContext)
type MiddlewareFunc func(*context.AppContext, func())

func (s *Server) registerRoutes() {
	homeController := controllers.NewHomeController()

	s.GET("/", homeController.GetIndex)
	g := s.Group("/test", m1, m2)
	g.GET("/", homeController.GetIndex)
}

func m1(ctx *context.AppContext, next func()) {

	next()
}
func m2(ctx *context.AppContext, next func()) {

	next()
}
