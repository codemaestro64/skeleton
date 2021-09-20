package controllers

import (
	"net/http"

	"github.com/codemaestro64/skeleton/web/context"
)

type HomeController struct {
}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (c *HomeController) GetIndex(ctx *context.AppContext) {
	res := map[string]interface{}{
		"id":   15,
		"name": "skeleton structure",
	}

	ctx.JSON(http.StatusOK, res)
}
