/**
 Converts the 'ConditionType' enum from the IndexerOrder protobuf (https://github.com/dydxprotocol/v4-proto/blob/4b721881fdfe99485336e221def03dc5b86eb0a1/dydxprotocol/indexer/protocol/v1/clob.proto#L131)
 to the 'OrderType' enum in postgres.
 */
CREATE OR REPLACE FUNCTION dydx_protocol_condition_type_to_order_type(condition_type jsonb) RETURNS text AS $$
DECLARE
    UNRECOGNIZED constant jsonb = '-1'::jsonb;
    CONDITION_TYPE_UNSPECIFIED constant jsonb = '0'::jsonb;
    CONDITION_TYPE_STOP_LOSS constant jsonb = '1'::jsonb;
    CONDITION_TYPE_TAKE_PROFIT constant jsonb = '2'::jsonb;
BEGIN
    CASE condition_type
    WHEN UNRECOGNIZED THEN
            RETURN 'LIMIT';
    WHEN CONDITION_TYPE_UNSPECIFIED THEN
        RETURN 'LIMIT';
    WHEN CONDITION_TYPE_STOP_LOSS THEN
        RETURN 'STOP_LIMIT';
    WHEN CONDITION_TYPE_TAKE_PROFIT THEN
        RETURN 'TAKE_PROFIT';
    ELSE
        RAISE EXCEPTION 'Unexpected ConditionType: %', condition_type;
    END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
