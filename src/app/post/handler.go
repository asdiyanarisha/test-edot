package post

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test-asset-fendr/src/dto"
	"test-asset-fendr/src/factory"
)

type handler struct {
	service Service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

func (h *handler) AddPostHandler(g *gin.Context) {
	var payload dto.AddPost
	if err := g.ShouldBindJSON(&payload); err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.service.AddPostService(g, payload); err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusCreated, dto.Response{
		Message: "post has been created",
	})
	return
}
