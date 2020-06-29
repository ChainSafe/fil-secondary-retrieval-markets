package client

type Client struct {
	host Host
}

func NewClient(host Host) *Client {
	return &Client{host: host}
}