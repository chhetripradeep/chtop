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
Run chtop pointing to prometheus stats endpoint of ClickHouse.

```
❯ chtop --url http://localhost:9363/metrics
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
    - name: ClickHouseProfileEvents_Query
    - name: ClickHouseProfileEvents_SelectQuery
    - name: ClickHouseProfileEvents_InsertQuery
    - name: ClickHouseMetrics_PartsActive
    - name: ClickHouseMetrics_TCPConnection
    - name: ClickHouseProfileEvents_FileOpen
```

## Sample Output

<img width="633" alt="Screenshot 2022-12-30 at 9 25 00 PM" src="https://user-images.githubusercontent.com/30620077/210074948-f453b33c-8158-47a3-8018-e6e59312f0a2.png">

## Thank you

This tool is built using [BubbleTea](https://github.com/charmbracelet/bubbletea), a very neat TUI Framework.

## Todos

- Allow to run CH queries to gather datapoints.
- Introduce more panel types.
- Allow to monitor CH clusters.
