/**
 Converts the TimeInForce field from an IndexerOrder proto (https://github.com/dydxprotocol/v4-proto/blob/437f6d8/dydxprotocol/indexer/protocol/v1/clob.proto#L95)
 to a TimeInForce enum in postgres.

  Raise an exception if the input TimeInForce enum is not in the known enum values for TimeInForce.
 */
CREATE OR REPLACE FUNCTION dydx_from_protocol_time_in_force(tif jsonb) RETURNS text AS $$
BEGIN
    CASE tif
        -- Default behavior with UNRECOGNIZED = GTT (Good-Til-Time)
        WHEN '-1'::jsonb THEN RETURN 'GTT';
        -- Default behavior with TIME_IN_FORCE_UNSPECIFIED = GTT (Good-Til-Time)
        WHEN '0'::jsonb THEN RETURN 'GTT';
        WHEN '1'::jsonb THEN RETURN 'IOC';
        WHEN '2'::jsonb THEN RETURN 'POST_ONLY';
        WHEN '3'::jsonb THEN RETURN 'FOK';
        ELSE RAISE EXCEPTION 'Unexpected TimeInForce from protocol %', tif;
        END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
