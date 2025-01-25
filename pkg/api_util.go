package pkg

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// FilterKey contains property and operator.
// filter is a map[FilterKey][]string
type FilterKey struct {
	Property string
	Operator string
}

// SortOptions holt sort parameter and direction. Direction can be DESC in  case of descending sorting and
// ASC or blank value in case of ascending sorting.
type SortOptions struct {
	Property  string
	Direction string
}

type Paging struct {
	All        string `query:"all" default:"no"`
	PageNumber string `query:"page_number" default:"1"`
	PageSize   string `query:"page_size"`
	Paging     string `query:"paging"`
}

type FilterAndSort struct {
	Filters map[FilterKey][]string
	Sort    []*SortOptions
	Paging  *Paging
}

func ParseUrlQuery(ctx *gin.Context) *FilterAndSort {

	m := &FilterAndSort{
		Filters: map[FilterKey][]string{},
		Sort:    []*SortOptions{},
		Paging:  &Paging{},
	}

	pageNumber, exist := ctx.GetQuery("page_number")
	if exist {
		m.Paging.PageNumber = pageNumber
	} else {
		m.Paging.PageNumber = "1"
	}

	pageSize, exist := ctx.GetQuery("page_size")
	if exist {

		m.Paging.PageSize = pageSize
	} else {
		m.Paging.PageSize = "10"
	}

	paging, exist := ctx.GetQuery("paging")
	if exist {
		m.Paging.Paging = paging
	}

	_, existAll := ctx.GetQuery("all")
	if existAll {
		m.Paging.All = "yes"
	}

	sort, exist := ctx.GetQuery("sort")

	if exist {
		sortKeys := strings.Split(sort, ",")
		for _, v := range sortKeys {
			if strings.HasPrefix(v, "-") {
				m.Sort = append(m.Sort, &SortOptions{
					Property:  v[1:],
					Direction: "DESC",
				})
			} else {
				m.Sort = append(m.Sort, &SortOptions{
					Property:  v,
					Direction: "ASC",
				})
			}
		}
	}

	// if exist {
	// 	start := strings.Index(sort, "(")
	// 	end := strings.Index(sort, ")")
	// 	substringKeySort := sort[start+1 : end]
	// 	sortKeys := strings.Split(substringKeySort, ",")
	// 	for _, v := range sortKeys {
	// 		if strings.HasPrefix(v, "-") {
	// 			m.Sort = append(m.Sort, &SortOptions{
	// 				Property:  v[1:],
	// 				Direction: "DESC",
	// 			})
	// 		} else {
	// 			m.Sort = append(m.Sort, &SortOptions{
	// 				Property:  v,
	// 				Direction: "ASC",
	// 			})
	// 		}
	// 	}
	// }

	// create maps
	keysWords := []string{"sort", "page_number", "page_size", "paging", "all"}

	queryParams := ctx.Request.URL.Query()

	for key, val := range queryParams {
		if InList(key, keysWords) {
			continue
		}
		//filters
		if !strings.Contains(key, "(") {
			m.Filters[FilterKey{
				Property: key,
				Operator: "eq",
			}] = val
		} else {
			start := strings.Index(key, "(")
			end := strings.Index(key, ")")
			substringKey := key[start+1 : end]
			substringOp := key[:start-1]

			m.Filters[FilterKey{
				Property: substringKey,
				Operator: substringOp,
			}] = val
		}
	}

	return m
}

// if key exists into slice keys
func InList(elem string, list []string) bool {
	for _, el := range list {
		if el == elem {
			return true
		}
	}
	return false
}
