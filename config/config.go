package config

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/caarlos0/env.v2"
)

type Config struct {
	Debug                      bool     `env:"KIT_OVERWATCH_DEBUG" envDefault:"true"`
	ListenAddress              string   `env:"KIT_OVERWATCH_LISTEN_ADDRESS" envDefault:":8080"`
	StatsDAddress              string   `env:"KIT_OVERWATCH_STATSD_ADDRESS" envDefault:"localhost:8125"`
	StatsDPrefix               string   `env:"KIT_OVERWATCH_STATSD_PREFIX" envDefault:"statsd.kit-overwatch.dev"`
	Namespace                  string   `env:"KIT_OVERWATCH_NAMESPACE" envDefault:"default"`
	InCluster                  bool     `env:"KIT_OVERWATCH_IN_CLUSTER" envDefault:"false"`
	ClusterName                string   `env:"KIT_OVERWATCH_CLUSTER_NAME" envDefault:"local"`
	ClusterHost                string   `env:"KIT_OVERWATCH_CLUSTER_HOST" envDefault:"http://127.0.0.1:8001"`
	NotificationLevel          string   `env:"KIT_OVERWATCH_NOTIFICATION_LEVEL" envDefault:"DEBUG"`
	MentionLabel               string   `env:"KIT_OVERWATCH_MENTION_LABEL" envDefault:""`
	MentionDefault             string   `env:"KIT_OVERWATCH_MENTION_DEFAULT" envDefault:"here"`
	NotifyLog                  bool     `env:"KIT_OVERWATCH_NOTIFY_LOG" envDefault:"true"`
	NotifySlack                bool     `env:"KIT_OVERWATCH_NOTIFY_SLACK" envDefault:"false"`
	NotifySlackToken           string   `env:"KIT_OVERWATCH_NOTIFY_SLACK_TOKEN" envDefault:""`
	NotifySlackAsUser          bool     `env:"KIT_OVERWATCH_NOTIFY_SLACK_AS_USER" envDefault:"false"`
	NotifySlackChannel         string   `env:"KIT_OVERWATCH_NOTIFY_SLACK_CHANNEL" envDefault:""`
}

func New() *Config {
	return &Config{}
}

func (c *Config) LoadEnvVars() error {
	if err := env.Parse(c); err != nil {
		return fmt.Errorf("Unable to fetch env vars: %v", err.Error())
	}

	var errorList []string

	// Verify we have a valid listen address
	matched, err := regexp.MatchString(`^(?:[^\s]+)?:\d+$`, c.ListenAddress)
	if err != nil {
		errorList = append(errorList, fmt.Sprintf("error parsing 'KIT_OVERWATCH_LISTEN_ADDRESS': %v", err.Error()))
	}

	if !matched {
		errorList = append(errorList, "invalid 'KIT_OVERWATCH_LISTEN_ADDRESS'")
	}

	if len(errorList) != 0 {
		return fmt.Errorf(strings.Join(errorList, "; "))
	}

	return nil
}
