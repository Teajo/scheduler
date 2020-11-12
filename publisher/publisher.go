package publisher

import (
	"fmt"
	"jpb/scheduler/config"
	"jpb/scheduler/utils"
	"time"
)

// Publisher interface
type Publisher interface {
	CheckConfig(map[string]string) error
	Publish(map[string]string) error
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
		pm.publish(scheduling)
	}
}

// Get retrieves a publisher according to provided id
func (pm *PubManager) Get(id string) (Publisher, bool) {
	pub, ok := pm.publishers[id]
	return pub, ok
}

func (pm *PubManager) publish(scheduling *utils.Scheduling) error {
	publisher, ok := pm.publishers[scheduling.Publisher]
	if ok {
		fmt.Println(fmt.Sprintf("publish to %s at %s", scheduling.Publisher, scheduling.Date.Format(time.RFC3339Nano)))
		err := publisher.Publish(scheduling.Settings)
		if err != nil {
			return err
		}
	}

	return nil
}
