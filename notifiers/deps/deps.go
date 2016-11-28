package deps

import (
	"k8s.io/kubernetes/pkg/api"
)

type Notification struct {
	Cluster string
	Event   api.Event
	Level   string
	Mention string
}
