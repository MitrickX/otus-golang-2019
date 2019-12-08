package grpc

import (
	"context"
	"fmt"
	"github.com/mitrickx/otus-golang-2019/25/calendar/internal/storage/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"regexp"
	"strconv"
	"testing"
	"time"
)

var bufConnSize = 4096

func TestCreateEventOK(t *testing.T) {

	service, client := RunTestGrpcPipe(t)

	request := &CreateEventRequest{
		Name:  "Do homework",
		Start: ts(2019, 10, 15, 20, 0),
		End:   ts(2019, 10, 15, 22, 0),
	}

	response, err := client.CreateEvent(context.Background(), request)

	if err != nil {
		t.Errorf("Create event must not return err %s", err)
	}

	re := regexp.MustCompile(`created (\d+)`)
	match := re.FindStringSubmatch(response.Result)
	if match == nil {
		t.Errorf("unexpected OkResponse.Result value `%s`", response.Result)
		return
	}

	id, err := strconv.Atoi(match[1])
	if err != nil {
		t.Fatalf("strconv return error %s", err)
	}

	if id <= 0 {
		t.Errorf("unexpected event id %d, must be > 0", id)
	}

	if service.getEventsTotalCount() != 1 {
		t.Errorf("unexpected count of events in entities, must be 1 instead of %d", service.getEventsTotalCount())
	}

}

func TestCreateEventInvalidName(t *testing.T) {
	_, client := RunTestGrpcPipe(t)

	request := &CreateEventRequest{
		Name:  "",
		Start: ts(2019, 10, 15, 20, 0),
		End:   ts(2019, 10, 15, 22, 0),
	}

	response, err := client.CreateEvent(context.Background(), request)

	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("expected status code %d (invalid argument) instread of %d", codes.InvalidArgument, status.Code(err))
	}

	if status.Convert(err).Message() != "name must not be empty" {
		t.Errorf("expected error message `name must not be empty` instread of `%s`", status.Convert(err).Message())
	}

	if response != nil {
		t.Errorf("response must be nil instread of %+v", response)
	}
}

func TestCreateEventInvalidStart(t *testing.T) {
	_, client := RunTestGrpcPipe(t)

	request := &CreateEventRequest{
		Name:  "Do homework",
		Start: nil,
		End:   ts(2019, 10, 15, 22, 0),
	}

	response, err := client.CreateEvent(context.Background(), request)

	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("expected status code %d (invalid argument) instread of %d", codes.InvalidArgument, status.Code(err))
	}

	if status.Convert(err).Message() != "start date must not be empty" {
		t.Errorf("expected error message `start date must not be empty` instread of `%s`", status.Convert(err).Message())
	}

	if response != nil {
		t.Errorf("response must be nil instread of %+v", response)
	}
}

func TestUpdateEventOK(t *testing.T) {
	service, client := RunTestGrpcPipe(t)

	event1 := &Event{
		Name:  "Do homework",
		Start: ts(2019, 10, 15, 2, 0),
		End:   ts(2019, 10, 15, 22, 0),
	}

	id := addEvent(t, &service.Calendar, event1, 1)
	if id <= 0 {
		return
	}

	request := &UpdateEventRequest{
		Id:    int32(id),
		Name:  "Watch movie",
		Start: ts(2019, 10, 15, 22, 0),
		End:   ts(2019, 10, 16, 01, 00),
	}

	response, err := client.UpdateEvent(context.Background(), request)

	if err != nil {
		t.Errorf("must not be error on update %s", err)
	}

	if response.Result != "updated" {
		t.Errorf("must be `updated` result on update")
	}

	event, err := service.GetEvent(id)
	if err != nil {
		t.Errorf("must not be error on get %s", err)
	}

	if event.Name != request.Name {
		t.Errorf("event is not udpated, Name must be `%s` not `%s`", request.Name, event.Name)
	}

	if !isTimestampEquals(event.Start, request.Start) {
		t.Errorf("event is not udpated, Start must be `%s` not `%s`", request.Start, event.Start)
	}

	if !isTimestampEquals(event.End, request.End) {
		t.Errorf("event is not udpated, End must be `%s` not `%s`", request.End, event.End)
	}
}

func TestUpdateEventNotExisting(t *testing.T) {
	service, client := RunTestGrpcPipe(t)

	event1 := &Event{
		Name:  "Do homework",
		Start: ts(2019, 10, 15, 2, 0),
		End:   ts(2019, 10, 15, 22, 0),
	}

	id := addEvent(t, &service.Calendar, event1, 1)
	if id <= 0 {
		return
	}

	request := &UpdateEventRequest{
		Id:    100,
		Name:  "Watch movie",
		Start: ts(2019, 10, 15, 22, 0),
		End:   ts(2019, 10, 16, 01, 00),
	}

	response, err := client.UpdateEvent(context.Background(), request)

	if response != nil {
		t.Errorf("error must be nil instread of %+v", response)
	}

	expected := "couldn't update event in storage: event not found"
	if status.Convert(err).Message() != expected {
		t.Errorf("expected error `%s` instread of `%s`", expected, status.Convert(err).Message())
	}
}

func TestUpdateEventInvalidId(t *testing.T) {

	_, client := RunTestGrpcPipe(t)

	request := &UpdateEventRequest{
		Id:    0,
		Name:  "Watch movie",
		Start: ts(2019, 10, 15, 22, 0),
		End:   ts(2019, 10, 16, 01, 00),
	}

	response, err := client.UpdateEvent(context.Background(), request)

	if response != nil {
		t.Errorf("error must be nil instread of %+v", response)
	}

	expected := "id must be greater 0"
	if status.Convert(err).Message() != expected {
		t.Errorf("expected error `%s` instread of `%s`", expected, status.Convert(err).Message())
	}

	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("expected status code %d (invalid argument) instread of %d", codes.InvalidArgument, status.Code(err))
	}
}

func TestUpdateEventInvalidName(t *testing.T) {

	service, client := RunTestGrpcPipe(t)

	event1 := &Event{
		Name:  "Do homework",
		Start: ts(2019, 10, 15, 20, 0),
		End:   ts(2019, 10, 15, 22, 0),
	}

	id := addEvent(t, &service.Calendar, event1, 1)
	if id <= 0 {
		return
	}

	request := &UpdateEventRequest{
		Id:    int32(id),
		Name:  "",
		Start: ts(2019, 10, 15, 22, 0),
		End:   ts(2019, 10, 16, 01, 00),
	}

	response, err := client.UpdateEvent(context.Background(), request)

	if response != nil {
		t.Errorf("error must be nil instread of %+v", response)
	}

	expected := "name must not be empty"
	if status.Convert(err).Message() != expected {
		t.Errorf("expected error `%s` instread of `%s`", expected, status.Convert(err).Message())
	}

	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("expected status code %d (invalid argument) instread of %d", codes.InvalidArgument, status.Code(err))
	}
}

func TestUpdateEventInvalidStart(t *testing.T) {

	service, client := RunTestGrpcPipe(t)

	event1 := &Event{
		Name:  "Do homework",
		Start: ts(2019, 10, 15, 2, 0),
		End:   ts(2019, 10, 15, 22, 0),
	}

	id := addEvent(t, &service.Calendar, event1, 1)
	if id <= 0 {
		return
	}

	request := &UpdateEventRequest{
		Id:    int32(id),
		Name:  "Watch movie",
		Start: nil,
		End:   ts(2019, 10, 16, 01, 00),
	}

	response, err := client.UpdateEvent(context.Background(), request)

	if response != nil {
		t.Errorf("error must be nil instread of %+v", response)
	}

	expected := "start date must not be empty"
	if status.Convert(err).Message() != expected {
		t.Errorf("expected error `%s` instread of `%s`", expected, status.Convert(err).Message())
	}

	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("expected status code %d (invalid argument) instread of %d", codes.InvalidArgument, status.Code(err))
	}
}

func TestDeleteEventOK(t *testing.T) {
	service, client := RunTestGrpcPipe(t)

	event1 := &Event{
		Name:  "Do homework",
		Start: ts(2019, 10, 15, 2, 0),
		End:   ts(2019, 10, 15, 22, 0),
	}

	id := addEvent(t, &service.Calendar, event1, 1)
	if id <= 0 {
		return
	}

	request := &DeleteEventRequest{
		Id: int32(id),
	}

	response, err := client.DeleteEvent(context.Background(), request)

	if err != nil {
		t.Errorf("must not be error on delete %s", err)
	}

	if response.Result != "deleted" {
		t.Errorf("result must be `deleted` instread of %s", response.Result)
	}

	_, err = service.GetEvent(id)
	if err != ErrorNotFound {
		t.Errorf("event might not deleted, expected error `%s` instread of `%s`", ErrorNotFound, err)
	}

}

func TestDeleteEventNotExisting(t *testing.T) {
	service, client := RunTestGrpcPipe(t)

	event1 := &Event{
		Name:  "Do homework",
		Start: ts(2019, 10, 15, 2, 0),
		End:   ts(2019, 10, 15, 22, 0),
	}

	id := addEvent(t, &service.Calendar, event1, 1)
	if id <= 0 {
		return
	}

	request := &DeleteEventRequest{
		Id: 100,
	}

	response, err := client.DeleteEvent(context.Background(), request)

	if response != nil {
		t.Errorf("error must be nil instread of %+v", response)
	}

	expected := "couldn't delete event from storage: event not found"
	if status.Convert(err).Message() != expected {
		t.Errorf("expected error `%s` instread of `%s`", expected, status.Convert(err).Message())
	}
}

func TestDeleteEventInvalidId(t *testing.T) {

	_, client := RunTestGrpcPipe(t)

	request := &DeleteEventRequest{
		Id: -10,
	}

	response, err := client.DeleteEvent(context.Background(), request)

	if response != nil {
		t.Errorf("error must be nil instread of %+v", response)
	}

	expected := "id must be greater 0"
	if status.Convert(err).Message() != expected {
		t.Errorf("expected error `%s` instread of `%s`", expected, status.Convert(err).Message())
	}

	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("expected status code %d (invalid argument) instread of %d", codes.InvalidArgument, status.Code(err))
	}
}

func TestGetEventsForDay(t *testing.T) {
	service, client := RunTestGrpcPipe(t)

	// set up fixes list of events
	addFixedListOfEvents(t, &service.Calendar)

	// set deterministic now time for test
	service.now = time.Date(2019, 11, 21, 8, 0, 0, 0, time.UTC)

	response, err := client.GetEventsForDay(context.Background(), &Nothing{})
	if err != nil {
		t.Errorf("must not be error instread of %s", err)
		return
	}

	if len(response.Events) != 1 {
		t.Errorf("event list must has one event instead of %d", len(response.Events))
		return
	}

	event := *response.Events[0]
	if event.Name != "Thursday" {
		t.Error("Must be Thursday event")
	}
}

func TestGetEventsForWeek(t *testing.T) {

	service, client := RunTestGrpcPipe(t)

	// set up fixes list of events
	addFixedListOfEvents(t, &service.Calendar)

	// set deterministic now time for test
	service.now = time.Date(2019, 11, 21, 8, 0, 0, 0, time.UTC)

	response, err := client.GetEventsForWeek(context.Background(), &Nothing{})
	if err != nil {
		t.Errorf("must not be error instread of %s", err)
		return
	}

	if len(response.Events) != 7 {
		t.Errorf("event list must has 7 events instead of %d", len(response.Events))
		return
	}
}

func TestGetEventsForMonth(t *testing.T) {
	service, client := RunTestGrpcPipe(t)

	// set up fixes list of events
	addFixedListOfEvents(t, &service.Calendar)

	addEvent(t, &service.Calendar, &Event{
		Name:  "First day",
		Start: ts(2019, 11, 01, 8, 0),
		End:   ts(2019, 11, 01, 10, 0),
	}, 8)

	addEvent(t, &service.Calendar, &Event{
		Name:  "Last day",
		Start: ts(2019, 11, 30, 8, 0),
		End:   ts(2019, 11, 30, 10, 0),
	}, 9)

	// set deterministic now time for test
	service.now = time.Date(2019, 11, 21, 8, 0, 0, 0, time.UTC)

	response, err := client.GetEventsForMonth(context.Background(), &Nothing{})
	if err != nil {
		t.Errorf("must not be error instread of %s", err)
		return
	}

	if len(response.Events) != 9 {
		t.Errorf("event list must has 9 events instead of %d", len(response.Events))
		return
	}
}

func RunTestGrpcPipe(t *testing.T) (*Service, ServiceClient) {

	listener := bufconn.Listen(bufConnSize)

	var client ServiceClient
	var clientResCh chan error

	service, serverResCh := RunTestService(listener)

	// If error return error right away not run client and close listener
	select {
	case err := <-serverResCh:
		if err != nil {
			t.Error(err)
		}
		_ = listener.Close()
		return nil, nil
	default:
		client, clientResCh = RunTestClient(listener)
	}

	// If server or client return error close listener
	go func() {
		select {
		case err := <-serverResCh:
			if err != nil {
				t.Error(err)
			}
			_ = listener.Close()
		case err := <-clientResCh:
			if err != nil {
				t.Error(err)
			}
			_ = listener.Close()
		}
	}()

	return service, client
}

func RunTestClient(listener *bufconn.Listener) (client ServiceClient, resultCh chan error) {
	resultCh = make(chan error, 1)

	bufDialer := func(_ context.Context, _ string) (net.Conn, error) {
		return listener.Dial()
	}
	conn, err := grpc.Dial("bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		resultCh <- fmt.Errorf("grpc Dial with bufconn connection return error %s\n", err)
		return
	}

	client = NewServiceClient(conn)
	return
}

func RunTestService(listener *bufconn.Listener) (service *Service, resultCh chan error) {

	resultCh = make(chan error, 1)

	storage := memory.NewStorage()
	service, err := NewService("", storage, nil)

	if err != nil {
		resultCh <- fmt.Errorf("test server exited with error %s", err)
		return
	}

	s := grpc.NewServer()
	RegisterServiceServer(s, service)

	go func() {
		err := s.Serve(listener)
		if err != nil {
			resultCh <- fmt.Errorf("test server exited with error %s", err)
		}
	}()

	return
}
