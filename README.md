# Music Analyze

[![CI](https://github.com/bryanwahyu/analyze/actions/workflows/ci.yml/badge.svg)](https://github.com/bryanwahyu/analyze/actions/workflows/ci.yml)

A small Go service that answers sales questions about the [Chinook](https://github.com/lerocha/chinook-database) music store dataset: which genres sell best, which countries buy the most, who the top support agents are, and the median and mode of invoice totals. MySQL, chi, Swagger.

## Run it

```bash
docker compose up --build
```

This starts MySQL with the Chinook data loaded and the API on `http://localhost:5000`. On macOS, port 5000 is often taken by AirPlay; use `APP_PORT=5050 docker compose up --build` instead.

Swagger UI: `http://localhost:5000/swagger/index.html`

## Endpoints

| Endpoint | Answers |
|---|---|
| `GET /analytics/top-genres?limit=10` | Genres by revenue, with share of total |
| `GET /analytics/top-countries?limit=5` | Countries by customer count, and by average order value |
| `GET /analytics/media-types` | Revenue per media type, with share of total |
| `GET /analytics/highest-sale` | The largest invoice(s) and who placed them |
| `GET /analytics/invoice-stats` | Median and mode of invoice totals |
| `GET /analytics/agents` | Revenue and transaction count per support agent |

Example, against the bundled data:

```bash
$ curl -s 'localhost:5000/analytics/top-genres?limit=3'
[{"genre":"Rock","genre_sales":826.65,"pct_of_total":35.499871},
 {"genre":"Latin","genre_sales":382.14,"pct_of_total":16.410719},
 {"genre":"Metal","genre_sales":261.36,"pct_of_total":11.223911}]

$ curl -s localhost:5000/analytics/invoice-stats
{"median_total":3.96,"mode_total":1.98}
```

Errors return `500 internal error`; the underlying database error is logged, not sent to the client.

## How the numbers are computed

All aggregation happens in SQL (`internal/infra/db/mysql/analyze_repo.go`):

- Shares of total use a CTE per group plus `SUM(...) OVER ()`, so one query returns both the group total and the grand total.
- The median orders invoices with `ROW_NUMBER()` and averages the middle one or two rows, which works for both odd and even counts.
- The mode is the most frequent invoice total, ties broken by the larger value.
- Highest sale returns every invoice that matches the maximum total, not just the first.

## Data

`internal/infra/db/mysql/DB.sql` is the Chinook MySQL script (original CamelCase tables) followed by snake_case views (`invoice`, `invoice_line`, `customer`, ...). The queries read the views. Docker Compose loads the script on first start; to load it into your own MySQL:

```bash
mysql -uroot -p < internal/infra/db/mysql/DB.sql
```

## Configuration

The service reads a YAML file from `CONFIG_PATH` (default `config.yaml`). `config.docker.yaml` is the one used in Docker; copy it to `config.yaml` and change `database.host` to run locally:

```yaml
server:
  port: 5000
database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: example
  name: music_analyze
```

```bash
go run ./cmd
```

## Tests

```bash
go test ./...                                     # handler tests; repository tests skip
docker compose up -d db
MYSQL_DSN='root:example@tcp(127.0.0.1:3306)/music_analyze?parseTime=true' go test ./...
```

The repository tests run every query against the real Chinook data and check known answers (Rock is the top genre, the USA has 13 customers, there are 412 invoices across three agents, and so on). CI runs them against a MySQL service container.

## Layout

```
cmd/main.go                    entrypoint, router, Swagger wiring
internal/domain/analyst        result types and the repository port
internal/application/analyst   use cases
internal/http/analytics        HTTP handlers
internal/infra/db/mysql        MySQL repository, connection, Chinook script
docs/                          generated Swagger files (swag init -g cmd/main.go -o docs)
```

## License

The Chinook script is by Luis Rocha under its own license (see the header of `DB.sql`).
