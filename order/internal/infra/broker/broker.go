package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"order_service/internal/config"
	"order_service/internal/infra"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	logger *zap.Logger
	conn   *amqp.Connection
	ch     *amqp.Channel
}

const OrderEventExchange = "order.events"

func New(log *zap.Logger, cfg *config.RabbitMQ) (*RabbitMQ, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("amqp.Dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("conn.Channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		OrderEventExchange,
		"topic",
		true,  // durable
		false, // auto-delete
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("ch.ExchangeDeclare: %w", err)
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
				true,  // durable
				false, // auto-delete
				false,
				false,
				nil,
			)
			if err != nil {
				return nil, err
			}

			// Привязка очереди к exchange и routing key
			if err := ch.QueueBind(
				queueName,
				routingKey,
				OrderEventExchange,
				false,
				nil,
			); err != nil {
				return nil, err
			}
		}
	}

	log.Info("Connect to RabbitMQ success")

	return &RabbitMQ{
		conn:   conn,
		ch:     ch,
		logger: log,
	}, nil
}

func (r *RabbitMQ) Close() error {
	return r.conn.Close()
}

func (r *RabbitMQ) PublishCreatedOrder(ctx context.Context, order *infra.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	return r.ch.Publish(
		OrderEventExchange,
		"order.created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
