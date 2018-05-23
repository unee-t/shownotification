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
