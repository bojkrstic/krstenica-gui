package handler

import "fmt"

const routePrefix = "api/v1"

func (h *httpHandler) addRoutes() {
	apiRouter := h.router.Group("", h.requireAPIAuth())
	adminRouter := apiRouter.Group("", h.requireRole(adminRoleDefault))

	// admin-only resources
	adminRouter.POST(pathWithAction("adminv2", "tamples"), h.createTample())
	adminRouter.GET(pathWithAction("adminv2", "tamples/:id"), h.getTample())
	adminRouter.GET(pathWithAction("adminv2", "tamples"), h.listTample())
	adminRouter.PUT(pathWithAction("adminv2", "tamples/:id"), h.updateTample())
	adminRouter.DELETE(pathWithAction("adminv2", "tamples/:id"), h.deleteTample())

	adminRouter.POST(pathWithAction("adminv2", "priests"), h.createPriest())
	adminRouter.GET(pathWithAction("adminv2", "priests/:id"), h.getPriest())
	adminRouter.GET(pathWithAction("adminv2", "priests"), h.listPriest())
	adminRouter.PUT(pathWithAction("adminv2", "priests/:id"), h.updatePriest())
	adminRouter.DELETE(pathWithAction("adminv2", "priests/:id"), h.deletePriest())

	adminRouter.POST(pathWithAction("adminv2", "eparhije"), h.createEparhije())
	adminRouter.GET(pathWithAction("adminv2", "eparhije/:id"), h.getEparhije())
	adminRouter.GET(pathWithAction("adminv2", "eparhije"), h.listEparhije())
	adminRouter.PUT(pathWithAction("adminv2", "eparhije/:id"), h.updateEparhije())
	adminRouter.DELETE(pathWithAction("adminv2", "eparhije/:id"), h.deleteEparhije())

	adminRouter.POST(pathWithAction("adminv2", "persons"), h.createPersons())
	adminRouter.GET(pathWithAction("adminv2", "persons/:id"), h.getPersons())
	adminRouter.GET(pathWithAction("adminv2", "persons"), h.listPersons())
	adminRouter.PUT(pathWithAction("adminv2", "persons/:id"), h.updatePersons())
	adminRouter.DELETE(pathWithAction("adminv2", "persons/:id"), h.deletePersons())

	adminRouter.GET(pathWithAction("adminv2", "users"), h.listUsers())
	adminRouter.POST(pathWithAction("adminv2", "users"), h.createUser())
	adminRouter.PUT(pathWithAction("adminv2", "users/:id"), h.updateUser())
	adminRouter.DELETE(pathWithAction("adminv2", "users/:id"), h.deleteUser())

	// krstenice routes available to any authenticated user (service enforces city/role)
	apiRouter.POST(pathWithAction("adminv2", "krstenice"), h.createKrstenice())
	apiRouter.GET(pathWithAction("adminv2", "krstenice/:id"), h.getKrstenice())
	apiRouter.GET(pathWithAction("adminv2", "krstenice"), h.listKrstenice())
	apiRouter.PUT(pathWithAction("adminv2", "krstenice/:id"), h.updateKrstenice())
	apiRouter.DELETE(pathWithAction("adminv2", "krstenice/:id"), h.deleteKrstenice())
	apiRouter.GET(pathWithAction("adminv2", "krstenice-print/:id"), h.getKrstenicePrint())
}

func pathWithAction(module string, action string) string {
	return fmt.Sprintf("%s/%s/%s", routePrefix, module, action)
}
