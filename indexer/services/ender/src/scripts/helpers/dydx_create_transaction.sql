CREATE OR REPLACE FUNCTION dydx_create_transaction(
    transaction_hash text, block_height int, transaction_index int
) RETURNS jsonb AS $$
/**
  Parameters:
    - transaction_hash: the hash of the transaction being processed.
    - block_height: the height of the block being processed.
    - transaction_index: the index of the transaction in the block.
  Returns: The inserted transaction.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    inserted_transaction jsonb;
BEGIN
    INSERT INTO transactions ("blockHeight", "transactionIndex", "transactionHash", "id")
    VALUES (block_height, transaction_index, transaction_hash,
        dydx_uuid_from_transaction_parts(block_height, transaction_index))
    RETURNING to_jsonb(transactions.*) INTO inserted_transaction;

    RETURN inserted_transaction;
END;
$$ LANGUAGE plpgsql;
