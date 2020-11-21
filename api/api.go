package api

type api interface {
	Listen()
}

// Scheduling represents a scheduling object
type scheduling struct {
	Date      string            `json:"date"`
	Publisher string            `json:"publisher"`
	Settings  map[string]string `json:"settings"`
}
