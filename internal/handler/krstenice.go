package handler

import (
	"fmt"
	"krstenica/internal/dto"
	"krstenica/internal/errorx"
	"krstenica/pkg"
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

		cx := ctx.Request.Context()

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

		cx := ctx.Request.Context()

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
		cx := ctx.Request.Context()

		filters := pkg.ParseUrlQuery(ctx)

		krstenica, totalCount, err := h.service.ListKrstenice(cx, filters)
		if err != nil {
			if err == errorx.ErrKrstenicaNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data":  krstenica,
			"total": totalCount,
		})
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

		cx := ctx.Request.Context()

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

		cx := ctx.Request.Context()

		err = h.service.DeleteKrstenica(cx, int64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}

//****************************************************end******Krstenica*************************************
