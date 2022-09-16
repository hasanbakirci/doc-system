package document

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Document struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Extension   string             `bson:"extension"`
	Path        string             `bson:"path"`
	MimeType    string             `bson:"mime_type"`
	CreatedAt   string             `bson:"created_at"`
	UpdatedAt   string             `bson:"updated_at"`
}

func (d *Document) Create() (*Document, error) {
	return &Document{
		ID:          primitive.NewObjectID(),
		Name:        d.Name,
		Description: d.Description,
		Extension:   d.Extension,
		Path:        d.Path,
		MimeType:    d.MimeType,
		CreatedAt:   time.Now().Format("2006-01-02-15-04-05"),
		UpdatedAt:   time.Now().Format("2006-01-02-15-04-05"),
	}, nil
}

type CreateDocumentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Extension   string `json:"extension"`
	Path        string `json:"path"`
	MimeType    string `json:"mime_type"`
}

func (receiver *CreateDocumentRequest) ToDocument() *Document {
	return &Document{
		ID:          primitive.NewObjectID(),
		Name:        receiver.Name,
		Description: receiver.Description,
		Extension:   receiver.Extension,
		Path:        receiver.Path,
		MimeType:    receiver.MimeType,
		CreatedAt:   time.Now().Format("2006-01-02-15-04-05"),
		UpdatedAt:   time.Now().Format("2006-01-02-15-04-05"),
	}
}

type UpdateDocumentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Extension   string `json:"extension"`
	Path        string `json:"path"`
	MimeType    string `json:"mime_type"`
}

func (receiver *UpdateDocumentRequest) ToDocument() *Document {
	return &Document{
		Name:        receiver.Name,
		Description: receiver.Description,
		Extension:   receiver.Extension,
		Path:        receiver.Path,
		MimeType:    receiver.MimeType,
		UpdatedAt:   time.Now().Format("2006-01-02-15-04-05"),
	}
}

type DocumentResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Extension   string `json:"extension"`
	Path        string `json:"path"`
	MimeType    string `json:"mime_type"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (d *Document) ToDocumentResponse() *DocumentResponse {
	return &DocumentResponse{
		ID:          d.ID.String(),
		Name:        d.Name,
		Description: d.Description,
		Extension:   d.Extension,
		Path:        d.Path,
		MimeType:    d.MimeType,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

type DocumentLog struct {
	DocumentId  string `json:"document_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Extension   string `json:"extension"`
	Path        string `json:"path"`
	MimeType    string `json:"mime_type"`
	UserId      string `json:"user_id"`
}

func CreateDocumentLog(doc *Document, uid string) *DocumentLog {
	return &DocumentLog{
		DocumentId:  doc.ID.String(),
		Name:        doc.Name,
		Description: doc.Description,
		Extension:   doc.Extension,
		Path:        doc.Path,
		MimeType:    doc.MimeType,
		UserId:      uid,
	}
}
