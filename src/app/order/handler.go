package order

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

func (h *handler) PaymentOrder(g *gin.Context) {
	userClaim := g.Value("userClaim").(dto.UserClaimJwt)
	orderId, err := strconv.Atoi(g.Param("order_id"))
	if err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: "product_id is not valid",
		})
		return
	}

	if err := h.service.PaymentOrder(g, userClaim, orderId); err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, dto.Response{
		Message: "order successfully payed",
	})
	return
}

func (h *handler) CreateOrder(g *gin.Context) {
	userClaim := g.Value("userClaim").(dto.UserClaimJwt)
	var payload dto.PayloadCreateOrder
	if err := g.ShouldBindJSON(&payload); err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.service.CreateOrder(g, userClaim, payload); err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusCreated, dto.Response{
		Message: "create order success",
	})
	return
}
