import {
  getInjectiveEvmAddress,
  getInjectiveMintRecipient,
  normalizeInjectiveAddress,
} from '../../src/lib/injective-address';
import {
  DYDX_ADDRESS,
  INJECTIVE_ADDRESS,
  INJECTIVE_EVM_ADDRESS,
  INJECTIVE_MINT_RECIPIENT,
} from '../helpers/deposit-policy-helpers';

describe('injective-address', () => {
  it('gives the account as a checksummed EVM address', () => {
    expect(getInjectiveEvmAddress(INJECTIVE_ADDRESS)).toEqual(INJECTIVE_EVM_ADDRESS);
  });

  it('gives the account as a 32-byte CCTP mint recipient', () => {
    expect(getInjectiveMintRecipient(INJECTIVE_ADDRESS)).toEqual(INJECTIVE_MINT_RECIPIENT);
  });

  it('leaves an address written the usual way as it is', () => {
    expect(normalizeInjectiveAddress(INJECTIVE_ADDRESS)).toEqual(INJECTIVE_ADDRESS);
  });

  it('reads an address written in upper case as the same account', () => {
    const upperCase: string = INJECTIVE_ADDRESS.toUpperCase();

    expect(normalizeInjectiveAddress(upperCase)).toEqual(INJECTIVE_ADDRESS);
    expect(getInjectiveEvmAddress(upperCase)).toEqual(INJECTIVE_EVM_ADDRESS);
    expect(getInjectiveMintRecipient(upperCase)).toEqual(INJECTIVE_MINT_RECIPIENT);
  });

  it.each([
    ['a dYdX address', DYDX_ADDRESS],
    ['a Noble address', 'noble1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5hpek7j'],
    ['an EVM address', INJECTIVE_EVM_ADDRESS],
    ['an Injective address with a bad checksum', `${INJECTIVE_ADDRESS.slice(0, -1)}q`],
    [
      'an Injective address that is not 20 bytes',
      'inj1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5qypqxpq9qcrsszg2pvxqy3rth0',
    ],
    ['an empty string', ''],
  ])('refuses %s', (_name: string, address: string) => {
    expect(() => normalizeInjectiveAddress(address)).toThrow();
    expect(() => getInjectiveEvmAddress(address)).toThrow();
    expect(() => getInjectiveMintRecipient(address)).toThrow();
  });
});
