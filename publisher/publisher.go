package publisher

import (
	"fmt"
	"jpb/scheduler/config"
	"jpb/scheduler/events"
	"jpb/scheduler/logger"
	"jpb/scheduler/retry"
	"jpb/scheduler/utils"
	"time"
)

// ValueType is config key value type
type ValueType string

const (
	STRING      ValueType = "STRING"
	JSON_STRING ValueType = "JSON_STRING"
	INT         ValueType = "INT"
	BOOL        ValueType = "BOOL"
)

// ConfigValueDef config value def
type ConfigValueDef struct {
	Type        ValueType   `json:"type"`
	Default     interface{} `json:"default"`
	Possible    interface{} `json:"possible"`
	Required    bool        `json:"required"`
	Placeholder string      `json:"placeholder"`
}

// Publisher interface
type Publisher interface {
	Publish(map[string]interface{}) *PublishError
	GetConfigDef() map[string]*ConfigValueDef
}

// PubManager is a publisher manager
type PubManager struct {
	Bus        *events.Bus `inject:""`
	publishers map[string]Publisher
}

// Start starts a publisher manager
func (pm *PubManager) Start() {
	pm.publishers = loadPublisherPlugins(config.Get().PluginDir)

	e := pm.Bus.Subscribe(events.TASKDONE)
	go pm.listen(e)
}

// Get retrieves a publisher according to provided id
func (pm *PubManager) Get(id string) (Publisher, bool) {
	pub, ok := pm.publishers[id]
	return pub, ok
}

// GetAvailable returns a list of available publishers
func (pm *PubManager) GetAvailable() map[string]map[string]*ConfigValueDef {
	pubs := make(map[string]map[string]*ConfigValueDef)
	for k, v := range pm.publishers {
		pubs[k] = v.GetConfigDef()
	}
	return pubs
}

// Listen listens for done tasks
func (pm *PubManager) listen(e chan interface{}) {
	logger.Info("publisher listening for done tasks")
	for {
		payload := <-e
		s, ok := payload.(*utils.Scheduling)
		if !ok {
			logger.Error("Impossible to cast publishing payload")
		}

		pm.publish(s)
	}
}

func (pm *PubManager) publish(scheduling *utils.Scheduling) {
	for _, pub := range scheduling.Publishers {
		pubName := pub.Publisher
		strat := pub.RetryStrat
		settings := pub.Settings
		publisher, ok := pm.publishers[pubName]

		if ok {
			go retry.Do(func() error {
				logger.Info(fmt.Sprintf("try to publish to %s at %s", pubName, time.Now().Format(time.RFC3339Nano)))
				err := publisher.Publish(settings)
				if err != nil {
					logger.Error(err.Error())
					if err.ShouldRetry() {
						return err
					}
				}
				return nil
			}, strat.Limit, strat.Timeout, strat.Exponential)
		}
	}
}
