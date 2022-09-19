package auth

import (
	"context"
	"github.com/hasanbakirci/doc-system/internal/config"
	"github.com/hasanbakirci/doc-system/pkg/errorHandler"
	"github.com/hasanbakirci/doc-system/pkg/helpers"
)

type Service interface {
	Create(ctx context.Context, request CreateUserRequest) (string, error)
	Update(ctx context.Context, id string, request UpdateUserRequest) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	GetAll(ctx context.Context) ([]UserResponse, error)
	GetById(ctx context.Context, id string) (*UserResponse, error)
	Login(ctx context.Context, request LoginUserRequest) (string, error)
}

type authService struct {
	repository Repository
	config     config.Configuration
}

func (a authService) Login(ctx context.Context, request LoginUserRequest) (string, error) {
	_, err := a.repository.CheckEmail(ctx, request.Email)
	if err != nil {
		errorHandler.Panic(400, "Service: e-mail address already exists")
	}

	user, e := a.repository.GetByEmail(ctx, request.Email)
	if e != nil {
		errorHandler.Panic(404, "Service: e-mail address was not found")
	}
	if !user.CheckPasswordHash(request.Password) {
		errorHandler.Panic(400, "Service: wrong password")
	}
	token := helpers.GenerateJwtToken(user.ID, user.Role, a.config.JwtSettings)
	return token, nil
}

func (a authService) Create(ctx context.Context, request CreateUserRequest) (string, error) {
	status, _ := a.repository.CheckEmail(ctx, request.Email)
	if status {
		errorHandler.Panic(400, "Service: email already exists")
	}
	user := request.ToUser()
	user.HashPassword()

	id, e := a.repository.Create(ctx, user)
	if e != nil {
		errorHandler.Panic(404, "Service: failed to create user")
	}
	return id, nil
}

func (a authService) Update(ctx context.Context, id string, request UpdateUserRequest) (bool, error) {
	user := request.ToUser()
	user.HashPassword()
	result, _ := a.repository.Update(ctx, id, user)
	if !result {
		errorHandler.Panic(404, "Service: failed to update user")
	}
	return result, nil
}

func (a authService) Delete(ctx context.Context, id string) (bool, error) {
	result, _ := a.repository.Delete(ctx, id)
	if !result {
		errorHandler.Panic(404, "Service: failed to delete user")
	}
	return true, nil
}

func (a authService) GetAll(ctx context.Context) ([]UserResponse, error) {
	users, err := a.repository.GetAll(ctx)
	if err != nil {
		errorHandler.Panic(404, err.Error())
	}
	userResponses := make([]UserResponse, 0)
	for i := 0; i < len(users); i++ {
		u := users[i].ToUserResponse()
		userResponses = append(userResponses, *u)
	}
	return userResponses, nil
}

func (a authService) GetById(ctx context.Context, id string) (*UserResponse, error) {
	user, err := a.repository.GetById(ctx, id)
	if err != nil {
		errorHandler.Panic(404, "Service: user id not found")
	}
	result := user.ToUserResponse()
	return result, nil
}

func NewAuthService(repo Repository, cfg config.Configuration) Service {
	return &authService{repository: repo, config: cfg}
}
