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
	logger := watermill.NewStdLoggerWithOut(os.Stdout, true, true)
	cqrsMarshaler := cqrs.ProtobufMarshaler{}
	commandsAMQPConfig := amqp.NewDurableQueueConfig("amqp://admin:admin@localhost:5672/")

	commandsPublisher, err := amqp.NewPublisher(commandsAMQPConfig, logger)
	if err != nil {
		panic(err)
	}

	commandsSubscriber, err := amqp.NewSubscriber(commandsAMQPConfig, logger)
	if err != nil {
		panic(err)
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
