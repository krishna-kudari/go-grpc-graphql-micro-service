package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	ErrNotFound = errors.New("Entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p *Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elasticsearch.Client
}

// Close implements [Repository].
func (e *elasticRepository) Close() {
	panic("unimplemented")
}

// GetProductByID implements [Repository].
func (e *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	res, err := e.client.Get("catalog", id, e.client.Get.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("Error getting doc: %s", res.String())
	}

	var r struct {
		Source ProductDocument `json:"_source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &Product{
		ID:          r.Source.Id,
		Name:        r.Source.Name,
		Description: r.Source.Description,
		Price:       r.Source.Price,
	}, nil
}

// ListProducts implements [Repository].
func (e *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	es := e.client
	res, err := e.client.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex("catalog"),
		es.Search.WithFrom(int(skip)),
		es.Search.WithSize(int(take)),
		es.Search.WithSort("id:asc"),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("Error fetching docs: %s", res.String())
	}

	var sr struct {
		Hits struct {
			Hits []struct {
				ID     string          `json:"_id"`
				Source ProductDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return nil, err
	}

	var products []Product
	for _, hit := range sr.Hits.Hits {
		products = append(products, Product{
			ID:          hit.Source.Id,
			Name:        hit.Source.Name,
			Description: hit.Source.Description,
			Price:       hit.Source.Price,
		})
	}
	return products, nil
}

// ListProductsWithIDs implements [Repository].
func (e *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	es := e.client
	body := map[string]interface{}{
		"ids": ids,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := es.Mget(
		bytes.NewReader(data),
		es.Mget.WithContext(ctx),
		es.Mget.WithIndex("catalog"),
	)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("mget error: %s", res.String())
	}

	var response struct {
		Docs []struct {
			ID     string          `json:"_id"`
			Found  bool            `json:"found"`
			Source ProductDocument `json:"_source"`
		} `json:"docs"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	var products []Product
	for _, doc := range response.Docs {
		if doc.Found {
			products = append(products, Product{
				ID:          doc.Source.Id,
				Name:        doc.Source.Name,
				Description: doc.Source.Description,
				Price:       doc.Source.Price,
			})
		}
	}
	return products, nil
}

// PutProduct implements [Repository].
func (e *elasticRepository) PutProduct(ctx context.Context, p *Product) error {
	data, err := json.Marshal(ProductDocument{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	})
	if err != nil {
		return err
	}

	res, err := e.client.Index(
		"catalog",
		bytes.NewReader(data),
		e.client.Index.WithDocumentID(p.ID),
		e.client.Index.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

// SearchProducts implements [Repository].
func (e *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	es := e.client
	body := map[string]interface{}{
		"from": skip,
		"size": take,
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  query,
				"fields": []string{"title", "description", "price"},
			},
		},
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	res, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex("catalog"),
		es.Search.WithBody(bytes.NewReader(data)),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}
	var sr struct {
		Hits struct {
			Hits []struct {
				ID     string          `json:"_id"`
				Source ProductDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return nil, err
	}
	var products []Product
	for _, hit := range sr.Hits.Hits {
		products = append(products, Product{
			ID:          hit.Source.Id,
			Name:        hit.Source.Name,
			Description: hit.Source.Description,
			Price:       hit.Source.Price,
		})
	}
	return products, nil
}

type ProductDocument struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{
				url,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client}, nil
}
