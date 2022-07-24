package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	userCtx             = "userId"
	adminCtx            = "adminId"
	trainerCtx          = "trainerId"
	authorizationHeader = "Authorization"
	domainCtx           = "domain"
)

func parseRequestHost(c *gin.Context) string {
	refererHeader := c.Request.Header.Get("Referer")
	refererParts := strings.Split(refererHeader, "/")

	// this logic is used to avoid crashes during integration testing
	if len(refererParts) < 3 {
		return c.Request.Host
	}

	hostParts := strings.Split(refererParts[2], ":")

	return hostParts[0]
}

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
}

func (h *Handler) authIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newErrResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, id)
}

func (h *Handler) trainerIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newErrResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(trainerCtx, id)
}

func (h *Handler) adminIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newErrResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(adminCtx, id)
}

func getUserId(c *gin.Context) (string, error) {
	return getIdByContext(c, userCtx)
}

func getTrainerId(c *gin.Context) (string, error) {
	return getIdByContext(c, trainerCtx)
}

func getAdminId(c *gin.Context) (string, error) {
	return getIdByContext(c, adminCtx)
}

func getIdByContext(c *gin.Context, context string) (string, error) {
	idFromCtx, ok := c.Get(context)
	if !ok {
		return "", errors.New("ctx in middleware not found")
	}

	id, ok := idFromCtx.(string)
	if !ok {
		return "", errors.New("ctx in middleware is of invalid type")
	}

	return id, nil
}

func getDomainFromContext(c *gin.Context) (string, error) {
	val, ex := c.Get(domainCtx)
	if !ex {
		return "", errors.New("domainCtx not found")
	}

	valStr, ok := val.(string)
	if !ok {
		return "", errors.New("domainCtx is of invalid type")
	}

	return valStr, nil
}
