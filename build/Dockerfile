FROM scratch

# Add the go binary
COPY kit-overwatch-linux ./kit-overwatch

# Add the certs needed to communicate with Slack
COPY /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Set default ENV vars as documented in README.md
ENV KIT_OVERWATCH_DEBUG=false\
    KIT_OVERWATCH_LISTEN_ADDRESS=:80\
    KIT_OVERWATCH_STATSD_ADDRESS=localhost:8125\
    KIT_OVERWATCH_STATSD_PREFIX=statsd.kit-overwatch.dev\
    KIT_OVERWATCH_NAMESPACE=default\
    KIT_OVERWATCH_IN_CLUSTER=true\
    KIT_OVERWATCH_CLUSTER_NAME=Kubernetes\
    KIT_OVERWATCH_CLUSTER_HOST=\
    KIT_OVERWATCH_NOTIFICATION_LEVEL=INFO\
    KIT_OVERWATCH_MENTION_LABEL=\
    KIT_OVERWATCH_MENTION_DEFAULT=here\
    KIT_OVERWATCH_NOTIFY_LOG=true\
    KIT_OVERWATCH_NOTIFY_SLACK=false\
    KIT_OVERWATCH_NOTIFY_SLACK_TOKEN=\
    KIT_OVERWATCH_NOTIFY_SLACK_AS_USER=false\
    KIT_OVERWATCH_NOTIFY_SLACK_CHANNEL=

CMD ["/kit-overwatch"]
