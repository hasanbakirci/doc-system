package auth

import (
	"context"
	"github.com/hasanbakirci/doc-system/internal/config"
	"github.com/hasanbakirci/doc-system/pkg/helpers"
	"github.com/hasanbakirci/doc-system/pkg/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Create(ctx context.Context, request CreateUserRequest) (string, error)
	Update(ctx context.Context, id primitive.ObjectID, request UpdateUserRequest) (bool, error)
	Delete(ctx context.Context, id primitive.ObjectID) (bool, error)
	GetAll(ctx context.Context) ([]UserResponse, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*UserResponse, error)
	Login(ctx context.Context, request LoginUserRequest) (string, error)
}

type authService struct {
	repository Repository
	config     config.Configuration
}

func (a authService) Login(ctx context.Context, request LoginUserRequest) (string, error) {
	status, err := a.repository.CheckEmail(ctx, request.Email)
	if !status {
		response.Panic(404, err.Error())
	}

	user, e := a.repository.GetByEmail(ctx, request.Email)
	if !user.CheckPasswordHash(request.Password) {
		response.Panic(400, e.Error())
	}
	token := helpers.GenerateJwtToken(user.ID.String(), user.Role, a.config.JwtSettings)
	return token, nil
}

func (a authService) Create(ctx context.Context, request CreateUserRequest) (string, error) {
	status, err := a.repository.CheckEmail(ctx, request.Email)
	if status {
		response.Panic(404, err.Error())
	}
	user := request.ToUser()
	user.HashPassword()

	id, e := a.repository.Create(ctx, user)
	if e != nil {
		response.Panic(404, e.Error())
	}
	return id.String(), nil
}

func (a authService) Update(ctx context.Context, id primitive.ObjectID, request UpdateUserRequest) (bool, error) {
	user := request.ToUser()
	user.HashPassword()
	result, err := a.repository.Update(ctx, id, user)
	if !result {
		response.Panic(404, err.Error())
	}
	return result, nil
}

func (a authService) Delete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	result, err := a.repository.Delete(ctx, id)
	if !result {
		response.Panic(404, err.Error())
	}
	return true, nil
}

func (a authService) GetAll(ctx context.Context) ([]UserResponse, error) {
	users, err := a.repository.GetAll(ctx)
	if err != nil {
		response.Panic(404, err.Error())
	}
	userResponses := make([]UserResponse, 0)
	for i := 0; i < len(users); i++ {
		u := users[i].ToUserResponse()
		userResponses = append(userResponses, *u)
	}
	return userResponses, nil
}

func (a authService) GetById(ctx context.Context, id primitive.ObjectID) (*UserResponse, error) {
	user, err := a.repository.GetById(ctx, id)
	if err == nil {
		response.Panic(404, err.Error())
	}
	result := user.ToUserResponse()
	return result, nil
}

func NewAuthService(repo Repository, cfg config.Configuration) Service {
	return &authService{repository: repo, config: cfg}
}
