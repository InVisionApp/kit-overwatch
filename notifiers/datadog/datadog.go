package notifiers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/InVisionApp/go-datadog-api"
	"github.com/InVisionApp/kit-overwatch/notifiers/deps"
	log "github.com/Sirupsen/logrus"
)

const EVENT_TYPE = "kubernetes"
const MD_PREFIX = "%%% \n"
const MD_SUFFIX = "\n %%%"

type NotifyDataDog struct {
	Client *datadog.Client
}

func New(apiKey string, appKey string) *NotifyDataDog {
	client := datadog.NewClient(apiKey, appKey)
	return &NotifyDataDog{
		Client: client,
	}
}

func (ndd *NotifyDataDog) Send(n *deps.Notification) error {

	event := &datadog.Event{
		SourceType: EVENT_TYPE,
		EventType:  EVENT_TYPE,
	}

	switch n.Level {
	case "INFO":
		event.AlertType = "info"
		event.Priority = "low"
	case "WARN":
		event.AlertType = "warning"
		event.Priority = "normal"
	case "ERROR":
		event.AlertType = "error"
		event.Priority = "normal"
	}

	title := fmt.Sprintf("`%s` event for `%s` on `%s`", n.Event.Reason, n.Event.ObjectMeta.Name, n.Cluster)
	if n.Mention != "" {
		title = fmt.Sprintf("Alerting [%s] concerning %s", n.Mention, title)
	}

	event.Title = title
	event.Tags = []string{
		"team:" + n.Mention,
		"cluster:" + n.Cluster,
		"type:" + n.Event.Type,
		"level:" + n.Level,
		"reason:" + n.Event.Reason,
		"node:" + n.Event.Source.Host,
		"name:" + n.Event.ObjectMeta.Name,
		"namespace:" + n.Event.ObjectMeta.Namespace,
		"component:" + n.Event.Source.Component,
		"count:" + fmt.Sprintf("%d", n.Event.Count),
		"involved-object-kind:" + n.Event.InvolvedObject.Kind,
		"involved-object-name:" + n.Event.InvolvedObject.Name,
		"mentioned:" + n.Mention,
	}

	message := `	%v
### Message Details
	%v
### Event Details
	%v
### Involved Object
	%v
`

	mDetails := &messageDetails{
		Cluster: n.Cluster,
		Reason:  n.Event.Reason,
		Type:    n.Event.Type,
		Level:   n.Level,
	}

	eDetails := &eventDetails{
		Node:           n.Event.Source.Host,
		Namespace:      n.Event.ObjectMeta.Namespace,
		Component:      n.Event.Source.Component,
		Count:          fmt.Sprintf("%d", n.Event.Count),
		FirstOccurance: n.Event.FirstTimestamp.Format(time.RFC1123),
		LastOccurance:  n.Event.LastTimestamp.Format(time.RFC1123),
	}

	iObject := &involvedObject{
		Kind: n.Event.InvolvedObject.Kind,
		Name: n.Event.InvolvedObject.Name,
	}

	mDetailsJson, err := json.MarshalIndent(mDetails, "", "  ")
	if err != nil {
		return fmt.Errorf("%s\n", err)
	}

	eDetailsJson, err := json.MarshalIndent(eDetails, "", "  ")
	if err != nil {
		return fmt.Errorf("%s\n", err)
	}

	iObjectJson, err := json.MarshalIndent(iObject, "", "  ")
	if err != nil {
		return fmt.Errorf("%s\n", err)
	}

	event.Text = MD_PREFIX + fmt.Sprintf(message, title, mDetailsJson, eDetailsJson, iObjectJson) + MD_SUFFIX

	newEvent, err := ndd.Client.PostEvent(event)
	if err != nil {
		return fmt.Errorf("Post Event error: %s\n", err)
	}

	log.Infof("NotifyDataDog: %s / %s / %s / %s", n.Event.Reason, n.Event.Message, newEvent.Id, newEvent.Title)
	return nil
}

type messageDetails struct {
	Cluster string `json:"cluster,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Type    string `json:"type,omitempty"`
	Level   string `json:"level,omitempty"`
}

type eventDetails struct {
	Node           string `json:"node,omitempty"`
	Namespace      string `json:"namespace,omitempty"`
	Component      string `json:"component,omitempty"`
	Count          string `json:"count,omitempty"`
	FirstOccurance string `json:"first_occurance,omitempty"`
	LastOccurance  string `json:"last_occurance,omitempty"`
}

type involvedObject struct {
	Kind string `json:"kind,omitempty"`
	Name string `json:"name,omitempty"`
}
