-- Calculates funding payments for all subaccounts between the most recent height up to which funding payments have been computed (exclusive) and the current height (inclusive).
INSERT INTO funding_payments (
    "subaccountId",
    "createdAt",
    "createdAtHeight",
    "perpetualId",
    ticker,
    "oraclePrice",
    size,
    side,
    rate,
    payment
)
WITH
    net AS (
        SELECT
            "subaccountId",
            "clobPairId", -- align the names
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
    -- figure out what the last funding payment was.
    last_funding_payment AS (
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
    paired AS (
        SELECT
            COALESCE(n."subaccountId", lfp."subaccountId") as "subaccountId",
            COALESCE(pm.id, lfp."perpetualId") AS "perpetualId",
            COALESCE(pm.ticker, lfp.ticker) AS ticker,
            COALESCE(n.net_size, 0) + COALESCE(lfp.last_snapshot_size, 0) AS net_size
        FROM
            net n
            -- left join pm here to processed clobPairId into perpetual_id.
            LEFT JOIN perpetual_markets pm ON pm."clobPairId" = n."clobPairId"
            -- full join here because we want the entries when either net or last_funding_payment is null.
            -- no match necessary. 
            FULL JOIN last_funding_payment lfp ON lfp."subaccountId" = n."subaccountId" 
                AND lfp."perpetualId" = pm.id
    ),
    funding AS (
        -- Grab the latest funding index update for each perpetual.
        SELECT DISTINCT
            ON (f."perpetualId") f."perpetualId" AS "perpetualId",
            f.rate,
            f."oraclePrice" AS "oraclePrice",
            f."effectiveAt" AS "effectiveAt"
        FROM
            funding_index_updates f
        WHERE f."effectiveAtHeight" = :current_height
        ORDER BY
            f."perpetualId",
            f."effectiveAtHeight" DESC
    )
SELECT
    p."subaccountId",
    CURRENT_TIMESTAMP as "createdAt",
    :current_height as "createdAtHeight",
    p."perpetualId",
    p.ticker,
    f."oraclePrice",
    p.net_size AS size,
    CASE
        WHEN p.net_size > 0 THEN 'LONG'
        ELSE 'SHORT'
    END AS side,
    f.rate,
    - p.net_size * f."oraclePrice" * f.rate AS payment
FROM
    paired p
    -- inner join here because we absolutely need a funding index to calculate funding payments. 
    -- if no funding index, no entry will be created.
    JOIN funding f ON f."perpetualId" = p."perpetualId"
WHERE
    p.net_size != 0
ORDER BY
    p."subaccountId",
    p."perpetualId";
