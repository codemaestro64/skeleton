package web

import (
	"github.com/codemaestro64/skeleton/web/controllers"
)

func (s *Server) registerRoutes() {
	homeController := controllers.NewHomeController()
	s.GET("/", homeController.GetIndex)
}
