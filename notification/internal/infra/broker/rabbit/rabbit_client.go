package rabbit

type Client struct {
}

func NewRabbitClient() *Client {
	return &Client{}
}

func (r *Client) StartConsuming() error {
	//TODO implement me
	panic("implement me")
}
