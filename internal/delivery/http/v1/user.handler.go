package v1

import (
	"errors"
	"github.com/IgorKravtsov/esport_server_go/internal/domain"
	"github.com/IgorKravtsov/esport_server_go/internal/service/user/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/user")
	{
		authenticated := users.Group("/", h.authIdentity)
		users.POST("/refresh-tokens", h.userRefresh)
		{
			authenticated.POST("/verify/:code", h.userVerify)
		}
	}
}

// @Summary User Verify Registration
// @Security UserAuth
// @Tags user
// @Description user verify registration
// @ModuleID userVerify
// @Accept  json
// @Produce  json
// @Param code path string true "verification code"
// @Success 200 {object} dto.TokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /user/verify/{code} [post]
func (h *Handler) userVerify(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		newErrResponse(c, http.StatusBadRequest, "code is empty")

		return
	}

	id, err := getUserId(c)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err = h.services.User.Verify(c.Request.Context(), id, code); err != nil {
		if errors.Is(err, domain.ErrVerificationCodeInvalid) {
			newErrResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newErrResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{"success"})
}

// @Summary User Refresh Tokens
// @Tags user
// @Description user refresh tokens
// @Accept  json
// @Produce  json
// @Param input body dto.RefreshToken true "register info"
// @Success 200 {object} dto.TokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /user/refresh-tokens [post]
func (h *Handler) userRefresh(c *gin.Context) {
	var inp dto.RefreshToken
	if err := c.BindJSON(&inp); err != nil {
		newErrResponse(c, http.StatusBadRequest, "Invalid input body")

		return
	}

	res, err := h.services.User.RefreshTokens(c.Request.Context(), inp.Token)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}
