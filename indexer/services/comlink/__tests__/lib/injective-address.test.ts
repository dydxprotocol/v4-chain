import {
  getInjectiveAddress,
  getInjectiveEvmAddress,
  getInjectiveMintRecipient,
} from '../../src/lib/injective-address';
import {
  DYDX_ADDRESS,
  INJECTIVE_ADDRESS,
  INJECTIVE_EVM_ADDRESS,
  INJECTIVE_MINT_RECIPIENT,
} from '../helpers/deposit-policy-helpers';

describe('injective-address', () => {
  it('gives the Injective address that shares its bytes with the dYdX address', () => {
    expect(getInjectiveAddress(DYDX_ADDRESS)).toEqual(INJECTIVE_ADDRESS);
  });

  it('gives the same account as a checksummed EVM address', () => {
    expect(getInjectiveEvmAddress(DYDX_ADDRESS)).toEqual(INJECTIVE_EVM_ADDRESS);
  });

  it('gives the same account as a 32-byte CCTP mint recipient', () => {
    expect(getInjectiveMintRecipient(DYDX_ADDRESS)).toEqual(INJECTIVE_MINT_RECIPIENT);
  });

  it.each([
    ['an Injective address', INJECTIVE_ADDRESS],
    ['a Noble address', 'noble1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5hpek7j'],
    ['an EVM address', INJECTIVE_EVM_ADDRESS],
    ['a dYdX address with a bad checksum', `${DYDX_ADDRESS.slice(0, -1)}q`],
    ['a dYdX address that is not 20 bytes', 'dydx1qqqsyqcyq5rqwzqfpg9scrgwpugpzysnzs23v9ccwdz36y'],
    ['an empty string', ''],
  ])('refuses %s', (_name: string, address: string) => {
    expect(() => getInjectiveAddress(address)).toThrow();
    expect(() => getInjectiveEvmAddress(address)).toThrow();
    expect(() => getInjectiveMintRecipient(address)).toThrow();
  });
});
