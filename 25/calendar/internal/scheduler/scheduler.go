package scheduler

import (
	"fmt"
	"github.com/mitrickx/otus-golang-2019/25/calendar/internal/domain/entities"
	"go.uber.org/zap"
	"time"
)

type Scheduler struct {
	scanTimeout  time.Duration    // frequency of scan
	storage      entities.Storage // calendar storage
	lastScanTime *time.Time       // time of last scan
	logger       *zap.SugaredLogger
	queue        Queue
}

func (s *Scheduler) Run() {
	ticker := time.NewTicker(s.scanTimeout)
	for range ticker.C {
		s.scan()
	}
}

func (s *Scheduler) scan() {
	var start, end *entities.DateTime

	// not first scan
	if s.lastScanTime != nil {
		scanTime := *s.lastScanTime
		t := scanTime.Add(-s.scanTimeout)
		dt := entities.ConvertFromTime(t)
		start = &dt
	}

	nowTime := time.Now()
	endTime := nowTime.Add(s.scanTimeout)

	dt := entities.ConvertFromTime(endTime)
	end = &dt

	events, err := s.storage.GetEventsForNotification(start, end)

	if err != nil {
		s.logErrorf("Scheduler.scan, storage.GetEventsForNotification return error %w", err)
	}

	err = s.enqueueEvents(events)
	if err != nil {
		nowTime := time.Now()
		s.lastScanTime = &nowTime
	}
}

func (s *Scheduler) enqueueEvents(events []entities.Event) error {
	if s.queue == nil {
		return errors.New("queue not inited")
	}

	for _, event := range events {
		err := s.queue.AddEvent(event)
		s.logErrorf("Scheduler.enqueEvents, queue.AddEvent return error %s", err)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) logErrorf(format string, err error) {
	if s.logger != nil {
		s.logger.Error(fmt.Errorf(format, err))
	}
}

func (s *Scheduler) logError(err error) {
	if s.logger != nil {
		s.logger.Error(err)
	}
}
