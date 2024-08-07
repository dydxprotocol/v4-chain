import {
    dbHelpers,
    testConstants,
    testMocks,
    YieldParamsTable,
  } from '@dydxprotocol-indexer/postgres';
  import { RequestMethod, YieldParamsResponseObject } from '../../../../src/types';
  import request from 'supertest';
  import { sendRequest } from '../../../helpers/helpers';
  
  describe('yield-params-controller#V4', () => {
    beforeAll(async () => {
      await dbHelpers.migrate();
    });
  
    afterAll(async () => {
      await dbHelpers.teardown();
    });
  
    describe('GET', () => {
      afterEach(async () => {
        await dbHelpers.clearData();
      });

      it('Get /yieldParams returns yieldParams when there are no yield params', async () => {
        await testMocks.seedData();

        const height: string = '5';

        const response: request.Response = await sendRequest({
            type: RequestMethod.GET,
            path: `/v4/yieldParams?createdAtOrBeforeHeight=${height}`,
          });

        expect(response.body.allYieldParams).toHaveLength(0);
      });

      it('Get /yieldParams returns all yieldParams when no height specified', async () => {
        await testMocks.seedData();

        await Promise.all([
            YieldParamsTable.create(testConstants.defaultYieldParams1),
            YieldParamsTable.create(testConstants.defaultYieldParams2),
        ]);
  
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/yieldParams?createdAtOrBeforeHeight`,
        });

        const expectedYieldParamsResponse1: YieldParamsResponseObject = {
          id: YieldParamsTable.uuid(testConstants.defaultYieldParams1.createdAtHeight),
          sDAIPrice: testConstants.defaultYieldParams1.sDAIPrice,
          assetYieldIndex: testConstants.defaultYieldParams1.assetYieldIndex,
          createdAt: testConstants.defaultYieldParams1.createdAt,
          createdAtHeight: testConstants.defaultYieldParams1.createdAtHeight,
      }

      const expectedYieldParamsResponse2: YieldParamsResponseObject = {
          id: YieldParamsTable.uuid(testConstants.defaultYieldParams2.createdAtHeight),
          sDAIPrice: testConstants.defaultYieldParams2.sDAIPrice,
          assetYieldIndex: testConstants.defaultYieldParams2.assetYieldIndex,
          createdAt: testConstants.defaultYieldParams2.createdAt,
          createdAtHeight: testConstants.defaultYieldParams2.createdAtHeight,
      }

      expect(response.body.allYieldParams).toHaveLength(2);
      expect(response.body.allYieldParams).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
              ...expectedYieldParamsResponse1,
          }),
          expect.objectContaining({
              ...expectedYieldParamsResponse2,
          }),
        ]),
      );
      });

      it('Get /yieldParams returns yieldParams at or before height for one row', async () => {
        await testMocks.seedData();

        await Promise.all([
            YieldParamsTable.create(testConstants.defaultYieldParams1),
            YieldParamsTable.create(testConstants.defaultYieldParams2),
        ]);
  
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/yieldParams?createdAtOrBeforeHeight=${testConstants.defaultYieldParams1.createdAtHeight}`,
        });

        const expectedYieldParamsResponse: YieldParamsResponseObject = {
            id: YieldParamsTable.uuid(testConstants.defaultYieldParams1.createdAtHeight),
            sDAIPrice: testConstants.defaultYieldParams1.sDAIPrice,
            assetYieldIndex: testConstants.defaultYieldParams1.assetYieldIndex,
            createdAt: testConstants.defaultYieldParams1.createdAt,
            createdAtHeight: testConstants.defaultYieldParams1.createdAtHeight,
        }
        
        expect(response.body.allYieldParams).toHaveLength(1);
        expect(response.body.allYieldParams).toEqual(
          expect.arrayContaining([
            expect.objectContaining({
                ...expectedYieldParamsResponse,
            }),
          ]),
        );
      });

      it('Get /yieldParams returns yieldParams at or before height for two rows', async () => {
        await testMocks.seedData();

        await Promise.all([
            YieldParamsTable.create(testConstants.defaultYieldParams1),
            YieldParamsTable.create(testConstants.defaultYieldParams2),
        ]);
  
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/yieldParams?createdAtOrBeforeHeight=${testConstants.defaultYieldParams2.createdAtHeight}`,
        });

        const expectedYieldParamsResponse1: YieldParamsResponseObject = {
            id: YieldParamsTable.uuid(testConstants.defaultYieldParams1.createdAtHeight),
            sDAIPrice: testConstants.defaultYieldParams1.sDAIPrice,
            assetYieldIndex: testConstants.defaultYieldParams1.assetYieldIndex,
            createdAt: testConstants.defaultYieldParams1.createdAt,
            createdAtHeight: testConstants.defaultYieldParams1.createdAtHeight,
        }

        const expectedYieldParamsResponse2: YieldParamsResponseObject = {
            id: YieldParamsTable.uuid(testConstants.defaultYieldParams2.createdAtHeight),
            sDAIPrice: testConstants.defaultYieldParams2.sDAIPrice,
            assetYieldIndex: testConstants.defaultYieldParams2.assetYieldIndex,
            createdAt: testConstants.defaultYieldParams2.createdAt,
            createdAtHeight: testConstants.defaultYieldParams2.createdAtHeight,
        }
        
        expect(response.body.allYieldParams).toHaveLength(2);
        expect(response.body.allYieldParams).toEqual(
          expect.arrayContaining([
            expect.objectContaining({
                ...expectedYieldParamsResponse1,
            }),
            expect.objectContaining({
                ...expectedYieldParamsResponse2,
            }),
          ]),
        );
      });

      it('Get /yieldParams/latestYieldParams throws error when no rows found', async () => {
        await testMocks.seedData();
        const response: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/yieldParams/latestYieldParams`,
          expectedStatus: 500,
        });

        const expectedErrorMsg: string = 'Internal Server Error';

        expect(response.body).toEqual(expect.objectContaining({
          errors: expect.arrayContaining([
            expect.objectContaining({
              msg: expectedErrorMsg,
            }),
          ]),
        }));
      });

      it('Get /yieldParams/latestYieldParams returns yieldParams succesfully gets latest yield params when one row is in DB', async () => {
        await testMocks.seedData();

        await Promise.all([
            YieldParamsTable.create(testConstants.defaultYieldParams1),
        ]);

        const response: request.Response = await sendRequest({
            type: RequestMethod.GET,
            path: `/v4/yieldParams/latestYieldParams`,
          });

        const expectedYieldParamsResponse1: YieldParamsResponseObject = {
            id: YieldParamsTable.uuid(testConstants.defaultYieldParams1.createdAtHeight),
            sDAIPrice: testConstants.defaultYieldParams1.sDAIPrice,
            assetYieldIndex: testConstants.defaultYieldParams1.assetYieldIndex,
            createdAt: testConstants.defaultYieldParams1.createdAt,
            createdAtHeight: testConstants.defaultYieldParams1.createdAtHeight,
        }

        expect(response.body.allYieldParams).toHaveLength(1);
        expect(response.body.allYieldParams).toEqual(
          expect.arrayContaining([
            expect.objectContaining({
                ...expectedYieldParamsResponse1,
            }),
          ]),
        );
      });

      it('Get /yieldParams/latestYieldParams returns yieldParams succesfully gets latest yield params when multiple rows are in DB', async () => {
        await testMocks.seedData();

        await Promise.all([
            YieldParamsTable.create(testConstants.defaultYieldParams1),
            YieldParamsTable.create(testConstants.defaultYieldParams2),
        ]);

        const response: request.Response = await sendRequest({
            type: RequestMethod.GET,
            path: `/v4/yieldParams/latestYieldParams`,
          });

        const expectedYieldParamsResponse2: YieldParamsResponseObject = {
            id: YieldParamsTable.uuid(testConstants.defaultYieldParams2.createdAtHeight),
            sDAIPrice: testConstants.defaultYieldParams2.sDAIPrice,
            assetYieldIndex: testConstants.defaultYieldParams2.assetYieldIndex,
            createdAt: testConstants.defaultYieldParams2.createdAt,
            createdAtHeight: testConstants.defaultYieldParams2.createdAtHeight,
        }
        
        expect(response.body.allYieldParams).toHaveLength(1);
        expect(response.body.allYieldParams).toEqual(
          expect.arrayContaining([
            expect.objectContaining({
                ...expectedYieldParamsResponse2,
            }),
          ]),
        );
      });
    });
  });
  