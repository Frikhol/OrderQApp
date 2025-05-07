package interfaces

type Service interface {
	HandleOrderCreatedMessages() error
	HandleOrderCancelledMessages() error
	HandleOrderCompletedMessages() error
}
