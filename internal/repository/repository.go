package repository

import "gorm.io/gorm"

type Repo interface {
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repo {
	return &repo{db: db}
}

func Paginate(db *gorm.DB, dest interface{}, limit int) *gorm.DB {
	return db.Limit(limit).Offset(0).Find(dest)
}
