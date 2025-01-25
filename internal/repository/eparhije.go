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

func (r *repo) GetEparhijeByID(ctx context.Context, id int64) (*model.Eparhija, error) {
	var eparhija model.Eparhija

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&eparhija).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrEparhijeNotFound
		}
		return nil, err
	}

	return &eparhija, nil
}

func (r *repo) ListEparhije(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Eparhija, int64, error) {

	var eparhija []model.Eparhija
	where, whereParams, err := pkg.FilterToSQL(filterAndSort.Filters, validateEparhijeFilterAttr)
	if err != nil {
		return nil, 0, err
	}

	if where == "" {
		where += "t.status != 'deleted' "
	} else {
		where += " AND t.status != 'deleted' "
	}

	orderBy, err := pkg.SortSQL(filterAndSort.Sort, transformEparhijeSortAttribute)
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
		Table("eparhije AS t").
		Where(where, whereParams...).
		Order(orderBy).
		Find(&eparhija).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, errorx.ErrEparhijeNotFound
		}
		return nil, 0, err
	}

	var totalCount int64

	//totalCount
	err = r.db.Table("eparhije AS t").
		Where("t.status != 'deleted' ").
		Where(where, whereParams...).
		Order(orderBy).
		Count(&totalCount).
		Error
	if err != nil {
		return nil, 0, err
	}

	return eparhija, totalCount, nil
}

var allowedAtributesInEparhijeFilters = []string{
	"id", "name", "status", "city", "created_at",
}

var allowedAtributesInEparhijeSort = []string{
	"id", "name", "status", "city", "created_at",
}

func transformEparhijeSortAttribute(p string) (string, error) {
	if !pkg.InList(p, allowedAtributesInEparhijeSort) {
		return "", fmt.Errorf("UNSUPPORTED_SORT_PROPERTY")
	}

	return "t." + p, nil
}

func validateEparhijeFilterAttr(p string, v []string) (string, error) {
	if !pkg.InList(p, allowedAtributesInEparhijeFilters) {
		return "", fmt.Errorf("UNSUPPORTED_FILTER_PROPERTY")
	}

	return "t." + p, nil
}

func (r *repo) CreateEparhije(ctx context.Context, eparhija *model.Eparhija) (*model.Eparhija, error) {
	err := r.db.WithContext(ctx).Create(eparhija).Error
	if err != nil {
		return nil, err
	}

	return eparhija, nil
}

func (r *repo) UpdateEparhije(ctx context.Context, id int64, updates map[string]interface{}) error {
	err := r.db.WithContext(ctx).
		Table("eparhije").
		Where("id = ? ", id).
		Updates(updates).Error
	if err != nil {
		return err
	}

	return nil
}
