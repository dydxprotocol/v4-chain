import {
  isRestrictedCountryHeaders,
  CountryHeaders,
  isWhitelistedAddress,
} from '../../src/geoblocking/restrict-countries';
import * as util from '../../src/geoblocking/util';
import config from '../../src/config';

const defaultHeaders: CountryHeaders = {
  'cf-ipcountry': 'US',
};

describe('isRestrictedCountryHeaders', () => {
  let isRestrictedCountrySpy: jest.SpyInstance;
  const defaultEnabled: boolean = config.INDEXER_LEVEL_GEOBLOCKING_ENABLED;

  beforeAll(() => {
    config.INDEXER_LEVEL_GEOBLOCKING_ENABLED = true;
  });

  afterAll(() => {
    config.INDEXER_LEVEL_GEOBLOCKING_ENABLED = defaultEnabled;
  });

  beforeEach(() => {
    isRestrictedCountrySpy = jest.spyOn(util, 'isRestrictedCountry');
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('does not reject headers from non-restricted countries', () => {
    // non-restricted country in header
    isRestrictedCountrySpy.mockReturnValue(false);

    expect(isRestrictedCountryHeaders(defaultHeaders)).toEqual(false);
  });

  it('does reject headers with restricted country', () => {
    // restricted country in header
    isRestrictedCountrySpy.mockReturnValue(true);

    expect(isRestrictedCountryHeaders(defaultHeaders)).toEqual(true);
  });

  it('does reject empty headers', () => {
    expect(isRestrictedCountryHeaders({})).toEqual(true);
  });
});

describe('isWhitelistedAddress', () => {
  it('returns true if address is whitelisted', () => {
    config.WHITELISTED_ADDRESSES = '0x123,0x456';

    expect(isWhitelistedAddress('0x123')).toEqual(true);
  });

  it('returns false if address is not whitelisted', () => {
    config.WHITELISTED_ADDRESSES = '0x123,0x456';

    expect(isWhitelistedAddress('0x789')).toEqual(false);
  });
});
