package notifiers

import (
	log "github.com/Sirupsen/logrus"

	"github.com/InVisionApp/kit-overwatch/config"
	"github.com/InVisionApp/kit-overwatch/notifiers/deps"
	notifyLog "github.com/InVisionApp/kit-overwatch/notifiers/log"
	notifySlack "github.com/InVisionApp/kit-overwatch/notifiers/slack"
)

type Notifiers struct {
	Config config.Config
}

func New(cfg *config.Config) *Notifiers {
	return &Notifiers{
		Config: *cfg,
	}
}

func (notifiers *Notifiers) SendAll(n *deps.Notification) {
	// Only send notification if it's a desired Level
	levels := [...]string{"DEBUG", "INFO", "WARN", "ERROR"}
	send := false
	switch notifiers.Config.NotificationLevel {
	case "DEBUG":
		if stringInSlice(n.Level, levels[0:]) {
			send = true
		}
	case "INFO":
		if stringInSlice(n.Level, levels[1:]) {
			send = true
		}
	case "WARN":
		if stringInSlice(n.Level, levels[2:]) {
			send = true
		}
	case "ERROR":
		if stringInSlice(n.Level, levels[3:]) {
			send = true
		}
	}

	if send {
		if notifiers.Config.NotifyLog {
			err := notifyLog.Send(n)
			if err != nil {
				log.Fatalf("NotifyLog Error: %v", err.Error())
			}
		}
		if notifiers.Config.NotifySlack {
			ns := notifySlack.New(notifiers.Config.NotifySlackToken, notifiers.Config.NotifySlackChannel, notifiers.Config.NotifySlackAsUser)
			err := ns.Send(n)
			if err != nil {
				log.Fatalf("NotifySlack Error: %v", err.Error())
			}
		}
	} else {
		log.Debugf("Skipping because %s is not within NotificationLevel: %s / %s / %s / %s", n.Level, n.Cluster, n.Event.Reason, n.Event.Message, n.Event.LastTimestamp)
	}
}

// For finding a string in an array
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
