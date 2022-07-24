package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"

	"github.com/IgorKravtsov/esport_server_go/internal/domain"
	"github.com/IgorKravtsov/esport_server_go/internal/service/user/dto"
)

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	users := api.Group("/auth")
	{
		users.POST("/register", h.register)
		users.POST("/login", h.login)
	}
}

// @Summary User Register
// @Tags auth
// @Description create user account
// @ModuleID register
// @Accept  json
// @Produce  json
// @Param input body dto.UserRegister true "register info"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/register [post]
func (h *Handler) register(c *gin.Context) {
	var inp dto.UserRegister
	if err := c.BindJSON(&inp); err != nil {
		newErrResponse(c, http.StatusBadRequest, "Invalid input body")
		return
	}
	foundUser, _ := h.services.User.GetByEmail(c.Request.Context(), inp.Email)

	if foundUser != nil {
		newErrResponse(c, http.StatusBadRequest, "User with this email already exists")
		return
	}

	ID, err := h.services.User.Register(c.Request.Context(), dto.UserRegister{
		Name:     inp.Name,
		Email:    strings.ToLower(inp.Email),
		Password: inp.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			newErrResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newErrResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, idResponse{
		ID: ID,
	})
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
// @Router /auth/login [post]
func (h *Handler) login(c *gin.Context) {
	var inp dto.UserLogin
	if err := c.BindJSON(&inp); err != nil {
		newErrResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.User.Login(c.Request.Context(), dto.UserLogin{
		Email:    inp.Email,
		Password: inp.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newErrResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newErrResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}
