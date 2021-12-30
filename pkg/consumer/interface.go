package consumer

type Consumer interface {
	Receive(msg string) error
}
