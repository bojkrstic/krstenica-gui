package repository

import (
	"context"
	"krstenica/internal/errorx"
	"krstenica/internal/model"

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

func (r *repo) ListEparhije(ctx context.Context) ([]model.Eparhija, error) {

	var eparhija []model.Eparhija

	err := r.db.WithContext(ctx).
		Where("status !=?", "deleted").
		Find(&eparhija).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrEparhijeNotFound
		}
		return nil, err
	}

	return eparhija, nil
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
