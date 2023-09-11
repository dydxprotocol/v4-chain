import { rejectRestrictedCountries } from '../../src/lib/restrict-countries';
import * as utils from '../../src/lib/utils';

const restrictedHeaders = {
  'cf-ipcountry': 'US',
};

const nonRestrictedHeaders = {
  'cf-ipcountry': 'SA',
};

describe('rejectRestrictedCountries', () => {
  let isRestrictedCountrySpy: jest.SpyInstance;
  let req: any;
  let res: any;
  let next: any;

  beforeEach(() => {
    isRestrictedCountrySpy = jest.spyOn(utils, 'isRestrictedCountry');
    req = {
      get: jest.fn().mockReturnThis(),
    };
    res = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn().mockReturnThis(),
      set: jest.fn().mockReturnThis(),
    };
    next = jest.fn();
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('does not reject requests without country headers', () => {
    // empty headers
    req.headers = {};

    rejectRestrictedCountries(req, res, next);
    expect(res.status).not.toHaveBeenCalled();
    expect(next).toHaveBeenCalled();
  });

  it('does not reject requests from non-restricted countries', () => {
    // non-restricted country in header
    req.headers = nonRestrictedHeaders;
    isRestrictedCountrySpy.mockReturnValueOnce(false);

    rejectRestrictedCountries(req, res, next);
    expect(res.status).not.toHaveBeenCalled();
    expect(next).toHaveBeenCalled();
  });

  it('rejects request from restricted countries with a 403', () => {
    // restricted ipcountry
    req.headers = restrictedHeaders;
    isRestrictedCountrySpy.mockReturnValueOnce(true);

    rejectRestrictedCountries(req, res, next);
    expect(res.status).toHaveBeenCalledWith(403);
    expect(next).not.toHaveBeenCalled();
  });
});
