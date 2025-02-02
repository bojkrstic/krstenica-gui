package repository

import (
	"context"
	"errors"
	"fmt"
	"krstenica/internal/errorx"
	"krstenica/internal/model"
	"krstenica/pkg"
	"strings"

	"gorm.io/gorm"
)

func (r *repo) GetKrstenicaByID(ctx context.Context, id int64) (*model.Krstenica, error) {
	var krstenica model.Krstenica
	if id <= 0 {
		return nil, errors.New("invalid ID provided")
	}

	eparhijaJoin := "LEFT JOIN eparhije as ep on ep.id = t.eparhija_id AND ep.status != 'deleted'"
	tampleJoin := "LEFT JOIN tamples as tm on tm.id = t.tample_id AND tm.status != 'deleted'"
	parentJoin := "LEFT JOIN persons as par on par.id = t.parent_id AND par.status != 'deleted'"
	godFatherJoin := "LEFT JOIN persons as fat on fat.id = t.godfather_id AND fat.status != 'deleted'"
	parohJoin := "LEFT JOIN persons as pa on pa.id = t.paroh_id AND pa.status != 'deleted'"
	priestJoin := "LEFT JOIN priests as pr on pr.id = t.priest_id AND pr.status != 'deleted'"

	err := r.db.WithContext(ctx).
		Debug().
		Table("krstenice AS t").
		Where("t.id = ?", id).
		Joins(eparhijaJoin).
		Joins(tampleJoin).
		Joins(parentJoin).
		Joins(godFatherJoin).
		Joins(parohJoin).
		Joins(priestJoin).
		Select(`t.*, ep.name as eparhija_name,
		tm.name as tample_name,
		tm.city as tample_city, 
		par.first_name as parent_first_name, 
		par.last_name as parent_last_name, 
		par.occupation as parent_occupation, 
		par.city as parent_city, 
		par.religion as parent_religion,
		fat.first_name as godfather_first_name, 
		fat.last_name as godfather_last_name, 
		fat.occupation as godfather_occupation, 
		fat.city as godfather_city, 
		fat.religion as godfather_religion, 
		pa.first_name as paroh_first_name,
		pa.last_name as paroh_last_name,
		pr.first_name as priest_first_name,
		pr.last_name as priest_last_name`).
		First(&krstenica).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // Umesto `err == gorm.ErrRecordNotFound`
			return nil, errorx.ErrKrstenicaNotFound
		}
		return nil, err
	}

	return &krstenica, nil
}

func (r *repo) ListKrstenice(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Krstenica, int64, error) {

	var krstenica []model.Krstenica

	where, whereParams, err := pkg.FilterToSQL(filterAndSort.Filters, validateKrstenicaFilterAttr)
	if err != nil {
		return nil, 0, err
	}

	if where == "" {
		where += "t.status != 'deleted' "
	} else {
		where += " AND t.status != 'deleted' "
	}

	orderBy, err := pkg.SortSQL(filterAndSort.Sort, transformKrstenicaSortAttribute)
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
		Table("krstenice AS t").
		Where(where, whereParams...).
		Order(orderBy).
		Find(&krstenica).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, errorx.ErrKrstenicaNotFound
		}
		return nil, 0, err
	}

	var totalCount int64

	//totalCount
	err = r.db.Table("krstenice AS t").
		Where("t.status != 'deleted' ").
		Where(where, whereParams...).
		Order(orderBy).
		Count(&totalCount).
		Error
	if err != nil {
		return nil, 0, err
	}

	return krstenica, totalCount, nil
}

var allowedAtributesInKrstenicaFilters = []string{
	"id", "book", "page", "current_number", "eparhija_id", "tample_id", "parent_id", "godfather_id", "paroh_id", "priest_id",
	"first_name", "last_name", "gender", "city", "country", "birth_date", "birth_order", "place_of_birthday", "municipality_of_birthday", "baptism",
	"is_church_married", "is_twin", "has_physical_disability", "anagrafa", "number_of_certificate", "town_of_certificate", "certificate",
	"comment", "status", "created_at",
}

var allowedAtributesInKrstenicaSort = []string{
	"id", "book", "page", "current_number", "eparhija_id", "tample_id", "parent_id", "godfather_id", "paroh_id", "priest_id",
	"first_name", "last_name", "gender", "city", "country", "birth_date", "birth_order", "place_of_birthday", "municipality_of_birthday", "baptism",
	"is_church_married", "is_twin", "has_physical_disability", "anagrafa", "number_of_certificate", "town_of_certificate", "certificate",
	"comment", "status", "created_at",
}

func transformKrstenicaSortAttribute(p string) (string, error) {
	if !pkg.InList(p, allowedAtributesInKrstenicaSort) {
		return "", fmt.Errorf("UNSUPPORTED_SORT_PROPERTY")
	}

	return "t." + p, nil
}

func validateKrstenicaFilterAttr(p string, v []string) (string, error) {
	if !pkg.InList(p, allowedAtributesInKrstenicaFilters) {
		return "", fmt.Errorf("UNSUPPORTED_FILTER_PROPERTY")
	}

	return "t." + p, nil
}

func (r *repo) CreateKrstenica(ctx context.Context, krstenica *model.Krstenica) (*model.Krstenica, error) {
	err := r.db.WithContext(ctx).Create(krstenica).Error
	if err != nil {
		return nil, err
	}

	return krstenica, nil
}

func (r *repo) UpdateKrstenica(ctx context.Context, id int64, updates map[string]interface{}) error {
	err := r.db.WithContext(ctx).
		Table("krstenice").
		Where("id = ? ", id).
		Updates(updates).Error
	if err != nil {
		return err
	}

	return nil
}
