package document

import (
	"context"
	"fmt"
	"github.com/hasanbakirci/doc-system/pkg/redisClient"
	"github.com/hasanbakirci/doc-system/pkg/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Create(ctx context.Context, request CreateDocumentRequest, uid string) (string, error)
	Update(ctx context.Context, id primitive.ObjectID, request UpdateDocumentRequest) (bool, error)
	Delete(ctx context.Context, id primitive.ObjectID) (bool, error)
	GetAll(ctx context.Context) ([]DocumentResponse, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*DocumentResponse, error)
}

type documentService struct {
	repository Repository
	redis      *redisClient.RedisClient
}

func (d documentService) Create(ctx context.Context, request CreateDocumentRequest, uid string) (string, error) {
	document := request.ToDocument()
	id, err := d.repository.Create(ctx, document)
	if err != nil {
		response.Panic(404, err.Error())
	}

	d.redis.Publish("doc-system", CreateDocumentLog(document, uid))
	return id.String(), nil
}

func (d documentService) Update(ctx context.Context, id primitive.ObjectID, request UpdateDocumentRequest) (bool, error) {
	document := request.ToDocument()
	result, err := d.repository.Update(ctx, id, document)
	if !result {
		response.Panic(404, err.Error())
	}
	return result, nil
}

func (d documentService) Delete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	result, err := d.repository.Delete(ctx, id)
	if !result {
		response.Panic(404, err.Error())
	}
	return result, nil
}

func (d documentService) GetAll(ctx context.Context) ([]DocumentResponse, error) {
	fmt.Println("service ---> ", ctx.Value("id"))
	documents, err := d.repository.GetAll(ctx)
	if len(documents) < 1 {
		response.Panic(404, err.Error())
	}
	documentResponses := make([]DocumentResponse, 0)
	for i := 0; i < len(documents); i++ {
		doc := documents[i].ToDocumentResponse()
		documentResponses = append(documentResponses, *doc)
	}
	return documentResponses, nil
}

func (d documentService) GetById(ctx context.Context, id primitive.ObjectID) (*DocumentResponse, error) {
	document, err := d.repository.GetById(ctx, id)
	if &document == nil {
		response.Panic(404, err.Error())
	}
	result := document.ToDocumentResponse()
	return result, nil
}

func NewDocumentService(repo Repository, redis *redisClient.RedisClient) Service {
	return &documentService{repository: repo, redis: redis}
}
