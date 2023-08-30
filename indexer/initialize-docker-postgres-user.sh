#!/bin/bash

# Initialize the postgres docker image.

set -e
set -u

# Grant access to datadog to monitor the DB instance.
# Define the function
grant_datadog_access_with_pword() {
     psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
create user datadog with password '$DATADOG_POSTGRES_PASSWORD';
grant pg_monitor to datadog;
grant SELECT ON pg_stat_database to datadog;
EOSQL
}

# Initialize extensions locally that are useful for development.
initialize_extensions() {
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
create extension if not exists "plpgsql";
create extension if not exists "plpgsql_check";
create extension if not exists "plprofiler";
create extension if not exists "pldbgapi";
EOSQL
}

initialize_extensions
grant_datadog_access_with_pword
