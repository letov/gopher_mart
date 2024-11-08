package queue

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopher_mart/internal/infrastructure/config"
	"log"
	"os"
)

type Rabbit struct {
	l    *zap.SugaredLogger
	conn *amqp.Connection
	ch   *amqp.Channel
	q    *amqp.Queue
}

func (r Rabbit) Publish(ctx context.Context, orderID string) error {
	return r.ch.PublishWithContext(ctx,
		"",
		r.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(orderID),
		},
	)
}

func (r Rabbit) RegisterHandler(h RequestAccrualHandler) error {
	messages, err := r.ch.Consume(
		r.q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			orderID := string(msg.Body)
			if err != nil {
				r.l.Error(err)
				return
			}
			err = h(orderID)
			if err != nil {
				r.l.Warn(err)
				return
			}
			err = msg.Ack(false)
			if err != nil {
				r.l.Error(err)
				return
			}
		}
	}()

	return nil
}

func NewRabbit(
	lc fx.Lifecycle,
	config *config.Config,
	l *zap.SugaredLogger,
) *Rabbit {
	conn, err := amqp.Dial(config.Ampq)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	q, err := ch.QueueDeclare(
		"request_accrual",
		true,
		os.Getenv("IS_TEST_ENV") == "true",
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	b := &Rabbit{
		l:    l,
		conn: conn,
		ch:   ch,
		q:    &q,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			_ = ch.Close()
			_ = conn.Close()
			return nil
		},
	})

	return b
}
