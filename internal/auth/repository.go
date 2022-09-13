package auth

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository interface {
	Create(ctx context.Context, user *User) (primitive.ObjectID, error)
	Update(ctx context.Context, id primitive.ObjectID, user *User) (bool, error)
	Delete(ctx context.Context, id primitive.ObjectID) (bool, error)
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	CheckEmail(ctx context.Context, email string) (bool, error)
}

type authRepository struct {
	collection *mongo.Collection
}

func (a authRepository) Create(ctx context.Context, user *User) (primitive.ObjectID, error) {
	result, err := a.collection.InsertOne(ctx, user)
	if result.InsertedID == nil {
		return primitive.NilObjectID, err
	}
	return user.ID, nil
}

func (a authRepository) Update(ctx context.Context, id primitive.ObjectID, user *User) (bool, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"user":       user.Username,
		"password":   user.Password,
		"email":      user.Email,
		"role":       user.Role,
		"updated_at": time.Now().Format("2006-01-02-15-04-05"),
	}}
	updateResult, err := a.collection.UpdateOne(ctx, filter, update)
	if updateResult.ModifiedCount < 1 {
		return false, err
	}
	return true, nil
}

func (a authRepository) Delete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	deleteResult, err := a.collection.DeleteOne(ctx, bson.M{"_id": id})
	if deleteResult.DeletedCount < 1 {
		return false, err
	}
	return true, nil
}

func (a authRepository) GetAll(ctx context.Context) ([]User, error) {
	cursor, err := a.collection.Find(ctx, bson.M{})
	users := make([]User, 0)
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (a authRepository) GetById(ctx context.Context, id primitive.ObjectID) (*User, error) {
	user := new(User)
	if err := a.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (a authRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	user := new(User)
	if err := a.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (a authRepository) CheckEmail(ctx context.Context, email string) (bool, error) {
	count, err := a.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count >= 1, nil
}

func NewAuthRepository(db *mongo.Database) Repository {
	col := db.Collection("users")
	return &authRepository{collection: col}
}
