package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
	userClaim := g.Value("userClaim").(dto.UserClaimJwt)

	fmt.Println(userClaim)

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

	g.JSON(http.StatusOK, dto.Response{
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

	res, err := h.service.Register(g, payload)
	if err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusCreated, dto.Response{
		Message: "user success created",
		Data:    res,
	})
	return
}
