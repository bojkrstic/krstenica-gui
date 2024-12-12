package handler

import (
	"krstenica/internal/config"
	"krstenica/internal/repository"
	"krstenica/internal/service"

	"github.com/gin-gonic/gin"
)

type HttpHandler interface {
	Init()
}

type httpHandler struct {
	conf    *config.Config
	router  *gin.Engine
	repo    repository.Repo
	service service.Service
}

func NewHttpHandler(s service.Service, c *config.Config, r repository.Repo) HttpHandler {
	return &httpHandler{service: s, conf: c, repo: r}
}

func (h *httpHandler) Init() {
	h.router = gin.New()
	h.router.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/api/v1/krstenica/ping"))

	h.addRoutes()

	// Spawing go-routine when starting server.
	go func() {
		h.router.Run(h.conf.HTTPPort)
	}()
}
