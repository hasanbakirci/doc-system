package document

import (
	"context"

	"github.com/hasanbakirci/doc-system/pkg/errorHandler"
	"github.com/hasanbakirci/doc-system/pkg/redisClient"
)

type Service interface {
	Create(ctx context.Context, request CreateDocumentRequest, uid string) (string, error)
	Update(ctx context.Context, id string, request UpdateDocumentRequest) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	GetAll(ctx context.Context) ([]DocumentResponse, error)
	GetById(ctx context.Context, id string) (*DocumentResponse, error)
}

type documentService struct {
	repository Repository
	redis      *redisClient.RedisClient
}

func (d documentService) Create(ctx context.Context, request CreateDocumentRequest, uid string) (string, error) {
	document := request.ToDocument()
	id, err := d.repository.Create(ctx, document)
	if err != nil {
		errorHandler.Panic(404, "Service: failed to create document")
	}

	d.redis.Publish("doc-system", CreateDocumentLog(document, uid))
	return id, nil
}

func (d documentService) Update(ctx context.Context, id string, request UpdateDocumentRequest) (bool, error) {
	document := request.ToDocument()
	result, _ := d.repository.Update(ctx, id, document)
	if !result {
		errorHandler.Panic(404, "Service: failed to update document")
	}
	return result, nil
}

func (d documentService) Delete(ctx context.Context, id string) (bool, error) {
	result, _ := d.repository.Delete(ctx, id)
	if !result {
		errorHandler.Panic(404, "Service: failed to delete document")
	}
	return result, nil
}

func (d documentService) GetAll(ctx context.Context) ([]DocumentResponse, error) {
	documents, err := d.repository.GetAll(ctx)
	if len(documents) < 1 {
		errorHandler.Panic(404, err.Error())
	}
	documentResponses := make([]DocumentResponse, 0)
	for i := 0; i < len(documents); i++ {
		doc := documents[i].ToDocumentResponse()
		documentResponses = append(documentResponses, *doc)
	}
	return documentResponses, nil
}

func (d documentService) GetById(ctx context.Context, id string) (*DocumentResponse, error) {
	document, err := d.repository.GetById(ctx, id)
	if err != nil {
		errorHandler.Panic(404, "Service: document id not found")
	}
	result := document.ToDocumentResponse()
	return result, nil
}

func NewDocumentService(repo Repository, redis *redisClient.RedisClient) Service {
	return &documentService{repository: repo, redis: redis}
}
