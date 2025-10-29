package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"krstenica/internal/dto"
)

type usersTableData struct {
	Items []*dto.User
}

func (h *httpHandler) renderUsersPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := h.service.ListUsers(ctx.Request.Context())
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{"Message": err.Error()})
			return
		}
		ctx.HTML(http.StatusOK, "users/index.html", gin.H{
			"Title":           "Корисници",
			"ContentTemplate": "users/content",
			"Users":           users,
		})
	}
}

func (h *httpHandler) renderUsersTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := h.service.ListUsers(ctx.Request.Context())
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{"Message": err.Error()})
			return
		}
		ctx.HTML(http.StatusOK, "users/table.html", usersTableData{Items: users})
	}
}

func (h *httpHandler) renderUsersNew() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "users/new.html", nil)
	}
}

func (h *httpHandler) handleUsersCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.UserCreateReq
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.Header("HX-Retarget", "closest dialog")
			ctx.HTML(http.StatusBadRequest, "users/new.html", gin.H{"Error": "Неисправан унос"})
			return
		}
		_, err := h.service.CreateUser(ctx.Request.Context(), &req)
		if err != nil {
			ctx.Header("HX-Retarget", "closest dialog")
			ctx.HTML(http.StatusBadRequest, "users/new.html", gin.H{"Error": err.Error()})
			return
		}
		users, err := h.service.ListUsers(ctx.Request.Context())
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{"Message": err.Error()})
			return
		}
		ctx.HTML(http.StatusOK, "users/table.html", usersTableData{Items: users})
	}
}

func (h *httpHandler) renderUsersEdit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{"Message": "Непознат корисник"})
			return
		}
		user, err := h.service.GetUser(ctx.Request.Context(), id)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{"Message": err.Error()})
			return
		}
		ctx.HTML(http.StatusOK, "users/edit.html", gin.H{"User": user})
	}
}

func (h *httpHandler) handleUsersUpdate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{"Message": "Непознат корисник"})
			return
		}
		user, err := h.service.GetUser(ctx.Request.Context(), id)
		if err != nil {
			ctx.Header("HX-Retarget", "closest dialog")
			ctx.HTML(http.StatusInternalServerError, "users/edit.html", gin.H{"Error": err.Error()})
			return
		}

		password := ctx.PostForm("password")
		if strings.TrimSpace(password) == "" {
			ctx.Header("HX-Retarget", "closest dialog")
			ctx.HTML(http.StatusBadRequest, "users/edit.html", gin.H{"Error": "Лозинка је обавезна", "User": user})
			return
		}

		if _, err = h.service.UpdateUserPassword(ctx.Request.Context(), id, password); err != nil {
			ctx.Header("HX-Retarget", "closest dialog")
			ctx.HTML(http.StatusBadRequest, "users/edit.html", gin.H{"Error": err.Error(), "User": user})
			return
		}
		users, err := h.service.ListUsers(ctx.Request.Context())
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{"Message": err.Error()})
			return
		}
		ctx.HTML(http.StatusOK, "users/table.html", usersTableData{Items: users})
	}
}

func (h *httpHandler) listUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := h.service.ListUsers(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": users})
	}
}

func (h *httpHandler) createUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.UserCreateReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		user, err := h.service.CreateUser(ctx.Request.Context(), &req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, user)
	}
}

func (h *httpHandler) updateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var payload struct {
			Password string `json:"password"`
		}
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		user, err := h.service.UpdateUserPassword(ctx.Request.Context(), id, payload.Password)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}
