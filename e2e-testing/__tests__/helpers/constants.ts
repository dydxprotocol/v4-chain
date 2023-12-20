export const DYDX_LOCAL_ADDRESS = 'dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4';
export const DYDX_LOCAL_MNEMONIC = 'merge panther lobster crazy road hollow amused security before critic about cliff exhibit cause coyote talent happy where lion river tobacco option coconut small';
export const DYDX_LOCAL_ADDRESS_2 = 'dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs';
export const DYDX_LOCAL_MNEMONIC_2 = 'color habit donor nurse dinosaur stable wonder process post perfect raven gold census inside worth inquiry mammal panic olive toss shadow strong name drum';

export const MNEMONIC_TO_ADDRESS: Record<string, string> = {
  [DYDX_LOCAL_MNEMONIC]: DYDX_LOCAL_ADDRESS,
  [DYDX_LOCAL_MNEMONIC_2]: DYDX_LOCAL_ADDRESS_2,
};

export const ADDRESS_TO_MNEMONIC: Record<string, string> = {
  [DYDX_LOCAL_ADDRESS]: DYDX_LOCAL_MNEMONIC,
  [DYDX_LOCAL_ADDRESS_2]: DYDX_LOCAL_MNEMONIC_2,
};
