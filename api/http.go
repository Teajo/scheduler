package api

import (
	"encoding/json"
	"fmt"
	"jpb/scheduler/controller"
	"jpb/scheduler/logger"
	"jpb/scheduler/utils"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	logger.Info(fmt.Sprintf("http api listening for schedules in %d", a.port))

	router := chi.NewRouter()

	FileServer(router, "/", "./ui/build/")
	router.Options("/*", a.options)
	router.Get("/ping", a.onPing)
	router.Get("/tasks", a.onGetTasks)
	router.Get("/tasks/publishers", a.onGetPublishers)
	router.Post("/tasks/schedule", a.onPostSchedule)
	router.Delete("/tasks/{id}", a.onDeleteTask)

	// TODO: check http port availability
	http.ListenAndServe(fmt.Sprintf(":%d", a.port), router)
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

func (a *HTTPApi) onDeleteTask(w http.ResponseWriter, r *http.Request) {
	cors(w)
	jsonResp(w)

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "No task id provided"})
		return
	}

	err := a.ctrl.RemoveTask(id)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "Error when removing task, check server logs"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse{Message: fmt.Sprintf("Task %s removed", id)})
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

// FileServer is serving static files
func FileServer(r chi.Router, public string, static string) {

	if strings.ContainsAny(public, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	root, _ := filepath.Abs(static)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		panic("Static Documents Directory Not Found")
	}

	fs := http.StripPrefix(public, http.FileServer(http.Dir(root)))

	if public != "/" && public[len(public)-1] != '/' {
		r.Get(public, http.RedirectHandler(public+"/", 301).ServeHTTP)
		public += "/"
	}

	r.Get(public+"*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file := strings.Replace(r.RequestURI, public, "/", 1)
		if _, err := os.Stat(root + file); os.IsNotExist(err) {
			http.ServeFile(w, r, path.Join(root, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	}))
}
