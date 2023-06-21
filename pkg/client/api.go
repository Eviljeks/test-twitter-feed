package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const saveMessageURI = "/messages"

var ErrResponseStatusCodeError = errors.New("response status code >= 400")

type APIClient struct {
	client      *http.Client
	apiBasePath string
}

func NewAPIClient(client *http.Client, apiBasePath string) *APIClient {
	return &APIClient{
		client:      client,
		apiBasePath: apiBasePath,
	}
}

func (c *APIClient) SaveMessage(content string) error {
	body, err := json.Marshal(struct {
		Content string `json:"content"`
	}{
		Content: content,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.path(saveMessageURI), bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return ErrResponseStatusCodeError
	}

	return nil
}

func (c *APIClient) path(uri string) string {
	return c.apiBasePath + uri
}
