CREATE OR REPLACE FUNCTION klyra_protocol_condition_type_to_order_type(condition_type jsonb) RETURNS text AS $$
/**
  Converts the 'ConditionType' enum from the IndexerOrder protobuf
  to the 'OrderType' enum in postgres.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    CASE condition_type
        WHEN '-1'::jsonb THEN RETURN 'LIMIT'; /** UNRECOGNIZED */
        WHEN '0'::jsonb THEN RETURN 'LIMIT'; /** CONDITION_TYPE_UNSPECIFIED */
        WHEN '1'::jsonb THEN RETURN 'STOP_LIMIT'; /** CONDITION_TYPE_STOP_LOSS */
        WHEN '2'::jsonb THEN RETURN 'TAKE_PROFIT'; /** CONDITION_TYPE_TAKE_PROFIT */
        ELSE RAISE EXCEPTION 'Unexpected ConditionType: %', condition_type;
    END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
