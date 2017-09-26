package vk

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
