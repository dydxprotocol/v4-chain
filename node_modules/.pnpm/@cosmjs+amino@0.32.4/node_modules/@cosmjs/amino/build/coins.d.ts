export interface Coin {
    readonly denom: string;
    readonly amount: string;
}
/**
 * Creates a coin.
 *
 * If your values do not exceed the safe integer range of JS numbers (53 bit),
 * you can use the number type here. This is the case for all typical Cosmos SDK
 * chains that use the default 6 decimals.
 *
 * In case you need to supportr larger values, use unsigned integer strings instead.
 */
export declare function coin(amount: number | string, denom: string): Coin;
/**
 * Creates a list of coins with one element.
 */
export declare function coins(amount: number | string, denom: string): Coin[];
/**
 * Takes a coins list like "819966000ucosm,700000000ustake" and parses it.
 *
 * Starting with CosmJS 0.32.3, the following imports are all synonym and support
 * a variety of denom types such as IBC denoms or tokenfactory. If you need to
 * restrict the denom to something very minimal, this needs to be implemented
 * separately in the caller.
 *
 * ```
 * import { parseCoins } from "@cosmjs/proto-signing";
 * // equals
 * import { parseCoins } from "@cosmjs/stargate";
 * // equals
 * import { parseCoins } from "@cosmjs/amino";
 * ```
 *
 * This function is not made for supporting decimal amounts and does not support
 * parsing gas prices.
 */
export declare function parseCoins(input: string): Coin[];
/**
 * Function to sum up coins with type Coin
 */
export declare function addCoins(lhs: Coin, rhs: Coin): Coin;
