package notifiers

import (
	"fmt"
	"github.com/InVisionApp/kit-overwatch/notifiers/deps"
	log "github.com/Sirupsen/logrus"
)

func Send(n *deps.Notification) error {
	message := fmt.Sprintf("NotifyLog: %s / %s / %s / %s", n.Cluster, n.Event.Reason, n.Event.Message, n.Event.LastTimestamp)

	// Add mention if one exists
	if n.Mention != "" {
		message = fmt.Sprintf("%s / @%s", message, n.Mention)
	}

	switch n.Level {
	case "DEBUG":
		log.Debugf(message)
	case "INFO":
		log.Infof(message)
	case "WARN":
		log.Warnf(message)
	case "ERROR":
		log.Errorf(message)
	default:
		return fmt.Errorf("Invalid Notification.Level provided")
	}
	return nil
}
