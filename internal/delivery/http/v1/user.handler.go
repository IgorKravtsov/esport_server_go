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

		authenticated := users.Group("/", h.userIdentity)
		{
			authenticated.POST("/verify/:code", h.userVerify)
		}
	}
}

type userRegisterInput struct {
	Name     string `json:"name" binding:"required,min=1,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=1,max=64"`
}

// @Summary User SignUp
// @Tags users-auth
// @Description create user account
// @ModuleID register
// @Accept  json
// @Produce  json
// @Param input body userRegisterInput true "sign up info"
// @Success 201 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/register [post]
func (h *Handler) register(c *gin.Context) {
	var inp userRegisterInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "Invalid input body")
		return
	}

	if err := h.services.Users.Register(c.Request.Context(), dto.UserRegister{
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

// @Summary User Verify Registration
// @Security UsersAuth
// @Tags users-auth
// @Description user verify registration
// @ModuleID userVerify
// @Accept  json
// @Produce  json
// @Param code path string true "verification code"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/verify/{code} [post]
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

	if err = h.services.Users.Verify(c.Request.Context(), id, code); err != nil {
		if errors.Is(err, domain.ErrVerificationCodeInvalid) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{"success"})
}
