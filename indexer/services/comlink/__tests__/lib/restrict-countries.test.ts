import { INDEXER_GEOBLOCKED_PAYLOAD, isRestrictedCountryHeaders } from '@dydxprotocol-indexer/compliance';
import { testConstants } from '@dydxprotocol-indexer/postgres';
import config from '../../src/config';
import { rejectRestrictedCountries } from '../../src/lib/restrict-countries';
import { BlockedCode } from '../../src/types';
import * as utils from '../../src/lib/utils';
import { matchedData } from 'express-validator';

jest.mock('@dydxprotocol-indexer/compliance');
jest.mock('express-validator');

const restrictedHeaders = {
  'cf-ipcountry': 'US',
};

const nonRestrictedHeaders = {
  'cf-ipcountry': 'SA',
};

const internalIp: string = '3.125.3.24';

describe('rejectRestrictedCountries', () => {
  let isRestrictedCountrySpy: jest.SpyInstance;
  let matchedDataSpy: jest.SpyInstance;
  let req: any;
  let res: any;
  let next: any;

  const defaultEnabled: boolean = config.INDEXER_LEVEL_GEOBLOCKING_ENABLED;

  beforeAll(() => {
    config.INDEXER_LEVEL_GEOBLOCKING_ENABLED = true;
  });

  afterAll(() => {
    config.INDEXER_LEVEL_GEOBLOCKING_ENABLED = defaultEnabled;
  });

  beforeEach(() => {
    isRestrictedCountrySpy = isRestrictedCountryHeaders as unknown as jest.Mock;
    matchedDataSpy = matchedData as unknown as jest.Mock;
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

  it('does not reject requests from non-restricted countries', async () => {
    // non-restricted country in header
    req.headers = nonRestrictedHeaders;
    isRestrictedCountrySpy.mockReturnValueOnce(false);
    matchedDataSpy.mockReturnValue({ address: testConstants.defaultAddress });

    await rejectRestrictedCountries(req, res, next);
    expect(res.status).not.toHaveBeenCalled();
    expect(next).toHaveBeenCalled();
  });

  it('rejects request from restricted countries with a 403', async () => {
    // restricted ipcountry
    req.headers = restrictedHeaders;
    isRestrictedCountrySpy.mockReturnValueOnce(true);
    matchedDataSpy.mockReturnValue({ address: testConstants.defaultAddress });

    await rejectRestrictedCountries(req, res, next);
    expect(res.status).toHaveBeenCalledWith(403);
    expect(res.json).toHaveBeenCalledWith(expect.objectContaining({
      errors: expect.arrayContaining([
        {
          msg: INDEXER_GEOBLOCKED_PAYLOAD,
          code: BlockedCode.GEOBLOCKED,
        },
      ]),
    }));
    expect(next).not.toHaveBeenCalled();
  });

  it('does not check headers for internal indexer ip address', async () => {
    // restricted ipcountry
    req.headers = restrictedHeaders;
    isRestrictedCountrySpy.mockReturnValueOnce(true);
    matchedDataSpy.mockReturnValue({ address: testConstants.defaultAddress });
    jest.spyOn(utils, 'getIpAddr').mockReturnValue(internalIp);
    jest.spyOn(utils, 'isIndexerIp').mockImplementation((ip: string): boolean => ip === internalIp);

    await rejectRestrictedCountries(req, res, next);
    expect(res.status).not.toHaveBeenCalled();
    expect(next).toHaveBeenCalled();
  });
});
