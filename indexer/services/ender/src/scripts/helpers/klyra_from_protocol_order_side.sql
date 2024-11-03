CREATE OR REPLACE FUNCTION klyra_from_protocol_order_side(order_side jsonb) RETURNS text AS $$
/**
  Converts the 'Side' enum from the IndexerOrder protobuf
  to the 'OrderSide' enum in postgres.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    CASE order_side
        WHEN '1'::jsonb THEN RETURN 'BUY'; /** SIDE_BUY */
        ELSE RETURN 'SELL';
    END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
