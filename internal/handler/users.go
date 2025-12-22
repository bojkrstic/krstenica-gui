package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"krstenica/internal/dto"
)

type usersTableData struct {
	Items   []*dto.User
	Error   string
	Success string
}

func (h *httpHandler) renderUsersPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := h.service.ListUsers(ctx.Request.Context())
		if err != nil {
			h.renderHTML(ctx, http.StatusInternalServerError, "partials/error.html", gin.H{"Message": err.Error()})
			return
		}
		h.renderHTML(ctx, http.StatusOK, "users/index.html", gin.H{
			"Title":           "Корисници",
			"ContentTemplate": "users/content",
			"Users":           users,
		})
	}
}

func (h *httpHandler) renderUsersTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.usersTableResponse(ctx, "", "")
	}
}

func (h *httpHandler) renderUsersNew() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.renderHTML(ctx, http.StatusOK, "users/new.html", nil)
	}
}

func (h *httpHandler) handleUsersCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.UserCreateReq
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.Header("HX-Retarget", "closest dialog")
			h.renderHTML(ctx, http.StatusBadRequest, "users/new.html", gin.H{"Error": "Неисправан унос"})
			return
		}
		created, err := h.service.CreateUser(ctx.Request.Context(), &req)
		if err != nil {
			ctx.Header("HX-Retarget", "closest dialog")
			h.renderHTML(ctx, http.StatusBadRequest, "users/new.html", gin.H{"Error": err.Error()})
			return
		}
		h.usersTableResponse(ctx, "Корисник '"+created.Username+"' је додат.", "")
	}
}

func (h *httpHandler) renderUsersEdit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			h.renderHTML(ctx, http.StatusBadRequest, "partials/error.html", gin.H{"Message": "Непознат корисник"})
			return
		}
		user, err := h.service.GetUser(ctx.Request.Context(), id)
		if err != nil {
			h.renderHTML(ctx, http.StatusInternalServerError, "partials/error.html", gin.H{"Message": err.Error()})
			return
		}
		h.renderHTML(ctx, http.StatusOK, "users/edit.html", gin.H{"User": user})
	}
}

func (h *httpHandler) handleUsersUpdate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			h.renderHTML(ctx, http.StatusBadRequest, "partials/error.html", gin.H{"Message": "Непознат корисник"})
			return
		}
		existing, err := h.service.GetUser(ctx.Request.Context(), id)
		if err != nil {
			ctx.Header("HX-Retarget", "closest dialog")
			h.renderHTML(ctx, http.StatusInternalServerError, "users/edit.html", gin.H{"Error": err.Error()})
			return
		}

		var req dto.UserUpdateReq
		if err := ctx.ShouldBind(&req); err != nil {
			ctx.Header("HX-Retarget", "closest dialog")
			h.renderHTML(ctx, http.StatusBadRequest, "users/edit.html", gin.H{"Error": "Неисправан унос", "User": existing})
			return
		}
		if strings.TrimSpace(req.Username) == "" {
			ctx.Header("HX-Retarget", "closest dialog")
			h.renderHTML(ctx, http.StatusBadRequest, "users/edit.html", gin.H{"Error": "Корисничко име је обавезно", "User": existing})
			return
		}

		updated, err := h.service.UpdateUser(ctx.Request.Context(), id, &req)
		if err != nil {
			ctx.Header("HX-Retarget", "closest dialog")
			h.renderHTML(ctx, http.StatusBadRequest, "users/edit.html", gin.H{"Error": err.Error(), "User": existing})
			return
		}
		h.usersTableResponse(ctx, "Корисник '"+updated.Username+"' је измењен.", "")
	}
}

func (h *httpHandler) handleUsersDelete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			h.renderHTML(ctx, http.StatusBadRequest, "partials/error.html", gin.H{"Message": "Непознат корисник"})
			return
		}
		user, err := h.service.GetUser(ctx.Request.Context(), id)
		if err != nil {
			h.usersTableResponse(ctx, "", err.Error())
			return
		}
		if err := h.service.DeleteUser(ctx.Request.Context(), id); err != nil {
			h.usersTableResponse(ctx, "", err.Error())
			return
		}
		h.usersTableResponse(ctx, "Корисник '"+user.Username+"' је обрисан.", "")
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
		var payload dto.UserUpdateReq
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		user, err := h.service.UpdateUser(ctx.Request.Context(), id, &payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}

func (h *httpHandler) deleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if err := h.service.DeleteUser(ctx.Request.Context(), id); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) usersTableResponse(ctx *gin.Context, successMsg, errorMsg string) {
	users, err := h.service.ListUsers(ctx.Request.Context())
	if err != nil {
		h.renderHTML(ctx, http.StatusInternalServerError, "partials/error.html", gin.H{"Message": err.Error()})
		return
	}
	h.renderHTML(ctx, http.StatusOK, "users/table.html", usersTableData{
		Items:   users,
		Success: successMsg,
		Error:   errorMsg,
	})
}
