package messagebroker

type Publisher struct {
	broker *Broker
	ID     string
}

func NewPublisher(broker *Broker, id string) *Publisher {
	return &Publisher{broker: broker, ID: id}
}

func (p *Publisher) Publish(topic string, payload interface{}) {
	if !p.broker.IsPublisherRegistered(p.ID) {
		return // Do nothing if the publisher is not registered
	}
	msg := Message{
		Topic:   topic,
		Payload: payload,
	}
	p.broker.Publish(msg)
}
