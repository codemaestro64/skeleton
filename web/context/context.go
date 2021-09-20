package context

import (
	"github.com/labstack/echo/v4"
	"github.com/codemaestro64/skeleton/lib/cache"
)

type AppContext struct {
	echo.Context
	Cache *cache.Cache
}
