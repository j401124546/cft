package monitor

import "cft/model"

type Monitor struct {
	containers map[string]*model.Container
}

var _monitor = NewMonitor()

func AddContainer(id string, name string, host string, state model.StateType) {
	if _, ok := _monitor.containers[id]; !ok {
		container := model.NewContainer(id, name, host, state)
		_monitor.containers[id] = container
	}
}

func GetContainer(id string) *model.Container {
	if _, ok := _monitor.containers[id]; ok {
		return _monitor.containers[id]
	}
	return nil
}

func NewMonitor() *Monitor {
	return &Monitor{
		containers: make(map[string]*model.Container),
	}
}
