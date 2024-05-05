package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/sikozonpc/commons/api"
	"github.com/sikozonpc/commons/broker"
	"github.com/sikozonpc/omsv2-kitchen/gateway"
	"go.opentelemetry.io/otel"
)

type Consumer struct {
	gateway gateway.KitchenGateway
}

func NewConsumer(gateway gateway.KitchenGateway) *Consumer {
	return &Consumer{gateway}
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
			// Create a new span
			tr := otel.Tracer("amqp")
			_, messageSpan := tr.Start(context.Background(), fmt.Sprintf("AMQP - consume - %s", q.Name))

			var o *pb.Order
			if err := json.Unmarshal(d.Body, &o); err != nil {
				log.Printf("Error unmarshalling order: %v", err)
				d.Nack(false, false)
				continue
			}

			if o.Status == "paid" {
				cookOrder() // let him cook

				messageSpan.AddEvent(fmt.Sprintf("Order Cooked: %v", o))

				if err := c.gateway.UpdateOrder(context.Background(), &pb.Order{
					Status:     "ready",
					ID:         o.ID,
					CustomerID: o.CustomerID,
				}); err != nil {
					log.Printf("error updating the order %v", o)

					if err := broker.HandleRetry(ch, &d); err != nil {
						log.Printf("error handling the retry: %v", err.Error())
					}
					continue
				}
			}

			messageSpan.AddEvent(fmt.Sprintf("order.updated: %v", o))
			messageSpan.End()

			d.Ack(false)
		}
	}()

	log.Printf("AMQP Listening. To exit press CTRL+C")
	<-forever
}

func cookOrder() {
	log.Println("Cooking order...")
	time.Sleep(5 * time.Second)
	log.Println("Order cooked!")
}
