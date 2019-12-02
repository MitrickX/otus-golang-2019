package scheduler

import (
	"github.com/mitrickx/otus-golang-2019/25/calendar/internal/domain/entities"
	"github.com/mitrickx/otus-golang-2019/25/calendar/internal/storage/memory"
	"testing"
	"time"
)

type channelQueue struct {
	ch chan entities.Event
}

func (c *channelQueue) AddEvent(event entities.Event) error {
	c.ch <- event
}

func (c *channelQueue) ReadEvent() (entities.Event, error) {
	return <-c.ch, nil
}

func newChannelQueue(size int) *channelQueue {
	ch := make(chan entities.Event, size)
	return &channelQueue{ch}
}

func newScheduler() *Scheduler {
	s := &Scheduler{
		scanTimeout: 200 * time.Millisecond,
		queue:       newChannelQueue(100),
		storage:     memory.NewStorage(),
	}
	return s
}

func TestSchedulerRun(t *testing.T) {

}
