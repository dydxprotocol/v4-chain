import { stats } from '@dydxprotocol-indexer/base';
import { dbHelpers, FirebaseNotificationTokenTable, testMocks } from '@dydxprotocol-indexer/postgres';

import config from '../../src/config';
import runTask from '../../src/tasks/delete-old-firebase-notification-tokens';
import { defaultWallet } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

describe('delete-old-firebase-notification-tokens', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    jest.spyOn(stats, 'timing');
  });

  beforeEach(async () => {
    await testMocks.seedData();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  it('deletes old Firebase notification tokens', async () => {
    // Create test data
    const currentDate = new Date();
    const oldDate = new Date(currentDate.getTime() - 40 * 24 * 60 * 60 * 1000); // 40 days ago
    const recentDate = new Date(currentDate.getTime() - 15 * 24 * 60 * 60 * 1000); // 15 days ago

    await FirebaseNotificationTokenTable.create({
      token: 'old_token',
      updatedAt: oldDate.toISOString(),
      address: defaultWallet.address,
      language: 'en',
    });
    await FirebaseNotificationTokenTable.create({
      token: 'recent_token',
      updatedAt: recentDate.toISOString(),
      address: defaultWallet.address,
      language: 'fr',
    });

    const initialTokens = await FirebaseNotificationTokenTable.findAll({}, []);
    expect(initialTokens.length).toBe(3);

    // Run the task
    await runTask();

    // Check if old token was deleted and recent token remains
    const remainingTokens = await FirebaseNotificationTokenTable.findAll({}, []);
    expect(remainingTokens.length).toBe(2);
    expect(remainingTokens[0].token).toBe('recent_token');

    // Check if stats.timing was called
    expect(stats.timing).toHaveBeenCalledWith(
      expect.stringContaining(`${config.SERVICE_NAME}.delete_old_firebase_notification_tokens`),
      expect.any(Number),
    );
  });
});
