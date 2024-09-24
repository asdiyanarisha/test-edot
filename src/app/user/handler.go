package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"test-edot/src/dto"
	"test-edot/src/factory"
)

type handler struct {
	service Service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

func (h *handler) UserMe(g *gin.Context) {
	userId, err := strconv.Atoi(g.Value("userId").(string))
	if err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	fmt.Println(userId)

	g.JSON(http.StatusCreated, dto.Response{
		Message: "user fetched",
	})
	return
}

func (h *handler) LoginUser(g *gin.Context) {
	var payload dto.LoginUser
	if err := g.ShouldBindJSON(&payload); err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	token, err := h.service.Login(g, payload)
	if err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusCreated, dto.Response{
		Message: "login successfully",
		Data:    dto.ResponseToken{Token: token},
	})
	return
}

func (h *handler) RegisterUser(g *gin.Context) {
	var payload dto.RegisterUser
	if err := g.ShouldBindJSON(&payload); err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.service.Register(g, payload); err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusCreated, dto.Response{
		Message: "user success created",
	})
	return
}
