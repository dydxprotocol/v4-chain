CREATE OR REPLACE FUNCTION dydx_protocol_convert_to_order_type(order_flags bigint, condition_type jsonb) RETURNS text AS $$
/**
  Converts the 'ConditionType' enum from the IndexerOrder protobuf (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/proto/dydxprotocol/indexer/protocol/v1/clob.proto#L130)
  to the 'OrderType' enum in postgres.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    CASE order_flags
        WHEN 0 THEN RETURN 'LIMIT';
        WHEN 32 THEN
            CASE condition_type
                WHEN '-1'::jsonb THEN RETURN 'LIMIT'; /** UNRECOGNIZED */
                WHEN '0'::jsonb THEN RETURN 'LIMIT'; /** CONDITION_TYPE_UNSPECIFIED */
                WHEN '1'::jsonb THEN RETURN 'STOP_LIMIT'; /** CONDITION_TYPE_STOP_LOSS */
                WHEN '2'::jsonb THEN RETURN 'TAKE_PROFIT'; /** CONDITION_TYPE_TAKE_PROFIT */
                ELSE RAISE EXCEPTION 'Unexpected ConditionType: %', condition_type;
            END CASE;
        WHEN 64 THEN RETURN 'LIMIT';
        WHEN 128 THEN RETURN 'TWAP';
        WHEN 256 THEN RETURN 'TWAP_SUBORDER';
        ELSE RAISE EXCEPTION 'Unexpected OrderFlags: %', order_flags;
    END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
