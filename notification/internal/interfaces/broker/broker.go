package broker

type Broker interface {
	StartConsuming() error
}
