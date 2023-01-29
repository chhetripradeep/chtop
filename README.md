# ClickHouse Top

Show live statistics for monitoring your ClickHouse node.

## Build

```
❯ make build
```

## Usage

It can populate graphs either by pulling metrics from the prometheus exporter 
endpoint of ClickHouse or by running sql queries on ClickHouse native protocol endpoint.

To enable prometheus exporter endpoint for ClickHouse, you will need to add following
ClickHouse server configuration:

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
  -c, --config string             path of the config file (default: $HOME/.chtop.yaml)
  -h, --help                      help for chtop
  -m, --metrics-url string        clickhouse url for pulling metrics in prometheus exposition format
  -d, --queries-database string   clickhouse database for connecting from clickhouse client (default "system")
  -p, --queries-password string   clickhouse password of the provided clickhouse user for running clickhouse queries
  -q, --queries-url string        clickhouse endpoint for running clickhouse queries via native protocol
  -u, --queries-username string   clickhouse username for running clickhouse queries (default "default")
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

<img width="1529" alt="demo" src="https://user-images.githubusercontent.com/30620077/214207781-577d75a0-e593-4b01-80cb-8228c2ee4c40.png">

## Thank you

This tool is built using [BubbleTea](https://github.com/charmbracelet/bubbletea), a very neat TUI Framework.

## Todos

- Introduce more panel types.
- Allow to monitor CH clusters.
