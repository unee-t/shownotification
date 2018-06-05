Originally this project set out to consume AWS SQS messages. This approach
doesn't work since you can't have more than one listener on a SQS queue. Please
see https://www.youtube.com/watch?v=bKFZVNzloiA for more details on the topic.

I've since accepted a contribution from this [SO
answer](https://stackoverflow.com/a/50498987/4534) which forms a better /
simpler example of how you could do it.

* [On Docker Hub](https://hub.docker.com/r/uneet/showhook/tags/)

# EventSource Server

Server Sent Events (SSE) allows clients to receive notifications.

### Build and Setup Network

```shell
make
make network
```

### Start Cluster of Two Linked Systems

You may set this using service `etcd` discovery.

```shell
PORT=9000 NEIGHBORS=192.168.0.1:9001 make start
PORT=9001 NEIGHBORS=192.168.0.1:9000 make start
```

You can add more neighbors with a `,` comma delimiter.
