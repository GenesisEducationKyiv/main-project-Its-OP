package rabbitmq

import (
	"btcRate/common/infrastructure/bus"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type Consumer struct {
	channel  *amqp.Channel
	exchange string
	handlers []bus.ICommandHandler
}

func NewConsumer(channel *amqp.Channel, exchange string, handlers []bus.ICommandHandler) *Consumer {
	return &Consumer{channel: channel, exchange: exchange, handlers: handlers}
}

func (c *Consumer) StartConsuming() {
	for _, handler := range c.handlers {
		go c.consume(handler)
	}
}

func (c *Consumer) consume(handler bus.ICommandHandler) {
	q, _ := c.channel.QueueDeclare(
		fmt.Sprintf("%s.%s", c.exchange, handler.GetName()), // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)

	c.channel.QueueBind(
		q.Name,            // queue name
		handler.GetName(), // routing key
		c.exchange,        // exchange
		false,
		nil)

	messages, _ := c.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for msg := range messages {
		var cmd bus.Command
		json.Unmarshal(msg.Body, &cmd)
		handler.Handle(cmd)
	}
}
