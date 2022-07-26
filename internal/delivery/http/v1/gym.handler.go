package v1

import (
	"github.com/IgorKravtsov/esport_server_go/internal/service/gym/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initGymsRoutes(api *gin.RouterGroup) {
	gyms := api.Group("/gym")
	{
		authintiticatedAdmin := gyms.Group("/", h.adminIdentity)
		{
			authintiticatedAdmin.POST("/create", h.createGym)
		}
	}
}

// @Summary Create Gym
// @Tags gym
// @Description creating a gym
// @Security AdminAuth
// @ModuleID createGym
// @Accept  json
// @Produce  json
// @Param input body dto.CreateGym true "gym info"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /gym/create [post]
func (h *Handler) createGym(c *gin.Context) {
	id, err := getAdminId(c)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	var inp dto.CreateGym
	if err = c.BindJSON(&inp); err != nil {
		newErrResponse(c, http.StatusBadRequest, "Invalid input body")

		return
	}

	ID, err := h.services.Gym.Create(c.Request.Context(), dto.CreateGym{
		Title:   inp.Title,
		Address: inp.Address,
	}, id)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusCreated, idResponse{
		ID: ID,
	})
}
