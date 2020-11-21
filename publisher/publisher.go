package publisher

import (
	"fmt"
	"jpb/scheduler/config"
	"jpb/scheduler/logger"
	"jpb/scheduler/retry"
	"jpb/scheduler/utils"
	"time"
)

// Publisher interface
type Publisher interface {
	CheckConfig(map[string]string) error
	Publish(map[string]string) *PublishError
}

// PubManager is a publisher manager
type PubManager struct {
	taskDone   chan *utils.Scheduling
	publishers map[string]Publisher
}

// New creates a publisher manager
func New(taskDone chan *utils.Scheduling) *PubManager {
	pubs := loadPublisherPlugins(config.Get().PluginDir)

	pm := &PubManager{
		taskDone:   taskDone,
		publishers: pubs,
	}

	go pm.listen()
	return pm
}

// Get retrieves a publisher according to provided id
func (pm *PubManager) Get(id string) (Publisher, bool) {
	pub, ok := pm.publishers[id]
	return pub, ok
}

// Listen listens for done tasks
func (pm *PubManager) listen() {
	logger.Info("publisher listening for done tasks")
	for {
		scheduling := <-pm.taskDone
		pm.publish(scheduling)
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
