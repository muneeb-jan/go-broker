package messagebroker

type Publisher struct {
	broker *Broker
}

func NewPublisher(broker *Broker) *Publisher {
	return &Publisher{broker: broker}
}

func (p *Publisher) Publish(topic string, payload interface{}) {
	msg := Message{
		Topic:   topic,
		Payload: payload,
	}
	p.broker.Publish(msg)
}
