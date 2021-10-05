package controllers

import (
	//"net/http"

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
		"name": "Michael",
	}

	ctx.WithStatusCode(200).JSON(res)

	//ctx.Cache.Put("name", "Michael", 86400)

	/**name, err := ctx.Cache.Get("name")
	if err != nil {
		res := map[string]interface{}{
			"error": err.Error(),
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	res := map[string]interface{}{
		"id":   15,
		"name": name.(string),
	}

	ctx.JSON(http.StatusOK, res)**/

	//ctx.Response.Write([]byte("hello world"))
}
