# beacon

A Laravel-inspired service-provider framework for Go services, built on
[uber/dig](https://github.com/uber-go/dig). Shared across multiple repositories.

`module github.com/iVampireSP/beacon`

## Core

- **container** — `Application` (Register → Boot → Run), DI container, `ProviderFactory` + `Adapt[P]` for compile-time-safe registration.
- **support** — the `ServiceProvider` interface.
- **console** — `ConsoleCommand` and the `CommandProvider` capability.

## Capability interfaces

A `ServiceProvider` may also implement any of these; the framework collects them automatically:

| Interface | Declares | Collected by |
| `console.CommandProvider` | CLI commands | `Application.Boot` |
| `queue/job.HandlerProvider` | job handlers | `worker` command |
| `bus.ListenerProvider` | event listeners | `eventbus` command |
| `schedule.CronProvider` | cron jobs | `scheduler` command |

## Modules

`db` (PostgreSQL/pgx, ORM-agnostic `*sql.DB`), `cache`, `lock`, `jwt`, `keystore`,
`bus` (Kafka), `queue` (Asynq/Redis), `schedule` (cron), `cron`, `tracing` (OTel),
`email`, `config`, `i18n`, `tmpl`, `cipher`, `logger`, `httpserver`, plus utilities
(`cerr`, `json`, `must`, `paginator`, `filter`, `namegen`, `ratelimit`, `validator`, `version`).

The ORM client (ent) is intentionally **not** here: it is generated per-consuming-repo
and binds app-side over `db`'s `*sql.DB`.
