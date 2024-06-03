# Environment Log Receiver

| Status        |           |
| ------------- |-----------|
| Stability     | [beta]   |


An extension that produces logs containing stats for the environment. A file can be used as an additional buffer.

## Configuration

| Field                              | Default | Description                                                                                               |
|------------------------------------|---------|-----------------------------------------------------------------------------------------------------------|
| `operators`                         | []                                   | An array of [operators](../../pkg/stanza/docs/operators/README.md#what-operators-are-available). See below for more details.                                                                                                                                    |
| `resource`                         | {}      | A map of `key: value` pairs to add to the entry's resource.                                               |
| `storage`                           | none                                 | The ID of a storage extension to be used to store file offsets. File offsets allow the receiver to pick up where it left off in the case of a collector restart. If no storage extension is used, the receiver will manage offsets in memory only.              |
| `retry_on_failure.enabled`         | `false`                              | If `true`, the receiver will pause reading a file and attempt to resend the current batch of logs if it encounters an error from downstream components.                                                                                                         |
| `retry_on_failure.initial_interval` | `1s`                                 | [Time](#time-parameters) to wait after the first failure before retrying.                                                                                                                                                                                       |
| `retry_on_failure.max_interval`    | `30s`                                | Upper bound on retry backoff [interval](#time-parameters). Once this value is reached the delay between consecutive retries will remain constant at the specified value.                                                                                        |
| `retry_on_failure.max_elapsed_time` | `5m`                                 | Maximum amount of [time](#time-parameters) (including retries) spent trying to send a logs batch to a downstream consumer. Once this value is reached, the data is discarded. Retrying never stops if set to `0`.                                               |
| `buffer`                           | None | A file buffer configuration |
| `log_samplers`                     | []      | A list of log samplers to be added to the file log receiver. For the moment only one sampler is supported |

# File buffer configuration

| Field           | Default  | Description                                                                                                                                           |
|-----------------|----------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| `attributes`                        | {}                                   | A map of `key: value` pairs to add to the entry's attributes.                                                                                                                                                                                                   |
| `delete_after_read`                 | `false`                              | If `true`, each log file will be read and then immediately deleted. Requires that the `filelog.allowFileDeletion` feature gate is enabled. Must be `false` when `start_at` is set to `end`.                                                                     |
| `encoding`                          | `utf-8`                              | The encoding of the file being read. See the list of [supported encodings below](#supported-encodings) for available options.                                                                                                                                   |
| `fingerprint_size`                  | `1kb`                                | The number of bytes with which to identify a file. The first bytes in the file are used as the fingerprint. Decreasing this value at any point will cause existing fingerprints to forgotten, meaning that all files will be read from the beginning (one time) |
| `force_flush_period`                | `500ms`                              | [Time](#time-parameters) since last time new data was found in the file, after which a partial log at the end of the file may be emitted.|
| `include_file_name`                 | `true`                               | Whether to add the file name as the attribute `log.file.name`.                                                                                                                                                                                                  |
| `include_file_path`                 | `false`                              | Whether to add the file path as the attribute `log.file.path`.                                                                                                                                                                                                  |
| `include_file_name_resolved`        | `false`                              | Whether to add the file name after symlinks resolution as the attribute `log.file.name_resolved`.                                                                                                                                                               |
| `include_file_path_resolved`        | `false`                              | Whether to add the file path after symlinks resolution as the attribute `log.file.path_resolved`.                                                                                                                                                               |
| `include_file_owner_name`           | `false`                              | Whether to add the file owner name as the attribute `log.file.owner.name`. Not supported for windows.                                                                                                                                                           |
| `include_file_owner_group_name`     | `false`                              | Whether to add the file group name as the attribute `log.file.owner.group.name`. Not supported for windows.                                                                                                                                                     |
| `max_log_size`                      | `1MiB`                               | The maximum size of a log entry to read. A log entry will be truncated if it is larger than `max_log_size`. Protects against reading large amounts of data into memory.                                                                                         |
| `max_concurrent_files`              | 1024                                 | The maximum number of log files from which logs will be read concurrently. If the number of files matched in the `include` pattern exceeds this number, then files will be processed in batches.                                                                |
| `max_batches`                       | 0                                    | Only applicable when files must be batched in order to respect `max_concurrent_files`. This value limits the number of batches that will be processed during a single poll interval. A value of 0 indicates no limit.                                           |
| `max_log_size`                      | `1MiB`                               | The maximum size of a log entry to read. A log entry will be truncated if it is larger than `max_log_size`. Protects against reading large amounts of data into memory.                                                                                         |

| `multiline`                         |                                      | A `multiline` configuration block. See [below](#multiline-configuration) for more details.                                                                                                                                                                      |
| `poll_interval`                     | 200ms                                | The [duration](#time-parameters) between filesystem polls.                                                                                                                                                                                                      |
| `preserve_leading_whitespaces`      | `false`                              | Whether to preserve leading whitespaces.                                                                                                                                                                                                                        |
| `preserve_trailing_whitespaces`     | `false`                              | Whether to preserve trailing whitespaces.                                                                                                                                                                                                                       |
| `start_at`                          | `end`                                | At startup, where to start reading logs from the file. Options are `beginning` or `end`.                                                                                                                                                                        |


## Log Sampler

| Field           | Default  | Description                                                                                                                                           |
|-----------------|----------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| `metric`        | Required | The metric to sample. Possible values [netstats]                                                                                                      |
| `poll_interval` | Optional | The interval for generating the metrics                                                                                                               |


## Examples

This will output netstats delta metrics using a file as a buffer.
The receiver will add a record to the filebuffer.log each 60s the receiver will read that file producing the logs to the pipeline each 10s.

```yaml
  envlogreceiver/metering:
    log_samplers:
      - metric: netstats
        poll_interval: 60s
    buffer:
      path: buffer.log
      include_file_name: false
      poll_interval: 10s
      fingerprint_size: 1kb
      start_at: beginning
    storage: file_storage/deltas
    retry_on_failure:
      enabled: true
      initial_interval: 1s
      max_interval: 10m
      max_elapsed_time: 1h
```

This will output netstats directly to the pipeline each 60s. The buffer will be done by the pipeline.

```yaml
  envlogreceiver/metering:
    log_samplers:
      - metric: netstats
        poll_interval: 60s
    storage: file_storage/deltas
    retry_on_failure:
      enabled: true
      initial_interval: 1s
      max_interval: 10m
      max_elapsed_time: 1h
```
