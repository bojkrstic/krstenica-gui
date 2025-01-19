package repository

import (
	"context"
	"krstenica/internal/errorx"
	"krstenica/internal/model"

	"gorm.io/gorm"
)

func (r *repo) GetKrstenicaByID(ctx context.Context, id int64) (*model.Krstenica, error) {
	var krstenica model.Krstenica

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&krstenica).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrKrstenicaNotFound
		}
		return nil, err
	}

	return &krstenica, nil
}

func (r *repo) ListKrstenice(ctx context.Context) ([]model.Krstenica, error) {

	var krstenica []model.Krstenica

	err := r.db.WithContext(ctx).
		Where("status !=?", "deleted").
		Find(&krstenica).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrKrstenicaNotFound
		}
		return nil, err
	}

	return krstenica, nil
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
