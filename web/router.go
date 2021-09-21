package web

import (
	"github.com/codemaestro64/skeleton/web/context"
	"github.com/codemaestro64/skeleton/web/controllers"
)

type Handler func(*context.AppContext)

func (s *Server) registerRoutes() {
	homeController := controllers.NewHomeController()

	s.GET("/", homeController.GetIndex)
}
