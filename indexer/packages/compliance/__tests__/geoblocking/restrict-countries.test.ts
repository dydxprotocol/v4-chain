import {
  isRestrictedCountryHeaders,
  CountryHeaders,
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
