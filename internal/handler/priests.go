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

//*************************************************************Priests*************************************

func (h *httpHandler) createPriest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &dto.PriestCreateReq{}

		if err := ctx.Bind(req); err != nil {
			fmt.Println("Error when parsing body", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "error when parsing request data"})
			return
		}

		cx := ctx.Request.Context()

		priest, err := h.service.CreatePriest(cx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, priest)
	}
}

func (h *httpHandler) getPriest() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cx := ctx.Request.Context()

		priest, err := h.service.GetPriestByID(cx, int64(id))
		if err != nil {
			if err == errorx.ErrPriestNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, priest)
	}
}

func (h *httpHandler) listPriest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cx := ctx.Request.Context()

		filters := pkg.ParseUrlQuery(ctx)

		priest, totalCount, err := h.service.ListPriests(cx, filters)
		if err != nil {
			if err == errorx.ErrPriestNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data":  priest,
			"total": totalCount,
		})
	}
}

func (h *httpHandler) updatePriest() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req := &dto.PriestUpdateReq{}

		if err := ctx.Bind(req); err != nil {
			fmt.Println("Error when parsing body", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "error when parsing request data"})
			return
		}

		cx := ctx.Request.Context()

		priest, err := h.service.UpdatePriest(cx, int64(id), req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, priest)
	}
}

func (h *httpHandler) deletePriest() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cx := ctx.Request.Context()

		err = h.service.DeletePriest(cx, int64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}
