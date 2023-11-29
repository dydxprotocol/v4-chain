CREATE OR REPLACE FUNCTION dydx_from_serializable_int(serializable_int jsonb) RETURNS numeric AS $$
/**
  Converts a JSON.stringify byte array representing a SerializableInt
  (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/protocol/dtypes/serializable_int.go#L84) to a numeric.
   Note that the underlying SerializableInt encoding format uses the big.Int GobEncoding
  (https://github.com/golang/go/blob/886fba5/src/math/big/intmarsh.go#L18)
  which is represented as [versionAndSignByte bigEndianByte0 bigEndianByte1 ... bigEndianByte2]
  byte array.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    rval numeric = 0;
    version_and_sign int;
    num_value_keys int;
    key jsonb;
    value jsonb;
BEGIN
    IF serializable_int->'0' IS NULL THEN
        RETURN 0;
    END IF;

    version_and_sign = (serializable_int->'0')::int;
    if (version_and_sign >> 1 != 1) THEN
        RAISE EXCEPTION 'Unsupported BigInt encoding format %', version_and_sign >> 1;
    END IF;

    SELECT COUNT(*) - 1 INTO num_value_keys FROM jsonb_object_keys(serializable_int);

    FOR key, value IN SELECT * FROM jsonb_each(serializable_int) LOOP
        IF key != '0'::jsonb THEN
            rval = rval + power(256::numeric, num_value_keys - key::numeric) * value::numeric;
        END IF;
    END LOOP;

    IF (version_and_sign & 1) != 0 THEN
        rval = -rval;
    END IF;
    RETURN rval;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
