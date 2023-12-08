CREATE OR REPLACE FUNCTION dydx_get_weighted_average(first_price numeric, first_weight numeric, second_price numeric, second_weight numeric) RETURNS numeric AS $$
/**
  Returns the weighted average between two prices.

  Note that since division is used the scale of the resulting number is limited to 20 which matches the division
  precision (DP) of the https://mikemcl.github.io/big.js/ library.

  Parameters:
    - first_price: The first price. Defaults to 0 if null.
    - first_weight: The weight of the first price.
    - second_price: The second price. Defaults to 0 if null.
    - second_weight: The weight of the second price.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    RETURN dydx_trim_scale((coalesce(first_price, 0::numeric) * first_weight +
                            coalesce(second_price, 0::numeric) * second_weight)::numeric(256, 20)
                               / (first_weight + second_weight)::numeric(256, 20));
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
