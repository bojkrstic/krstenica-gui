package service

import (
	"context"
	"errors"
	"strings"

	"krstenica/internal/dto"
	"krstenica/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *service) EnsureDefaultUser(ctx context.Context) error {
	username := strings.TrimSpace(s.conf.Auth.Username)
	if username == "" {
		username = "admin"
	}
	password := strings.TrimSpace(s.conf.Auth.Password)
	if password == "" {
		password = "admin"
	}

	user, err := s.repo.GetUserByUsername(ctx, username)
	switch {
	case err == nil:
		if password == "" {
			return nil
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		return s.repo.UpdateUser(ctx, user.ID, map[string]interface{}{"password_hash": string(hash)})
	case errors.Is(err, gorm.ErrRecordNotFound):
		return s.createUserInternal(ctx, username, password)
	default:
		return err
	}
}

func (s *service) AuthenticateUser(ctx context.Context, username, password string) (bool, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return false, nil
	}

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return false, nil
	}
	return true, nil
}

func (s *service) ListUsers(ctx context.Context) ([]*dto.User, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*dto.User, 0, len(users))
	for _, u := range users {
		user := u
		res = append(res, &dto.User{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		})
	}
	return res, nil
}

func (s *service) CreateUser(ctx context.Context, req *dto.UserCreateReq) (*dto.User, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" {
		return nil, errors.New("username is required")
	}

	if len(username) > 255 {
		return nil, errors.New("username can not be longer than 255 characters")
	}

	if strings.TrimSpace(req.Password) == "" {
		return nil, errors.New("password is required")
	}

	if _, err := s.repo.GetUserByUsername(ctx, username); err == nil {
		return nil, errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := s.createUserInternal(ctx, username, req.Password); err != nil {
		return nil, err
	}

	created, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &dto.User{
		ID:        created.ID,
		Username:  created.Username,
		CreatedAt: created.CreatedAt,
	}, nil
}

func (s *service) GetUser(ctx context.Context, id int64) (*dto.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.User{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *service) UpdateUser(ctx context.Context, id int64, req *dto.UserUpdateReq) (*dto.User, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}

	current, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	username := strings.TrimSpace(req.Username)
	if username != "" && username != current.Username {
		if len(username) > 255 {
			return nil, errors.New("username can not be longer than 255 characters")
		}
		if _, err := s.repo.GetUserByUsername(ctx, username); err == nil {
			return nil, errors.New("username already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		updates["username"] = username
	}

	password := strings.TrimSpace(req.Password)
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updates["password_hash"] = string(hash)
	}

	if len(updates) == 0 {
		return &dto.User{
			ID:        current.ID,
			Username:  current.Username,
			CreatedAt: current.CreatedAt,
		}, nil
	}

	if err := s.repo.UpdateUser(ctx, id, updates); err != nil {
		return nil, err
	}

	return s.GetUser(ctx, id)
}

func (s *service) DeleteUser(ctx context.Context, id int64) error {
	if _, err := s.repo.GetUserByID(ctx, id); err != nil {
		return err
	}

	count, err := s.repo.CountUsers(ctx)
	if err != nil {
		return err
	}
	if count <= 1 {
		return errors.New("не може се обрисати последњи корисник")
	}

	return s.repo.DeleteUser(ctx, id)
}

func (s *service) createUserInternal(ctx context.Context, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &model.User{
		Username:     username,
		PasswordHash: string(hash),
	}
	_, err = s.repo.CreateUser(ctx, user)
	return err
}
