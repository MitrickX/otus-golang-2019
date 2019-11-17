package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type OkResponse struct {
	Result string `json:"result"`
}

type EventListResponse struct {
	Result []*Event `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Service struct {
	Calendar
	logger *zap.SugaredLogger
	port   string
}

func NewService(port string, logger *zap.SugaredLogger) *Service {
	service := NewCalendar()
	return &Service{
		Calendar: *service,
		logger:   logger,
		port:     port,
	}
}

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

func (service *Service) Run() {

	router := mux.NewRouter()
	router.HandleFunc("/create_event", service.CreatEvent).Methods("POST")
	router.HandleFunc("/update_event", service.UpdateEvent).Methods("POST")
	router.HandleFunc("/delete_event", service.DeleteEvent).Methods("POST")

	handler := service.requestLogMiddleware(router)

	if service.logger != nil {
		service.logger.Infof("start server at %s", service.port)
	}

	err := http.ListenAndServe(":"+service.port, handler)
	if err != nil && service.logger != nil {
		service.logger.Errorf("Service.Run, http listen and serve failed, return error %s", err)
	}
}

func RunService(port string, logger *zap.SugaredLogger) {
	NewService(port, logger).Run()
}

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

func (service *Service) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	service.parseForm(r)

	now := time.Now()
	startTime := now.Format(dateLayout) + " 00:00"
	endTime := now.Format(dateLayout) + " 23:59"
	events, err := service.Calendar.GetEventsByPeriod(startTime, endTime)

	if err != nil {
		service.writeErrorResponse(w, "internal server error", 500)
		if service.logger != nil {
			service.logger.Errorf("Service.GetEventsForDay, error Calendar.GetEventsByPeriod %s", err)
		}
	}

	service.writeEventListResponse(w, events, 200)
}

func (service *Service) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	service.parseForm(r)

	now := time.Now()
	nowWeek := now.Weekday()
	nowWeekInt := int(nowWeek)
	if nowWeekInt == 0 {
		nowWeekInt = 7
	}
	//-int()
	startTime := now.Format(dateLayout) + " 00:00"
	endTime := now.Format(dateLayout) + " 23:59"
	events, err := service.Calendar.GetEventsByPeriod(startTime, endTime)

	if err != nil {
		service.writeErrorResponse(w, "internal server error", 500)
		if service.logger != nil {
			service.logger.Errorf("Service.GetEventsForDay, error Calendar.GetEventsByPeriod %s", err)
		}
	}

	service.writeEventListResponse(w, events, 200)
}

func (service *Service) parseForm(r *http.Request) {
	err := r.ParseForm()
	if err != nil && service.logger != nil {
		service.logger.Errorf("Service.parseForm error %s", err)
	}
}

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
