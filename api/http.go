package api

import (
	"encoding/json"
	"fmt"
	"jpb/scheduler/controller"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// HTTPApi represents http api
type HTTPApi struct {
	port int
	ctrl *controller.Ctrl
}

type scheduling struct {
	Date string `json:"date"`
}

// NewHTTPApi creates a new http api struct
func NewHTTPApi(port int) *HTTPApi {
	return &HTTPApi{
		port: port,
		ctrl: controller.New(),
	}
}

// Listen starts listening
func (a *HTTPApi) Listen() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Post("/schedule", func(w http.ResponseWriter, r *http.Request) {
		var scheduling scheduling
		err := json.NewDecoder(r.Body).Decode(&scheduling)
		if err != nil {
			fmt.Println(err.Error())
			w.Write([]byte("not ok"))
			return
		}

		layout := "2006-01-02T15:04:05.999Z"
		t, err := time.Parse(layout, scheduling.Date)
		if err != nil {
			fmt.Println(err.Error())
			w.Write([]byte("not ok"))
			return
		}

		id, err := a.ctrl.Schedule(&t)
		if err != nil {
			w.Write([]byte("failed to create task"))
			return
		}

		w.Write([]byte(fmt.Sprintf("task %s created", id)))
	})

	http.ListenAndServe(fmt.Sprintf(":%d", a.port), r)
}
