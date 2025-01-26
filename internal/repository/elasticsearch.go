package repository

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"irm_backend/internal/config"
	"irm_backend/internal/models"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
)

type ElasticsearchDB struct {
	client *elasticsearch.Client
}

func NewElasticsearchDB(cfg *config.Config) (*ElasticsearchDB, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.ElasticsearchURL},
		Username:  cfg.ElasticsearchUsername,
		Password:  cfg.ElasticsearchPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating elasticsearch client: %v", err)
	}

	return &ElasticsearchDB{client: es}, nil
}

// Helper function to create index if not exists
func (es *ElasticsearchDB) createIndexIfNotExists(index string) error {
	res, err := es.client.Indices.Exists([]string{index})
	if err != nil {
		return err
	}
	if res.StatusCode == 404 {
		_, err = es.client.Indices.Create(index)
		if err != nil {
			return err
		}
	}
	return nil
}

// User repository methods
func (es *ElasticsearchDB) CreateUser(user *models.User) error {
	if err := es.createIndexIfNotExists("users"); err != nil {
		return err
	}

	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = es.client.Index(
		"users",
		strings.NewReader(string(data)),
		es.client.Index.WithDocumentID(user.ID.String()),
		es.client.Index.WithRefresh("true"),
	)
	return err
}

func (es *ElasticsearchDB) GetUserByEmail(email string) (*models.User, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"email": email,
			},
		},
	}

	data, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	res, err := es.client.Search(
		es.client.Search.WithIndex("users"),
		es.client.Search.WithBody(strings.NewReader(string(data))),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		Hits struct {
			Hits []struct {
				Source models.User `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Hits.Hits) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &result.Hits.Hits[0].Source, nil
}

// Post repository methods
func (es *ElasticsearchDB) CreatePost(post *models.Post) error {
	if err := es.createIndexIfNotExists("posts"); err != nil {
		return err
	}

	post.ID = uuid.New()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	data, err := json.Marshal(post)
	if err != nil {
		return err
	}

	_, err = es.client.Index(
		"posts",
		strings.NewReader(string(data)),
		es.client.Index.WithDocumentID(post.ID.String()),
		es.client.Index.WithRefresh("true"),
	)
	return err
}

func (es *ElasticsearchDB) GetUserPosts(userID string) ([]models.Post, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"user_id": userID,
			},
		},
	}

	data, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	res, err := es.client.Search(
		es.client.Search.WithIndex("posts"),
		es.client.Search.WithBody(strings.NewReader(string(data))),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		Hits struct {
			Hits []struct {
				Source models.Post `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	posts := make([]models.Post, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		posts[i] = hit.Source
	}

	return posts, nil
}

func (es *ElasticsearchDB) UpdatePost(post *models.Post) error {
	post.UpdatedAt = time.Now()
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}

	_, err = es.client.Index(
		"posts",
		strings.NewReader(string(data)),
		es.client.Index.WithDocumentID(post.ID.String()),
		es.client.Index.WithRefresh("true"),
	)
	return err
}

func (es *ElasticsearchDB) DeletePost(postID string) error {
	_, err := es.client.Delete("posts", postID)
	return err
}
