package notifiers

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	dd "github.com/zorkian/go-datadog-api"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"

	"github.com/InVisionApp/kit-overwatch/fakes/depsfakes"
	"github.com/InVisionApp/kit-overwatch/notifiers/deps"
)

var _ = Describe("NewNotifyDataDog", func() {

	Context("when called", func() {
		It("should set the client", func() {
			var fakeDataDogClient *depsfakes.FakeIDataDogClient = &depsfakes.FakeIDataDogClient{}
			notifier := New(fakeDataDogClient)
			Expect(notifier.Client).ToNot(BeNil())
		})
	})

})

var _ = Describe("NotifyDataDogSend", func() {

	var (
		fakeDataDogClient *depsfakes.FakeIDataDogClient = &depsfakes.FakeIDataDogClient{}
		notifier          *NotifyDataDog
		expectedNotifier  *deps.Notification
		actualEvent       *dd.Event
	)

	BeforeEach(func() {

		now := time.Date(2010, 1, 1, 1, 00, 00, 0, time.UTC)
		later := now.AddDate(0, 0, 1)

		notifier = New(fakeDataDogClient)
		expectedNotifier = &deps.Notification{
			Cluster: "local",
			Event: api.Event{
				Reason:  "Scheduled",
				Message: "Scheduled event message from k8s",
				Source: api.EventSource{
					Component: "default-scheduler",
					Host:      "Host",
				},
				ObjectMeta: api.ObjectMeta{
					Name:      "joebob-service-deployment-12345-fd34b.123456789",
					Namespace: "default",
				},
				FirstTimestamp: unversioned.NewTime(now),
				LastTimestamp:  unversioned.NewTime(later),
				Count:          12345,
				Type:           "Normal",
				InvolvedObject: api.ObjectReference{
					Kind: "Pod",
					Name: "joebob-service-deployment-12345-fd34b.123456789",
				},
			},
			Level:   "INFO",
			Mention: "here",
		}

		// common stub
		fakeDataDogClient.PostEventStub = func(event *dd.Event) (*dd.Event, error) {
			actualEvent = event
			actualEvent.Id = 19191
			return actualEvent, nil
		}
	})

	Context("when called", func() {
		It("should return error when notification is nil", func() {
			var n *deps.Notification
			err := notifier.Send(n)
			Expect(err).ToNot(BeNil())
		})

		It("should set alert/priority correctly", func() {
			type alert struct {
				atype    string
				priority string
			}

			levels := map[string]alert{
				"INFO":  alert{atype: "Info", priority: "low"},
				"WARN":  alert{atype: "Warning", priority: "high"},
				"ERROR": alert{atype: "Error", priority: "high"},
			}

			for k, v := range levels {
				expectedNotifier.Level = k
				err := notifier.Send(expectedNotifier)
				Expect(err).To(BeNil())
				Expect(actualEvent).ToNot(BeNil())
				Expect(actualEvent.AlertType).To(Equal(v.atype))
				Expect(actualEvent.Priority).To(Equal(v.priority))
			}
		})

		It("should format serviceName correctly when `deployment` occurs and Kind is `Pod`", func() {
			err := notifier.Send(expectedNotifier)
			Expect(err).To(BeNil())
			Expect(actualEvent.Title).To(Equal("k8s Event for [here] concerning `Scheduled` event for `joebob-service` on `local`"))
		})

		It("should format serviceName correctly when `deployment` is missing and Kind is `Pod`", func() {
			expectedNotifier.Event.ObjectMeta.Name = "joebob-service-12345-fd34b.123456789"
			expectedNotifier.Event.InvolvedObject.Name = expectedNotifier.Event.ObjectMeta.Name
			err := notifier.Send(expectedNotifier)
			Expect(err).To(BeNil())
			Expect(actualEvent.Title).To(Equal("k8s Event for [here] concerning `Scheduled` event for `joebob-service` on `local`"))
		})

		It("should pass name through serviceName when Kind is not `Pod`", func() {
			expectedNotifier.Event.InvolvedObject.Kind = "Replicator"
			expectedNotifier.Event.ObjectMeta.Name = "joebob-service-deployment-12345-fd34b.123456789"
			expectedNotifier.Event.InvolvedObject.Name = expectedNotifier.Event.ObjectMeta.Name
			err := notifier.Send(expectedNotifier)
			Expect(err).To(BeNil())
			Expect(actualEvent.Title).To(Equal("k8s Event for [here] concerning `Scheduled` event for `joebob-service-deployment-12345-fd34b.123456789` on `local`"))
		})

		It("should not panic when formatting serviceName and split array is out of bounds", func() {
			expectedNotifier.Event.ObjectMeta.Name = "joebob-service12345fd34b.123456789"
			expectedNotifier.Event.InvolvedObject.Name = expectedNotifier.Event.ObjectMeta.Name
			err := notifier.Send(expectedNotifier)
			Expect(err).To(BeNil())
			Expect(actualEvent.Title).To(Equal("k8s Event for [here] concerning `Scheduled` event for `joebob-service12345fd34b.123456789` on `local`"))
		})

		It("should format title differently if no mention", func() {
			expectedNotifier.Mention = ""
			err := notifier.Send(expectedNotifier)
			Expect(err).To(BeNil())
			Expect(actualEvent.Title).To(Equal("`Scheduled` event for `joebob-service` on `local`"))
		})

		It("should error when datadog errors", func() {
			fakeDataDogClient.PostEventStub = func(event *dd.Event) (*dd.Event, error) {
				actualEvent = nil
				return nil, fmt.Errorf("This is an error.")
			}
			err := notifier.Send(expectedNotifier)
			Expect(err).ToNot(BeNil())
			Expect(actualEvent).To(BeNil())
		})

	})

})
