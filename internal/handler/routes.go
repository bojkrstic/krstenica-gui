package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const routePrefix = "api/v1"

func (h *httpHandler) addRoutes() {
	adminRouter := h.router.Group("", h.needAdminAccess())

	adminRouter.POST(pathWithAction("adminv2", "birth-certificates"), h.needAdminAccess(), h.createBirthCertificate())
	adminRouter.GET(pathWithAction("adminv2", "birth-certificates/:id"), h.needAdminAccess(), h.getBirthCertificate())
	adminRouter.PUT(pathWithAction("adminv2", "birth-certificates/:id"), h.needAdminAccess(), h.updateBirthCertificate())
	adminRouter.DELETE(pathWithAction("adminv2", "birth-certificates/:id"), h.needAdminAccess(), h.deleteBirthCertificate())

	//tamples
	adminRouter.POST(pathWithAction("adminv2", "tamples"), h.needAdminAccess(), h.createTample())
	adminRouter.GET(pathWithAction("adminv2", "tamples/:id"), h.needAdminAccess(), h.getTample())
	adminRouter.GET(pathWithAction("adminv2", "tamples"), h.needAdminAccess(), h.listTample())
	adminRouter.PUT(pathWithAction("adminv2", "tamples/:id"), h.needAdminAccess(), h.updateTample())
	adminRouter.DELETE(pathWithAction("adminv2", "tamples/:id"), h.needAdminAccess(), h.deleteTample())
	//priests
	adminRouter.POST(pathWithAction("adminv2", "priests"), h.needAdminAccess(), h.createPriest())
	adminRouter.GET(pathWithAction("adminv2", "priests/:id"), h.needAdminAccess(), h.getPriest())
	adminRouter.GET(pathWithAction("adminv2", "priests"), h.needAdminAccess(), h.listPriest())
	adminRouter.PUT(pathWithAction("adminv2", "priests/:id"), h.needAdminAccess(), h.updatePriest())
	adminRouter.DELETE(pathWithAction("adminv2", "priests/:id"), h.needAdminAccess(), h.deletePriest())
	//eparhije
	adminRouter.POST(pathWithAction("adminv2", "eparhije"), h.needAdminAccess(), h.createEparhije())
	adminRouter.GET(pathWithAction("adminv2", "eparhije/:id"), h.needAdminAccess(), h.getEparhije())
	adminRouter.GET(pathWithAction("adminv2", "eparhije"), h.needAdminAccess(), h.listEparhije())
	adminRouter.PUT(pathWithAction("adminv2", "eparhije/:id"), h.needAdminAccess(), h.updateEparhije())
	adminRouter.DELETE(pathWithAction("adminv2", "eparhije/:id"), h.needAdminAccess(), h.deleteEparhije())

	//persons
	adminRouter.POST(pathWithAction("adminv2", "persons"), h.needAdminAccess(), h.createPersons())
	adminRouter.GET(pathWithAction("adminv2", "persons/:id"), h.needAdminAccess(), h.getPersons())
	adminRouter.GET(pathWithAction("adminv2", "persons"), h.needAdminAccess(), h.listPersons())
	adminRouter.PUT(pathWithAction("adminv2", "persons/:id"), h.needAdminAccess(), h.updatePersons())
	adminRouter.DELETE(pathWithAction("adminv2", "persons/:id"), h.needAdminAccess(), h.deletePersons())
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
