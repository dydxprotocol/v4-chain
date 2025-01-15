# Postgres

Postgres package holds all postgres knex migrations, postgres models and helper functions.

## Knex migration
Add a knex migration by running `pnpm run migrate:make <create_fake_table>`

Run the migration with `pnpm run migrate`

In `__tests__/db/migrations.test.ts`, test cases may be expected to fail (and hence should be commented out) if a model is modified due to the latest migration.

In order to migrate in dev and staging, you must redeploy and run bazooka.

TODO(CORE-512): Add info/resources around bazooka. [Doc](https://www.notion.so/dydx/Engineering-Runbook-15064661da9643188ce33e341b68e7bb#cb2283d80ef14a51924f3bd1a538fd82).
