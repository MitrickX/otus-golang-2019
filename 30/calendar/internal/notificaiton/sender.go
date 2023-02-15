package notificaiton

import (
	"context"
	"fmt"
	"time"

	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/monitoring"

	"go.uber.org/zap"
)

// Sender interface
type Sender interface {
	Run() error
	Stop() error
	Send(info EventInfo) error
}

// Sender, that send info into log file
type LogSender struct {
	queue    Queue
	logger   *zap.SugaredLogger
	eventsCh <-chan EventInfo
	cancel   context.CancelFunc
	metrics  *monitoring.SenderMetrics
}

// Constructor
func NewLogSender(queue Queue, logger zap.SugaredLogger, metrics *monitoring.SenderMetrics) *LogSender {
	return &LogSender{
		queue:   queue,
		logger:  &logger,
		metrics: metrics,
	}
}

// Run sender
func (s *LogSender) Run() error {
	eventsCh, err := s.queue.Consume()
	if err != nil {
		return err
	}
	s.eventsCh = eventsCh

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	s.run(ctx)

	return nil
}

// inner run helper that read from input channel of events and also take into account context Done channel
func (s *LogSender) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Debug("LogSender stopped")
			return
		case eventInfo := <-s.eventsCh:
			err := s.Send(eventInfo)
			if err != nil {
				s.logger.Error(fmt.Errorf("LogSender.Run error while send event info %w", err))
			}
		}
	}
}

// Stop Sender
func (s *LogSender) Stop() error {
	if s.cancel != nil {
		s.cancel()
	}
	if s.metrics != nil {
		s.metrics.Stop()
	}
	return nil
}

// Send method, for this sender just print into log
func (s *LogSender) Send(info EventInfo) error {

	if s.metrics != nil {
		s.metrics.IncSenderCounter()
	}

	t := time.Now()
	slzInfo, _ := serializeEvent(info)
	s.logger.Desugar().Debug("LogSender.Send",
		zap.String("datetime", t.String()),
		zap.String("event", fmt.Sprintf("%#v", info)),
		zap.String("eventSerialized", string(slzInfo)),
	)
	return nil
}