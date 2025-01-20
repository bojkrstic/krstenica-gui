package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const routePrefix = "api/v1"

func (h *httpHandler) addRoutes() {
	adminRouter := h.router.Group("", h.needAdminAccess()) // checking admin prililges at this point, applied for all routes in admin router

	//tamples
	adminRouter.POST(pathWithAction("adminv2", "tamples"), h.createTample())
	adminRouter.GET(pathWithAction("adminv2", "tamples/:id"), h.getTample())
	adminRouter.GET(pathWithAction("adminv2", "tamples"), h.listTample())
	adminRouter.PUT(pathWithAction("adminv2", "tamples/:id"), h.updateTample())
	adminRouter.DELETE(pathWithAction("adminv2", "tamples/:id"), h.deleteTample())
	//priests
	adminRouter.POST(pathWithAction("adminv2", "priests"), h.createPriest())
	adminRouter.GET(pathWithAction("adminv2", "priests/:id"), h.getPriest())
	adminRouter.GET(pathWithAction("adminv2", "priests"), h.listPriest())
	adminRouter.PUT(pathWithAction("adminv2", "priests/:id"), h.updatePriest())
	adminRouter.DELETE(pathWithAction("adminv2", "priests/:id"), h.deletePriest())
	//eparhije
	adminRouter.POST(pathWithAction("adminv2", "eparhije"), h.createEparhije())
	adminRouter.GET(pathWithAction("adminv2", "eparhije/:id"), h.getEparhije())
	adminRouter.GET(pathWithAction("adminv2", "eparhije"), h.listEparhije())
	adminRouter.PUT(pathWithAction("adminv2", "eparhije/:id"), h.updateEparhije())
	adminRouter.DELETE(pathWithAction("adminv2", "eparhije/:id"), h.deleteEparhije())

	//persons
	adminRouter.POST(pathWithAction("adminv2", "persons"), h.createPersons())
	adminRouter.GET(pathWithAction("adminv2", "persons/:id"), h.getPersons())
	adminRouter.GET(pathWithAction("adminv2", "persons"), h.listPersons())
	adminRouter.PUT(pathWithAction("adminv2", "persons/:id"), h.updatePersons())
	adminRouter.DELETE(pathWithAction("adminv2", "persons/:id"), h.deletePersons())

	//krstenice
	adminRouter.POST(pathWithAction("adminv2", "krstenice"), h.createKrstenice())
	adminRouter.GET(pathWithAction("adminv2", "krstenice/:id"), h.getKrstenice())
	adminRouter.GET(pathWithAction("adminv2", "krstenice"), h.listKrstenice())
	adminRouter.PUT(pathWithAction("adminv2", "krstenice/:id"), h.updateKrstenice())
	adminRouter.DELETE(pathWithAction("adminv2", "krstenice/:id"), h.deleteKrstenice())
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
