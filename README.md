# prometheus-pulsar-adapter
Use Pulsar as a remote storage database for Prometheus  (remote write only)

Prometheus-pulsar-adapter is a service which receives [Prometheus](https://github.com/prometheus) metrics through [`remote_write`](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write), marshal into JSON and sends them into [Pulsar](https://github.com/apache/pulsar).  

`Prometheus-pulsar-adapter`  uses [Prometheus-kafka-adapter](https://github.com/Telefonica/prometheus-kafka-adapter) code to a large extent, thanks to `Prometheus-kafka-adapter team`  

## output

It is able to write JSON or Avro-JSON messages in a kafka topic, depending on the `SERIALIZATION_FORMAT` configuration variable.

### JSON

```json
{
  "timestamp": "1970-01-01T00:00:00Z",
  "value": "9876543210",
  "name": "up",

  "labels": {
    "__name__": "up",
    "label1": "value1",
    "label2": "value2"
  }
}
```

`timestamp` and `value` are reserved values, and can't be used as label names. `__name__` is a special label that defines the name of the metric and is copied as `name` to the top level for convenience.  
  

## development

```
go test
go build
```

## configuration

### prometheus

Prometheus needs to have a `remote_write` url configured, pointing to the '/receive' endpoint of the host and port where the prometheus-pulsar-adapter service is running. For example:

```yaml
remote_write:
  - url: "http://prometheus-pulsar-adapter:8080/receive"
```

## contributing

With issues:
  - Use the search tool before opening a new issue.
  - Please provide source code and commit sha if you found a bug.
  - Review existing issues and provide feedback or react to them.

With pull requests:
  - Open your pull request against `master`
  - You should add/modify tests to cover your proposed code changes.
  - If your pull request contains a new feature, please document it on the README.


## license

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
