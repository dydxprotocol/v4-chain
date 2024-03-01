# Scripts

Scripts to be run from devbox.

### Available Commands
Check `package.json` for all available commands

### Args

In order to submit args with these scripts you may need to separate the script invocation and the args with `--`.

Example:

```sh
pnpm run validate-pnl -- \
--s 964e4f37-1314-5181-bcf9-d3b0a30d86ed \
--p d6f67aac-749d-5f62-8a55-ef825915b575 d6f67aac-749d-5f62-8a55-ef825915b575


KAFKA_BROKER_URLS=...kafka.ap-northeast-1.amazonaws.com:9092 \
SERVICE_NAME=scripts pnpm run print-block -- --h 9265388

SERVICE_NAME=script pnpm run print-candle-time-boundaries -- \
--t 2024-02-28T10:01:36.17+00:00
```

### EnvVars

In order to properly run all scripts, add these lines to your `~/.bashrc`:

```
export DB_PORT=5432
export DB_NAME=dydx
export DB_USERNAME=dydx
export DB_HOSTNAME=staging-indexer-apne1-db.cv0zh5lkgpcw.ap-northeast-1.rds.amazonaws.com
export DB_READONLY_HOSTNAME=staging-indexer-apne1-db-read-replica.cv0zh5lkgpcw.ap-northeast-1.rds.amazonaws.com
export DB_PASSWORD=<insert DB password from ~/.pgpass here>
```
