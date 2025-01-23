package main

import (
	"fmt"
	"listener-service/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)


func main() {
	// try to connect to rabbitmq


	rabbitConn, err := connect()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}


	defer rabbitConn.Close()

	
	
	// strat listening for messages
	log.Println("Listening for messages and consuming RabbitMQ messages ...")


	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)

	if err != nil {
		panic(err)
	}

	


	// watch the queue and consume event
	err = consumer.Listen([]string{"log.Info", "log.Warning", "log.ERROR"})

	if err != nil {
		log.Println(err)
	}
}


func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection


	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitmMq not redey....")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break 
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off connection")
		time.Sleep(backoff)
		continue
	}

	return connection, nil
}