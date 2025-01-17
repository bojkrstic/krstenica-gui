package handler

import (
	"context"
	"fmt"
	"krstenica/internal/dto"
	"krstenica/internal/errorx"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// *************************************************************Eparhije*************************************
func (h *httpHandler) createEparhije() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &dto.EparhijeCreateReq{}

		if err := ctx.Bind(req); err != nil {
			fmt.Println("Error when parsing body", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "error when parsing request data"})
			return
		}

		cx := context.Background()

		eparhija, err := h.service.CreateEparhije(cx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, eparhija)
	}
}

func (h *httpHandler) getEparhije() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cx := context.Background()

		eparhija, err := h.service.GetEparhijeByID(cx, int64(id))
		if err != nil {
			if err == errorx.ErrEparhijeNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, eparhija)
	}
}

func (h *httpHandler) listEparhije() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		fmt.Println("Path ", path)
		fmt.Println("Query ", query)
		fmt.Printf("Path: %s, Query: %s\n", path, query)

		cx := context.Background()
		eparhija, err := h.service.ListEparhije(cx)
		if err != nil {
			if err == errorx.ErrEparhijeNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, eparhija)
	}
}

func (h *httpHandler) updateEparhije() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req := &dto.EparhijeUpdateReq{}

		if err := ctx.Bind(req); err != nil {
			fmt.Println("Error when parsing body", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "error when parsing request data"})
			return
		}

		cx := context.Background()

		eparhija, err := h.service.UpdateEparhije(cx, int64(id), req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, eparhija)
	}
}

func (h *httpHandler) deleteEparhije() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = h.service.DeleteEparhije(ctx, int64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}

//****************************************************end******Eparhije*************************************
