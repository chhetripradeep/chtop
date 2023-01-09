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
      --config string             config file (default: $HOME/.chtop.yaml)
  -h, --help                      help for chtop
      --metrics-url string        clickhouse url for pulling metrics in prometheus exposition format
      --queries-database string   clickhouse database for connecting clickhouse client (default "system")
      --queries-password string   clickhouse password for running clickhouse queries
      --queries-url string        clickhouse url for running clickhouse queries (native protocol port)
      --queries-username string   clickhouse username for running clickhouse queries (default "default")
```

Run chtop pointing to prometheus stats endpoint & http endpoint of ClickHouse.

Sample Run:
```
❯ chtop --metrics-url http://localhost:9363/metrics --queries-url localhost:9000 --config chtop.yaml
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
    - name: Number of Databases
      sql: "select count(*) from system.databases"
    - name: Number of Tables
      sql: "select count(*) from system.tables"
    - name: Number of Parts
      sql: "select count(*) from system.parts"
```

## Sample Output


<img width="1515" alt="Screenshot 2023-01-10 at 3 32 20 AM" src="https://user-images.githubusercontent.com/30620077/211392656-7a8a261d-5e4f-4ed4-9f3d-107e376eba34.png">

## Thank you

This tool is built using [BubbleTea](https://github.com/charmbracelet/bubbletea), a very neat TUI Framework.

## Todos

- Introduce more panel types.
- Allow to monitor CH clusters.
