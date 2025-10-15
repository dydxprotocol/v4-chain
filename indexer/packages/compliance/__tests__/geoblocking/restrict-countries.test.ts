import { GeoOriginHeaders, GeoOriginStatus } from '../../src/types';
import {
  isRestrictedCountryHeaders,
  isWhitelistedAddress,
} from '../../src/geoblocking/restrict-countries';
import * as restrictCountries from '../../src/geoblocking/restrict-countries';
import config from '../../src/config';

const defaultHeaders: GeoOriginHeaders = {
  'geo-origin-country': 'FR',
  'geo-origin-region': 'FR-75', // Paris
  'geo-origin-status': GeoOriginStatus.OK,
};

describe('isRestrictedCountryHeaders', () => {
  let isRestrictedCountryHeadersSpy: jest.SpyInstance;
  const defaultEnabled: boolean = config.INDEXER_LEVEL_GEOBLOCKING_ENABLED;

  beforeAll(() => {
    config.INDEXER_LEVEL_GEOBLOCKING_ENABLED = true;
  });

  afterAll(() => {
    config.INDEXER_LEVEL_GEOBLOCKING_ENABLED = defaultEnabled;
  });

  beforeEach(() => {
    isRestrictedCountryHeadersSpy = jest.spyOn(restrictCountries, 'isRestrictedCountryHeaders');
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('does not reject headers from non-restricted countries', () => {
    // non-restricted country in header
    isRestrictedCountryHeadersSpy.mockReturnValue(false);

    expect(isRestrictedCountryHeaders(defaultHeaders)).toEqual(false);
  });

  it('does reject headers with restricted country', () => {
    // restricted country in header
    isRestrictedCountryHeadersSpy.mockReturnValue(true);

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
