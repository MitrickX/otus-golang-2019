package notificaiton

import (
	"fmt"
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/domain/entities"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"io"
	"strconv"
	"time"
)

// Config for rabbit-mq Queue implementation
type Config struct {
	Host           string
	Port           string
	QName          string
	User           string
	Password       string
	ConnectRetries int
}

// Config constructor
func NewConfig(m map[string]string) (*Config, error) {
	keys := []string{"host", "port", "qname", "user", "password"}
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return nil, fmt.Errorf("`%s` key is missing", key)
		}
	}

	connectRetries, err := strconv.Atoi(m["connect_retries"])
	if err != nil {
		return nil, fmt.Errorf("rabbit connect_retries key error %w", err)
	}

	return &Config{
		Host:           m["host"],
		Port:           m["port"],
		QName:          m["qname"],
		User:           m["user"],
		Password:       m["password"],
		ConnectRetries: connectRetries,
	}, nil

}

// Rabbit mq Queue implementation
type Rabbit struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	eventsCh chan EventInfo
	q        *amqp.Queue
	logger   *zap.SugaredLogger
}

// Constructor
func NewRabbitQueue(cfg Config) (*Rabbit, error) {

	rabbit := &Rabbit{}

	var conn *amqp.Connection
	var connError error

	for i := 0; i < cfg.ConnectRetries; i++ {
		conn, connError = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.Host, cfg.Port))
		if connError == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if connError != nil {
		return nil, connError
	}

	rabbit.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		_ = rabbit.Close()
		return nil, err
	}

	rabbit.ch = ch

	q, err := ch.QueueDeclare(
		"events",
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		_ = rabbit.Close()
		return nil, err
	}

	rabbit.q = &q

	rabbit.eventsCh = make(chan EventInfo, 100)

	return rabbit, nil
}

// Push main info about biz event entity into Queue
func (r *Rabbit) Push(event entities.Event) error {

	msg, err := serializeEvent(extractEventInfo(event))
	if err != nil {
		return err
	}

	err = r.ch.Publish(
		"",       // exchange
		r.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// Subscribe on event info channel, from where event info items
func (r *Rabbit) Consume() (<-chan EventInfo, error) {
	messages, err := r.ch.Consume(
		r.q.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)

	if err != nil {
		return nil, err
	}

	go func() {
		for msg := range messages {
			eventInfo := &EventInfo{}
			err := unSerializeEvent(msg.Body, eventInfo)
			if err != nil {
				r.logErrorf("Rabbit.Consume, error when unserialize event info: %w", err)
			} else {
				r.eventsCh <- *eventInfo
			}
		}
	}()

	return r.eventsCh, nil
}

// Close Rabbit Queue implementation
func (r *Rabbit) Close() error {
	conn := r.conn
	ch := r.ch
	r.conn = nil
	r.ch = nil
	return closeAll(ch, conn)
}

// Try close closers one by one until all done of first error happened
func closeAll(closers ...io.Closer) error {
	for _, closer := range closers {
		err := closer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// log formatted error into log
func (r *Rabbit) logErrorf(format string, err error) {
	if r.logger != nil {
		r.logger.Error(fmt.Errorf(format, err))
	}
}
