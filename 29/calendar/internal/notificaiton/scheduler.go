package notificaiton

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mitrickx/otus-golang-2019/29/calendar/internal/domain/entities"
	"go.uber.org/zap"
)

var ErrorQueueNotInitialized = errors.New("queue not initialized")
var ErrorStorageNotInitialized = errors.New("storage not initialized")

// Notification scheduler
// Scan storage with some freq and put events info into queue
// Once event info pushed into queue even mark as notified
type Scheduler struct {
	scanTimeout time.Duration    // frequency of scan
	storage     entities.Storage // calendar storage
	start       *time.Time       // start of interval for get events by interval for notification
	logger      *zap.SugaredLogger
	queue       Queue
	nowTimeFn   func() time.Time // for possibility to redeclare in tests
	cancelFn    context.CancelFunc
}

// Constructor
func NewScheduler(scanTimeout time.Duration, storage entities.Storage, queue Queue, logger *zap.SugaredLogger) *Scheduler {
	return &Scheduler{
		scanTimeout: scanTimeout,
		storage:     storage,
		logger:      logger,
		queue:       queue,
	}
}

// Run scheduler
func (s *Scheduler) Run() error {

	if s.queue == nil {
		return ErrorQueueNotInitialized
	}

	if s.storage == nil {
		return ErrorStorageNotInitialized
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	s.cancelFn = cancelFn

	s.run(ctx)

	return nil
}

// Stop scheduler
func (s *Scheduler) Stop() {
	if s.cancelFn != nil {
		s.cancelFn()
	}
}

// Inner helper that run scan process. Take into account context.Done()
func (s *Scheduler) run(ctx context.Context) {
	ticker := time.NewTicker(s.scanTimeout)
	s.scan()
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			err := s.queue.Close()
			if err != nil {
				s.logErrorf("error happened on close queue %w\n", err)
			}
			return
		case <-ticker.C:
			s.scan()
		}
	}
}

// scan db to find events to notify about
// push event info into queue
// once event info pushed into queue mark event as notified
func (s *Scheduler) scan() {
	var start, end *entities.DateTime

	// not first scan
	if s.start != nil {
		startTime := *s.start
		dt := entities.ConvertFromTime(startTime)
		start = &dt
	}

	endTime := s.now()

	dt := entities.ConvertFromTime(endTime)
	end = &dt

	events, err := s.storage.GetEventsToNotify(start, end)

	s.logInfof("%d event(s) push into queue (%s, %s)", len(events), start, end)

	if err != nil {
		s.logErrorf("Scheduler.scan, storage.GetEventsToNotify return error %w", err)
	}

	s.enqueueEvents(events)

	s.start = &endTime

}

// push event info into queue and mark as notified
func (s *Scheduler) enqueueEvents(events []entities.Event) {

	for _, event := range events {
		err := s.queue.Push(event)
		if err != nil {
			s.logErrorf("Scheduler.enqueueEvents, queue.Push return error %s", err)
		} else {
			err = s.storage.MarkEventAsNotified(event.Id(), s.now())
			if err != nil {
				s.logErrorf("Scheduler.enqueueEvents, storage.MarkEventAsNotified return error %s", err)
			}
		}
	}
}

// now helper, call nowTimeFn, that could be redefined in test
// aka template method pattern (with go specific)
func (s *Scheduler) now() time.Time {
	if s.nowTimeFn == nil {
		return time.Now()
	} else {
		return s.nowTimeFn()
	}
}

// log formatted error
func (s *Scheduler) logErrorf(format string, err error) {
	if s.logger != nil {
		s.logger.Error(fmt.Errorf(format, err))
	}
}

// print formatted info level message into log
func (s *Scheduler) logInfof(format string, args ...interface{}) {
	if s.logger != nil {
		///tmp/calendar/scheduler/output
		s.logger.Infof(format, args...)
	}
}
