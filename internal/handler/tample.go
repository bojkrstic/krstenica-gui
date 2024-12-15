package handler

import (
	"context"
	"krstenica/internal/errorx"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *httpHandler) createTample() gin.HandlerFunc {
	return func(ctx *gin.Context) {
	}
}

func (h *httpHandler) getTample() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tample, err := h.service.GetTampleByID(context.Background(), int64(id))
		if err != nil {
			if err == errorx.ErrTampleNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, tample)
	}
}

func (h *httpHandler) updateTample() gin.HandlerFunc {
	return func(ctx *gin.Context) {
	}
}

func (h *httpHandler) deleteTample() gin.HandlerFunc {
	return func(ctx *gin.Context) {
	}
}
