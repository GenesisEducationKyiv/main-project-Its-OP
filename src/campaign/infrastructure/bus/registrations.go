package bus

import (
	"btcRate/common/infrastructure"
	"btcRate/common/infrastructure/bus/command_handlers"
	"btcRate/common/infrastructure/bus/commands"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
	"os"
	"time"
)

func AddCommandBus(messageBusHost string, consumerGroup string) (*cqrs.CommandBus, *message.Router) {
	cqrsMarshaler := cqrs.JSONMarshaler{}
	logger := watermill.NewStdLoggerWithOut(os.Stdout, false, false)
	commandsKafkaConfig := kafka.DefaultSaramaSyncPublisherConfig()

	var commandsPublisher *kafka.Publisher
	var commandsSubscriber *kafka.Subscriber
	var err error

	for i := 0; i < 10; i++ {
		var err error

		commandsPublisher, err = kafka.NewPublisher(
			kafka.PublisherConfig{
				Brokers:   []string{messageBusHost},
				Marshaler: kafka.DefaultMarshaler{},
			},
			logger,
		)
		if err == nil {
			commandsSubscriber, err = kafka.NewSubscriber(
				kafka.SubscriberConfig{
					Brokers:               []string{messageBusHost},
					Unmarshaler:           kafka.DefaultMarshaler{},
					OverwriteSaramaConfig: commandsKafkaConfig,
					ConsumerGroup:         consumerGroup,
				},
				logger,
			)
		}

		if err != nil {
			log.Printf("Failed to connect to Kafka: %s. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	if commandsPublisher == nil || commandsSubscriber == nil {
		panic("Failed to connect to Kafka after several attempts")
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