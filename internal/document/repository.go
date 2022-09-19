package document

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, document *Document) (string, error)
	Update(ctx context.Context, id string, document *Document) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	GetAll(ctx context.Context) ([]Document, error)
	GetById(ctx context.Context, id string) (*Document, error)
}
