package deps

import (
	dd "github.com/zorkian/go-datadog-api"
)

//go:generate counterfeiter -o ../fakes/depsfakes/fake_idatadogclient.go . IDataDogClient

// Interface for faking the DataDog Client
type IDataDogClient interface {
	PostEvent(event *dd.Event) (*dd.Event, error)
}
