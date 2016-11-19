package notifiers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/InVisionApp/go-datadog-api"
	"github.com/InVisionApp/kit-overwatch/notifiers/deps"
	log "github.com/Sirupsen/logrus"
	"strings"
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
		event.AlertType = "Info"
		event.Priority = "low"
	case "WARN":
		event.AlertType = "Warning"
		event.Priority = "high"
	case "ERROR":
		event.AlertType = "Error"
		event.Priority = "high"
	}

	serviceName := n.Event.ObjectMeta.Name
	splitName := strings.Split(n.Event.ObjectMeta.Name, "-")

	switch n.Event.InvolvedObject.Kind {
	case "Pod":
		// Default Rules to ignore `service-name-[deploynumber-podnumber.number]`
		indexToIgnore := 2

		if len(splitName) > 3 && splitName[len(splitName)-3] == "deployment" {
			// Rules to ignore `service-name-[deployment-deploynumber-podnumber.number]`
			indexToIgnore = 3
		}
		serviceName = strings.Join(splitName[:len(splitName)-indexToIgnore], "-")

	}

	title := fmt.Sprintf("`%s` event for `%s` on `%s`", n.Event.Reason, serviceName, n.Cluster)
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
		"object-kind:" + n.Event.InvolvedObject.Kind,
		"object-name:" + n.Event.InvolvedObject.Name,
		"mentioned:" + n.Mention,
		"service:" + serviceName,
	}

	message := `#### Message Details
	%v
#### Event Details
	%v
#### Involved Object
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

	mDetailsJson, err := json.Marshal(mDetails)
	if err != nil {
		return fmt.Errorf("%s\n", err)
	}
	mDetailsFmt := string(mDetailsJson)

	eDetailsJson, err := json.Marshal(eDetails)
	if err != nil {
		return fmt.Errorf("%s\n", err)
	}
	eDetailsFmt := string(eDetailsJson)

	iObjectJson, err := json.Marshal(iObject)
	if err != nil {
		return fmt.Errorf("%s\n", err)
	}
	iObjectFmt := string(iObjectJson)

	event.Text = MD_PREFIX + fmt.Sprintf(message, mDetailsFmt, eDetailsFmt, iObjectFmt) + MD_SUFFIX

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
