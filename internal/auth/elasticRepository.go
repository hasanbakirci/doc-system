package auth

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

// CheckEmail implements Repository
func (e *elasticRepository) CheckEmail(ctx context.Context, email string) (bool, error) {
	shouldFilter := map[string]interface{}{"term": map[string]interface{}{"Email.keyword": email}}
	query := map[string]interface{}{
		"query": shouldFilter,
	}
	dataBytes, err := json.Marshal(&query)
	if err != nil {
		return false, err
	}
	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(e.index),
		e.client.Search.WithBody(bytes.NewReader(dataBytes)),
		e.client.Search.WithPretty(),
		e.client.Search.WithHuman(),
	)
	if err != nil {
		return false, err
	}

	if res.IsError() {
		return false, errors.Wrap(errors.New(res.String()), "esClient.Search error")
	}
	defer res.Body.Close()

	hits := elasticclient.ElasticResponse[User]{}
	if err := json.NewDecoder(res.Body).Decode(&hits); err != nil {
		errorHandler.Panic(400, err.Error())
	}
	if hits.Hits.Total.Value > 0 {
		return true, nil
	}
	return false, nil
}

// Create implements Repository
func (e *elasticRepository) Create(ctx context.Context, user *User) (string, error) {
	exists, err := e.client.Indices.Exists([]string{e.index})
	if exists.StatusCode != 200 {
		//index, err := e.client.Indices.Create(e.index)
		//fmt.Println(index.Body)
		//if err != nil {
		//	errorHandler.Panic(400, err.Error())
		//}
		//log.Info("created index:", index.Body)
		//alias, err := e.client.Indices.PutAlias([]string{e.index}, uuid.New().String())
		//fmt.Println(alias.Body)
		//if err != nil {
		//	errorHandler.Panic(400, err.Error())
		//}
		//log.Info("created alias:", alias.StatusCode)
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
	user.Create()
	u, _ := json.Marshal(user)
	req := esapi.IndexRequest{Index: e.index, DocumentID: user.ID, Body: bytes.NewReader(u)}
	res, err := req.Do(ctx, e.client)
	if err != nil || res.StatusCode > 201 {
		errorHandler.Panic(400, err.Error())
	}
	defer res.Body.Close()

	return user.ID, nil
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
func (e *elasticRepository) GetAll(ctx context.Context) ([]User, error) {
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

	hits := elasticclient.ElasticResponse[User]{}
	if err := json.NewDecoder(res.Body).Decode(&hits); err != nil {
		errorHandler.Panic(400, err.Error())
	}

	hitsize := len(hits.Hits.Hits)
	if hitsize > 0 {
		resList := make([]User, hitsize)
		for i, source := range hits.Hits.Hits {
			resList[i] = source.Source
		}
		return resList, nil
	}
	return nil, errors.New("Elastic repository: Users index is null")
}

// GetByEmail implements Repository
func (e *elasticRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	shouldFilter := map[string]interface{}{"term": map[string]interface{}{"Email.keyword": email}}
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

	hits := elasticclient.ElasticResponse[User]{}
	if err := json.NewDecoder(res.Body).Decode(&hits); err != nil {
		errorHandler.Panic(400, err.Error())
	}
	hitsize := len(hits.Hits.Hits)
	if hitsize < 1 {
		return nil, errors.New("Elastic repository: email not found")
	}
	resList := make([]User, hitsize)
	for i, source := range hits.Hits.Hits {
		resList[i] = source.Source
	}
	return &resList[0], nil
}

// GetById implements Repository
func (e *elasticRepository) GetById(ctx context.Context, id string) (*User, error) {
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

	hits := elasticclient.ElasticResponse[User]{}
	if err := json.NewDecoder(res.Body).Decode(&hits); err != nil {
		errorHandler.Panic(400, err.Error())
	}
	hitsize := len(hits.Hits.Hits)
	if hitsize < 1 {
		return nil, errors.New("id not found")
	}
	resList := make([]User, hitsize)
	for i, source := range hits.Hits.Hits {
		resList[i] = source.Source
	}
	return &resList[0], nil
}

// Update implements Repository
func (e *elasticRepository) Update(ctx context.Context, id string, user *User) (bool, error) {
	updateField := map[string]interface{}{
		"source": fmt.Sprintf(
			"ctx._source['Username'] = '%s';ctx._source['Password'] = '%s';ctx._source['Email'] = '%s';ctx._source['Role'] = '%s';ctx._source['UpdatedAt'] = '%s'",
			user.Username, user.Password, user.Email, user.Role, user.UpdatedAt)}

	shouldFilter := map[string]interface{}{"match_phrase": map[string]interface{}{"ID.keyword": id}}

	updateRequest := map[string]interface{}{
		"script": updateField,
		"query":  shouldFilter,
	}
	dataBytes, err := json.Marshal(&updateRequest)
	if err != nil {
		return false, err
	}

	reader := bytes.NewReader(dataBytes)
	res, err := e.client.UpdateByQuery([]string{e.index},
		e.client.UpdateByQuery.WithBody(reader),
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
	return &elasticRepository{client: elastic, index: "users_19092022", alias: "users"}
}
