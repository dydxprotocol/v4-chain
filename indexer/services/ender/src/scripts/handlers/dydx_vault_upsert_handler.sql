CREATE OR REPLACE FUNCTION dydx_vault_upsert_handler(
  block_time timestamp, event_data jsonb) RETURNS jsonb AS $$
/**
  Parameters:
    - block_time: the time of the block being processed.
    - event_data: The 'data' field of the IndexerTendermintEvent (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/proto/dydxprotocol/indexer/indexer_manager/event.proto#L25)
        converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.
  Returns: JSON object containing fields:
    - vault: The upserted vault in vault-model format

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    vault_status text;
    vault_record vaults%ROWTYPE;
BEGIN
    vault_status = dydx_protocol_vault_status_to_vault_status(event_data->'status');

    vault_record."address" = jsonb_extract_path_text(event_data, 'address');
    vault_record."clobPairId" = (event_data->'clobPairId')::bigint;
    vault_record."status" = vault_status;
    vault_record."createdAt" = block_time;
    vault_record."updatedAt" = block_time;

    INSERT INTO vaults VALUES (vault_record.*)
    ON CONFLICT ("address") DO
      UPDATE
        SET 
          "status" = vault_status,
          "updatedAt" = block_time;

    RETURN jsonb_build_object(
      'vault',
      dydx_to_jsonb(vault_record)
    );
END;
$$ LANGUAGE plpgsql;
