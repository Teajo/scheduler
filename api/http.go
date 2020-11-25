package api

import (
	"encoding/json"
	"fmt"
	"jpb/scheduler/controller"
	"jpb/scheduler/logger"
	"jpb/scheduler/utils"
	"net/http"
	"net/url"
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

type dataResponse struct {
	Data interface{} `json:"data"`
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

	r.Options("/*", a.options)
	r.Get("/", a.onPing)
	r.Get("/tasks", a.onGetTasks)
	r.Get("/tasks/publishers", a.onGetPublishers)
	r.Post("/tasks/schedule", a.onPostSchedule)

	// TODO: check http port availability
	http.ListenAndServe(fmt.Sprintf(":%d", a.port), r)
}

func (a *HTTPApi) options(w http.ResponseWriter, r *http.Request) {
	cors(w)
	jsonResp(w)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse{Message: "OK"})
}

func (a *HTTPApi) onPing(w http.ResponseWriter, r *http.Request) {
	cors(w)
	jsonResp(w)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse{Message: "OK"})
}

func (a *HTTPApi) onGetTasks(w http.ResponseWriter, r *http.Request) {
	cors(w)
	jsonResp(w)

	startStr := getQueryValue(r.URL.Query(), "startDate", []string{utils.FirstDate.Format(time.RFC3339Nano)})
	endStr := getQueryValue(r.URL.Query(), "endDate", []string{utils.LastDate.Format(time.RFC3339Nano)})

	end, err := time.Parse(time.RFC3339Nano, endStr[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "End date query is not ISO formatted"})
		return
	}

	start, err := time.Parse(time.RFC3339Nano, startStr[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Start date query is not ISO formatted"})
		return
	}

	tasks := a.ctrl.GetTasks(start, end)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dataResponse{Data: tasks})
}

func (a *HTTPApi) onGetPublishers(w http.ResponseWriter, r *http.Request) {
	cors(w)
	jsonResp(w)

	pubs := a.ctrl.GetPublishers()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dataResponse{Data: pubs})
}

func (a *HTTPApi) onPostSchedule(w http.ResponseWriter, r *http.Request) {
	cors(w)
	jsonResp(w)

	var scheduling Scheduling

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

	id, err := a.ctrl.Schedule(utils.NewScheduling(t, scheduling.Publishers))
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse{Message: fmt.Sprintf("Task %s created", id)})
}

func getQueryValue(query url.Values, key string, dflt []string) []string {
	v, ok := query[key]
	if !ok {
		v = dflt
	}
	return v
}

func cors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Origin, content-type")
}

func jsonResp(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
