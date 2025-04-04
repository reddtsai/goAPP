package http

import (
	"context"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HttpOption func(*HttpOptions)

type HttpOptions struct {
	Middlewares  []gin.HandlerFunc
	Logging      bool
	CustomLogger gin.HandlerFunc
	CORSConfig   *cors.Config
}

func DefaultOptions() *HttpOptions {
	return &HttpOptions{
		Logging: true,
		CORSConfig: &cors.Config{
			AllowOrigins:     []string{"*"}, // 默認允許所有來源
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
	}
}

func WithMiddleware(middleware gin.HandlerFunc) HttpOption {
	return func(o *HttpOptions) {
		o.Middlewares = append(o.Middlewares, middleware)
	}
}

func WithLogging(enabled bool) HttpOption {
	return func(o *HttpOptions) {
		o.Logging = enabled
	}
}

func WithCustomLogger(logger gin.HandlerFunc) HttpOption {
	return func(o *HttpOptions) {
		o.Logging = true
		o.CustomLogger = logger
	}
}

func WithCORS(config cors.Config) HttpOption {
	return func(o *HttpOptions) {
		o.CORSConfig = &config
	}
}

func NewRouter(ctx context.Context, opts ...HttpOption) *gin.Engine {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	g := gin.New()
	g.Use(gin.Recovery())

	if options.Logging {
		if options.CustomLogger != nil {
			g.Use(options.CustomLogger)
		} else {
			g.Use(gin.Logger())
		}
	}

	if options.CORSConfig != nil {
		g.Use(cors.New(*options.CORSConfig))
	}

	for _, middleware := range options.Middlewares {
		g.Use(middleware)
	}

	return g
}
