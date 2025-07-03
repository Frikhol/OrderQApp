package rabbit

import (
	"context"
	"github.com/streadway/amqp"
	"log"
	"notification_service/internal/config"
)

type Client struct {
	url   string
	query string
}

func NewRabbitClient(cfg config.RabbitMQ) *Client {
	return &Client{cfg.URL, cfg.Query}
}

func (r *Client) StartConsuming(ctx context.Context) error {
	conn, err := amqp.Dial(r.url)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(r.query, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-msgs:
				log.Printf("Received a message: %s", msg.Body)
				//TODO: unmarshal and notice users
			case <-ctx.Done():
				//FIXME:??
				//ch.Close()
				//conn.Close()
				return
			}
		}
	}()

	return nil
}
