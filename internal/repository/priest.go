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

func (r *repo) GetPriestByID(ctx context.Context, id int64) (*model.Priest, error) {
	var priest model.Priest

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&priest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrPriestNotFound
		}
		return nil, err
	}

	return &priest, nil
}

func (r *repo) ListPriests(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Priest, int64, error) {

	var priest []model.Priest

	where, whereParams, err := pkg.FilterToSQL(filterAndSort.Filters, validatePriestFilterAttr)
	if err != nil {
		return nil, 0, err
	}

	if where == "" {
		where += "t.status != 'deleted' "
	} else {
		where += " AND t.status != 'deleted' "
	}

	orderBy, err := pkg.SortSQL(filterAndSort.Sort, transformPriestSortAttribute)
	if err != nil {
		return nil, 0, err
	}

	if orderBy != "" {
		if !strings.Contains(orderBy, "t.id") {
			orderBy += ", t.id DESC"
		}
	} else {
		orderBy = "t.id DESC"
	}

	query := r.db.WithContext(ctx).
		Table("priests AS t").
		Where(where, whereParams...).
		Order(orderBy)

	query = applyPagination(query, filterAndSort)

	err = query.Find(&priest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, errorx.ErrPriestNotFound
		}
		return nil, 0, err
	}

	var totalCount int64

	//totalCount
	err = r.db.Table("priests AS t").
		Where("t.status != 'deleted' ").
		Where(where, whereParams...).
		Order(orderBy).
		Count(&totalCount).
		Error
	if err != nil {
		return nil, 0, err
	}

	return priest, totalCount, nil
}

var allowedAtributesInPriestFilters = []string{
	"id", "first_name", "last_name", "city", "title", "status", "created_at",
}

var allowedAtributesInPriestSort = []string{
	"id", "first_name", "last_name", "city", "title", "status", "created_at",
}

func transformPriestSortAttribute(p string) (string, error) {
	if !pkg.InList(p, allowedAtributesInPriestSort) {
		return "", fmt.Errorf("UNSUPPORTED_SORT_PROPERTY")
	}

	return "t." + p, nil
}

func validatePriestFilterAttr(p string, v []string) (string, error) {
	if !pkg.InList(p, allowedAtributesInPriestFilters) {
		return "", fmt.Errorf("UNSUPPORTED_FILTER_PROPERTY")
	}

	return "t." + p, nil
}

func (r *repo) CreatePriest(ctx context.Context, priest *model.Priest) (*model.Priest, error) {
	err := r.db.WithContext(ctx).Create(priest).Error
	if err != nil {
		return nil, err
	}

	return priest, nil
}

func (r *repo) UpdatePriest(ctx context.Context, id int64, updates map[string]interface{}) error {
	err := r.db.WithContext(ctx).
		Table("priests").
		Where("id = ? ", id).
		Updates(updates).Error
	if err != nil {
		return err
	}

	return nil
}
