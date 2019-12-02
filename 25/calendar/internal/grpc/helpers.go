package grpc

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mitrickx/otus-golang-2019/25/calendar/internal/domain/entities"
	"time"
)

// Helper for create new timestamp by 5 int components make sense for this package and application
func NewTimestamp(year, month, day, hour, minute int) (*timestamp.Timestamp, error) {
	t := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
	return ptypes.TimestampProto(t)
}

// Inner helper that convert grpc timestamp.Timestamp to entities.EventTime
func convertToCalendarEventTime(tspb *timestamp.Timestamp) (*entities.EventTime, error) {
	t, err := ptypes.Timestamp(tspb)
	if err != nil {
		return nil, err
	}
	eventTime := entities.ConvertFromTime(t)
	return &eventTime, nil
}

// Inner Helper that helps convert grpc.Event to entities.Event
func convertToCalendarEvent(event *Event) (*entities.Event, error) {

	startTime, err := convertToCalendarEventTime(event.Start)
	if err != nil {
		return nil, err
	}

	endTime, err := convertToCalendarEventTime(event.End)
	if err != nil {
		return nil, err
	}

	calendarEvent := entities.NewEventWithId(int(event.Id), event.Name, *startTime, *endTime)

	return &calendarEvent, nil
}

// Convert from inner Event entity (entities.Event) to grpc.Event
func convertFromCalendarEvent(calendarEvent entities.Event) (*Event, error) {
	start, err := ptypes.TimestampProto(calendarEvent.Start().Time())
	if err != nil {
		return nil, err
	}

	end, err := ptypes.TimestampProto(calendarEvent.End().Time())
	if err != nil {
		return nil, err
	}

	event := &Event{
		Id:    int32(calendarEvent.Id()),
		Name:  calendarEvent.Name(),
		Start: start,
		End:   end,
	}

	return event, nil
}

// Is proto timestamps equals
func isTimestampEquals(tspb1 *timestamp.Timestamp, tspb2 *timestamp.Timestamp) bool {
	// as pointers
	if tspb1 == tspb2 {
		return true
	}
	return tspb1.GetSeconds() == tspb2.GetSeconds() && tspb1.GetNanos() == tspb2.GetNanos()
}

// Is grpc events are equals
func isEventEquals(event1 *Event, event2 *Event, deep bool) bool {
	// as pointers
	if event1 == event2 {
		return true
	}
	nameEquals := event1.Name == event2.Name
	startEquals := isTimestampEquals(event1.Start, event2.Start)
	endEquals := isTimestampEquals(event2.End, event2.End)
	if deep {
		return nameEquals && startEquals && endEquals && event1.Id == event2.Id
	} else {
		return nameEquals && startEquals && endEquals
	}
}

// Same as NewTimestamp but with suppressing error
// Recommended for tests only
func ts(year, month, day, hour, minute int) *timestamp.Timestamp {
	t := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
	ts, _ := ptypes.TimestampProto(t)
	return ts
}
