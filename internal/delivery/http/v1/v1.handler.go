package v1

import (
	"errors"
	"github.com/IgorKravtsov/esport_server_go/internal/service"
	"github.com/IgorKravtsov/esport_server_go/pkg/auth"
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

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUserRoutes(v1)
		h.initGymsRoutes(v1)
		h.initAuthRoutes(v1)

		//v1.GET("/settings", h.setSchoolFromRequest, h.getSchoolSettings)
		//v1.GET("/promocodes/:code", h.setSchoolFromRequest, h.getPromo)
		//v1.GET("/offers/:id", h.setSchoolFromRequest, h.getOffer)
	}
}

func parseIdFromPath(c *gin.Context, param string) (string, error) {
	idParam := c.Param(param)
	if idParam == "" {
		return "", errors.New("empty id param")
	}

	return idParam, nil
}
