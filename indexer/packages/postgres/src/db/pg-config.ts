import pg from 'pg';
import { DateTime } from 'luxon';

/**
 * we need to add this line, because the default type parser
 * changes all datetime objects to javascript dates, when we
 * need all dates returned with iso strings
 */

const utcZone = {
  zone: "utc",
}

pg.types.setTypeParser(
  pg.types.builtins.TIMESTAMPTZ,
  (val) => (val === null ? null : DateTime.fromSQL(val, utcZone).toISO()),
);
