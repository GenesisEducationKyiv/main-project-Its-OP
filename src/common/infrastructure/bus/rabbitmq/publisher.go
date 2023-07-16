package rabbitmq

import (
	"btcRate/common/infrastructure/bus"
	"encoding/json"
	"github.com/streadway/amqp"
)

type IPublisher interface {
	Publish(queue string, cmd bus.Command) error
}

type Publisher struct {
	channel  *amqp.Channel
	exchange string
}

func NewPublisher(channel *amqp.Channel, exchange string) IPublisher {
	return &Publisher{channel: channel, exchange: exchange}
}

func (p *Publisher) Publish(routingKey string, cmd bus.Command) error {
	body, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		p.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
