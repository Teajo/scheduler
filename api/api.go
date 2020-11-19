package api

type api interface {
	Listen()
}

// Scheduling represents a scheduling object
type scheduling struct {
	Date       string            `json:"date"`
	Publishers []string          `json:"publishers"`
	Settings   map[string]string `json:"settings"`
}
