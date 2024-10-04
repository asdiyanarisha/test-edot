package product

import (
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
		Message: "success fetch product list",
		Data:    res,
	})
	return
}

func (h *handler) DetailProduct(g *gin.Context) {
	userClaim := g.Value("userClaim").(dto.UserClaimJwt)

	productId, err := strconv.Atoi(g.Param("product_id"))
	if err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: "product_id is not valid",
		})
		return
	}

	res, err := h.service.GetProductDetail(g, userClaim, productId)
	if err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, dto.Response{
		Message: "success fetch detail product",
		Data:    res,
	})
}

func (h *handler) TransferProduct(g *gin.Context) {
	userClaim := g.Value("userClaim").(dto.UserClaimJwt)

	productId, err := strconv.Atoi(g.Param("product_id"))
	if err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: "product_id is not valid",
		})
		return
	}

	var payload dto.TransferProductWarehouse
	if err := g.ShouldBindJSON(&payload); err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.service.TransferProductWarehouse(g, payload, userClaim, productId); err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, dto.Response{
		Message: "success transfer product",
	})
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
