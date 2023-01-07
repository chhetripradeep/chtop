# ClickHouse Top

Show live statistics for monitoring your ClickHouse node.

## Build

```
❯ make build
```

## Usage

It pulls metrics from the prometheus exporter endpoint of ClickHouse, so
make sure that you have enabled prometheus exporter endpoint for ClickHouse.

```
<clickhouse>
    <prometheus>
        <endpoint>/metrics</endpoint>
        <port>9363</port>
        <metrics>true</metrics>
        <events>true</events>
        <asynchronous_metrics>true</asynchronous_metrics>
        <status_info>true</status_info>
    </prometheus>
</clickhouse>
```

```
❯ chtop --help
Monitor your ClickHouse clusters without ever leaving your terminal

Usage:
  chtop [flags]

Flags:
      --config string        config file (default: $HOME/.chtop.yaml)
  -h, --help                 help for chtop
      --metrics-url string   clickhouse url for metrics in promql format
      --queries-url string   clickhouse url for running clickhouse queries
```

Run chtop pointing to prometheus stats endpoint & http endpoint of ClickHouse.
```
❯ chtop --metrics-url http://localhost:9363/metrics --queries-url http://localhost:8123
```

## Themes

You can configure the theme (default path: $HOME/.chtop.yaml) 

```
theme:
  graph:
    color: red
    height: 10
    precision: 1
```
## Metrics

You can configure the metrics to plot (default path: $HOME/.chtop.yaml)

```
clickhousemetrics:
  metrics:
    - alias: Total Queries
      name: ClickHouseProfileEvents_Query
    - alias: Total Select Queries
      name: ClickHouseProfileEvents_SelectQuery
    - alias: Total Insert Queries
      name: ClickHouseProfileEvents_InsertQuery
    - alias: Number of Active Parts
      name: ClickHouseMetrics_PartsActive
    - alias: Number of TCP Connections
      name: ClickHouseMetrics_TCPConnection
    - alias: Number of Open File Descriptors
      name: ClickHouseProfileEvents_FileOpen
```

You can configure to run sql queries to populate metrics to plot (default path: $HOME/.chtop.yaml)
```
clickhousequeries:
  queries:
    - name: Number of Running Queries
      sql: "select count(*) from system.processes"
```

## Sample Output

<img width="633" alt="Screenshot 2022-12-30 at 9 25 00 PM" src="https://user-images.githubusercontent.com/30620077/210074948-f453b33c-8158-47a3-8018-e6e59312f0a2.png">

## Thank you

This tool is built using [BubbleTea](https://github.com/charmbracelet/bubbletea), a very neat TUI Framework.

## Todos

- Allow to run CH queries to gather datapoints.
- Introduce more panel types.
- Allow to monitor CH clusters.
