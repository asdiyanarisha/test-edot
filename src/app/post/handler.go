package post

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func (h *handler) GetPostHandler(g *gin.Context) {
	var payload dto.ParamGetPost
	if err := g.ShouldBindQuery(&payload); err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	res, err := h.service.GetPostService(g, payload)
	if err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, dto.Response{
		Message: "post has been resulted",
		Data:    res,
	})
	return
}

func (h *handler) GetPostByIdHandler(g *gin.Context) {
	postIdParam := g.Param("id")

	postId, err := strconv.Atoi(postIdParam)
	if err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	res, err := h.service.GetPostByIdService(g, postId)
	if err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, dto.Response{
		Message: "post has been resulted",
		Data:    res,
	})
	return
}

func (h *handler) DeletePostHandler(g *gin.Context) {
	postIdParam := g.Param("id")

	postId, err := strconv.Atoi(postIdParam)
	if err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.service.DeletePostService(g, postId); err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, dto.Response{
		Message: "post has been deleted",
	})
	return
}

func (h *handler) UpdatePostHandler(g *gin.Context) {
	postIdParam := g.Param("id")

	var payload dto.AddPost
	if err := g.ShouldBindJSON(&payload); err != nil {
		g.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	postId, err := strconv.Atoi(postIdParam)
	if err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := h.service.UpdatePostService(g, postId, payload); err != nil {
		g.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, dto.Response{
		Message: "post has been updated",
	})
	return
}
