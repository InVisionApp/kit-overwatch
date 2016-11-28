<p align="center">
  <a href="https://invisionapp.github.io/kit/">
    <img src="https://github.com/InVisionApp/kit/raw/master/media/kit-logo-horz-sm.png">
  </a>
</p>

# kit-overwatch

[![Docker Repository on Quay](https://quay.io/repository/invision/kit-overwatch/status "Docker Repository on Quay")](https://quay.io/repository/invision/kit-overwatch)
[ ![Codeship Status for InVisionApp/kit-overwatch](https://codeship.com/projects/1232c860-4089-0134-3306-1a3f66148744/status?branch=develop)](https://codeship.com/projects/167729)

Monitors events within a Kubernetes cluster and sends notifications to Slack. You can adjust the events you want to be notified about and it can automatically handle @mentions based on Kubernetes labels.

To avoid flooding channels, duplicate events will be throttled back automatically by 1 minute * the number of notifications already sent.

## API Usage

### API Endpoints

#### `GET /version`
+ **Description**: Display version of `auth-api` service
+ **On success**:
  * Status: `200`
  * Response: text string
+ **On failure**:
  * No failure case

----------------------------------------------------

#### `GET /healthcheck`
+ **Description**: Check if service is healthy/up
+ **On success**:
  * Status: `200`
  * Response: text string
+ **On failure**:
  * Status: `400`
  * Response: JSON error blob

----------------------------------------------------

## Expected environment variables

The following environment variables are used by this service.

| Variable | Description | Required | Default |
| :--- | :--- | :--- | :--- |
| `KIT_OVERWATCH_DEBUG` | Enables debug logging | yes | `false` |
| `KIT_OVERWATCH_LISTEN_ADDRESS` | The port the service listens on | yes | `:80` |
| `KIT_OVERWATCH_STATSD_ADDRESS` | The statsd address | yes | `localhost:8125` |
| `KIT_OVERWATCH_STATSD_PREFIX` | The statsd prefix | yes | `statsd.kit-overwatch.dev` |
| `KIT_OVERWATCH_NAMESPACE` | The namespace to watch events on | yes | `default` |
| `KIT_OVERWATCH_IN_CLUSTER` | Enable when deployed in a Kubernetes cluster to automatically watch events in that cluster | yes | `true` |
| `KIT_OVERWATCH_CLUSTER_NAME` | This name is displayed in all the notifications generated | false | `Kubernetes` |
| `KIT_OVERWATCH_CLUSTER_HOST` | The address to the cluster. Only needed when using KIT_OVERWATCH_IN_CLUSTER=false | false | *empty* |
| `KIT_OVERWATCH_NOTIFICATION_LEVEL` | Determines what level of events you want to be notified about. Goes from `DEBUG` -> `INFO` -> `WARN` -> `ERROR` | false | `INFO` |
| `KIT_OVERWATCH_MENTION_LABEL` | Will use this label found on a resource as a mention in the notification | false | *empty* |
| `KIT_OVERWATCH_MENTION_DEFAULT` | If no KIT_OVERWATCH_MENTION_LABEL is found, it will default to using this as a mention in the notification | false | `here` |
| `KIT_OVERWATCH_NOTIFY_LOG` | Enable to send a notification to stdout | true | `true` |
| `KIT_OVERWATCH_NOTIFY_SLACK` | Enable to send a notification to slack | true | `false` |
| `KIT_OVERWATCH_NOTIFY_SLACK_TOKEN` | The auth token for Slack. Required if KIT_OVERWATCH_NOTIFY_SLACK=true | false | *empty* |
| `KIT_OVERWATCH_NOTIFY_SLACK_AS_USER` | Enable to send a notification to slack as the given user associated with the token | true | `false` |
| `KIT_OVERWATCH_NOTIFY_SLACK_CHANNEL` | The Slack channel to post notifications to. Required if KIT_OVERWATCH_NOTIFY_SLACK=true | false | *empty* |
| `KIT_OVERWATCH_NOTIFY_DATADOG` | Enable to send an event to DataDog | true | `false` |
| `KIT_OVERWATCH_NOTIFY_DATADOG_APIKEY` | The apikey for DataDog. Required if KIT_OVERWATCH_NOTIFY_DATADOG=true | false | *empty* |
| `KIT_OVERWATCH_NOTIFY_DATADOG_APPKEY` | The appkey for DataDog. Required if KIT_OVERWATCH_NOTIFY_DATADOG=true | false | *empty* |


## How to run locally

This requires that you have `Go` installed locally.

1. Run `make run`
1. Test by running `curl -I -X GET localhost:8080/healthcheck` and you should get a 200 response

If you would like to change the ENV var settings from their defaults, you can copy `.env.dist` to `.env` and modify them there. Then just re-run `make run`.

## How to test Codeship build

Install Codeship's `jet` tool then run:

```
$ jet steps
```

## Testing

1. Install [minikube](https://github.com/kubernetes/minikube)
1. Install [kubectl](https://github.com/kubernetes/minikube)
1. Run `minikube start`
1. Start a proxy `kubectl proxy`
1. Run `make run`
1. Now perform actions on your local kubernetes cluster that will generate events

## Limitations

- You cannot run more than one instance of this service within a Cluster or you'll end up with duplication notifications
- To avoid duplicate notifications being sent, the service will only send notifications for past events that have happened 1 minute before the service was started

## TODO

- [x] Use Minikube for local functional testing
- [x] Allow setting name of cluster to be included in notifications
- [x] Get events from cluster
- [x] Log events to stdout
- [x] Add Slack notification of events
- [x] Add DataDog notification of events
- [ ] Add Unit Tests for Watcher
- [ ] Add Unit Tests for Notifiers
- [ ] Add Unit Tests for Slack Notifier
- [ ] Add Unit Tests for Log Notifier
- [ ] Use datastore so we can make this service stateless
- [ ] Split watcher and notifiers to allow multiple instances using task queue
