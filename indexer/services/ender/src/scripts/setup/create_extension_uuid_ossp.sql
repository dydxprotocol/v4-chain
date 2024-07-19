/**
  Loads the 'uuid-ossp' extension which provides UUID generation functions, specifically
  `uuid_generate_v5` to generate version 5.0 UUIDs.

  See https://www.postgresql.org/docs/current/uuid-ossp.html for additional details.
*/
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
