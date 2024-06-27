import Big from 'big.js';

export const DENOM_TO_COIN_CONVERSION: number = 1e-18;
export const DENOM_COIN_SCALE: number = 18;

export function denomToHumanReadableConversion(denom: number): string {
  return Big(denom).times(DENOM_TO_COIN_CONVERSION).toFixed(DENOM_COIN_SCALE);
}

export function convertToDenomScale(num: string): string {
  return Big(num).toFixed(DENOM_COIN_SCALE);
}
