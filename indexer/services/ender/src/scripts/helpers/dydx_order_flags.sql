-- Create a separate schema for constants
CREATE SCHEMA IF NOT EXISTS constants;

-- Create functions in the constants schema
CREATE OR REPLACE FUNCTION constants.order_flag_long_term() RETURNS bigint AS $$
BEGIN
    RETURN 64;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;

CREATE OR REPLACE FUNCTION constants.order_flag_twap() RETURNS bigint AS $$
BEGIN
    RETURN 128;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;

CREATE OR REPLACE FUNCTION constants.order_flag_twap_suborder() RETURNS bigint AS $$
BEGIN
    RETURN 256;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;