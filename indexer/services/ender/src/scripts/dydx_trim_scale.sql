CREATE OR REPLACE FUNCTION dydx_trim_scale(value numeric) RETURNS numeric AS $$
/**
  Returns a numeric with the zeros after the decimal point removed. Note that this function should be replaced by
  trim_scale which has become available with Postgres 13 (https://www.postgresql.org/docs/current/functions-math.html).

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    trimmed_text text;
    trimmed_num numeric;
BEGIN
    trimmed_text = rtrim(value::text, '0');
    /** Check that we didn't trim all the digits in the case of the value being '0'. */
    IF length(trimmed_text) = 0 THEN
        RETURN 0;
    END IF;
    trimmed_num = trimmed_text::numeric;
    /** Check that the trimmed values are equivalent in case we trim values that are important, for example '10'. */
    IF trimmed_num = value THEN
        RETURN trimmed_num;
    END IF;
    RETURN value;
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
