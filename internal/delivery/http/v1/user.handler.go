package v1

import (
	"errors"
	"github.com/IgorKravtsov/esport_server_go/internal/domain"
	"github.com/IgorKravtsov/esport_server_go/internal/service/user/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/user")
	{
		users.POST("/register", h.register)
		users.POST("/login", h.login)
		users.POST("/refresh-tokens", h.userRefresh)

		authenticated := users.Group("/", h.userIdentity)
		{
			authenticated.POST("/verify/:code", h.userVerify)
		}
	}
}

// @Summary User Register
// @Tags auth
// @Description create user account
// @ModuleID register
// @Accept  json
// @Produce  json
// @Param input body dto.UserRegister true "register info"
// @Success 201 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /api/v1/user/register [post]
func (h *Handler) register(c *gin.Context) {
	var inp dto.UserRegister
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "Invalid input body")
		return
	}
	foundUser, _ := h.services.User.GetByEmail(c.Request.Context(), inp.Email)

	if foundUser != nil {
		newResponse(c, http.StatusBadRequest, "User with this email already exists")
		return
	}

	if err := h.services.User.Register(c.Request.Context(), dto.UserRegister{
		Name:     inp.Name,
		Email:    inp.Email,
		Password: inp.Password,
	}); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary Learner Login
// @Tags auth
// @Description user sign in
// @ModuleID login
// @Accept  json
// @Produce  json
// @Param input body dto.UserLogin true "login info"
// @Success 200 {object} dto.TokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /api/v1/user/login [post]
func (h *Handler) login(c *gin.Context) {
	var inp dto.UserLogin
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.User.Login(c.Request.Context(), dto.UserLogin{
		Email:    inp.Email,
		Password: inp.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

// @Summary User Verify Registration
// @Security UserAuth
// @Tags auth
// @Description user verify registration
// @ModuleID userVerify
// @Accept  json
// @Produce  json
// @Param code path string true "verification code"
// @Success 200 {object} dto.TokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /api/v1/user/verify/{code} [post]
func (h *Handler) userVerify(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		newResponse(c, http.StatusBadRequest, "code is empty")

		return
	}

	id, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err = h.services.User.Verify(c.Request.Context(), id, code); err != nil {
		if errors.Is(err, domain.ErrVerificationCodeInvalid) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{"success"})
}

// @Summary User Refresh Tokens
// @Tags auth
// @Description user refresh tokens
// @Accept  json
// @Produce  json
// @Param input body dto.RefreshToken true "register info"
// @Success 200 {object} dto.TokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /api/v1/user/refresh-tokens [post]
func (h *Handler) userRefresh(c *gin.Context) {
	var inp dto.RefreshToken
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "Invalid input body")

		return
	}

	res, err := h.services.User.RefreshTokens(c.Request.Context(), inp.Token)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}
