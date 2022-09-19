package document

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type documentRepository struct {
	collection *mongo.Collection
}

func (d documentRepository) Create(ctx context.Context, document *Document) (string, error) {
	result, err := d.collection.InsertOne(ctx, document)
	if result.InsertedID == nil {
		return "", err
	}
	return document.ID, nil
}

func (d documentRepository) Update(ctx context.Context, id string, document *Document) (bool, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"name":        document.Name,
		"description": document.Description,
		"extension":   document.Extension,
		"path":        document.Path,
		"mime_type":   document.Path,
		"updated_at":  time.Now().Format("2006-01-02-15-04-05"),
	}}
	updateResult, err := d.collection.UpdateOne(ctx, filter, update)
	if updateResult.ModifiedCount < 1 {
		return false, err
	}
	return true, nil
}

func (d documentRepository) Delete(ctx context.Context, id string) (bool, error) {
	deleteResult, err := d.collection.DeleteOne(ctx, bson.M{"_id": id})
	if deleteResult.DeletedCount < 1 {
		return false, err
	}
	return true, nil
}

func (d documentRepository) GetAll(ctx context.Context) ([]Document, error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	documents := make([]Document, 0)
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, err
	}
	return documents, nil
}

func (d documentRepository) GetById(ctx context.Context, id string) (*Document, error) {
	document := new(Document)
	if err := d.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&document); err != nil {
		return nil, err
	}
	return document, nil
}

func NewDocumentRepository(db *mongo.Database) Repository {
	col := db.Collection("documents")
	return &documentRepository{collection: col}
}
