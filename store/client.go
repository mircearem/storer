package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Maybe this should be a new function in the storer package
type Client struct {
	remoteAddr string
	client     *http.Client
}

func NewClient(remoteAddr string) *Client {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	return &Client{
		remoteAddr: remoteAddr,
		client:     client,
	}
}

// Get values from the store using keys
func (c *Client) Get(collection string, key string) (string, error) {
	uri := fmt.Sprintf("http://%s/api/%s", c.remoteAddr, collection)
	data := map[string]string{"key": key}
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	body := bytes.NewBuffer(b)
	req, err := http.NewRequest("GET", uri, body)
	if err != nil {
		return "", err
	}
	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result := make(map[string]string)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result["value"], nil
}

// Insert key value pairs into the store
func (c *Client) Post(collection string, data []byte) error {
	// Check if data is valid json format
	jmap := make(map[string]string)
	if err := json.Unmarshal(data, &jmap); err != nil {
		return err
	}
	// Data is valit json format: {"key":"value"}
	uri := fmt.Sprintf("http://%s/api/%s", c.remoteAddr, collection)
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return err
	}
	// Execute request
	if _, err = c.client.Do(req); err != nil {
		return err
	}
	return nil
}
