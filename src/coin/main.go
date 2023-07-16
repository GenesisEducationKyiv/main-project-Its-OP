package main

import (
	"btcRate/coin/web"
	"btcRate/common/infrastructure/bus/command_handlers"
	"context"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
	"os"
	"os/signal"
)

func main() {
	commandBus := addCommandBus()
	server := web.NewServerManager()
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
}

func addCommandBus() *cqrs.CommandBus {
	cqrsMarshaler := cqrs.ProtobufMarshaler{}
	commandsAMQPConfig := amqp.NewDurableQueueConfig("amqp://localhost:5672")

	commandsPublisher, err := amqp.NewPublisher(commandsAMQPConfig, nil)
	if err != nil {
		panic(err)
	}

	commandsSubscriber, err := amqp.NewSubscriber(commandsAMQPConfig, nil)
	if err != nil {
		panic(err)
	}

	router, err := message.NewRouter(message.RouterConfig{}, nil)
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

	// processors are based on router, so they will work when router will start
	if err := router.Run(context.Background()); err != nil {
		panic(err)
	}

	return commandBus
}
