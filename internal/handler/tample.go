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

// *************************************************************Tamples*************************************
func (h *httpHandler) createTample() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &dto.TampleCreateReq{}

		if err := ctx.Bind(req); err != nil {
			fmt.Println("Error when parsing body", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "error when parsing request data"})
			return
		}

		cx := context.Background()

		tample, err := h.service.CreateTample(cx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, tample)
	}
}

func (h *httpHandler) getTample() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cx := context.Background()

		tample, err := h.service.GetTampleByID(cx, int64(id))
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

func (h *httpHandler) listTample() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		fmt.Println("Path ", path)
		fmt.Println("Query ", query)
		fmt.Printf("Path: %s, Query: %s\n", path, query)

		cx := context.Background()
		tample, err := h.service.ListTamples(cx)
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

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req := &dto.TampleUpdateReq{}

		if err := ctx.Bind(req); err != nil {
			fmt.Println("Error when parsing body", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "error when parsing request data"})
			return
		}

		cx := context.Background()

		tample, err := h.service.UpdateTample(cx, int64(id), req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, tample)
	}
}

func (h *httpHandler) deleteTample() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = h.service.DeleteTample(ctx, int64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}

//****************************************************end******Tamples*************************************

//*************************************************************Priests*************************************

func (h *httpHandler) createPriest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &dto.PriestCreateReq{}

		if err := ctx.Bind(req); err != nil {
			fmt.Println("Error when parsing body", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "error when parsing request data"})
			return
		}

		cx := context.Background()

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

		cx := context.Background()

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

		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		fmt.Println("Path ", path)
		fmt.Println("Query ", query)
		fmt.Printf("Path: %s, Query: %s\n", path, query)

		cx := context.Background()
		priest, err := h.service.ListPriests(cx)
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

		cx := context.Background()

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

		err = h.service.DeletePriest(ctx, int64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}
