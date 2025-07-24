package rabbit

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"notification_service/internal/config"
	"notification_service/internal/domain/types"
	"notification_service/internal/interfaces/app"
)

type Client struct {
	url      string
	notifier app.Notifier
}

const OrderEventExchange = "order.events"

func NewRabbitClient(cfg config.RabbitMQ, notifier app.Notifier) *Client {
	return &Client{cfg.URL, notifier}
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

	if err := ch.ExchangeDeclare(
		OrderEventExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("ch.ExchangeDeclare: %w", err)
	}

	var orderQueues = map[string][]string{
		"order.created":   {"queue_order_created"},
		"order.assigned":  {"queue_order_assigned"},
		"order.accepted":  {"queue_order_accepted"},
		"order.cancelled": {"queue_order_cancelled"},
		"order.completed": {"queue_order_completed"},
		"order.updated":   {"queue_order_updated"},
	}

	for routingKey, queues := range orderQueues {
		for _, queueName := range queues {
			_, err := ch.QueueDeclare(
				queueName,
				true,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				return err
			}

			if err := ch.QueueBind(
				queueName,
				routingKey,
				OrderEventExchange,
				false,
				nil,
			); err != nil {
				return err
			}
		}
	}

	for _, queues := range orderQueues {
		for _, queueName := range queues {
			msgs, err := ch.Consume(
				queueName,
				"",
				true,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				return fmt.Errorf("consume from %s: %w", queueName, err)
			}

			go func(queue string, m <-chan amqp.Delivery) {
				for msg := range m {
					log.Printf("[queue: %s] received message: %s", queue, msg.Body)
					newMsg := types.WSMessage{
						Type:    "test type",
						Payload: msg.Body,
					}
					log.Printf("[queue: %s] sending message: %s to user: %v", queue, newMsg, msg.UserId)
					r.notifier.Notify(msg.UserId, newMsg)
				}
			}(queueName, msgs)
		}
	}

	return nil
}
