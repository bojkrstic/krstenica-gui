package repository

import (
	"context"
	"strconv"
	"strings"

	"krstenica/internal/model"
	"krstenica/pkg"

	"gorm.io/gorm"
)

type Repo interface {
	GetTampleByID(ctx context.Context, id int64) (*model.Tample, error)
	CreateTample(ctx context.Context, tample *model.Tample) (*model.Tample, error)
	UpdateTample(ctx context.Context, id int64, updates map[string]interface{}) error
	ListTamples(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Tample, int64, error)

	GetPriestByID(ctx context.Context, id int64) (*model.Priest, error)
	CreatePriest(ctx context.Context, priest *model.Priest) (*model.Priest, error)
	UpdatePriest(ctx context.Context, id int64, updates map[string]interface{}) error
	ListPriests(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Priest, int64, error)

	GetEparhijeByID(ctx context.Context, id int64) (*model.Eparhija, error)
	CreateEparhije(ctx context.Context, eparhija *model.Eparhija) (*model.Eparhija, error)
	UpdateEparhije(ctx context.Context, id int64, updates map[string]interface{}) error
	ListEparhije(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Eparhija, int64, error)

	GetPersonByID(ctx context.Context, id int64) (*model.Person, error)
	CreatePerson(ctx context.Context, person *model.Person) (*model.Person, error)
	UpdatePerson(ctx context.Context, id int64, updates map[string]interface{}) error
	ListPersons(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Person, int64, error)

	GetKrstenicaByID(ctx context.Context, id int64) (*model.Krstenica, error)
	CreateKrstenica(ctx context.Context, krstenica *model.KrstenicaPost) (*model.Krstenica, error)
	UpdateKrstenica(ctx context.Context, id int64, updates map[string]interface{}) error
	ListKrstenice(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Krstenica, int64, error)

	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
	CountUsers(ctx context.Context) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, id int64, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, id int64) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repo {
	return &repo{db: db}
}

func (r *repo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repo) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repo) ListUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := r.db.WithContext(ctx).Order("username ASC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repo) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repo) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repo) UpdateUser(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repo) DeleteUser(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func Paginate(db *gorm.DB, dest interface{}, limit int) *gorm.DB {
	return db.Limit(limit).Offset(0).Find(dest)
}

const defaultPageSize = 10

func applyPagination(db *gorm.DB, fas *pkg.FilterAndSort) *gorm.DB {
	if fas == nil || fas.Paging == nil {
		return db
	}

	paging := fas.Paging
	if strings.EqualFold(strings.TrimSpace(paging.All), "yes") {
		return db
	}
	if strings.EqualFold(strings.TrimSpace(paging.Paging), "no") {
		return db
	}

	pageSize := parsePositiveInt(paging.PageSize, defaultPageSize)
	if pageSize <= 0 {
		return db
	}

	pageNumber := parsePositiveInt(paging.PageNumber, 1)
	if pageNumber <= 0 {
		pageNumber = 1
	}

	offset := (pageNumber - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	return db.Limit(pageSize).Offset(offset)
}

func parsePositiveInt(raw string, fallback int) int {
	value := strings.TrimSpace(raw)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}
