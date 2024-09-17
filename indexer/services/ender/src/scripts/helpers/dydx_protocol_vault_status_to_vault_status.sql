CREATE OR REPLACE FUNCTION dydx_protocol_vault_status_to_vault_status(vaultStatus jsonb)
    RETURNS text AS $$

BEGIN
    CASE vaultStatus
        WHEN '1'::jsonb THEN RETURN 'DEACTIVATED'; /** VAULT_STATUS_DEACTIVATED **/
        WHEN '2'::jsonb THEN RETURN 'STAND_BY'; /** VAULT_STATUS_STAND_BY **/
        WHEN '3'::jsonb THEN RETURN 'QUOTING'; /** VAULT_STATUS_QUOTING **/
        WHEN '4'::jsonb THEN RETURN 'CLOSE_ONLY'; /** VAULT_STATUS_CLOSE_ONLY **/
        ELSE RAISE EXCEPTION 'Invalid vault status: %', vaultStatus;
    END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
