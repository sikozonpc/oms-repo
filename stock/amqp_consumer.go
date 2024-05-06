package main

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sikozonpc/commons/broker"
	"go.opentelemetry.io/otel"
)

type Consumer struct{}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(
		q.Name,                // queue name
		"",                    // routing key
		broker.OrderPaidEvent, // exchange
		false,                 // no-wait
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			// Extract headers
			ctx := broker.ExtractAMQPHeader(context.Background(), d.Headers)

			// Create a new span
			tr := otel.Tracer("amqp")
			_, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - consume - %s", q.Name))

			log.Printf("Received a message: %s", d.Body)

			orderID := string(d.Body)

			d.Ack(false)

			messageSpan.End()
			log.Printf("Order received: %s", orderID)
		}
	}()

	log.Printf("AMQP Listening. To exit press CTRL+C")
	<-forever
}
