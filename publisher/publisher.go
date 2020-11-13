package publisher

import (
	"fmt"
	"jpb/scheduler/config"
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

	return &PubManager{
		taskDone:   taskDone,
		publishers: pubs,
	}
}

// Listen listens for done tasks
func (pm *PubManager) Listen() {
	fmt.Println("publisher listening for done tasks")
	for {
		scheduling := <-pm.taskDone
		go pm.publish(scheduling)
	}
}

// Get retrieves a publisher according to provided id
func (pm *PubManager) Get(id string) (Publisher, bool) {
	pub, ok := pm.publishers[id]
	return pub, ok
}

func (pm *PubManager) publish(scheduling *utils.Scheduling) {
	publisher, ok := pm.publishers[scheduling.Publisher]
	strat := scheduling.RetryStrat

	if ok {
		retry.Do(func() error {
			fmt.Println(fmt.Sprintf("try publish to %s at %s", scheduling.Publisher, time.Now().Format(time.RFC3339Nano)))
			err := publisher.Publish(scheduling.Settings)
			if err != nil {
				fmt.Println(err.Error())
				if err.ShouldRetry() {
					return err
				}
			}
			return nil
		}, strat.Limit, strat.Timeout, strat.Exponential)
	}
}
