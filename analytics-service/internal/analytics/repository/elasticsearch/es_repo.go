package elasticsearch

import (
	"analytics-service/domain/entities"
	"analytics-service/domain/interfaces"
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

type analyticsRepository struct {
	esClient *elasticsearch.Client
	index    string
}

func NewAnalyticsRepository(esClient *elasticsearch.Client, index string) interfaces.AnalyticsRepository {
	return &analyticsRepository{
		esClient: esClient,
		index:    index,
	}
}

func (r *analyticsRepository) Store(ctx context.Context, in entities.URLClicked) error {
	body, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("error marshaling document: %s", err)
	}
	res, err := r.esClient.Index(
		r.index,
		bytes.NewReader(body),
		r.esClient.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}
	return nil
}

func (r *analyticsRepository) GetAllByUrlUUID(ctx context.Context, urlUUID, ownerUUID string) ([]entities.URLClicked, error) {
	var urls []entities.URLClicked

	// Build the query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"match": map[string]interface{}{"url_uuid": urlUUID}},
					{"match": map[string]interface{}{"owner_uuid": ownerUUID}},
				},
			},
		},
		"sort": []map[string]interface{}{
			{"timestamp": map[string]string{"order": "desc"}},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	// Perform the search request
	res, err := r.esClient.Search(
		r.esClient.Search.WithContext(ctx),
		r.esClient.Search.WithIndex(r.index),
		r.esClient.Search.WithBody(&buf),
		r.esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %w", err)
	}
	defer res.Body.Close()

	// Check if response is an error
	if res.IsError() {
		return nil, fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	// Define response struct to parse ES hits
	var esRes struct {
		Hits struct {
			Hits []struct {
				Source entities.URLClicked `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&esRes); err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	// Extract _source into result slice
	for _, hit := range esRes.Hits.Hits {
		urls = append(urls, hit.Source)
	}

	return urls, nil
}
