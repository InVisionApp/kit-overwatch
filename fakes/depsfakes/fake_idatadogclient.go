// This file was generated by counterfeiter
package depsfakes

import (
	"sync"

	"github.com/InVisionApp/kit-overwatch/deps"
	go_datadog_api "github.com/zorkian/go-datadog-api"
)

type FakeIDataDogClient struct {
	PostEventStub        func(event *go_datadog_api.Event) (*go_datadog_api.Event, error)
	postEventMutex       sync.RWMutex
	postEventArgsForCall []struct {
		event *go_datadog_api.Event
	}
	postEventReturns struct {
		result1 *go_datadog_api.Event
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeIDataDogClient) PostEvent(event *go_datadog_api.Event) (*go_datadog_api.Event, error) {
	fake.postEventMutex.Lock()
	fake.postEventArgsForCall = append(fake.postEventArgsForCall, struct {
		event *go_datadog_api.Event
	}{event})
	fake.recordInvocation("PostEvent", []interface{}{event})
	fake.postEventMutex.Unlock()
	if fake.PostEventStub != nil {
		return fake.PostEventStub(event)
	} else {
		return fake.postEventReturns.result1, fake.postEventReturns.result2
	}
}

func (fake *FakeIDataDogClient) PostEventCallCount() int {
	fake.postEventMutex.RLock()
	defer fake.postEventMutex.RUnlock()
	return len(fake.postEventArgsForCall)
}

func (fake *FakeIDataDogClient) PostEventArgsForCall(i int) *go_datadog_api.Event {
	fake.postEventMutex.RLock()
	defer fake.postEventMutex.RUnlock()
	return fake.postEventArgsForCall[i].event
}

func (fake *FakeIDataDogClient) PostEventReturns(result1 *go_datadog_api.Event, result2 error) {
	fake.PostEventStub = nil
	fake.postEventReturns = struct {
		result1 *go_datadog_api.Event
		result2 error
	}{result1, result2}
}

func (fake *FakeIDataDogClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.postEventMutex.RLock()
	defer fake.postEventMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeIDataDogClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ deps.IDataDogClient = new(FakeIDataDogClient)