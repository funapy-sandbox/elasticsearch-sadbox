package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/uuid"
)

const (
	indexName = "user"
)

// User represents user.
type User struct {
	ID        string `json:"id"`
	ChatbotID string `json:"chatbot_id"`
	Tags      []Tag  `json:"tags" elastic:"type:nested"`
}

// Tag --
type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func testData() []User {
	return []User{
		{
			ID:        "u1",
			ChatbotID: "ch1",
			Tags: []Tag{
				{
					ID:   "t1",
					Name: "t1_name",
				},
			},
		},
		{
			ID:        "u2",
			ChatbotID: "ch1",
			Tags: []Tag{
				{
					ID:   "t1",
					Name: "t1_name",
				},
				{
					ID:   "t2",
					Name: "t2_name",
				},
			},
		},
		{
			ID:        "u3",
			ChatbotID: "ch3",
			Tags: []Tag{
				{
					ID:   "t3",
					Name: "t3_name",
				},
			},
		},
	}
}

func run(ctx context.Context) error {
	c, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	indexing := func() error {
		for _, user := range testData() {
			b, err := json.Marshal(user)
			if err != nil {
				return fmt.Errorf("failed to marshal object: %w", err)
			}

			id := uuid.NewString()
			req := esapi.IndexRequest{
				Index:      indexName,
				DocumentID: id,
				Refresh:    "true",
				Body:       bytes.NewReader(b),
			}

			resp, err := req.Do(ctx, c)
			if err != nil {
				return fmt.Errorf("failed to create index: %w", err)
			}
			defer resp.Body.Close()

			if resp.IsError() {
				return fmt.Errorf("failed to create index id: %s  status: %s", id, resp.Status())
			}

		}

		return nil
	}

	search := func() error {
		var buf bytes.Buffer
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"chatbot_id": "ch1",
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

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Println(string(b))

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
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
