package repository

import (
	"context"
	"krstenica/internal/errorx"
	"krstenica/internal/model"

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

func (r *repo) ListPersons(ctx context.Context) ([]model.Person, error) {

	var person []model.Person

	err := r.db.WithContext(ctx).
		Where("status !=?", "deleted").
		Find(&person).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrPersonNotFound
		}
		return nil, err
	}

	return person, nil
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
