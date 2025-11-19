package handler

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

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

	staticDir := resolveDir("web/static")
	h.router.Static("/static", staticDir)
		h.router.SetFuncMap(template.FuncMap{
			"formatDate": func(t time.Time) string {
				formatted := formatSerbianDate(t)
				if formatted == "" {
					return "-"
				}
				return formatted
			},
		"int64Value": func(v *int64) string {
			if v == nil {
				return ""
			}
			return strconv.FormatInt(*v, 10)
		},
	})
	templateDir := resolveDir("web/templates")
	h.mustLoadTemplates(templateDir)

	h.addAuthRoutes()
	h.addRoutes()
	h.addGuiRoutes()

	if err := h.service.EnsureDefaultUser(context.Background()); err != nil {
		log.Fatalf("failed to ensure default user: %v", err)
	}

	// Spawing go-routine when starting server.
	go func() {
		h.router.Run(h.conf.HTTPPort)
	}()
}

func resolveDir(relative string) string {
	if path, ok := searchUpward(relative, true); ok {
		return path
	}
	return relative
}

func resolveFile(relative string) string {
	if path, ok := searchUpward(relative, false); ok {
		return path
	}
	return relative
}

func searchUpward(relative string, wantDir bool) (string, bool) {
	if wd, err := os.Getwd(); err == nil {
		if path, ok := searchFromBase(wd, relative, wantDir); ok {
			return path, true
		}
	}

	if execPath, err := os.Executable(); err == nil {
		base := filepath.Dir(execPath)
		if path, ok := searchFromBase(base, relative, wantDir); ok {
			return path, true
		}
	}

	return "", false
}

func searchFromBase(start, relative string, wantDir bool) (string, bool) {
	current := start
	for {
		candidate := filepath.Join(current, relative)
		if info, err := os.Stat(candidate); err == nil {
			if wantDir && info.IsDir() {
				return candidate, true
			}
			if !wantDir && !info.IsDir() {
				return candidate, true
			}
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", false
}

func (h *httpHandler) mustLoadTemplates(root string) {
	patterns := []string{
		"*.html",
		filepath.Join("*", "*.html"),
		filepath.Join("*", "*", "*.html"),
	}

	templates := make(map[string]struct{})
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(root, pattern))
		if err != nil {
			continue
		}
		for _, match := range matches {
			templates[match] = struct{}{}
		}
	}

	if len(templates) == 0 {
		panic(fmt.Sprintf("no templates found in %s", root))
	}

	paths := make([]string, 0, len(templates))
	for path := range templates {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	h.router.LoadHTMLFiles(paths...)
}
