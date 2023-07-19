package bus

import (
	"btcRate/common/infrastructure"
	"btcRate/common/infrastructure/bus/command_handlers"
	"btcRate/common/infrastructure/bus/commands"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
	"os"
	"time"
)

type RabbitMQConfig struct {
	Host     string
	User     string
	Password string
}

func AddCommandBus(busConfig *RabbitMQConfig) (*cqrs.CommandBus, *message.Router) {
	cqrsMarshaler := cqrs.JSONMarshaler{}
	logger := watermill.NewStdLoggerWithOut(os.Stdout, true, true)

	commandsAMQPConfig := amqp.NewDurableQueueConfig(fmt.Sprintf("amqp://%s:%s@%s/", busConfig.User, busConfig.Password, busConfig.Host))
	commandsAMQPConfig.Exchange.GenerateName = func(topic string) string {
		return "btc-rate_topic"
	}
	commandsAMQPConfig.Exchange.Type = "topic"
	commandsAMQPConfig.QueueBind.GenerateRoutingKey = func(topic string) string {
		return topic
	}

	var commandsPublisher *amqp.Publisher
	var commandsSubscriber *amqp.Subscriber
	var err error

	commandsPublisher, err = amqp.NewPublisher(commandsAMQPConfig, logger)
	if err == nil {
		commandsSubscriber, err = amqp.NewSubscriber(commandsAMQPConfig, logger)
	}

	for i := 0; i < 10; i++ {
		var err error

		commandsPublisher, err = amqp.NewPublisher(commandsAMQPConfig, logger)
		if err == nil {
			commandsSubscriber, err = amqp.NewSubscriber(commandsAMQPConfig, logger)
		}

		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %s. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	if commandsPublisher == nil || commandsSubscriber == nil {
		panic("Failed to connect to RabbitMQ after several attempts")
	}

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	commandBus, err := cqrs.NewCommandBusWithConfig(
		commandsPublisher,
		cqrs.CommandBusConfig{
			GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
				// Custom routing for LogCommands
				if logCommand, ok := params.Command.(*commands.LogCommand); ok {
					return fmt.Sprintf("%s.%s", params.CommandName, logCommand.LogLevel), nil
				}
				return params.CommandName, nil
			},
			Marshaler: cqrsMarshaler,
		})
	if err != nil {
		panic(err)
	}

	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(
		router,
		cqrs.CommandProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
				switch params.CommandHandler.(type) {
				case command_handlers.LogCommandHandler:
					return fmt.Sprintf("%s.*", params.CommandName), nil
				case command_handlers.ErrorCommandHandler:
					return fmt.Sprintf("%s.%s", params.CommandName, infrastructure.LogLevelError), nil
				default:
					return params.CommandName, nil
				}
			},
			SubscriberConstructor: func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return commandsSubscriber, nil
			},
			Marshaler: cqrsMarshaler,
		},
	)
	if err != nil {
		panic(err)
	}

	err = commandProcessor.AddHandlers(
		command_handlers.ErrorCommandHandler{},
	)
	if err != nil {
		panic(err)
	}

	return commandBus, router
}
