package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"krstenica/internal/dto"
	"krstenica/pkg"
)

const (
	refreshEparhijeEvent   = "{\"refresh-eparhije-table\": true}"
	refreshHramoviEvent    = "{\"refresh-hramovi-table\": true}"
	refreshSvesteniciEvent = "{\"refresh-svestenici-table\": true}"
	refreshOsobeEvent      = "{\"refresh-osobe-table\": true}"
)

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

	h.router.GET("/ui/hramovi", h.renderHramoviPage())
	h.router.GET("/ui/hramovi/table", h.renderHramoviTable())
	h.router.GET("/ui/hramovi/new", h.renderHramoviNew())
	h.router.GET("/ui/hramovi/:id/edit", h.renderHramoviEdit())
	h.router.POST("/ui/hramovi", h.handleHramoviCreate())
	h.router.PUT("/ui/hramovi/:id", h.handleHramoviUpdate())
	h.router.DELETE("/ui/hramovi/:id", h.handleHramoviDelete())

	h.router.GET("/ui/svestenici", h.renderSvesteniciPage())
	h.router.GET("/ui/svestenici/table", h.renderSvesteniciTable())
	h.router.GET("/ui/svestenici/new", h.renderSvesteniciNew())
	h.router.GET("/ui/svestenici/:id/edit", h.renderSvesteniciEdit())
	h.router.POST("/ui/svestenici", h.handleSvesteniciCreate())
	h.router.PUT("/ui/svestenici/:id", h.handleSvesteniciUpdate())
	h.router.DELETE("/ui/svestenici/:id", h.handleSvesteniciDelete())

	h.router.GET("/ui/osobe", h.renderOsobePage())
	h.router.GET("/ui/osobe/table", h.renderOsobeTable())
	h.router.GET("/ui/osobe/new", h.renderOsobeNew())
	h.router.GET("/ui/osobe/:id/edit", h.renderOsobeEdit())
	h.router.POST("/ui/osobe", h.handleOsobeCreate())
	h.router.PUT("/ui/osobe/:id", h.handleOsobeUpdate())
	h.router.DELETE("/ui/osobe/:id", h.handleOsobeDelete())

	h.router.GET("/ui/osobe/picker", h.renderOsobePicker())
	h.router.GET("/ui/osobe/picker/table", h.renderOsobePickerTable())
	h.router.GET("/ui/osobe/picker/select/:id", h.handleOsobePickerSelect())
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

type hramoviTableData struct {
	Items      []*dto.Tample
	Pagination paginationData
	Total      int64
	Filters    map[string]string
}

type svesteniciTableData struct {
	Items      []*dto.Priest
	Pagination paginationData
	Total      int64
	Filters    map[string]string
}

type osobeTableData struct {
	Items      []*dto.Person
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
		cx := context.Background()

		eparhije, err := h.listActiveEparhijeForForm(cx)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		hramovi, err := h.listActiveHramoviForForm(cx)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.HTML(http.StatusOK, "krstenice/new.html", gin.H{
			"Eparhije": eparhije,
			"Hramovi":  hramovi,
		})
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

func (h *httpHandler) renderHramoviPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "hramovi/index.html", gin.H{
			"Title":           "Hramovi",
			"ContentTemplate": "hramovi/content",
		})
	}
}

func (h *httpHandler) renderHramoviTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := h.buildHramoviTable(ctx.Request.URL.Query(), ctx.Request.URL.Path)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.HTML(http.StatusOK, "hramovi/table.html", data)
	}
}

func (h *httpHandler) renderHramoviNew() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "hramovi/new.html", gin.H{})
	}
}

func (h *httpHandler) renderHramoviEdit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Nepostojeci identifikator hrama",
			})
			return
		}

		cx := context.Background()
		hram, err := h.service.GetTampleByID(cx, int64(id))
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.HTML(http.StatusOK, "hramovi/edit.html", gin.H{
			"Hram": hram,
		})
	}
}

func (h *httpHandler) handleHramoviCreate() gin.HandlerFunc {
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
			ctx.HTML(http.StatusBadRequest, "hramovi/new.html", gin.H{
				"Error": "Naziv hrama je obavezan",
				"Form":  formState,
			})
			return
		}

		cx := context.Background()
		_, err := h.service.CreateTample(cx, &dto.TampleCreateReq{
			Name: name,
			City: city,
		})
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "hramovi/new.html", gin.H{
				"Error": err.Error(),
				"Form":  formState,
			})
			return
		}

		ctx.Header("HX-Trigger", refreshHramoviEvent)
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) handleHramoviUpdate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Nepostojeci identifikator hrama",
			})
			return
		}

		if err := ctx.Request.ParseForm(); err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Neuspesno parsiranje forme",
			})
			return
		}

		rawName := strings.TrimSpace(ctx.PostForm("name"))
		rawCity := strings.TrimSpace(ctx.PostForm("city"))
		rawStatus := strings.TrimSpace(ctx.PostForm("status"))

		if rawName == "" {
			ctx.HTML(http.StatusBadRequest, "hramovi/edit.html", gin.H{
				"Error": "Naziv hrama je obavezan",
				"Hram": &dto.Tample{
					ID:     int64(id),
					Name:   rawName,
					City:   rawCity,
					Status: rawStatus,
				},
			})
			return
		}

		req := &dto.TampleUpdateReq{}
		nameCopy := rawName
		req.Name = &nameCopy
		req.City = &rawCity
		if rawStatus != "" {
			statusCopy := rawStatus
			req.Status = &statusCopy
		}

		cx := context.Background()
		if _, err := h.service.UpdateTample(cx, int64(id), req); err != nil {
			ctx.HTML(http.StatusBadRequest, "hramovi/edit.html", gin.H{
				"Error": err.Error(),
				"Hram": &dto.Tample{
					ID:     int64(id),
					Name:   rawName,
					City:   rawCity,
					Status: rawStatus,
				},
			})
			return
		}

		ctx.Header("HX-Trigger", refreshHramoviEvent)
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) handleHramoviDelete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Nepostojeci identifikator hrama",
			})
			return
		}

		cx := context.Background()
		if err := h.service.DeleteTample(cx, int64(id)); err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.Header("HX-Trigger", refreshHramoviEvent)
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) buildHramoviTable(values url.Values, basePath string) (*hramoviTableData, error) {
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

	items, total, err := h.service.ListTamples(context.Background(), filters)
	if err != nil {
		return nil, err
	}

	queryCopy := cloneValues(values)

	data := &hramoviTableData{
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

func (h *httpHandler) renderSvesteniciPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "svestenici/index.html", gin.H{
			"Title":           "Svestenici",
			"ContentTemplate": "svestenici/content",
		})
	}
}

func (h *httpHandler) renderSvesteniciTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := h.buildSvesteniciTable(ctx.Request.URL.Query(), ctx.Request.URL.Path)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.HTML(http.StatusOK, "svestenici/table.html", data)
	}
}

func (h *httpHandler) renderSvesteniciNew() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "svestenici/new.html", gin.H{})
	}
}

func (h *httpHandler) renderSvesteniciEdit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Nepostojeci identifikator svestenika",
			})
			return
		}

		cx := context.Background()
		svestenik, err := h.service.GetPriestByID(cx, int64(id))
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.HTML(http.StatusOK, "svestenici/edit.html", gin.H{
			"Svestenik": svestenik,
		})
	}
}

func (h *httpHandler) handleSvesteniciCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.Request.ParseForm(); err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Neuspesno parsiranje forme",
			})
			return
		}

		firstName := strings.TrimSpace(ctx.PostForm("first_name"))
		lastName := strings.TrimSpace(ctx.PostForm("last_name"))
		city := strings.TrimSpace(ctx.PostForm("city"))
		title := strings.TrimSpace(ctx.PostForm("title"))

		formState := gin.H{
			"FirstName": firstName,
			"LastName":  lastName,
			"City":      city,
			"Title":     title,
		}

		if firstName == "" || lastName == "" {
			ctx.HTML(http.StatusBadRequest, "svestenici/new.html", gin.H{
				"Error": "Ime i prezime su obavezni",
				"Form":  formState,
			})
			return
		}

		cx := context.Background()
		_, err := h.service.CreatePriest(cx, &dto.PriestCreateReq{
			FirstName: firstName,
			LastName:  lastName,
			City:      city,
			Title:     title,
		})
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "svestenici/new.html", gin.H{
				"Error": err.Error(),
				"Form":  formState,
			})
			return
		}

		ctx.Header("HX-Trigger", refreshSvesteniciEvent)
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) handleSvesteniciUpdate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Nepostojeci identifikator svestenika",
			})
			return
		}

		if err := ctx.Request.ParseForm(); err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Neuspesno parsiranje forme",
			})
			return
		}

		firstName := strings.TrimSpace(ctx.PostForm("first_name"))
		lastName := strings.TrimSpace(ctx.PostForm("last_name"))
		city := strings.TrimSpace(ctx.PostForm("city"))
		title := strings.TrimSpace(ctx.PostForm("title"))
		status := strings.TrimSpace(ctx.PostForm("status"))

		if firstName == "" || lastName == "" {
			ctx.HTML(http.StatusBadRequest, "svestenici/edit.html", gin.H{
				"Error": "Ime i prezime su obavezni",
				"Svestenik": &dto.Priest{
					ID:        int64(id),
					FirstName: firstName,
					LastName:  lastName,
					City:      city,
					Title:     title,
					Status:    status,
				},
			})
			return
		}

		req := &dto.PriestUpdateReq{}
		firstNameCopy := firstName
		req.FirstName = &firstNameCopy
		lastNameCopy := lastName
		req.LastName = &lastNameCopy
		cityCopy := city
		req.City = &cityCopy
		titleCopy := title
		req.Title = &titleCopy
		if status != "" {
			statusCopy := status
			req.Status = &statusCopy
		}

		cx := context.Background()
		if _, err := h.service.UpdatePriest(cx, int64(id), req); err != nil {
			ctx.HTML(http.StatusBadRequest, "svestenici/edit.html", gin.H{
				"Error": err.Error(),
				"Svestenik": &dto.Priest{
					ID:        int64(id),
					FirstName: firstName,
					LastName:  lastName,
					City:      city,
					Title:     title,
					Status:    status,
				},
			})
			return
		}

		ctx.Header("HX-Trigger", refreshSvesteniciEvent)
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) handleSvesteniciDelete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Nepostojeci identifikator svestenika",
			})
			return
		}

		cx := context.Background()
		if err := h.service.DeletePriest(cx, int64(id)); err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.Header("HX-Trigger", refreshSvesteniciEvent)
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) buildSvesteniciTable(values url.Values, basePath string) (*svesteniciTableData, error) {
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

	items, total, err := h.service.ListPriests(context.Background(), filters)
	if err != nil {
		return nil, err
	}

	queryCopy := cloneValues(values)

	data := &svesteniciTableData{
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

func (h *httpHandler) renderOsobePage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "osobe/index.html", gin.H{
			"Title":           "Osobe",
			"ContentTemplate": "osobe/content",
		})
	}
}

func (h *httpHandler) renderOsobeTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := h.buildOsobeTable(ctx.Request.URL.Query(), ctx.Request.URL.Path)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.HTML(http.StatusOK, "osobe/table.html", data)
	}
}

func (h *httpHandler) renderOsobeNew() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "osobe/new.html", gin.H{})
	}
}

func (h *httpHandler) renderOsobeEdit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Nepostojeci identifikator osobe",
			})
			return
		}

		cx := context.Background()
		osoba, err := h.service.GetPersonByID(cx, int64(id))
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.HTML(http.StatusOK, "osobe/edit.html", gin.H{
			"Osoba": osoba,
		})
	}
}

func (h *httpHandler) renderOsobePicker() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		field := strings.TrimSpace(ctx.Query("field"))
		if field == "" {
			field = "parent_id"
		}

		ctx.HTML(http.StatusOK, "osobe/picker.html", gin.H{
			"Field": field,
		})
	}
}

func (h *httpHandler) renderOsobePickerTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		field := strings.TrimSpace(ctx.Query("field"))
		if field == "" {
			field = "parent_id"
		}

		values := cloneValues(ctx.Request.URL.Query())
		if values.Get("status") == "" {
			values.Set("status", "active")
		}

		data, err := h.buildOsobeTable(values, ctx.Request.URL.Path)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.HTML(http.StatusOK, "osobe/picker-table.html", gin.H{
			"Field": field,
			"Data":  data,
		})
	}
}

func (h *httpHandler) handleOsobePickerSelect() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Nepostojeci identifikator osobe",
			})
			return
		}

		field := strings.TrimSpace(ctx.Query("field"))
		if field == "" {
			field = "parent_id"
		}

		person, err := h.service.GetPersonByID(context.Background(), int64(id))
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		name := strings.TrimSpace(strings.Join([]string{
			strings.TrimSpace(person.FirstName),
			strings.TrimSpace(person.LastName),
		}, " "))
		label := name
		if person.Role != "" {
			label = label + " (" + person.Role + ")"
		}

		payload := map[string]interface{}{
			"person-selected": map[string]interface{}{
				"field": field,
				"id":    person.ID,
				"label": label,
			},
			"close-picker": true,
		}

		bytes, err := json.Marshal(payload)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.Header("HX-Trigger", string(bytes))
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) handleOsobeCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.Request.ParseForm(); err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Neuspesno parsiranje forme",
			})
			return
		}

		firstName := strings.TrimSpace(ctx.PostForm("first_name"))
		lastName := strings.TrimSpace(ctx.PostForm("last_name"))
		briefName := strings.TrimSpace(ctx.PostForm("brief_name"))
		role := strings.TrimSpace(ctx.PostForm("role"))
		occupation := strings.TrimSpace(ctx.PostForm("occupation"))
		religion := strings.TrimSpace(ctx.PostForm("religion"))
		address := strings.TrimSpace(ctx.PostForm("address"))
		country := strings.TrimSpace(ctx.PostForm("country"))
		city := strings.TrimSpace(ctx.PostForm("city"))
		birthDateRaw := strings.TrimSpace(ctx.PostForm("birth_date"))

		formState := gin.H{
			"FirstName":  firstName,
			"LastName":   lastName,
			"BriefName":  briefName,
			"Role":       role,
			"Occupation": occupation,
			"Religion":   religion,
			"Address":    address,
			"Country":    country,
			"City":       city,
			"BirthDate":  birthDateRaw,
		}

		if firstName == "" || lastName == "" || role == "" {
			ctx.HTML(http.StatusBadRequest, "osobe/new.html", gin.H{
				"Error": "Ime, prezime i uloga su obavezni",
				"Form":  formState,
			})
			return
		}

		var birthDate time.Time
		if birthDateRaw != "" {
			parsed, err := time.Parse("2006-01-02", birthDateRaw)
			if err != nil {
				ctx.HTML(http.StatusBadRequest, "osobe/new.html", gin.H{
					"Error": "Neispravan format datuma. Koristite YYYY-MM-DD.",
					"Form":  formState,
				})
				return
			}
			birthDate = parsed
		}

		cx := context.Background()
		_, err := h.service.CreatePerson(cx, &dto.PersonCreateReq{
			FirstName:  firstName,
			LastName:   lastName,
			BriefName:  briefName,
			Occupation: occupation,
			Religion:   religion,
			Address:    address,
			Country:    country,
			Role:       role,
			BirthDate:  birthDate,
			City:       city,
		})
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "osobe/new.html", gin.H{
				"Error": err.Error(),
				"Form":  formState,
			})
			return
		}

		ctx.Header("HX-Trigger", refreshOsobeEvent)
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) handleOsobeUpdate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Nepostojeci identifikator osobe",
			})
			return
		}

		if err := ctx.Request.ParseForm(); err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Neuspesno parsiranje forme",
			})
			return
		}

		firstName := strings.TrimSpace(ctx.PostForm("first_name"))
		lastName := strings.TrimSpace(ctx.PostForm("last_name"))
		briefName := strings.TrimSpace(ctx.PostForm("brief_name"))
		role := strings.TrimSpace(ctx.PostForm("role"))
		occupation := strings.TrimSpace(ctx.PostForm("occupation"))
		religion := strings.TrimSpace(ctx.PostForm("religion"))
		address := strings.TrimSpace(ctx.PostForm("address"))
		country := strings.TrimSpace(ctx.PostForm("country"))
		city := strings.TrimSpace(ctx.PostForm("city"))
		birthDateRaw := strings.TrimSpace(ctx.PostForm("birth_date"))
		status := strings.TrimSpace(ctx.PostForm("status"))

		if firstName == "" || lastName == "" || role == "" {
			ctx.HTML(http.StatusBadRequest, "osobe/edit.html", gin.H{
				"Error": "Ime, prezime i uloga su obavezni",
				"Osoba": &dto.Person{
					ID:         int64(id),
					FirstName:  firstName,
					LastName:   lastName,
					BriefName:  briefName,
					Role:       role,
					Occupation: occupation,
					Religion:   religion,
					Address:    address,
					Country:    country,
					City:       city,
					Status:     status,
				},
			})
			return
		}

		var birthDatePtr *time.Time
		if birthDateRaw != "" {
			parsed, err := time.Parse("2006-01-02", birthDateRaw)
			if err != nil {
				ctx.HTML(http.StatusBadRequest, "osobe/edit.html", gin.H{
					"Error": "Neispravan format datuma. Koristite YYYY-MM-DD.",
					"Osoba": &dto.Person{
						ID:         int64(id),
						FirstName:  firstName,
						LastName:   lastName,
						BriefName:  briefName,
						Role:       role,
						Occupation: occupation,
						Religion:   religion,
						Address:    address,
						Country:    country,
						City:       city,
						Status:     status,
					},
				})
				return
			}
			birthDatePtr = &parsed
		}

		req := &dto.PersonUpdateReq{}
		firstNameCopy := firstName
		req.FirstName = &firstNameCopy
		lastNameCopy := lastName
		req.LastName = &lastNameCopy
		briefNameCopy := briefName
		req.BriefName = &briefNameCopy
		roleCopy := role
		req.Role = &roleCopy
		occupationCopy := occupation
		req.Occupation = &occupationCopy
		religionCopy := religion
		req.Religion = &religionCopy
		addressCopy := address
		req.Address = &addressCopy
		countryCopy := country
		req.Country = &countryCopy
		cityCopy := city
		req.City = &cityCopy
		if birthDatePtr != nil {
			req.BirthDate = birthDatePtr
		}
		if status != "" {
			statusCopy := status
			req.Status = &statusCopy
		}

		cx := context.Background()
		if _, err := h.service.UpdatePerson(cx, int64(id), req); err != nil {
			osoba := &dto.Person{
				ID:         int64(id),
				FirstName:  firstName,
				LastName:   lastName,
				BriefName:  briefName,
				Role:       role,
				Occupation: occupation,
				Religion:   religion,
				Address:    address,
				Country:    country,
				City:       city,
				Status:     status,
			}
			if birthDatePtr != nil {
				osoba.BirthDate = *birthDatePtr
			}
			ctx.HTML(http.StatusBadRequest, "osobe/edit.html", gin.H{
				"Error": err.Error(),
				"Osoba": osoba,
			})
			return
		}

		ctx.Header("HX-Trigger", refreshOsobeEvent)
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) handleOsobeDelete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "partials/error.html", gin.H{
				"Message": "Nepostojeci identifikator osobe",
			})
			return
		}

		cx := context.Background()
		if err := h.service.DeletePerson(cx, int64(id)); err != nil {
			ctx.HTML(http.StatusInternalServerError, "partials/error.html", gin.H{
				"Message": err.Error(),
			})
			return
		}

		ctx.Header("HX-Trigger", refreshOsobeEvent)
		ctx.Status(http.StatusNoContent)
	}
}

func (h *httpHandler) buildOsobeTable(values url.Values, basePath string) (*osobeTableData, error) {
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

	items, total, err := h.service.ListPersons(context.Background(), filters)
	if err != nil {
		return nil, err
	}

	queryCopy := cloneValues(values)

	data := &osobeTableData{
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

func (h *httpHandler) listActiveEparhijeForForm(ctx context.Context) ([]*dto.Eparhije, error) {
	filters := &pkg.FilterAndSort{
		Filters: map[pkg.FilterKey][]string{},
		Sort:    []*pkg.SortOptions{},
		Paging: &pkg.Paging{
			All: "yes",
		},
	}

	filters.Filters[pkg.FilterKey{Property: "status", Operator: "eq"}] = []string{"active"}
	filters.Sort = append(filters.Sort, &pkg.SortOptions{Property: "name", Direction: "ASC"})

	items, _, err := h.service.ListEparhije(ctx, filters)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (h *httpHandler) listActiveHramoviForForm(ctx context.Context) ([]*dto.Tample, error) {
	filters := &pkg.FilterAndSort{
		Filters: map[pkg.FilterKey][]string{},
		Sort:    []*pkg.SortOptions{},
		Paging: &pkg.Paging{
			All: "yes",
		},
	}

	filters.Filters[pkg.FilterKey{Property: "status", Operator: "eq"}] = []string{"active"}
	filters.Sort = append(filters.Sort, &pkg.SortOptions{Property: "name", Direction: "ASC"})

	items, _, err := h.service.ListTamples(ctx, filters)
	if err != nil {
		return nil, err
	}

	return items, nil
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
