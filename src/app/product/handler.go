package product

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

func (h *handler) ProductList(g *gin.Context) {
	var payload dto.ParameterQuery
	if err := g.ShouldBind(&payload); err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	res, err := h.service.ProductList(g, payload)
	if err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, dto.Response{
		Message: "success create shop",
		Data:    res,
	})
	return
}

func (h *handler) AddProduct(g *gin.Context) {
	userClaim := g.Value("userClaim").(dto.UserClaimJwt)
	var payload dto.PayloadAddProduct
	if err := g.ShouldBindJSON(&payload); err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.service.AddProduct(g, payload, userClaim); err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusCreated, dto.Response{
		Message: "success create product",
	})
	return
}
