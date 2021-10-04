package context

import (
	"encoding/json"
	"net/http"

	"github.com/codemaestro64/skeleton/lib/cache"
	"github.com/codemaestro64/skeleton/lib/logger"
	"github.com/codemaestro64/skeleton/web/models"
	"github.com/go-chi/chi/v5"
)

type AppContext struct {
	response http.ResponseWriter
	request  *http.Request
	cache    *cache.Cache
	logger   *logger.Logger
	db       *models.Database
}

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json; charset=UTF-8"
)

func New() *AppContext {
	return &AppContext{}
}

func (c *AppContext) Setup(w http.ResponseWriter, req *http.Request) {
	c.response = w
	c.request = req
}

func (c *AppContext) SetCache(cache *cache.Cache) {
	c.cache = cache
}

func (c *AppContext) SetDB(db *models.Database) {
	c.db = db
}

func (c *AppContext) SetLogger(logger *logger.Logger) {
	c.logger = logger
}

func (c *AppContext) Request() *http.Request {
	return c.request
}

func (c *AppContext) Response() http.ResponseWriter {
	return c.response
}

func (c *AppContext) SetContentType(contentTypeValue string) {
	header := c.Response().Header()
	if header.Get(ContentType) == "" {
		header.Set(ContentType, contentTypeValue)
	}
}

func (c *AppContext) WithStatusCode(code int) *AppContext {
	c.Response().WriteHeader(code)
	return c
}

func (c *AppContext) WithContentType(contentTypeValue string) *AppContext {
	c.SetContentType(contentTypeValue)
	return c
}

func (c *AppContext) JSON(data interface{}) error {
	enc := json.NewEncoder(c.Response())
	c.SetContentType(ApplicationJSON)
	return enc.Encode(data)
}

func (c *AppContext) Param(name string) string {
	return chi.URLParam(c.Request(), name)
}
