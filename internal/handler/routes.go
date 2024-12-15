package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const routePrefix = "api/v1"

func (h *httpHandler) addRoutes() {
	adminRouter := h.router.Group("", h.needAdminAccess())

	adminRouter.POST(pathWithAction("adminv2", "krstenica"), h.needAdminAccess(), h.createKrstenica())
	adminRouter.GET(pathWithAction("adminv2", "krstenica/:id"), h.needAdminAccess(), h.getKrstenica())
	adminRouter.PUT(pathWithAction("adminv2", "krstenica/:id"), h.needAdminAccess(), h.updateKrstenica())
	adminRouter.DELETE(pathWithAction("adminv2", "krstenica/:id"), h.needAdminAccess(), h.deleteKrstenica())

	adminRouter.POST(pathWithAction("adminv2", "tample"), h.needAdminAccess(), h.createTample())
	adminRouter.GET(pathWithAction("adminv2", "tample/:id"), h.needAdminAccess(), h.getTample())
	adminRouter.PUT(pathWithAction("adminv2", "tample/:id"), h.needAdminAccess(), h.updateTample())
	adminRouter.DELETE(pathWithAction("adminv2", "tample/:id"), h.needAdminAccess(), h.deleteTample())
}

func (h *httpHandler) needAdminAccess() gin.HandlerFunc {
	//return middleware.VerifyToken("PFG-BO-MEMBER", h.conf.AdminJWTSecret)

	//will be added later
	return func(ctx *gin.Context) {

	}
}

func (h *httpHandler) needUserAccess() gin.HandlerFunc {
	//return middleware.VerifyToken(token.PFGCasinoRoleUser, h.conf.JWTSecret)

	//will be added later
	return func(ctx *gin.Context) {

	}
}

func pathWithAction(module string, action string) string {
	return fmt.Sprintf("%s/%s/%s", routePrefix, module, action)
}
