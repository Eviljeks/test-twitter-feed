package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const saveMessageURI = "/messages"

var ErrResponseStatusCodeError = errors.New("response status code >= 400")

type ApiClient struct {
	client      *http.Client
	apiBasePath string
}

func NewApiClient(client *http.Client, apiBasePath string) *ApiClient {
	return &ApiClient{
		client:      client,
		apiBasePath: apiBasePath,
	}
}

func (c *ApiClient) SaveMessage(content string) error {
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

func (c *ApiClient) path(uri string) string {
	return c.apiBasePath + uri
}
