package pkg

import (
	"fmt"
	"log"
	"strings"
)

type FilterPropertyValidator func(attributeName string, attributeValue []string) (string, error)

type SortPropertyTransformer func(attributeName string) (string, error)

func FilterToSQL(filters map[FilterKey][]string, fn FilterPropertyValidator) (string, []interface{}, error) {
	var where []string
	params := []interface{}{}
	log.Printf("filters %v", filters)

	for k, val := range filters {
		log.Printf("val  %v %v", k, val)
		w, p, err := filterToSQL(k.Property, k.Operator, val, fn)
		if err != nil {
			return "", nil, err
		}
		if w != "" {
			where = append(where, w)
			if p != nil {
				params = append(params, p...)
			}
		}
	}

	return strings.Join(where, " AND "), params, nil
}

func filterToSQL(attr, op string, value []string, fn FilterPropertyValidator) (string, []interface{}, error) {

	val, err := filterValueParser(op, value)
	if err != nil {
		if err == fmt.Errorf("UNSUPPORTED_OPERATOR") {
			return "", nil, nil
		}
		return "", nil, err
	}

	attribute, err := fn(attr, val)
	if err != nil {
		// when attribute us unsupported, skip , ignore unsupported attributes
		if err == fmt.Errorf("UNSUPPORTED_FILTER_PROPERTY") {
			return "", nil, nil
		}
		return "", nil, err
	}

	if attribute == "" {
		return "", nil, nil
	}
	switch op {
	case "isnull":
		return fmt.Sprintf("%s IS NULL", attribute), nil, nil
	case "isnotnull":
		return fmt.Sprintf("%s IS NOT NULL", attribute), nil, nil
	case "isempty":
		return fmt.Sprintf("%s = ''", attribute), nil, nil
	case "isnotempty":
		return fmt.Sprintf("%s <> ''", attribute), nil, nil
	case "eq":
		return fmt.Sprintf("%s = ?", attribute), stringValueToInterface(val[0]), nil
	case "neq":
		return fmt.Sprintf("%s <> ?", attribute), stringValueToInterface(val[0]), nil
	case "startswith":
		return fmt.Sprintf("%s LIKE ?", attribute), stringValueToInterface(val[0] + "%"), nil
	case "contains":
		return fmt.Sprintf("%s LIKE ?", attribute), stringValueToInterface("%" + val[0] + "%"), nil
	case "endswith":
		return fmt.Sprintf("%s LIKE ?", attribute), stringValueToInterface("%" + val[0]), nil
	case "doesnotcontain":
		return fmt.Sprintf("%s NOT LIKE ?", attribute), stringValueToInterface("%" + val[0] + "%"), nil
	case "lt":
		return fmt.Sprintf("%s < ?", attribute), stringValueToInterface(val[0]), nil
	case "lte":
		return fmt.Sprintf("%s <= ?", attribute), stringValueToInterface(val[0]), nil
	case "gt":
		return fmt.Sprintf("%s > ?", attribute), stringValueToInterface(val[0]), nil
	case "gte":
		return fmt.Sprintf("%s >= ?", attribute), stringValueToInterface(val[0]), nil
	case "in":
		return fmt.Sprintf("%s IN  (", attribute) + strings.TrimSuffix(strings.Repeat("?,", len(val)), ",") + ")", stringSliceToInterface(val), nil
	case "notin":
		return fmt.Sprintf("%s NOT IN  (", attribute) + strings.TrimSuffix(strings.Repeat("?,", len(val)), ",") + ")", stringSliceToInterface(val), nil
	case "between":
		return fmt.Sprintf("%s BETWEEN ? AND ?", attribute), stringSliceToInterface(val), nil
	case "notbetween":
		return fmt.Sprintf("%s NOT BETWEEN ? AND ?", attribute), stringSliceToInterface(val), nil
	default:
		return "", nil, nil
	}

}

func stringValueToInterface(s string) []interface{} {
	return []interface{}{s}
}

func stringSliceToInterface(s []string) []interface{} {
	par := []interface{}{}
	for _, e := range s {
		par = append(par, e)
	}
	return par
}

func filterValueParser(op string, value []string) ([]string, error) {
	switch op {
	case "isnull", "isnotnull", "isempty", "isnotempty":
		return nil, nil
	case "eq", "neq", "lt", "lte", "gt", "gte", "startswith", "endswith", "contains", "doesnotcontain":
		if len(value) != 1 {
			return nil, fmt.Errorf("BAD_PARAM")
		}
		return value[:1], nil
	case "in", "notin":
		if len(value) == 0 {
			return nil, fmt.Errorf("BAD_PARAM")
		}
		par := []string{}
		for _, e := range value {
			sl := strings.Split(e, ",")
			for _, e := range sl {
				par = append(par, e)
			}
		}
		return par, nil
	case "between", "notbetween":
		if len(value) != 1 {
			return nil, fmt.Errorf("BAD_PARAM")
		}
		vals := strings.Split(value[0], ",")
		if len(vals) != 2 {
			return nil, fmt.Errorf("BAD_PARAM")
		}
		return []string{vals[0], vals[1]}, nil
	default:
		return nil, fmt.Errorf("UNSUPPORTED_OPERATOR")
	}
}

// SortSQL returns SQL Order BY clause based on list of sort options suitable to query embeding
func SortSQL(sorting []*SortOptions, fn SortPropertyTransformer) (string, error) {
	var sorts []string
	//params := []interface{}{}

	for _, val := range sorting {
		s, err := sortToSQL(val, fn)
		if err != nil {
			return "", err
		}
		if s != "" {
			sorts = append(sorts, s)
		}
	}

	return strings.Join(sorts, ","), nil
}

func sortToSQL(attr *SortOptions, fn SortPropertyTransformer) (string, error) {
	property, err := fn(attr.Property)
	if err != nil {
		// when attribute us unsupported, ignore it
		if err == fmt.Errorf("UNSUPPORTED_SORT_PROPERTY") {
			return "", nil
		}
		return "", err
	}

	if property == "" {
		return "", nil
	}

	return property + " " + attr.Direction, nil

}
