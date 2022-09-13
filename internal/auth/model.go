package auth

import (
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	Email     string             `bson:"email"`
	Role      string             `bson:"role"`
	CreatedAt string             `bson:"created_at"`
	UpdatedAt string             `bson:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Username  string             `json:"username"`
	Password  string             `json:"password"`
	Email     string             `json:"email"`
	Role      string             `json:"role"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (receiver *CreateUserRequest) ToUser() *User {
	return &User{
		ID:        primitive.NewObjectID(),
		Username:  receiver.Username,
		Password:  receiver.Password,
		Email:     receiver.Email,
		Role:      receiver.Role,
		CreatedAt: time.Now().Format("2006-01-02-15-04-05"),
		UpdatedAt: time.Now().Format("2006-01-02-15-04-05"),
	}
}

func (receiver *UpdateUserRequest) ToUser() *User {
	return &User{
		Username:  receiver.Username,
		Password:  receiver.Password,
		Email:     receiver.Email,
		Role:      receiver.Role,
		UpdatedAt: time.Now().Format("2006-01-02-15-04-05"),
	}
}

func (receiver *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:        receiver.ID,
		Username:  receiver.Username,
		Password:  receiver.Password,
		Email:     receiver.Email,
		Role:      receiver.Role,
		CreatedAt: receiver.CreatedAt,
		UpdatedAt: receiver.UpdatedAt,
	}
}

func (receiver *User) HashPassword() {
	bytes, err := bcrypt.GenerateFromPassword([]byte(receiver.Password), 14)
	if err != nil {
		log.Error("Hash Password Error : %v", err)
		panic(err)
	}
	receiver.Password = string(bytes)
}

func (receiver *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(receiver.Password), []byte(password))
	return err == nil
}
