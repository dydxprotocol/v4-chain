# Postgres

Postgres package holds all postgres knex migrations, seeders, postgres models and helper functions.

### Running the seeder

On a machine with access to the database instance for the Indexer, update `.env` with the required
environment variables for accessing the database instance. (See `.env.test` for an example of the
required values).

Then run:

```
pnpm run build && pnpm run seed
```

## Knex migration
Add a knex migration by running `pnpm run migrate:make <create_fake_table>`

Run the migration with `pnpm run migrate`

In order to migrate in v4 dev and staging, you must redeploy and run bazooka following the instructions [here](https://www.notion.so/dydx/Engineering-Runbook-15064661da9643188ce33e341b68e7bb#cb2283d80ef14a51924f3bd1a538fd82).
