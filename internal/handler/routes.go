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

	adminRouter.POST(pathWithAction("adminv2", "hram"), h.needAdminAccess(), h.createHram())
	adminRouter.GET(pathWithAction("adminv2", "hram/:id"), h.needAdminAccess(), h.getHram())
	adminRouter.PUT(pathWithAction("adminv2", "hram/:id"), h.needAdminAccess(), h.updateHram())
	adminRouter.DELETE(pathWithAction("adminv2", "hram/:id"), h.needAdminAccess(), h.deleteHram())
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
