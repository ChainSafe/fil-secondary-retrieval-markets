package provider

type Provider struct {
	host Host
}

func NewClient(host Host) *Provider {
	return &Provider{host: host}
}