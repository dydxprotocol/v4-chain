CREATE OR REPLACE FUNCTION dydx_from_protocol_time_in_force(tif jsonb) RETURNS text AS $$
/**
  Converts the TimeInForce field from an IndexerOrder proto (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/proto/dydxprotocol/indexer/protocol/v1/clob.proto#L94)
  to a TimeInForce enum in postgres.

  Raise an exception if the input TimeInForce enum is not in the known enum values for TimeInForce.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    CASE tif
        WHEN '-1'::jsonb THEN RETURN 'GTT'; /** Default behavior with UNRECOGNIZED = GTT (Good-Til-Time) */
        WHEN '0'::jsonb THEN RETURN 'GTT'; /** Default behavior with TIME_IN_FORCE_UNSPECIFIED = GTT (Good-Til-Time) */
        WHEN '1'::jsonb THEN RETURN 'IOC'; /** TIME_IN_FORCE_IOC */
        WHEN '2'::jsonb THEN RETURN 'POST_ONLY'; /** TIME_IN_FORCE_POST_ONLY */
        WHEN '3'::jsonb THEN RETURN 'FOK'; /** TIME_IN_FORCE_FILL_OR_KILL */
        ELSE RAISE EXCEPTION 'Unexpected TimeInForce from protocol %', tif;
    END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
