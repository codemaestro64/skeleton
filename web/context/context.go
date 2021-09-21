package context

import (
	"github.com/codemaestro64/skeleton/lib/cache"
	"github.com/codemaestro64/skeleton/lib/logger"
	"github.com/labstack/echo/v4"
)

type AppContext struct {
	echo.Context
	Cache  *cache.Cache
	Logger *logger.Logger
}
