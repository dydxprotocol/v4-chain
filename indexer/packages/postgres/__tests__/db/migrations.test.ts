import * as OrderTable from '../../src/stores/order-table';
import { knexPrimary } from '../../src/helpers/knex';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import {
  defaultOrder,
} from '../helpers/constants';
import { seedData } from '../helpers/mock-generators';

describe('Test new migration', () => {
  beforeEach(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('test UP and DOWN for most recent migration without seed data', async () => {
    // remove latest migration
    await multiDown(1);

    // re-add latest migration
    await knexPrimary.migrate.latest({ loadExtensions: ['.js'] });

    // re-remove latest migration
    await multiDown(1);
  });

  // NOTE: If a model is modified for a migration then these
  // tests must be skipped until the following migration
  it.skip('[Will fail if a model is modified for migration - see README] test adding most recent migration', async () => {
    // remove latest migration
    await multiDown(1);

    // add data to verify you can roll up and then later roll down
    await seedData();

    // re-add latest migration
    await knexPrimary.migrate.latest({ loadExtensions: ['.js'] });

    // re-remove latest migration
    await multiDown(1);
  });

  // NOTE: If a model is modified for a migration then these
  // tests must be skipped until the following migration
  it.skip('[Will fail if a model is modified for migration - see README] test adding most recent migration with rows that fail index that should only be applied going forward', async () => {
    // remove latest migration
    await multiDown(1);

    // add data to verify you can roll up and then later roll down
    await seedData();
    await OrderTable.create(defaultOrder);

    // re-add latest migration
    await knexPrimary.migrate.latest({ loadExtensions: ['.js'] });

    // re-remove latest migration
    await multiDown(1);
  });
});

/* ------- Helpers ------- */

async function multiDown(downCount: number = 3) {
  for (let i = 0; i < downCount; i += 1) {
    await knexPrimary.migrate.down({ loadExtensions: ['.js'] });
  }
}
