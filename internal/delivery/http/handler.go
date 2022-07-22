package http

import (
	v1 "github.com/IgorKravtsov/esport_server_go/internal/delivery/http/v1"
	"net/http"

	"github.com/IgorKravtsov/esport_server_go/internal/config"
	"github.com/IgorKravtsov/esport_server_go/internal/service"
	"github.com/IgorKravtsov/esport_server_go/pkg/auth"
	"github.com/IgorKravtsov/esport_server_go/pkg/limiter"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	// Init gin handler
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		limiter.Limit(cfg.Limiter.RPS, cfg.Limiter.Burst, cfg.Limiter.TTL),
		corsMiddleware,
	)

	//docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	//if cfg.Environment != configs.EnvLocal {
	//  docs.SwaggerInfo.Host = cfg.HTTP.Host
	//}

	//if cfg.Environment != configs.Prod {
	//  router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//}

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
