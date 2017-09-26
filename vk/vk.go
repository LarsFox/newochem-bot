package vk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Client works with VK API
type Client interface {
}

type client struct {
	token   string
	version string
}

// NewClient returns a new client to work with VK API
func NewClient(token, version string) Client {
	return &client{token: token, version: version}
}

func (c *client) GetWall(offset int) (*GetWallResponse, error) {
	url := c.url("wall.get", map[string]interface{}{
		"owner_id": newochemID,
		"offset":   offset,
		"count":    vkCount,
	})
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	target := &GetWallResponse{}
	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return target, nil
}

func (c *client) url(method string, params map[string]interface{}) string {
	url := fmt.Sprintf(vkMethod, method, c.token, c.version)
	for key, value := range params {
		url += fmt.Sprintf("&%s=%v", key, value)
	}
	return url
}
