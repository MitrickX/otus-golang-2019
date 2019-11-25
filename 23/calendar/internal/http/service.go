package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mitrickx/otus-golang-2019/23/calendar/internal/domain/entities"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

// Ok json response
type OkResponse struct {
	Result string `json:"result"`
}

// Ok json response with list of events
type EventListResponse struct {
	Result []*Event `json:"result"`
}

// Error json response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Http entities service itself
// Clean architecture approach - not working with inner biz logic layer directly
type Service struct {
	Calendar
	logger *zap.SugaredLogger
	port   string
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
	}, nil
}

// Middleware to log requests
func (service *Service) requestLogMiddleware(next http.Handler) http.Handler {
	// if not logger - no middleware
	if service.logger == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service.logger.Infof("Request %s %s %s %s", r.Method, r.RemoteAddr, r.URL.Path, time.Now())
		next.ServeHTTP(w, r)
	})
}

// Run http entities service
func (service *Service) Run() {

	router := mux.NewRouter()
	router.HandleFunc("/create_event", service.CreatEvent).Methods("POST")
	router.HandleFunc("/update_event", service.UpdateEvent).Methods("POST")
	router.HandleFunc("/delete_event", service.DeleteEvent).Methods("POST")
	router.HandleFunc("/events_for_day", service.GetEventsForDay).Methods("GET")
	router.HandleFunc("/events_for_week", service.GetEventsForWeek).Methods("GET")
	router.HandleFunc("/events_for_month", service.GetEventsForMonth).Methods("GET")

	handler := service.requestLogMiddleware(router)

	if service.logger != nil {
		service.logger.Infof("start server at %s", service.port)
	}

	err := http.ListenAndServe(":"+service.port, handler)
	if err != nil && service.logger != nil {
		service.logger.Errorf("Service.Run, http listen and serve failed, return error %s", err)
	}
}

// Run new http entities service
func RunService(port string, storage entities.Storage, logger *zap.SugaredLogger) error {
	service, err := NewService(port, storage, logger)
	if err != nil {
		return err
	}
	service.Run()
	return nil
}

// Create event handler
// On success response by ok json response with "create %d" result string
func (service *Service) CreatEvent(w http.ResponseWriter, r *http.Request) {
	service.parseForm(r)

	event, err := NewEvent(r.Form.Get("name"), r.Form.Get("start"), r.Form.Get("end"))
	if err != nil {
		service.writeErrorResponse(w, err.Error(), 400)
		return
	}

	id, err := service.AddEvent(event)
	if err != nil {
		service.writeErrorResponse(w, err.Error(), 200)
		return
	}

	service.writeOkResponse(w, fmt.Sprintf("created %d", id), 200)
}

// Update event handler
// On success response by ok json response with "updated" result string
func (service *Service) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	service.parseForm(r)

	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil || id <= 0 {
		service.writeErrorResponse(w, "invalid id parameter, must be int greater than 0", 400)
		return
	}

	event, err := NewEvent(r.Form.Get("name"), r.Form.Get("start"), r.Form.Get("end"))
	if err != nil {
		service.writeErrorResponse(w, err.Error(), 400)
		return
	}

	err = service.Calendar.UpdateEvent(id, event)
	if err != nil {
		service.writeErrorResponse(w, err.Error(), 200)
		return
	}

	service.writeOkResponse(w, "updated", 200)
}

// Delete event handler
// On success response by ok json response with "deleted" result string
func (service *Service) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	service.parseForm(r)

	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil || id <= 0 {
		service.writeErrorResponse(w, "invalid id parameter, must be int greater than 0", 400)
		return
	}

	err = service.Calendar.DeleteEvent(id)
	if err != nil {
		service.writeErrorResponse(w, err.Error(), 200)
		return
	}

	service.writeOkResponse(w, "deleted", 200)
}

// Get events for current day handler
// response by ok json response with list of events
func (service *Service) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	service.getEventsForDay(time.Now(), w, r)
}

// Inner method for testing, in test we want pass own 'now'
func (service *Service) getEventsForDay(now time.Time, w http.ResponseWriter, r *http.Request) {
	startTime, endTime := GetDayPeriod(now)
	service.getEventsForPeriod(startTime, endTime, w, r)
}

// Get events for current week handler
// response by ok json response with list of events
func (service *Service) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	service.getEventsForWeek(time.Now(), w, r)
}

// Inner method for testing, in test we want pass own 'now'
func (service *Service) getEventsForWeek(now time.Time, w http.ResponseWriter, r *http.Request) {
	startTime, endTime := GetWeekPeriod(now)
	service.getEventsForPeriod(startTime, endTime, w, r)
}

// Get events for current month handler
// response by ok json response with list of events
func (service *Service) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	service.getEventsForMonth(time.Now(), w, r)
}

// Inner method for testing, in test we want pass own 'now'
func (service *Service) getEventsForMonth(now time.Time, w http.ResponseWriter, r *http.Request) {
	startTime, endTime := GetMonthPeriod(now)
	service.getEventsForPeriod(startTime, endTime, w, r)
}

// Helper for GetEventsFor* methods to reduce code duplication
func (service *Service) getEventsForPeriod(start, end string, w http.ResponseWriter, r *http.Request) {
	events, err := service.Calendar.GetEventsByPeriod(start, end)

	if err != nil {
		service.writeErrorResponse(w, "internal server error", 500)
		if service.logger != nil {
			service.logger.Errorf("Service.GetEventsForDay, error Calendar.GetEventsByTimestampsPeriod %s", err)
		}
	}

	service.writeEventListResponse(w, events, 200)
}

// inner helper for parse form
func (service *Service) parseForm(r *http.Request) {
	err := r.ParseForm()
	if err != nil && service.logger != nil {
		service.logger.Errorf("Service.parseForm error %s", err)
	}
}

// inner helper for write ok json response
func (service *Service) writeOkResponse(w http.ResponseWriter, result string, code int) {
	response := &OkResponse{result}
	data, err := json.Marshal(response)

	if err != nil {
		if service.logger != nil {
			service.logger.Errorf("Service.writeOkResponse, marshal response error %s", err)
		}
		w.WriteHeader(500)
		_, writeErr := w.Write([]byte("Server error"))
		if writeErr != nil && service.logger != nil {
			service.logger.Errorf("Service.writeOkResponse, write `Server error` error %s", err)
		}
		return
	}

	w.WriteHeader(code)
	_, writeErr := w.Write(data)
	if writeErr != nil && service.logger != nil {
		service.logger.Errorf("Service.writeOkResponse, write `OkResponse` error %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
}

// inner helper for write error json response
func (service *Service) writeErrorResponse(w http.ResponseWriter, result string, code int) {
	response := &ErrorResponse{result}
	data, err := json.Marshal(response)

	if err != nil {
		if service.logger != nil {
			service.logger.Errorf("Service.writeErrorResponse, marshal response error %s", err)
		}
		w.WriteHeader(500)
		_, writeErr := w.Write([]byte("Server error"))
		if writeErr != nil && service.logger != nil {
			service.logger.Errorf("Service.writeErrorResponse, write `Server error` error %s", err)
		}
		return
	}

	w.WriteHeader(code)
	_, writeErr := w.Write(data)
	if writeErr != nil && service.logger != nil {
		service.logger.Errorf("Service.writeErrorResponse, write `OkResponse` error %s", err)
	}
}

// inner helper for write ok json response with list of events
func (service *Service) writeEventListResponse(w http.ResponseWriter, evens []*Event, code int) {
	response := &EventListResponse{evens}
	data, err := json.Marshal(response)

	if err != nil {
		if service.logger != nil {
			service.logger.Errorf("Service.writeEventListResponse, marshal response error %s", err)
		}
		w.WriteHeader(500)
		_, writeErr := w.Write([]byte("internal server error"))
		if writeErr != nil && service.logger != nil {
			service.logger.Errorf("Service.writeEventListResponse, write `internal server error` error %s", err)
		}
		return
	}

	w.WriteHeader(code)
	_, writeErr := w.Write(data)
	if writeErr != nil && service.logger != nil {
		service.logger.Errorf("Service.writeEventListResponse, write `EventListResponse` error %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
}
