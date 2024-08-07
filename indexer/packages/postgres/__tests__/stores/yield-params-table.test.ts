import {
    Ordering,
    YieldParamsColumns,
    YieldParamsFromDatabase,
  } from '../../src/types';
  import * as YieldParamsTable from '../../src/stores/yield-params-table';
  import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
  import { seedData } from '../helpers/mock-generators';
  import {
    defaultYieldParams1,
    defaultYieldParams2,
  } from '../helpers/constants';
  import { DateTime } from 'luxon';
  
  describe('Yield params store', () => {
    beforeEach(async () => {
      await seedData();
    });
  
    beforeAll(async () => {
      await migrate();
    });
  
    afterEach(async () => {
      await clearData();
    });
  
    afterAll(async () => {
      await teardown();
    });
  
    it('Successfully creates new yield params', async () => {
      await YieldParamsTable.create(defaultYieldParams1);
    });

    it('Successfully creates multiple new yield params', async () => {
        await YieldParamsTable.create(defaultYieldParams1);
        await YieldParamsTable.create(defaultYieldParams2);
    });

    it('Succesfully creates yield params and finds it', async () => {
        await YieldParamsTable.create(defaultYieldParams1);
        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll({}, [], {});
    
        expect(yieldParams.length).toEqual(1);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
    })

    it('Succesfully creates multiple yield params and finds them', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll({}, [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(2);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
        expect(yieldParams[1]).toEqual(expect.objectContaining(defaultYieldParams2));
    })
  
    it('Successfully finds all yield params with id', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { id: [YieldParamsTable.uuid(defaultYieldParams1.createdAtHeight)] }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(1);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
    });

    it('Successfully finds all yield params with sDAI price', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { sDAIPrice: defaultYieldParams1.sDAIPrice }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(1);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
    });

    it('Successfully finds all yield params with asset yield index', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { assetYieldIndex: defaultYieldParams1.assetYieldIndex }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(1);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
    });

    it('Successfully finds all yield params at height: Finds all yield params', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdAtHeight: [defaultYieldParams1.createdAtHeight, defaultYieldParams2.createdAtHeight] }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(2);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
        expect(yieldParams[1]).toEqual(expect.objectContaining(defaultYieldParams2));
    });

    it('Successfully finds all yield params at height: Find one set of yield params', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdAtHeight: [defaultYieldParams2.createdAtHeight] }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(1);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams2));
    });

    it('Successfully finds all yield params at height: No yield params found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const latestHeightPlusOne: string = (parseInt(defaultYieldParams2.createdAtHeight) + 1).toString();

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdAtHeight: [latestHeightPlusOne] }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(0);
    });

    it('Successfully finds all yield params before or at height: All yield params found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdBeforeOrAtHeight: defaultYieldParams2.createdAtHeight }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(2);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
        expect(yieldParams[1]).toEqual(expect.objectContaining(defaultYieldParams2));
    });

    it('Successfully finds all yield params before or at height: One set of yield params found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdBeforeOrAtHeight: defaultYieldParams1.createdAtHeight }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(1);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
    });

    it('Successfully finds all yield params before or at height: No yield parms to be found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const firsHeightMinusOne: string = (parseInt(defaultYieldParams1.createdAtHeight) - 1).toString();

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdBeforeOrAtHeight: firsHeightMinusOne }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(0);
    });

    it('Successfully finds all yield params after height: All yield params found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const firsHeightMinusOne: string = (parseInt(defaultYieldParams1.createdAtHeight) - 1).toString();

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdAfterHeight: firsHeightMinusOne }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(2);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
        expect(yieldParams[1]).toEqual(expect.objectContaining(defaultYieldParams2));
    });

    it('Successfully finds all yield params after height: One set of yield params found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdAfterHeight: defaultYieldParams1.createdAtHeight }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(1);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams2));
    });

    it('Successfully finds all yield params after height: No yield params found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdAfterHeight: defaultYieldParams2.createdAtHeight }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(0);
    });

    it('Successfully finds all yield params at time: One set of yield params found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdAt: defaultYieldParams1.createdAt }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(1);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
    });

    it('Successfully finds all yield params at time: No yield params found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const latestCreatedAtPlusOne = DateTime.fromISO(defaultYieldParams2.createdAt).plus({ days: 1 }).toISO();

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdAt: latestCreatedAtPlusOne }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(0);
    });

    it('Successfully finds all yield params before or at time: All yield params found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdBeforeOrAt: defaultYieldParams2.createdAt }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(2);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
        expect(yieldParams[1]).toEqual(expect.objectContaining(defaultYieldParams2));
    });

    it('Successfully finds all yield params before or at time: One set of yield params found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdBeforeOrAt: defaultYieldParams1.createdAt }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(1);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
    });

    it('Successfully finds all yield params before or at time: No yield parms found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const latestCreatedAtMinusOne = DateTime.fromISO(defaultYieldParams1.createdAt).minus({ days: 1 }).toISO();

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdBeforeOrAt: latestCreatedAtMinusOne }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(0);
    });

    it('Successfully finds all yield params after time: All yield params found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);
        const latestCreatedAtMinusOne = DateTime.fromISO(defaultYieldParams1.createdAt).minus({ days: 1 }).toISO();

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdAfter: latestCreatedAtMinusOne }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(2);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
        expect(yieldParams[1]).toEqual(expect.objectContaining(defaultYieldParams2));
    });

    it('Successfully finds all yield params after time: One set of yield params found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdAfter: defaultYieldParams1.createdAt }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(1);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams2));
    });

    it('Successfully finds all yield params : No yield params found', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdAfter: defaultYieldParams2.createdAt }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(0);
    });

    it('Successfully finds yield params before or at height and with specific assetYieldIndex', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdBeforeOrAtHeight: defaultYieldParams2.createdAtHeight,
              assetYieldIndex: defaultYieldParams1.assetYieldIndex,
            }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(1);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
    });

    it('Successfully finds yield params before or at height with asset yield index and sDAI price', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdBeforeOrAtHeight: defaultYieldParams2.createdAtHeight,
              assetYieldIndex: defaultYieldParams1.assetYieldIndex,
              sDAIPrice: defaultYieldParams1.sDAIPrice,
            }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(1);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams1));
    });

    it('Successfully finds yield params at time with asset yield index and sDAI price', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { createdAt: defaultYieldParams2.createdAt,
              assetYieldIndex: defaultYieldParams2.assetYieldIndex,
              sDAIPrice: defaultYieldParams2.sDAIPrice,
            }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(1);
        expect(yieldParams[0]).toEqual(expect.objectContaining(defaultYieldParams2));
    });

    it('Finds no yield params on findAll parameters mismatch', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase[] = await YieldParamsTable.findAll(
            { assetYieldIndex: defaultYieldParams1.assetYieldIndex,
              sDAIPrice: defaultYieldParams2.sDAIPrice,
            }, 
            [], {
            orderBy: [[YieldParamsColumns.createdAtHeight, Ordering.ASC]],
        });
      
        expect(yieldParams.length).toEqual(0);
    });
  
    it('Successfully finds yield params by from Id', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase | undefined = await YieldParamsTable.findById(
            YieldParamsTable.uuid(
                defaultYieldParams1.createdAtHeight,
            )
        );

        expect(yieldParams).toBeDefined();
        expect(yieldParams).toEqual(expect.objectContaining(defaultYieldParams1));
    });

    it('Successfully finds gets latest yield params: Multiple sets of yield params stored', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
            YieldParamsTable.create(defaultYieldParams2),
        ]);

        const yieldParams: YieldParamsFromDatabase | undefined = await YieldParamsTable.getLatest();

        expect(yieldParams).toBeDefined();
        expect(yieldParams).toEqual(expect.objectContaining(defaultYieldParams2));
    });

    it('Successfully finds gets latest yield params: One set of yield params stored', async () => {
        await Promise.all([
            YieldParamsTable.create(defaultYieldParams1),
        ]);

        const yieldParams: YieldParamsFromDatabase | undefined = await YieldParamsTable.getLatest();

        expect(yieldParams).toBeDefined();
        expect(yieldParams).toEqual(expect.objectContaining(defaultYieldParams1));
    });

    it('Successfully finds gets latest yield params: No yield params stored', async () => {
        await expect(YieldParamsTable.getLatest()).rejects.toThrow('Unable to find latest yield params');
    });
  })
  