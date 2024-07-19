CREATE OR REPLACE FUNCTION dydx_uuid_from_transaction_parts(block_height int, transaction_index int) RETURNS uuid AS $$
/**
  Returns a UUID using the parts of a transaction.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    return dydx_uuid(concat(block_height, '-', transaction_index));
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
