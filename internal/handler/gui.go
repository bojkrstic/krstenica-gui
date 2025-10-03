package handler

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"

	"krstenica/internal/dto"
	"krstenica/pkg"
)

func (h *httpHandler) addGuiRoutes() {
	h.router.GET("/ui", h.renderDashboard())
	h.router.GET("/ui/", h.renderDashboard())
	h.router.GET("/ui/krstenice", h.renderKrstenicePage())
	h.router.GET("/ui/krstenice/table", h.renderKrsteniceTable())
	h.router.GET("/ui/krstenice/new", h.renderKrsteniceNew())
}

func (h *httpHandler) renderDashboard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "dashboard/index.html", gin.H{
			"Title": "Kontrolna tabla",
		})
	}
}

type krsteniceTableData struct {
	Items      []*dto.Krstenica
	Pagination paginationData
	Total      int64
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
}

func (h *httpHandler) renderKrstenicePage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "krstenice/index.html", gin.H{
			"Title": "Krstenice",
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

		data := &krsteniceTableData{
			Items: items,
			Total: total,
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

		data.Pagination.PrevLink = buildPageLink(ctx, data.Pagination.PrevPage, pageSize)
		data.Pagination.NextLink = buildPageLink(ctx, data.Pagination.NextPage, pageSize)

		ctx.HTML(http.StatusOK, "krstenice/table.html", data)
	}
}

func (h *httpHandler) renderKrsteniceNew() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "krstenice/new.html", gin.H{})
	}
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

func buildPageLink(ctx *gin.Context, page int, pageSize int) string {
	values := url.Values{}
	for key, v := range ctx.Request.URL.Query() {
		for _, item := range v {
			if key == "page_number" || key == "page_size" {
				continue
			}
			values.Add(key, item)
		}
	}

	values.Set("page_number", strconv.Itoa(page))
	values.Set("page_size", strconv.Itoa(pageSize))

	return "/ui/krstenice/table?" + values.Encode()
}
