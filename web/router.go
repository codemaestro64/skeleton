package web

import (
	"net/http"

	appContext "github.com/codemaestro64/skeleton/web/context"
)

type HandlerFunc func(*appContext.AppContext)
type MiddlewareFunc func(*appContext.AppContext) bool

func (s *Server) USE(middleware ...MiddlewareFunc) {
	for _, v := range middleware {
		s.router.Use(s.resolveMiddlewareFunc(v))
	}
}

func (s *Server) GET(path string, handler HandlerFunc) {
	s.router.Get(path, s.resolveHandlerFunc(handler))
}

func (s *Server) POST(path string, handler HandlerFunc) {
	s.router.Post(path, s.resolveHandlerFunc(handler))
}

func (s *Server) PUT(path string, handler HandlerFunc) {
	s.router.Put(path, s.resolveHandlerFunc(handler))
}

func (s *Server) PATCH(path string, handler HandlerFunc) {
	s.router.Patch(path, s.resolveHandlerFunc(handler))
}

func (s *Server) DELETE(path string, handler HandlerFunc) {
	s.router.Delete(path, s.resolveHandlerFunc(handler))
}

func (s *Server) resolveHandlerFunc(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := s.resolveContext(w, r)
		handler(ctx)
		ctx.Reset()
		s.pool.Put(ctx)
	}
}

func (s *Server) resolveMiddlewareFunc(middleware MiddlewareFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if middleware(s.resolveContext(w, r)) {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func (s *Server) resolveContext(w http.ResponseWriter, r *http.Request) *appContext.AppContext {
	// get context from pool
	ctx := s.pool.Get().(*appContext.AppContext)
	ctx.Setup(w, r)
	ctx.SetDB(s.db)
	ctx.SetLogger(s.logger)
	ctx.SetCache(s.cache)

	return ctx
}

/**import (
	appContext "github.com/codemaestro64/skeleton/web/context"
	"github.com/labstack/echo/v4"
)

type Group struct {
	group             *echo.Group
	resolveHandler    func(HandlerFunc) echo.HandlerFunc
	resolveMiddleware func([]MiddlewareFunc) []echo.MiddlewareFunc
}

func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.group.Use(g.resolveMiddleware(middleware)...)
}

func (g *Group) GET(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.group.GET(path, g.resolveHandler(handler), g.resolveMiddleware(middleware)...)
}

func (g *Group) POST(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.group.POST(path, g.resolveHandler(handler), g.resolveMiddleware(middleware)...)
}

func (g *Group) PUT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.group.PUT(path, g.resolveHandler(handler), g.resolveMiddleware(middleware)...)
}

func (g *Group) PATCH(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.group.PATCH(path, g.resolveHandler(handler), g.resolveMiddleware(middleware)...)
}

func (g *Group) DELETE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.group.DELETE(path, g.resolveHandler(handler), g.resolveMiddleware(middleware)...)
}

func (s *Server) Group(prefix string, middleware ...MiddlewareFunc) *Group {
	g := &Group{
		resolveHandler:    s.resolveHandlerFunc,
		resolveMiddleware: s.resolveMiddlewareFuncs,
	}

	g.group = s.echo.Group(prefix, g.resolveMiddleware(middleware)...)
	return g
}

func (s *Server) Use(middleware ...MiddlewareFunc) {
	s.echo.Use(s.resolveMiddlewareFuncs(middleware)...)
}

func (s *Server) GET(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	s.echo.GET(path, s.resolveHandlerFunc(handler), s.resolveMiddlewareFuncs(middleware)...)
}

func (s *Server) POST(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	s.echo.POST(path, s.resolveHandlerFunc(handler), s.resolveMiddlewareFuncs(middleware)...)
}

func (s *Server) PUT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	s.echo.PUT(path, s.resolveHandlerFunc(handler), s.resolveMiddlewareFuncs(middleware)...)
}

func (s *Server) PATCH(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	s.echo.PATCH(path, s.resolveHandlerFunc(handler), s.resolveMiddlewareFuncs(middleware)...)
}

func (s *Server) DELETE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	s.echo.DELETE(path, s.resolveHandlerFunc(handler), s.resolveMiddlewareFuncs(middleware)...)
}

func (s *Server) resolveContext(ctx echo.Context) *appContext.AppContext {
	s.db.NewSession()

	c := &appContext.AppContext{
		Context: ctx,
		Cache:   s.cache,
		Logger:  s.logger,
		DB:      s.db,
	}

	return c
}

func (s *Server) resolveHandlerFunc(handler HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		handler(s.resolveContext(ctx))
		return nil
	}
}

func (s *Server) resolveMiddlewareFuncs(middleware []MiddlewareFunc) []echo.MiddlewareFunc {
	m := make([]echo.MiddlewareFunc, len(middleware))

	for index := range middleware {
		i := index
		m[i] = func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ctx echo.Context) error {
				nextFunc := func() {
					next(ctx)
				}

				middleware[i](s.resolveContext(ctx), nextFunc)
				return nil
			}
		}
	}

	return m
}
**/
