package watcher

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/types"
	"k8s.io/kubernetes/pkg/watch"

	"github.com/InVisionApp/kit-overwatch/config"
	"github.com/InVisionApp/kit-overwatch/notifiers"
	"github.com/InVisionApp/kit-overwatch/notifiers/deps"
)

type Watcher struct {
	Client client.Client
	ClientConfig restclient.Config
	Config config.Config
}

type WatcherEvent struct {
	Event   api.Event
	WatchEvent watch.Event
}

type SentEvent struct {
	LastSent time.Time
	Count int
}

func New(cfg *config.Config) *Watcher {
	var c *client.Client
	var cErr error
	var clientConfig *restclient.Config

	if cfg.InCluster {
		var confErr error
		clientConfig, confErr = restclient.InClusterConfig()
		if confErr != nil {
			log.Fatalf("Unable to instantiate in cluster config: %v", cErr.Error())
		}
		c, cErr = client.New(clientConfig)
		if cErr != nil {
			log.Fatalf("Unable to instantiate kube client in cluster: %v", cErr.Error())
		}
	} else {
		clientConfig = &restclient.Config{
			Host: cfg.ClusterHost,
		}
		c, cErr = client.New(clientConfig)
		if cErr != nil {
			log.Fatalf("Unable to instantiate kube client: %v", cErr.Error())
		}
	}

	return &Watcher{
		Client: *c,
		ClientConfig: *clientConfig,
		Config: *cfg,
	}
}

func (w *Watcher) Watch() {
	startTime := time.Now()
	pastEvents := make(map[types.UID]WatcherEvent)
	sentEvents := make(map[types.UID]SentEvent)

	opts := api.ListOptions{
		ResourceVersion: "0",
	}
	cw, err := w.Client.Events(w.Config.Namespace).Watch(opts)
	if err != nil {
		log.Fatalf("Unable to instantiate events watcher: %v", err.Error())
	}

	// Get the events channel
	ec := cw.ResultChan()
	log.Info("Watching for events...")

	// Process event channel
	for we := range ec {
		log.Infof("%s event detected", we.Type)

		// When an event occurs, get list of events
		list, err := w.Client.Events(w.Config.Namespace).List(api.ListOptions{
			ResourceVersion: "0",
		})
		if err != nil {
			log.Fatalf("Unable to get events: %v", err.Error())
		}
		for _, e := range list.Items {
			// Only log if we haven't logged before
			past, ok := pastEvents[e.ObjectMeta.UID]
			if ok {
				// Only skip if the count hasn't increased
				if e.Count == past.Event.Count {
					log.Debugf("Skip: already notified for %s / %s / %s", e.ObjectMeta.UID, e.Reason, e.Message)
					continue
				}
			}

			// Remember this event so we don't send duplicate notifications
			pastEvents[e.ObjectMeta.UID] = WatcherEvent{
				Event: e,
				WatchEvent:   we,
			}

			// Only log events that have happened since the service started
			diff := startTime.Sub(e.LastTimestamp.Time)
			if int(diff.Minutes()) > 1 {
				log.Debugf("Skip: %s / %s / %s - %s happened more than a minute before service started", e.ObjectMeta.UID, e.Reason, e.Message, e.LastTimestamp)
				continue
			}

			// Throttle duplicate events so we don't notify too many times
			sent, ok := sentEvents[e.ObjectMeta.UID]
			var count int
			if ok {
				canSendAfter := sent.LastSent.Add(time.Minute * time.Duration(sent.Count))
				if time.Now().After(canSendAfter) {
					count = sent.Count + 1
				} else {
					log.Debugf("Skip: throttle back notifications %v minutes for %s / %s / %s", sent.Count, e.ObjectMeta.UID, e.Reason, e.Message)
					continue
				}
			}
			sentEvents[e.ObjectMeta.UID] = SentEvent{
				LastSent: time.Now(),
				Count: count,
			}

			// Generate and send the notification
			go w.notify(e)
		}
	}

	log.Fatalf("Event watching has ended")
}

func (w *Watcher) getLevel(e api.Event) string {
	reasonLevels := map[string]string{
		"SuccessfulCreate": "INFO",
		"SuccessfulDelete": "INFO",
		"ContainerCreating": "INFO",
		"Pulled": "INFO",
		"Pulling": "INFO",
		"Created": "INFO",
		"Starting": "INFO",
		"Started": "INFO",
		"Killing": "INFO",
		"NodeReady": "INFO",
		"ScalingReplicaSet": "INFO",
		"Scheduled": "INFO",
		"NodeNotReady": "WARN",
		"MAPPING": "WARN",
		"UPDATE": "INFO",
		"DELETE": "INFO",
		"NodeOutOfDisk": "ERROR",
		"BackOff": "ERROR",
		"ImagePullBackOff": "ERROR",
		"FailedSync": "ERROR",
		"FreeDiskSpaceFailed": "WARN",
		"MissingClusterDNS": "ERROR",
		"RegisteredNode": "INFO",
		"TerminatingEvictedPod": "WARN",
		"RemovingNode": "WARN",
		"TerminatedAllPods": "WARN",
		"CreatedLoadBalancer": "INFO",
		"CreatingLoadBalancer": "INFO",
		"NodeHasSufficientDisk": "INFO",
		"NodeHasSufficientMemory": "INFO",
		"NodeNotSchedulable": "ERROR",
		"DeletingAllPods": "WARN",
		"DeletingNode": "WARN",
		"UpdatedLoadBalancer": "INFO",
	}

	var ok bool
	if _, ok = reasonLevels[e.Reason]; !ok {
		return "ERROR"
	}

	return reasonLevels[e.Reason]
}

func (w *Watcher) notify(e api.Event) {
	// Determine notification level
	level := w.getLevel(e)

	// Get label to use as mention in notification
	var mention string
	var rErr error
	var rOk bool
	ec, err := client.NewExtensions(&w.ClientConfig)
	if err != nil {
		log.Fatalf("Unable to instantiate new ExtensionsClient: %v", err.Error())
	}
	switch e.InvolvedObject.Kind {
	case "Pod":
		var resource *api.Pod
		resource, rErr = w.Client.Pods(w.Config.Namespace).Get(e.InvolvedObject.Name)
		mention, rOk = resource.ObjectMeta.Labels[w.Config.MentionLabel]
	case "Service":
		var resource *api.Service
		resource, rErr = w.Client.Services(w.Config.Namespace).Get(e.InvolvedObject.Name)
		mention, rOk = resource.ObjectMeta.Labels[w.Config.MentionLabel]
	case "Node":
		var resource *api.Node
		resource, rErr = w.Client.Nodes().Get(e.InvolvedObject.Name)
		mention, rOk = resource.ObjectMeta.Labels[w.Config.MentionLabel]
	case "Deployment":
		var resource *extensions.Deployment
		resource, rErr = ec.Deployments(w.Config.Namespace).Get(e.InvolvedObject.Name)
		mention, rOk = resource.ObjectMeta.Labels[w.Config.MentionLabel]
	case "ReplicaSet":
		var resource *extensions.ReplicaSet
		resource, rErr = ec.ReplicaSets(w.Config.Namespace).Get(e.InvolvedObject.Name)
		mention, rOk = resource.ObjectMeta.Labels[w.Config.MentionLabel]
	case "Job":
		var resource *batch.Job
		resource, rErr = ec.Jobs(w.Config.Namespace).Get(e.InvolvedObject.Name)
		mention, rOk = resource.ObjectMeta.Labels[w.Config.MentionLabel]
	case "DaemonSet":
		var resource *extensions.DaemonSet
		resource, rErr = ec.DaemonSets(w.Config.Namespace).Get(e.InvolvedObject.Name)
		mention, rOk = resource.ObjectMeta.Labels[w.Config.MentionLabel]
	default:
		log.Debugf("Cannot retrieve label for unsported Kind: %s", e.InvolvedObject.Kind)
	}
	if rErr != nil {
		log.Warnf("Unable to get %s: %v", e.InvolvedObject.Kind, rErr.Error())
	}
	if !rOk {
		log.Warnf("Mention label not found for %s: %s, using default: %s", e.InvolvedObject.Kind, e.InvolvedObject.Name, w.Config.MentionDefault)
		mention = w.Config.MentionDefault
	}

	// Send notifications
	n := notifiers.New(&w.Config)
	notification := deps.Notification{
		Cluster: w.Config.ClusterName,
		Event:   e,
		Level:   level,
		Mention: mention,
	}
	n.SendAll(&notification)
}
