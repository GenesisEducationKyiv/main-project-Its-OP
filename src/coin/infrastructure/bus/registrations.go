package bus

import (
	"btcRate/common/infrastructure"
	"btcRate/common/infrastructure/bus/command_handlers"
	"btcRate/common/infrastructure/bus/commands"
	"btcRate/common/infrastructure/bus/decorators"
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

func AddCommandBus(busConfig *RabbitMQConfig) (*cqrs.CommandBus, *message.Router, error) {
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
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ after several attempts")
	}

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(
		router,
		cqrs.CommandProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
				switch params.CommandHandler.HandlerName() {
				case command_handlers.LogCommandHandlerName:
					return fmt.Sprintf("%s.*", params.CommandName), nil
				case command_handlers.ErrorLogCommandHandlerName:
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
		return nil, nil, err
	}

	err = commandProcessor.AddHandlers(
		decorators.NewLoggedCommandHandler(command_handlers.LogCommandHandler{}, cqrsMarshaler.Name),
	)
	if err != nil {
		return nil, nil, err
	}

	return commandBus, router, nil
}
