package notifiers

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"

	"github.com/InVisionApp/kit-overwatch/notifiers/deps"
)

type NotifySlack struct {
	Token   string
	Channel string
	AsUser  bool
}

func New(token string, channel string, asUser bool) *NotifySlack {
	return &NotifySlack{
		Token:   token,
		Channel: channel,
		AsUser:  asUser,
	}
}

func (ns *NotifySlack) Send(n *deps.Notification) error {
	api := slack.New(ns.Token)
	params := slack.PostMessageParameters{
		Username:  "kit-overwatch",
		LinkNames: 1,
		AsUser:    ns.AsUser,
	}
	eventAttachment := slack.Attachment{
		Fallback: n.Event.Message,
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Message",
				Value: n.Event.Message,
				Short: false,
			},
			slack.AttachmentField{
				Title: "Cluster",
				Value: n.Cluster,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Reason",
				Value: n.Event.Reason,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Type",
				Value: n.Event.Type,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Level",
				Value: n.Level,
				Short: true,
			},
		},
	}

	// Determine slack color to use for event attachment based on Level
	switch n.Level {
	case "INFO":
		eventAttachment.Color = "good"
	case "WARN":
		eventAttachment.Color = "warning"
	case "ERROR":
		eventAttachment.Color = "danger"
	}

	eventDetailsAttachment := slack.Attachment{
		Fallback: n.Event.Message,
		Title:    "Details",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Name",
				Value: n.Event.ObjectMeta.Name,
				Short: false,
			},
			slack.AttachmentField{
				Title: "Node",
				Value: n.Event.Source.Host,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Namespace",
				Value: n.Event.ObjectMeta.Namespace,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Component",
				Value: n.Event.Source.Component,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Count",
				Value: fmt.Sprintf("%d", n.Event.Count),
				Short: true,
			},
			slack.AttachmentField{
				Title: "First Occurrence",
				Value: n.Event.FirstTimestamp.Format(time.RFC1123),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Lastest Occurrence",
				Value: n.Event.LastTimestamp.Format(time.RFC1123),
				Short: true,
			},
		},
	}

	involvedObjectAttachment := slack.Attachment{
		Fallback: fmt.Sprintf("%s %s", n.Event.InvolvedObject.Kind, n.Event.InvolvedObject.Name),
		Title:    "Involved Object",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Kind",
				Value: n.Event.InvolvedObject.Kind,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Name",
				Value: n.Event.InvolvedObject.Name,
				Short: true,
			},
		},
	}

	params.Attachments = []slack.Attachment{eventAttachment, eventDetailsAttachment, involvedObjectAttachment}
	message := fmt.Sprintf("`%s` event for `%s` on `%s`", n.Event.Reason, n.Event.ObjectMeta.Name, n.Cluster)

	// Add mention handle if there is one
	if n.Mention != "" {
		message = fmt.Sprintf("Alerting @%s concerning %s", n.Mention, message)
	}

	channelID, timestamp, err := api.PostMessage(ns.Channel, message, params)
	if err != nil {
		return fmt.Errorf("%s\n", err)
	}
	log.Infof("NotifySlack: %s / %s / %s / %s", n.Event.Reason, n.Event.Message, channelID, timestamp)
	return nil
}
