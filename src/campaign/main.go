package main

import (
	"btcRate/campaign/web"
	"btcRate/common/infrastructure"
	"btcRate/common/infrastructure/bus/command_handlers"
	"btcRate/common/infrastructure/bus/commands"
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	commandBus, router := addCommandBus(os.Getenv("KAFKA_HOST"))
	go func() {
		_ = router.Run(context.Background())
	}()

	server := web.NewServerManager()

	fc := &web.FileConfiguration{EmailStorageFile: "./data/emails.json"}
	sc := &web.SendgridConfiguration{ApiKey: os.Getenv("SENDGRID_KEY"), SenderName: os.Getenv("SENDGRID_SENDER_NAME"), SenderEmail: os.Getenv("SENDGRID_SENDER_EMAIL")}
	pc := &web.ProviderConfiguration{Hostname: os.Getenv("COIN_HOST"), Schema: os.Getenv("COIN_SCHEMA")}

	stop, err := server.RunServer(fc, sc, pc, commandBus)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	if err := stop(); err != nil {
		log.Fatal("Failed to stop the server: ", err)
	}
}

func addCommandBus(messageBusHost string) (*cqrs.CommandBus, *message.Router) {
	cqrsMarshaler := cqrs.JSONMarshaler{}
	logger := watermill.NewStdLoggerWithOut(os.Stdout, false, false)
	commandsKafkaConfig := kafka.DefaultSaramaSyncPublisherConfig()

	commandsKafkaConfig.Producer.Return.Successes = true
	commandsKafkaConfig.Net.SASL.Enable = true // add additional Kafka configuration as needed

	var commandsPublisher *kafka.Publisher
	var commandsSubscriber *kafka.Subscriber
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
				ConsumerGroup:         "campaign-consumer-group", // replace with your Kafka consumer group
			},
			logger,
		)
	}

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
					ConsumerGroup:         "campaign-consumer-group", // replace with your Kafka consumer group
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
		command_handlers.LogCommandHandler{},
	)
	if err != nil {
		panic(err)
	}

	return commandBus, router
}
