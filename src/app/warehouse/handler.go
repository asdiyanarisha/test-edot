package warehouse

import (
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

func (h *handler) AddWarehouse(g *gin.Context) {
	userClaim := g.Value("userClaim").(dto.UserClaimJwt)

	var payload dto.PayloadAddWarehouse
	if err := g.ShouldBindJSON(&payload); err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.service.AddWarehouse(g, payload, userClaim); err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusCreated, dto.Response{
		Message: "success create warehouse",
	})
	return
}

func (h *handler) GetWarehouses(g *gin.Context) {
	userClaim := g.Value("userClaim").(dto.UserClaimJwt)

	var payload dto.ParameterQueryWarehouse
	if err := g.ShouldBind(&payload); err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	res, err := h.service.GetWarehouses(g, userClaim, payload)
	if err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusCreated, dto.Response{
		Message: "success list warehouse",
		Data:    res,
	})
	return
}
