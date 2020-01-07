package notificaiton

import (
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/domain/entities"
	"github.com/mitrickx/otus-golang-2019/30/calendar/internal/storage/memory"
	"reflect"
	"testing"
	"time"
)

type testQueue struct {
	ch          chan EventInfo
	readTimeout time.Duration
}

func (c *testQueue) Push(event entities.Event) error {
	c.ch <- extractEventInfo(event)
	return nil
}

func (c *testQueue) ReadEvent() (EventInfo, error) {
	events := c.ReadAllEvents()
	if len(events) > 0 {
		return events[0], nil
	}
	return EventInfo{}, ErrQueueEmpty
}

func (c *testQueue) Consume() (<-chan EventInfo, error) {
	return c.ch, nil
}

func (c *testQueue) Close() error {
	close(c.ch)
	return nil
}

func (c *testQueue) ReadAllEvents() []EventInfo {
	var events []EventInfo

	ticker := time.NewTicker(c.readTimeout)

LOOP:
	for {
		select {
		case <-ticker.C:
			break LOOP
		case event := <-c.ch:
			events = append(events, event)
		}
	}

	return events
}

func newTestQueue() *testQueue {
	return &testQueue{
		ch:          make(chan EventInfo, 100),
		readTimeout: 50 * time.Millisecond,
	}
}

func newScheduler() *Scheduler {
	s := &Scheduler{
		scanTimeout: 1 * time.Hour,
		queue:       newTestQueue(),
		storage:     memory.NewStorage(),
	}
	return s
}

func TestSchedulerScan1(t *testing.T) {

	scheduler := newScheduler()

	originalEvents := []entities.Event{
		entities.NewDetailedEvent(
			"TestEvent",
			entities.NewDateTime(2019, 11, 18, 8, 0),
			entities.NewDateTime(2019, 11, 18, 10, 0),
			true,
			30,
			false,
			time.Now(),
		),
		entities.NewDetailedEvent(
			"TestEvent2",
			entities.NewDateTime(2019, 11, 18, 8, 0),
			entities.NewDateTime(2019, 11, 18, 10, 0),
			true,
			10,
			false,
			time.Now(),
		),
	}

	for _, event := range originalEvents {
		_, err := scheduler.storage.AddEvent(event)
		if err != nil {
			t.Errorf("Error while insert event %s: %s", event, err)
			return
		}
	}

	scheduler.scan()

	queue := scheduler.queue.(*testQueue)
	events := queue.ReadAllEvents()

	if len(events) != 2 {
		t.Errorf("Must be scaned 2 events")
		return
	}

	_, err := scheduler.queue.(*testQueue).ReadEvent()
	if err != ErrQueueEmpty {
		t.Errorf("Queue must be empty after read all")
	}

	scheduler.scan()

	_, err = scheduler.queue.(*testQueue).ReadEvent()
	if err != ErrQueueEmpty {
		t.Errorf("Queue must be empty after next scan")
	}

}

func TestSchedulerScan2(t *testing.T) {
	scheduler := newScheduler()

	deviation := 1 * time.Minute

	scheduler.scan()

	lastScanTime := *scheduler.start

	type eventPrepItem struct {
		name     string
		duration time.Duration
		before   int
	}

	items := []eventPrepItem{
		// way beyond start boundary of interval
		{
			name:     "A",
			duration: -10 * deviation,
			before:   30,
		},
		// a bit beyond start boundary of interval
		{
			name:     "B",
			duration: -1 * deviation,
			before:   0,
		},
		// on start boundary of interval
		{
			name:     "C",
			duration: 0,
			before:   5,
		},
		// inside interval
		{
			name:     "D",
			duration: deviation,
			before:   10,
		},
		// on end boundary of interval
		{
			name:     "E",
			duration: scheduler.scanTimeout,
			before:   6,
		},
		// a bit beyond of end of boundary of interval
		{
			name:     "F",
			duration: scheduler.scanTimeout + deviation,
			before:   7,
		},
		// a 2 scans beyond of end boundary of interval
		{
			name:     "G",
			duration: 2*scheduler.scanTimeout + deviation,
			before:   0,
		},
		// a 2 scans beyond of end boundary of interval (right on the end edge)
		{
			name:     "H",
			duration: 3 * scheduler.scanTimeout,
			before:   0,
		},
	}

	var originalEvents []entities.Event

	for _, item := range items {

		startTime := lastScanTime.Add(item.duration + time.Duration(item.before)*time.Minute)
		end := entities.NewDateTime(2030, 11, 18, 10, 0)

		event := entities.NewDetailedEvent(
			item.name,
			entities.ConvertFromTime(startTime),
			end,
			true,
			item.before,
			false,
			time.Now(),
		)

		originalEvents = append(originalEvents, event)

	}

	for _, event := range originalEvents {
		_, err := scheduler.storage.AddEvent(event)
		if err != nil {
			t.Errorf("Error while insert event %s: %s", event, err)
			return
		}
	}

	// way to emulate next scan run after exactly scanTimeout
	scheduler.nowTimeFn = func() time.Time {
		return scheduler.start.Add(scheduler.scanTimeout)
	}

	scheduler.scan()

	queue := scheduler.queue.(*testQueue)
	events := queue.ReadAllEvents()
	names := extractNames(events)

	expectedNames := []string{"C", "D", "E"}
	if !reflect.DeepEqual(expectedNames, names) {
		t.Errorf("Expected names = %v, instreadof %v", expectedNames, names)
		return
	}

	scheduler.scan()

	queue = scheduler.queue.(*testQueue)
	events = queue.ReadAllEvents()
	names = extractNames(events)

	expectedNames = []string{"F"}
	if !reflect.DeepEqual(expectedNames, names) {
		t.Errorf("Expected names = %v, instreadof %v", expectedNames, names)
		return
	}

	scheduler.scan()

	queue = scheduler.queue.(*testQueue)
	events = queue.ReadAllEvents()
	names = extractNames(events)

	expectedNames = []string{"G", "H"}
	if !reflect.DeepEqual(expectedNames, names) {
		t.Errorf("Expected names = %v, instreadof %v", expectedNames, names)
		return
	}

}

func TestSchedulerScan3(t *testing.T) {
	scheduler := newScheduler()

	_, _ = scheduler.storage.AddEvent(
		entities.NewDetailedEvent(
			"TestEvent",
			entities.NewDateTime(2019, 11, 18, 8, 0),
			entities.NewDateTime(2019, 11, 18, 10, 0),
			true,
			30,
			false,
			time.Now(),
		),
	)

	scheduler.scan()

	eventInfo, err := scheduler.queue.(*testQueue).ReadEvent()

	if err != nil {
		t.Errorf("Must not be error %s", err)
		return
	}

	id := eventInfo.Id
	event, err := scheduler.storage.GetEvent(id)

	if err != nil {
		t.Errorf("Must not be error %s", err)
		return
	}

	if !event.IsNotified() {
		t.Errorf("Event must be mark as notified")
	}
}

func extractNames(events []EventInfo) []string {
	var names []string
	for _, event := range events {
		names = append(names, event.Name)
	}
	return names
}
