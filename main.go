package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/uuid"
)

const (
	indexName = "user_tag"
)

// User represents user.
type User struct {
	Name string
	Age  int
}

func run() error {
	c, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	// indexing

	indexing := func() error {
		id := uuid.NewString()

		req := esapi.IndexRequest{
			Index:      indexName,
			DocumentID: uuid.NewString(),
			Refresh:    "true",
			Body:       strings.NewReader(`{"name":"hlts2", "age":25}`),
		}

		resp, err := req.Do(context.Background(), c)
		if err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
		defer resp.Body.Close()

		if resp.IsError() {
			return fmt.Errorf("failed to create index id: %s  status: %s", id, resp.Status())
		}

		return nil
	}

	// search

	search := func() error {
		var buf bytes.Buffer
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"name": "hlts2",
				},
			},
		}

		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			return err
		}

		resp, err := c.Search(
			c.Search.WithContext(context.Background()),
			c.Search.WithIndex(indexName),
			c.Search.WithBody(&buf),
			c.Search.WithTrackTotalHits(true),
			c.Search.WithPretty(),
		)
		if err != nil {
			return fmt.Errorf("failed to search object: %w", err)
		}
		defer resp.Body.Close()

		if resp.IsError() {
			return fmt.Errorf("failed to search object status: %s", resp.Status())
		}

		// TODO: decode

		return nil
	}

	for _, ope := range []func() error{
		indexing,
		search,
	} {
		if err := ope(); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
