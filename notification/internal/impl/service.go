package impl

import (
	"notification_service/internal/infra/broker"
	"notification_service/internal/interfaces"

	"go.uber.org/zap"
)

type service struct {
	logger *zap.Logger
	broker *broker.RabbitMQ
}

func New(logger *zap.Logger, broker *broker.RabbitMQ) interfaces.Service {
	return &service{logger: logger, broker: broker}
}

func (s *service) HandleOrderCreatedMessages() error {
	msgs, err := s.broker.GetChannel().Consume(
		"queue_order_created",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		s.logger.Error("failed to consume messages", zap.Error(err))
		return err
	}

	for msg := range msgs {
		s.logger.Info("received created order", zap.String("message", string(msg.Body)))
	}
	return nil
}

func (s *service) HandleOrderCancelledMessages() error {
	msgs, err := s.broker.GetChannel().Consume(
		"queue_order_cancelled",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		s.logger.Error("failed to consume messages", zap.Error(err))
		return err
	}

	for msg := range msgs {
		s.logger.Info("received cancelled order", zap.String("message", string(msg.Body)))
	}
	return nil
}

func (s *service) HandleOrderCompletedMessages() error {
	msgs, err := s.broker.GetChannel().Consume(
		"queue_order_completed",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		s.logger.Error("failed to consume messages", zap.Error(err))
		return err
	}

	for msg := range msgs {
		s.logger.Info("received completed order", zap.String("message", string(msg.Body)))
	}
	return nil
}
