package handler

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"krstenica/internal/dto"
	"krstenica/pkg"
)

const refreshEparhijeEvent = "{\"refresh-eparhije-table\": true}"

func (h *httpHandler) addGuiRoutes() {
	h.router.GET("/ui", h.renderDashboard())
	h.router.GET("/ui/", h.renderDashboard())
	h.router.GET("/ui/krstenice", h.renderKrstenicePage())
	h.router.GET("/ui/krstenice/table", h.renderKrsteniceTable())
	h.router.GET("/ui/krstenice/new", h.renderKrsteniceNew())

	h.router.GET("/ui/eparhije", h.renderEparhijePage())
	h.router.GET("/ui/eparhije/table", h.renderEparhijeTable())
	h.router.GET("/ui/eparhije/new", h.renderEparhijeNew())
	h.router.GET("/ui/eparhije/:id/edit", h.renderEparhijeEdit())
	h.router.POST("/ui/eparhije", h.handleEparhijeCreate())
	h.router.PUT("/ui/eparhije/:id", h.handleEparhijeUpdate())
	h.router.DELETE("/ui/eparhije/:id", h.handleEparhijeDelete())
}

func (h *httpHandler) renderDashboard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "dashboard/index.html", gin.H{
			"Title":           "Kontrolna tabla",
			"ContentTemplate": "dashboard/content",
		})
	}
}

type krsteniceTableData struct {
	Items      []*dto.Krstenica
	Pagination paginationData
	Total      int64
	Filters    map[string]string
}

type eparhijeTableData struct {
	Items      []*dto.Eparhije
	Pagination paginationData
	Total      int64
	Filters    map[string]string
}

type paginationData struct {
	Page       int
	PageSize   int
	Total      int64
	TotalPages int
	HasPrev    bool
	HasNext    bool
	PrevPage   int
	NextPage   int
	PrevLink   string
	NextLink   string
	Query      string
}

func (h *httpHandler) renderKrstenicePage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "krstenice/index.html", gin.H{
			"Title":           "Krstenice",
			"ContentTemplate": "krstenice/content",
		})
	}
}

func (h *httpHandler) renderKrsteniceTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		filters := pkg.ParseUrlQuery(ctx)
		pageSize := parsePositiveInt(filters.Paging.PageSize, 10)
		pageNumber := parsePositiveInt(filters.Paging.PageNumber, 1)

		cx := context.Background()

		items, total, err := h.service.ListKrstenice(cx, filters)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		queryValues := cloneValues(ctx.Request.URL.Query())

		data := &krsteniceTableData{
			Items:   items,
			Total:   total,
			Filters: buildFilterMap(queryValues),
			Pagination: paginationData{
				Page:       pageNumber,
				PageSize:   pageSize,
				Total:      total,
				TotalPages: calculateTotalPages(total, pageSize),
				HasPrev:    pageNumber > 1,
				HasNext:    int64(pageNumber*pageSize) < total,
				PrevPage:   max(pageNumber-1, 1),
				NextPage:   pageNumber + 1,
			},
		}

		data.Pagination.Query = queryValues.Encode()
		data.Pagination.PrevLink = buildPageLink(ctx.Request.URL.Path, queryValues, data.Pagination.PrevPage, pageSize)
		data.Pagination.NextLink = buildPageLink(ctx.Request.URL.Path, queryValues, data.Pagination.NextPage, pageSize)

		ctx.HTML(http.StatusOK, "krstenice/table.html", data)
	}
}

func (h *httpHandler) renderKrsteniceNew() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "krstenice/new.html", gin.H{})
	}
}

func (h *httpHandler) renderEparhijePage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "eparhije/index.html", gin.H{
			"Title":           "Eparhije",
			"ContentTemplate": "eparhije/content",
		})
	}
}

func (h *httpHandler) renderEparhijeTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := h.buildEparhijeTable(ctx.Request.URL.Query(), ctx.Request.URL.Path)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.HTML(http.StatusOK, "eparhije/table.html", data)
	}
}

func (h *httpHandler) renderEparhijeNew() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "eparhije/new.html", gin.H{})
	}
}

func (h *httpHandler) renderEparhijeEdit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Nepostojeci identifikator eparhije",
			})
			return
		}

		cx := context.Background()
		eparhija, err := h.service.GetEparhijeByID(cx, int64(id))
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.HTML(http.StatusOK, "eparhije/edit.html", gin.H{
			"Eparhija": eparhija,
		})
	}
}

func (h *httpHandler) handleEparhijeCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.Request.ParseForm(); err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Neuspesno parsiranje forme",
			})
			return
		}

		name := strings.TrimSpace(ctx.PostForm("name"))
		city := strings.TrimSpace(ctx.PostForm("city"))
		formState := gin.H{
			"Name": name,
			"City": city,
		}

		if name == "" {
			ctx.HTML(http.StatusBadRequest, "eparhije/new.html", gin.H{
				"Error": "Naziv eparhije je obavezan",
				"Form":  formState,
			})
			return
		}

		cx := context.Background()
		_, err := h.service.CreateEparhije(cx, &dto.EparhijeCreateReq{
			Name: name,
			City: city,
		})
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "eparhije/new.html", gin.H{
				"Error": err.Error(),
				"Form":  formState,
			})
			return
		}

		ctx.Header("HX-Trigger", refreshEparhijeEvent)
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) handleEparhijeUpdate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Nepostojeci identifikator eparhije",
			})
			return
		}

		if err := ctx.Request.ParseForm(); err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Neuspesno parsiranje forme",
			})
			return
		}

		rawName := ctx.PostForm("name")
		rawCity := ctx.PostForm("city")
		rawStatus := ctx.PostForm("status")

		req := &dto.EparhijeUpdateReq{}
		if name := strings.TrimSpace(rawName); name != "" {
			req.Name = &name
		}

		if city := strings.TrimSpace(rawCity); city != "" {
			req.City = &city
		}

		if status := strings.TrimSpace(rawStatus); status != "" {
			req.Status = &status
		}

		cx := context.Background()
		if _, err := h.service.UpdateEparhije(cx, int64(id), req); err != nil {
			ctx.HTML(http.StatusBadRequest, "eparhije/edit.html", gin.H{
				"Error": err.Error(),
				"Eparhija": &dto.Eparhije{
					ID:     int64(id),
					Name:   rawName,
					City:   rawCity,
					Status: rawStatus,
				},
			})
			return
		}

		ctx.Header("HX-Trigger", refreshEparhijeEvent)
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) handleEparhijeDelete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Nepostojeci identifikator eparhije",
			})
			return
		}

		cx := context.Background()
		if err := h.service.DeleteEparhije(cx, int64(id)); err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.Header("HX-Trigger", refreshEparhijeEvent)
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) buildEparhijeTable(values url.Values, basePath string) (*eparhijeTableData, error) {
	filters := &pkg.FilterAndSort{
		Filters: map[pkg.FilterKey][]string{},
		Sort:    []*pkg.SortOptions{},
		Paging:  &pkg.Paging{},
	}

	pageNumber := parsePositiveInt(values.Get("page_number"), 1)
	pageSize := parsePositiveInt(values.Get("page_size"), 10)
	filters.Paging.PageNumber = strconv.Itoa(pageNumber)
	filters.Paging.PageSize = strconv.Itoa(pageSize)

	for key, val := range values {
		if isPagingKey(key) {
			continue
		}

		trimmed := make([]string, 0, len(val))
		for _, item := range val {
			if strings.TrimSpace(item) != "" {
				trimmed = append(trimmed, item)
			}
		}
		if len(trimmed) == 0 {
			continue
		}

		filters.Filters[pkg.FilterKey{Property: key, Operator: "eq"}] = trimmed
	}

	items, total, err := h.service.ListEparhije(context.Background(), filters)
	if err != nil {
		return nil, err
	}

	queryCopy := cloneValues(values)

	data := &eparhijeTableData{
		Items:   items,
		Total:   total,
		Filters: buildFilterMap(queryCopy),
		Pagination: paginationData{
			Page:       pageNumber,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: calculateTotalPages(total, pageSize),
			HasPrev:    pageNumber > 1,
			HasNext:    int64(pageNumber*pageSize) < total,
			PrevPage:   max(pageNumber-1, 1),
			NextPage:   pageNumber + 1,
			Query:      queryCopy.Encode(),
		},
	}

	data.Pagination.PrevLink = buildPageLink(basePath, queryCopy, data.Pagination.PrevPage, pageSize)
	data.Pagination.NextLink = buildPageLink(basePath, queryCopy, data.Pagination.NextPage, pageSize)

	return data, nil
}

func parsePositiveInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	v, err := strconv.Atoi(value)
	if err != nil || v <= 0 {
		return defaultValue
	}

	return v
}

func calculateTotalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}

	if total == 0 {
		return 0
	}

	d := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		d++
	}

	return int(d)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func buildPageLink(basePath string, baseValues url.Values, page int, pageSize int) string {
	values := url.Values{}
	for key, v := range baseValues {
		if isPagingKey(key) {
			continue
		}
		for _, item := range v {
			values.Add(key, item)
		}
	}

	values.Set("page_number", strconv.Itoa(page))
	values.Set("page_size", strconv.Itoa(pageSize))

	encoded := values.Encode()
	if encoded == "" {
		return basePath
	}
	return basePath + "?" + encoded
}

func cloneValues(values url.Values) url.Values {
	copyValues := url.Values{}
	for key, v := range values {
		for _, item := range v {
			copyValues.Add(key, item)
		}
	}
	return copyValues
}

func buildFilterMap(values url.Values) map[string]string {
	filters := map[string]string{}
	for key, val := range values {
		if isPagingKey(key) {
			continue
		}
		if len(val) == 0 {
			continue
		}
		trimmed := strings.TrimSpace(val[0])
		if trimmed == "" {
			continue
		}
		filters[key] = trimmed
	}
	return filters
}

func isPagingKey(key string) bool {
	switch key {
	case "page_number", "page_size", "paging", "all", "sort":
		return true
	default:
		return false
	}
}
