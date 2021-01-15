package delay_server

type (
	// delayClient 延迟client
	delayClient struct {
		brokers []string
	}
	// COption delayClient Func
	COption func(*delayClient)
)

func NewDelayClient(options []COption) (*delayClient, error) {
	client := new(delayClient)

	for _, f := range options {
		f(client)
	}
	return client, client.validate()
}

// WithBrokerList 发送到的broker地址
func (c *delayClient) WithBrokerList(brokers []string) COption {
	return func(c *delayClient) {
		c.brokers = brokers
	}
}

func (c *delayClient) validate() error {
	return nil
}
