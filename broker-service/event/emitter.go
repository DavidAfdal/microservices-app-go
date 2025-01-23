package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)


type Emmitter struct {
	connection *amqp.Connection
}

func (e *Emmitter) setup() error {
	channel, err := e.connection.Channel()

	if err != nil {
		return err
	}

	defer channel.Close()

	return declareEnxchange(channel)
}


func (e *Emmitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()

	if err != nil {
		return err
	}

	defer channel.Close()
	

	log.Println("Pubhlising to channel")


   err = channel.Publish(
	 	"logs_topics",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(event),
		},
   )

   if err != nil {
	return err
   }

   return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emmitter, error) {
	emitter := Emmitter{
		connection: conn,
	}

	err := emitter.setup()

	if err != nil {
		return Emmitter{}, err
	}

	return emitter, nil
}