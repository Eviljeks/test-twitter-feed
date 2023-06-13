package amqp

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type SingleQueueAMQPPublisher struct {
	ch        *amqp.Channel
	queueName string
}

func NewSingleQueueAMQPPublisher(ch *amqp.Channel, queueName string) *SingleQueueAMQPPublisher {
	return &SingleQueueAMQPPublisher{
		ch:        ch,
		queueName: queueName,
	}
}

func (p *SingleQueueAMQPPublisher) Publish(data interface{}) error {
	mData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = p.ch.Publish(
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        mData,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *SingleQueueAMQPPublisher) Setup() error {
	_, err := p.ch.QueueDeclare(
		p.queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
