-- Computes funding payments that occurred at current_height by aggregating fills that occurred between
-- last_height (exclusive) and current_height (inclusive) on top of the state of all open positions at
-- last_height.
INSERT INTO funding_payments (
    "subaccountId",
    "createdAt",
    "createdAtHeight",
    "perpetualId",
    "ticker",
    "oraclePrice",
    "size",
    "fundingIndex",
    "side",
    "rate",
    "payment"
)
WITH
    -- net computes the net size of each (subaccount, perpetual) pair.
    net AS (
        SELECT
            "subaccountId",
            "clobPairId",
            SUM(
                CASE
                    WHEN side = 'BUY' THEN size
                    WHEN side = 'SELL' THEN - size
                END
            ) AS net_size
        FROM
            fills
        -- we only want fills in range [last processed height + 1, current height]
        WHERE "createdAtHeight" > :last_height AND "createdAtHeight" <= :current_height
        GROUP BY
            "subaccountId",
            "clobPairId"
    ),
    -- position_snapshot computes the last snapshot size of each (subaccount, perpetual) pair.
    -- this is retrieved from the funding_payments table.
    position_snapshot AS (
        SELECT DISTINCT ON ("subaccountId", "perpetualId")
            "subaccountId",
            "perpetualId",
            ticker,
            size as last_snapshot_size,
            "createdAtHeight"
        FROM funding_payments
        -- we only want the last funding payment for each subaccount and perpetual_id
        -- this way, we can calculate current position by positions at last height
        -- plus fills from last height + 1 to current height.
        WHERE "createdAtHeight" = :last_height
        ORDER BY "subaccountId", "perpetualId", "createdAtHeight" DESC
    ),
    -- paired computes the net size of each (subaccount, perpetual) pair joined with the snapshot
    -- to figure out the current open positions of everyone.
    paired AS (
        SELECT
            COALESCE(n."subaccountId", ps."subaccountId") as "subaccountId",
            COALESCE(pm.id, ps."perpetualId") AS "perpetualId",
            COALESCE(pm.ticker, ps.ticker) AS ticker,
            COALESCE(n.net_size, 0) + COALESCE(ps.last_snapshot_size, 0) AS net_size
        FROM
            net n
            -- left join pm here to processed clobPairId into perpetual_id.
            LEFT JOIN perpetual_markets pm ON pm."clobPairId" = n."clobPairId"
            -- full join here because we want the entries when either net or position_snapshot is null.
            -- no match necessary. 
            FULL JOIN position_snapshot ps ON ps."subaccountId" = n."subaccountId" 
                AND ps."perpetualId" = pm.id
    ),
    -- funding computes the funding index update for each perpetual.
    new_funding AS (
        SELECT DISTINCT
            ON (f."perpetualId") f."perpetualId" AS "perpetualId",
            f.rate,
            f."oraclePrice" AS "oraclePrice",
            f."effectiveAt" AS "effectiveAt",
            f."fundingIndex" AS "fundingIndex"
        FROM
            funding_index_updates f
        WHERE f."effectiveAtHeight" = :current_height
        ORDER BY
            f."perpetualId",
            f."effectiveAtHeight" DESC
    ),
    last_funding AS (
        SELECT DISTINCT
            ON (f."perpetualId") f."perpetualId" AS "perpetualId",
            f."fundingIndex" AS "fundingIndex"
        FROM
            funding_index_updates f
        WHERE f."effectiveAtHeight" = :last_height
        ORDER BY
            f."perpetualId",
            f."effectiveAtHeight" DESC
    ),
    overall_funding AS (
        SELECT
            nf."perpetualId" AS "perpetualId",
            nf.rate AS rate,
            nf."oraclePrice" AS "oraclePrice",
            nf."effectiveAt" AS "effectiveAt",
            nf."fundingIndex" - COALESCE(lf."fundingIndex", 0) AS "fundingIndexDelta",
            nf."fundingIndex" AS "fundingIndex"
        FROM
            new_funding nf
            LEFT JOIN last_funding lf ON nf."perpetualId" = lf."perpetualId"
    )
SELECT
    p."subaccountId",
    f."effectiveAt" as "createdAt",
    :current_height as "createdAtHeight",
    p."perpetualId",
    p."ticker",
    f."oraclePrice",
    p.net_size AS "size",
    f."fundingIndex" AS "fundingIndex",
    CASE
        WHEN p.net_size > 0 THEN 'LONG'
        ELSE 'SHORT'
    END AS side,
    f."rate",
    - p.net_size * f."fundingIndexDelta" AS "payment"
FROM
    paired p
    -- inner join here because we absolutely need a funding index to calculate funding payments. 
    -- if no funding index, no entry will be created.
    JOIN overall_funding f ON f."perpetualId" = p."perpetualId"
WHERE
    p.net_size != 0
ORDER BY
    p."subaccountId",
    p."perpetualId";
