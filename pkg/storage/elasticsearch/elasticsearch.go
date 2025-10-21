package elasticsearch

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"yourapp/pkg/config"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var (
	client *elasticsearch.Client
)

// Init initializes the Elasticsearch connection
func Init(ctx context.Context, cfg config.ElasticsearchConfig) error {
	esConfig := elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)},
		Username:  cfg.Username,
		Password:  cfg.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
		},
	}

	var err error
	client, err = elasticsearch.NewClient(esConfig)
	if err != nil {
		return fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	// Test the connection
	res, err := client.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Elasticsearch ping failed with status: %s", res.Status())
	}

	return nil
}

// GetClient returns the Elasticsearch client
func GetClient() *elasticsearch.Client {
	return client
}

// Close closes the Elasticsearch connection
func Close() error {
	// Elasticsearch client doesn't need explicit closing
	return nil
}

// Health checks the health of the Elasticsearch connection
func Health(ctx context.Context) error {
	if client == nil {
		return fmt.Errorf("Elasticsearch client not initialized")
	}

	res, err := client.Cluster.Health()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Elasticsearch health check failed with status: %s", res.Status())
	}

	return nil
}

// CreateIndex creates an index with the given name and mapping
func CreateIndex(ctx context.Context, indexName string, mapping string) error {
	if client == nil {
		return fmt.Errorf("Elasticsearch client not initialized")
	}

	res, err := client.Indices.Create(
		indexName,
		client.Indices.Create.WithBody(strings.NewReader(mapping)),
		client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to create index %s: %w", indexName, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create index %s with status: %s", indexName, res.Status())
	}

	return nil
}

// DeleteIndex deletes an index
func DeleteIndex(ctx context.Context, indexName string) error {
	if client == nil {
		return fmt.Errorf("Elasticsearch client not initialized")
	}

	res, err := client.Indices.Delete(
		[]string{indexName},
		client.Indices.Delete.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to delete index %s: %w", indexName, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to delete index %s with status: %s", indexName, res.Status())
	}

	return nil
}

// IndexDocument indexes a document
func IndexDocument(ctx context.Context, indexName, documentID string, document string) error {
	if client == nil {
		return fmt.Errorf("Elasticsearch client not initialized")
	}

	res, err := client.Index(
		indexName,
		strings.NewReader(document),
		client.Index.WithDocumentID(documentID),
		client.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to index document with status: %s", res.Status())
	}

	return nil
}

// Search performs a search query
func Search(ctx context.Context, indexName string, query string) (*esapi.Response, error) {
	if client == nil {
		return nil, fmt.Errorf("Elasticsearch client not initialized")
	}

	res, err := client.Search(
		client.Search.WithIndex(indexName),
		client.Search.WithBody(strings.NewReader(query)),
		client.Search.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to perform search: %w", err)
	}

	if res.IsError() {
		return nil, fmt.Errorf("search failed with status: %s", res.Status())
	}

	return res, nil
}
