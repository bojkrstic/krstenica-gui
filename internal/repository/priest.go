package repository

import (
	"context"
	"krstenica/internal/errorx"
	"krstenica/internal/model"

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

func (r *repo) ListPriests(ctx context.Context) ([]model.Priest, error) {

	var priest []model.Priest

	err := r.db.WithContext(ctx).
		Where("status !=?", "deleted").
		Find(&priest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.ErrPriestNotFound
		}
		return nil, err
	}

	return priest, nil
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
