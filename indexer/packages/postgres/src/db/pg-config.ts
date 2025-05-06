import { DateTime } from 'luxon';
import pg from 'pg';

/**
 * we need to add this line, because the default type parser
 * changes all datetime objects to javascript dates, when we
 * need all dates returned with iso strings
 */

const utcZone = {
  zone: 'utc',
};

pg.types.setTypeParser(
  pg.types.builtins.TIMESTAMPTZ,
  (val) => (val === null ? null : DateTime.fromISO(val.replace(' ', 'T'), utcZone).toISO()),
);
