package repository

import (
	"context"
	"fmt"
	"krstenica/internal/errorx"
	"krstenica/internal/model"
	"krstenica/pkg"
	"strings"

	"gorm.io/gorm"
)

func (r *repo) GetTampleByID(ctx context.Context, id int64) (*model.Tample, error) {
	var tample model.Tample

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&tample).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrTampleNotFound
		}
		return nil, err
	}

	return &tample, nil
}

func (r *repo) CreateTample(ctx context.Context, tample *model.Tample) (*model.Tample, error) {
	err := r.db.WithContext(ctx).Create(tample).Error
	if err != nil {
		return nil, err
	}

	return tample, nil
}

func (r *repo) UpdateTample(ctx context.Context, id int64, updates map[string]interface{}) error {
	err := r.db.WithContext(ctx).
		Table("tamples").
		Where("id = ? ", id).
		Updates(updates).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) ListTamples(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Tample, int64, error) {

	var tample []model.Tample

	where, whereParams, err := pkg.FilterToSQL(filterAndSort.Filters, validateTampleFilterAttr)
	if err != nil {
		return nil, 0, err
	}

	if where == "" {
		where += "t.status != 'deleted' "
	} else {
		where += " AND t.status != 'deleted' "
	}

	orderBy, err := pkg.SortSQL(filterAndSort.Sort, transformTamplesSortAttribute)
	if err != nil {
		return nil, 0, err
	}

	if orderBy != "" {
		if !strings.Contains(orderBy, "t.id") {
			orderBy += ", t.id"
		}
	} else {
		orderBy = "t.id"
	}

	err = r.db.WithContext(ctx).
		Table("tamples AS t").
		Where(where, whereParams...).
		Order(orderBy).
		Find(&tample).Error
	if err != nil {
		return nil, 0, err
	}

	var totalCount int64

	//totalCount
	err = r.db.Table("tamples AS t").
		Where("t.status != 'deleted' ").
		Where(where, whereParams...).
		Order(orderBy).
		Count(&totalCount).
		Error
	if err != nil {
		return nil, 0, err
	}

	return tample, totalCount, nil
}

var allowedAtributesInTampleFilters = []string{
	"id", "name", "status", "city", "created_at",
}

var allowedAtributesInTampleSort = []string{
	"id", "name", "status", "city", "created_at",
}

func transformTamplesSortAttribute(p string) (string, error) {
	if !inList(p, allowedAtributesInTampleSort) {
		return "", fmt.Errorf("UNSUPPORTED_SORT_PROPERTY")
	}

	return "t." + p, nil
}

func validateTampleFilterAttr(p string, v []string) (string, error) {
	if !inList(p, allowedAtributesInTampleFilters) {
		return "", fmt.Errorf("UNSUPPORTED_FILTER_PROPERTY")
	}

	return "t." + p, nil
}

func inList(elem string, list []string) bool {
	for _, el := range list {
		if el == elem {
			return true
		}
	}
	return false
}
