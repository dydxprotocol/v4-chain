CREATE OR REPLACE FUNCTION dydx_from_jsonlib_long(long_value jsonb) RETURNS numeric AS $$
/**
  Converts JSON objects of the form (https://www.npmjs.com/package/long):
    {
      "low": 10000000,
      "high": 0,
      "unsigned": false
    }
  and converts it to a numeric. Note that this is the format used to convert Long types when converted using
  JSON.stringify.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    POWER_2_32 constant numeric = power(2::numeric, 32::numeric);
BEGIN
    /*
    We use ::int4::bit(32)::int8::numeric to convert a signed 4-byte integer to an unsigned integer.
    This is equivalent to `number >>> 0` being used in the json long package.
     */
    IF (long_value->'unsigned')::bool THEN
        RETURN (long_value->'high')::int4::bit(32)::int8::numeric * POWER_2_32
            + (long_value->'low')::numeric::int4::bit(32)::int8::numeric;
    END IF;
    RETURN (long_value->'high')::numeric * POWER_2_32 + (long_value->'low')::int4::bit(32)::int8::numeric;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
