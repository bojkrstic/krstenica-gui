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
