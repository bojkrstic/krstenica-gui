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

// *************************************************************Krstenica*************************************
func (h *httpHandler) createKrstenice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &dto.KrstenicaCreateReq{}

		if err := ctx.Bind(req); err != nil {
			fmt.Println("Error when parsing body", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "error when parsing request data"})
			return
		}

		cx := context.Background()

		krstenica, err := h.service.CreateKrstenica(cx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, krstenica)
	}
}

func (h *httpHandler) getKrstenice() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cx := context.Background()

		krstenica, err := h.service.GetKrstenicaByID(cx, int64(id))
		if err != nil {
			if err == errorx.ErrKrstenicaNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, krstenica)
	}
}

func (h *httpHandler) listKrstenice() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		fmt.Println("Path ", path)
		fmt.Println("Query ", query)
		fmt.Printf("Path: %s, Query: %s\n", path, query)

		cx := context.Background()
		krstenica, err := h.service.ListKrstenice(cx)
		if err != nil {
			if err == errorx.ErrKrstenicaNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, krstenica)
	}
}

func (h *httpHandler) updateKrstenice() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req := &dto.KrstenicaUpdateReq{}

		if err := ctx.Bind(req); err != nil {
			fmt.Println("Error when parsing body", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "error when parsing request data"})
			return
		}

		cx := context.Background()

		krstenica, err := h.service.UpdateKrstenica(cx, int64(id), req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, krstenica)
	}
}

func (h *httpHandler) deleteKrstenice() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = h.service.DeleteKrstenica(ctx, int64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}

//****************************************************end******Krstenica*************************************
