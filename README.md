# ccloud-tele

A Confluent Cloud Telemetry CLI Tool for issuing queries against the Telemetry API


### Installing the CLI

```shell
go get github.com/nerdynick/cccloud-tele
```

### Usage

The CLI is coded to be interacted with similar to GIT.
It leverages the project [Cobra](https://github.com/spf13/cobra) to provide this experience.

Command Tree

* list
  * resources - Query Available Resources
  * metrics - Query Available Metrics
  * topics - Query Available Topics for a Metric
* query
  * metric - Query Data for a metric
  * metrics - Query Data for multipule metrics
  * topic - Query Data for a metric and topic
  * topics - Query Data for a metric and a list of topics
    * all - Query Data for a metric and all available topics (Queries run in parallel for each topic)

**List Available Resources**

```shell
./ccloud-tele list resources --apikey MY-KEY --apisecret MY-SECRET
```

**List Available Metrics for Kafka Cluster**

```shell
./ccloud-tele list metrics --apikey MY-KEY --apisecret MY-SECRET --kafka
```

**List Available Topics for a given Metric**

```shell
./ccloud-tele list topics --apikey MY-KEY --apisecret MY-SECRET --cluster MY-CLUSTER-ID  io.confluent.kafka.server/retained_bytes
```

**Query Cluster level results for a given Metric**

```shell
./ccloud-tele query metric --apikey MY-KEY --apisecret MY-SECRET --cluster MY-CLUSTER-ID  io.confluent.kafka.server/retained_bytes
```

**Query Topic level results for a given Metric & Topic**

```shell
./ccloud-tele query topic --apikey MY-KEY --apisecret MY-SECRET --cluster MY-CLUSTER-ID --metric io.confluent.kafka.server/retained_bytes  MY-TOPIC
```

**Query Topic Parition level results for a given Metric & Topic**

```shell
./ccloud-tele query topic --apikey MY-KEY --apisecret MY-SECRET --cluster MY-CLUSTER-ID --metric io.confluent.kafka.server/retained_bytes MY-TOPIC --partitions
```