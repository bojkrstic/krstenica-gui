package repository

import (
	"context"
	"krstenica/internal/errorx"
	"krstenica/internal/model"

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
