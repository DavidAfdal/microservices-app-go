package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {

	rabbitConn, err := connect()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}


	defer rabbitConn.Close()
	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s", webPort)

	srv := &http.Server{
		Addr : fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
   
	if err :=  srv.ListenAndServe(); err != nil {
		log.Panicf("Failed to listen :%s", webPort)
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