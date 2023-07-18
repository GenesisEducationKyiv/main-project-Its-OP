package main

import (
	"btcRate/coin/web"
	"btcRate/common/infrastructure/bus/command_handlers"
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	commandBus, router := addCommandBus()
	go func() {
		_ = router.Run(context.Background())
	}()

	server := web.NewServerManager(commandBus)
	stop, err := server.RunServer("./logs/coin-logs.csv")
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	if err := stop(); err != nil {
		log.Fatal("Failed to stop the server: ", err)
	}

	if err := router.Close(); err != nil {
		log.Fatal("Failed to stop the Router: ", err)
	}
}

func addCommandBus() (*cqrs.CommandBus, *message.Router) {
	cqrsMarshaler := cqrs.ProtobufMarshaler{}
	logger := watermill.NewStdLoggerWithOut(os.Stdout, false, false)
	commandsAMQPConfig := amqp.NewDurableQueueConfig("amqp://admin:admin@rabbitmq:5672/")

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
				return params.CommandName, nil
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
