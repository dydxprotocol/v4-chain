/**
  Parameters:
    - transaction_hash: the hash of the transaction being processed.
    - block_height: the height of the block being processed.
    - transaction_index: the index of the transaction in the block.
  Returns: The inserted transaction.
*/
CREATE OR REPLACE FUNCTION dydx_create_transaction(
    transaction_hash text, block_height text, transaction_index int
) RETURNS jsonb AS $$
DECLARE
    inserted_transaction jsonb;
BEGIN
    INSERT INTO transactions ("blockHeight", "transactionIndex", "transactionHash", "id")
    VALUES (block_height::bigint, transaction_index, transaction_hash,
        dydx_uuid_from_transaction_parts(block_height, transaction_index::text))
    RETURNING to_jsonb(transactions.*) INTO inserted_transaction;

    RETURN inserted_transaction;
END;
$$ LANGUAGE plpgsql;
