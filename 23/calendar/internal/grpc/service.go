package grpc

import (
	"context"
	"fmt"
	"github.com/mitrickx/otus-golang-2019/23/calendar/internal/domain/entities"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

// grpc entities service itself
// Clean architecture approach - not working with inner biz logic layer directly
type Service struct {
	Calendar
	logger *zap.SugaredLogger
	port   string

	// inject now time for getEventsForDay/getEventsForWeek/getEventsForPeriod
	// need to tests
	now time.Time
}

// Constructor
func NewService(port string, storage entities.Storage, logger *zap.SugaredLogger) (*Service, error) {
	service, err := NewCalendar(storage)
	if err != nil {
		return nil, err
	}
	return &Service{
		Calendar: *service,
		logger:   logger,
		port:     port,
		now:      time.Now(),
	}, nil
}

// Run grpc entities service
func (service *Service) Run() {
	s := grpc.NewServer()
	l, err := net.Listen("tcp", ":"+service.port)
	if err != nil && service.logger != nil {
		service.logger.Errorf("Service.Run, net listen, return error %s", err)
	}
	reflection.Register(s)
	RegisterServiceServer(s, service)
	err = s.Serve(l)
	if err != nil && service.logger != nil {
		service.logger.Errorf("Service.Run, grpc.Serve return error %s", err)
	}
}

// Run new grpc entities service
func RunService(port string, storage entities.Storage, logger *zap.SugaredLogger) error {
	service, err := NewService(port, storage, logger)
	if err != nil {
		return err
	}
	service.Run()
	return nil
}

// Create event service method (grpc remote call)
// On success result is  "created %d" string
// On invalid argument return error with codes.InvalidArgument code
// On other cases return some another error
func (service *Service) CreateEvent(ctx context.Context, request *CreateEventRequest) (*SimpleResponse, error) {
	if request.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name must not be empty")
	}
	if request.Start == nil {
		return nil, status.Error(codes.InvalidArgument, "start date must not be empty")
	}
	if request.End == nil {
		return nil, status.Error(codes.InvalidArgument, "end date must not be empty")
	}
	event := &Event{
		Name:  request.Name,
		Start: request.Start,
		End:   request.End,
	}
	id, err := service.AddEvent(event)
	if err != nil {
		return nil, err
	}
	return &SimpleResponse{
		Result: fmt.Sprintf("created %d", id),
	}, nil
}

// Update event service method (grpc remote call)
// On success result is "updated" string
// On invalid argument return error with codes.InvalidArgument code
// On other cases return some another error
func (service *Service) UpdateEvent(ctx context.Context, request *UpdateEventRequest) (*SimpleResponse, error) {
	id := request.GetId()
	if id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be greater 0")
	}
	if request.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name must not be empty")
	}
	if request.Start == nil {
		return nil, status.Error(codes.InvalidArgument, "start date must not be empty")
	}
	if request.End == nil {
		return nil, status.Error(codes.InvalidArgument, "end date must not be empty")
	}
	event := &Event{
		Name:  request.Name,
		Start: request.Start,
		End:   request.End,
	}
	err := service.Calendar.UpdateEvent(int(id), event)
	if err != nil {
		return nil, err
	}
	return &SimpleResponse{
		Result: "updated",
	}, nil
}

// Delete event service method (grpc remote call)
// On success result is "deleted" string
// On invalid argument return error with codes.InvalidArgument code
// On other cases return some another error
func (service *Service) DeleteEvent(ctx context.Context, request *DeleteEventRequest) (*SimpleResponse, error) {
	id := request.GetId()
	if id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be greater 0")
	}
	err := service.Calendar.DeleteEvent(int(id))
	if err != nil {
		return nil, err
	}
	return &SimpleResponse{
		Result: "deleted",
	}, nil
}

// Get events for current day service method (grpc remote call)
// On full success result is list of events
// On partial success (if only some events could be received) return as list as error about other events
// Otherwise return some another error
func (service *Service) GetEventsForDay(ctx context.Context, _ *Nothing) (*EventListResponse, error) {
	period, err := NewDayPeriod(service.now)
	if err != nil {
		return nil, err
	}
	return service.getEventsForPeriod(period)
}

// Get events for current week service method (grpc remote call)
// On full success result is list of events
// On partial success (if only some events could be received) return as list as error about other events
// Otherwise return some another error
func (service *Service) GetEventsForWeek(ctx context.Context, _ *Nothing) (*EventListResponse, error) {
	period, err := NewWeekPeriod(service.now)
	if err != nil {
		return nil, err
	}
	return service.getEventsForPeriod(period)
}

// Get events for current week service method (grpc remote call)
// On full success result is list of events
// On partial success (if only some events could be received) return as list as error about other events
// Otherwise return some another error
func (service *Service) GetEventsForMonth(ctx context.Context, _ *Nothing) (*EventListResponse, error) {
	period, err := NewMonthPeriod(service.now)
	if err != nil {
		return nil, err
	}
	return service.getEventsForPeriod(period)
}

// Helper for GetEventsFor* methods to reduce code duplication
func (service *Service) getEventsForPeriod(period *Period) (*EventListResponse, error) {
	events, err := service.Calendar.GetEventsByPeriod(period)
	if events == nil {
		return nil, err
	}
	response := &EventListResponse{
		Events: events,
	}
	return response, err

}
