/**
  Loads the 'pg_stat_statements' extension which captures timing metrics for statements and
  functions being executed within the database.

  See https://www.postgresql.org/docs/current/pgstatstatements.html for additional details.
*/
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";
