package document

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	elasticclient "github.com/hasanbakirci/doc-system/pkg/elasticClient"
	log "github.com/sirupsen/logrus"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/hasanbakirci/doc-system/pkg/errorHandler"
	"github.com/pkg/errors"
)

type elasticRepository struct {
	client *elasticsearch.Client
	index  string
	alias  string
}

// Create implements Repository
func (e *elasticRepository) Create(ctx context.Context, document *Document) (string, error) {
	exists, err := e.client.Indices.Exists([]string{e.index})
	if exists.StatusCode != 200 {
		settings := map[string]interface{}{
			"aliases": map[string]interface{}{
				e.alias: map[string]interface{}{},
			},
			"settings": map[string]interface{}{
				"number_of_shards":   3,
				"number_of_replicas": 2,
			},
		}
		dataBytes, err := json.Marshal(&settings)
		if err != nil {
			return "", err
		}
		req := esapi.IndicesCreateRequest{
			Index: e.index,
			Body:  bytes.NewReader(dataBytes),
		}
		res, err := req.Do(ctx, e.client)
		if err != nil || res.StatusCode > 201 {
			errorHandler.Panic(400, err.Error())
		}
		log.Infof("Elastic repositroy: %s index created, %s alias added.", e.index, e.alias)
	}
	if err != nil {
		errorHandler.Panic(404, err.Error())
	}
	document.Create()
	doc, _ := json.Marshal(document)
	req := esapi.IndexRequest{Index: e.index, DocumentID: document.ID, Body: bytes.NewReader(doc)}
	res, err := req.Do(ctx, e.client)
	if err != nil || res.StatusCode > 201 {
		errorHandler.Panic(400, err.Error())
	}
	defer res.Body.Close()

	return document.ID, err
}

// Delete implements Repository
func (e *elasticRepository) Delete(ctx context.Context, id string) (bool, error) {
	shouldFilter := map[string]interface{}{"match": map[string]interface{}{"ID.keyword": id}}
	query := map[string]interface{}{
		"query": shouldFilter,
	}
	dataBytes, err := json.Marshal(&query)
	if err != nil {
		return false, err
	}
	reader := bytes.NewReader(dataBytes)
	res, err := e.client.DeleteByQuery([]string{e.index}, reader,
		e.client.DeleteByQuery.WithContext(ctx),
		e.client.DeleteByQuery.WithPretty(),
		e.client.DeleteByQuery.WithHuman(),
		e.client.DeleteByQuery.WithTimeout(5*time.Second))
	if err != nil {
		return false, err
	}

	if res.IsError() {
		return false, errors.Wrap(errors.New(res.String()), "esClient.Search error")
	}

	defer res.Body.Close()

	deleted := elasticclient.ElasticResultResponse{}
	if err := json.NewDecoder(res.Body).Decode(&deleted); err != nil {
		errorHandler.Panic(400, err.Error())
	}
	if deleted.Deleted < 1 {
		return false, errors.New("id not found")
	}

	return true, nil
}

// GetAll implements Repository
func (e *elasticRepository) GetAll(ctx context.Context) ([]Document, error) {
	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(e.index),
		e.client.Search.WithPretty(),
		e.client.Search.WithHuman(),
		e.client.Search.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, errors.Wrap(errors.New(res.String()), "esClient.Search error")
	}
	defer res.Body.Close()

	hits := elasticclient.ElasticResponse[Document]{}
	if err := json.NewDecoder(res.Body).Decode(&hits); err != nil {
		errorHandler.Panic(400, err.Error())
	}
	hitsize := len(hits.Hits.Hits)
	if hitsize > 0 {
		resList := make([]Document, hitsize)
		for i, source := range hits.Hits.Hits {
			resList[i] = source.Source
		}
		return resList, nil
	}
	return nil, errors.New("Elastic repository: Documents index is null")
}

// GetById implements Repository
func (e *elasticRepository) GetById(ctx context.Context, id string) (*Document, error) {
	shouldFilter := map[string]interface{}{"match": map[string]interface{}{"ID.keyword": id}}
	query := map[string]interface{}{
		"query": shouldFilter,
	}
	dataBytes, err := json.Marshal(&query)
	if err != nil {
		return nil, err
	}
	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(e.index),
		e.client.Search.WithBody(bytes.NewReader(dataBytes)),
		e.client.Search.WithPretty(),
		e.client.Search.WithHuman(),
		e.client.Search.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, errors.Wrap(errors.New(res.String()), "esClient.Search error")
	}
	defer res.Body.Close()

	hits := elasticclient.ElasticResponse[Document]{}
	if err := json.NewDecoder(res.Body).Decode(&hits); err != nil {
		errorHandler.Panic(404, err.Error())
	}
	hitsize := len(hits.Hits.Hits)
	if hitsize < 1 {
		return nil, errors.New("Elastic repository: id not found")
	}
	responseList := make([]Document, len(hits.Hits.Hits))
	for i, source := range hits.Hits.Hits {
		responseList[i] = source.Source
	}
	return &responseList[0], nil
}

// Update implements Repository
func (e *elasticRepository) Update(ctx context.Context, id string, document *Document) (bool, error) {
	updateField := map[string]interface{}{
		"source": fmt.Sprintf(
			"ctx._source['Name'] = '%s';ctx._source['Description'] = '%s';ctx._source['Extension'] = '%s';ctx._source['Path'] = '%s';ctx._source['Mimetype'] = '%s';ctx._source['UpdatedAt'] = '%s'",
			document.Name, document.Description, document.Extension, document.Path, document.MimeType, document.UpdatedAt)}

	shouldFilter := map[string]interface{}{"match_phrase": map[string]interface{}{"ID.keyword": id}}

	updateRequest := map[string]interface{}{
		"script": updateField,
		"query":  shouldFilter,
	}
	dataBytes, err := json.Marshal(&updateRequest)
	if err != nil {
		return false, err
	}

	res, err := e.client.UpdateByQuery([]string{e.index},
		e.client.UpdateByQuery.WithBody(bytes.NewReader(dataBytes)),
		e.client.UpdateByQuery.WithContext(ctx),
		e.client.UpdateByQuery.WithPretty(),
		e.client.UpdateByQuery.WithHuman(),
		e.client.UpdateByQuery.WithTimeout(5*time.Second))
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	updated := elasticclient.ElasticResultResponse{}
	if err := json.NewDecoder(res.Body).Decode(&updated); err != nil {
		errorHandler.Panic(400, err.Error())
	}
	if updated.Updated < 1 {
		return false, errors.New("id not found")
	}
	return true, nil

}

func NewElasticRepository(elastic *elasticsearch.Client) Repository {
	return &elasticRepository{client: elastic, index: "documents_19092022", alias: "documents"}
}
