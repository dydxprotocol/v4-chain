DO $$
DECLARE
    bh_limit bigint := 200;  -- REPLACE THIS
BEGIN
    UPDATE transactions
    SET
        "transactionIndex" = "transactionIndex" - 1,
        id = dydx_uuid_from_transaction_parts("blockHeight"::text, ("transactionIndex" - 1)::text)
    WHERE
        "blockHeight" < bh_limit AND
        "transactionIndex" >= 1;
END $$;
