package api

import (
	"encoding/json"
	"fmt"
	"jpb/scheduler/controller"
	"jpb/scheduler/logger"
	"jpb/scheduler/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// HTTPApi represents http api
type HTTPApi struct {
	port int
	ctrl *controller.Ctrl
}

type successResponse struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type tasksResponse struct {
	Data []*utils.Scheduling `json:"data"`
}

// NewHTTPApi creates a new http api struct
func NewHTTPApi(port int, ctrl *controller.Ctrl) *HTTPApi {
	return &HTTPApi{
		port: port,
		ctrl: ctrl,
	}
}

// Listen starts listening
func (a *HTTPApi) Listen() {
	logger.Info("http api listening for schedules")

	r := chi.NewRouter()

	r.Get("/", a.onPing)
	r.Get("/tasks", a.onGetTasks)
	r.Post("/schedule", a.onPostSchedule)

	http.ListenAndServe(fmt.Sprintf(":%d", a.port), r)
}

func (a *HTTPApi) onPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse{Message: "OK"})
}

func (a *HTTPApi) onGetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dateString := r.URL.Query().Get("endDate")
	d, err := time.Parse(time.RFC3339Nano, dateString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Date query is not ISO formatted"})
		return
	}

	tasks := a.ctrl.GetTasks(d)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasksResponse{Data: tasks})
}

func (a *HTTPApi) onPostSchedule(w http.ResponseWriter, r *http.Request) {
	var scheduling scheduling
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&scheduling)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Error with body format"})
		return
	}

	layout := "2006-01-02T15:04:05.999Z"
	t, err := time.Parse(layout, scheduling.Date)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Date is not ISO formatted"})
		return
	}

	id, err := a.ctrl.Schedule(utils.NewScheduling(t, scheduling.Publisher, scheduling.Settings))
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse{Message: fmt.Sprintf("Task %s created", id)})
}
