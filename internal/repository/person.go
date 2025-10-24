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

func (r *repo) GetPersonByID(ctx context.Context, id int64) (*model.Person, error) {
	var person model.Person

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&person).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrPersonNotFound
		}
		return nil, err
	}

	return &person, nil
}

func (r *repo) ListPersons(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Person, int64, error) {

	var person []model.Person

	where, whereParams, err := pkg.FilterToSQL(filterAndSort.Filters, validatePersonFilterAttr)
	if err != nil {
		return nil, 0, err
	}

	if where == "" {
		where += "t.status != 'deleted' "
	} else {
		where += " AND t.status != 'deleted' "
	}

	orderBy, err := pkg.SortSQL(filterAndSort.Sort, transformPersonSortAttribute)
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
	query := r.db.WithContext(ctx).
		Table("persons AS t").
		Where(where, whereParams...).
		Order(orderBy)

	query = applyPagination(query, filterAndSort)

	err = query.Find(&person).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, errorx.ErrPersonNotFound
		}
		return nil, 0, err
	}

	var totalCount int64

	//totalCount
	err = r.db.Table("persons AS t").
		Where("t.status != 'deleted' ").
		Where(where, whereParams...).
		Order(orderBy).
		Count(&totalCount).
		Error
	if err != nil {
		return nil, 0, err
	}

	return person, totalCount, nil
}

var allowedAtributesInPersonFilters = []string{
	"id", "first_name", "last_name", "brief_name", "occupation", "religion", "address", "country", "role", "status", "city", "birth_date", "created_at",
}

var allowedAtributesInPersonSort = []string{
	"id", "first_name", "last_name", "brief_name", "occupation", "religion", "address", "country", "role", "status", "city", "birth_date", "created_at",
}

func transformPersonSortAttribute(p string) (string, error) {
	if !pkg.InList(p, allowedAtributesInPersonSort) {
		return "", fmt.Errorf("UNSUPPORTED_SORT_PROPERTY")
	}

	return "t." + p, nil
}

func validatePersonFilterAttr(p string, v []string) (string, error) {
	if !pkg.InList(p, allowedAtributesInPersonFilters) {
		return "", fmt.Errorf("UNSUPPORTED_FILTER_PROPERTY")
	}

	return "t." + p, nil
}

func (r *repo) CreatePerson(ctx context.Context, person *model.Person) (*model.Person, error) {
	err := r.db.WithContext(ctx).Create(person).Error
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (r *repo) UpdatePerson(ctx context.Context, id int64, updates map[string]interface{}) error {
	err := r.db.WithContext(ctx).
		Table("persons").
		Where("id = ? ", id).
		Updates(updates).Error
	if err != nil {
		return err
	}

	return nil
}
