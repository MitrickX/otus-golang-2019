package notificaiton

import (
	"errors"
	"github.com/mitrickx/otus-golang-2019/25/calendar/internal/domain/entities"
	"io"
)

var ErrQueueEmpty = errors.New("queue is empty")

// Simple queue interface
type Queue interface {
	Push(event entities.Event) error    // push main info about biz event entity into queue
	Consume() (<-chan EventInfo, error) // subscribe on event info channel, from where event info items will be read
	io.Closer
}
