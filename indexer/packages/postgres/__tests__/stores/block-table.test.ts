import { BlockFromDatabase } from 'packages/postgres/src/types';
import * as BlockTable from '../../src/stores/block-table';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import { DateTime } from 'luxon';

const testBlock1 = {
  blockHeight: '1',
  time: DateTime.utc(2025, 3, 5).toISO(),
};

const testBlock2 = {
  blockHeight: '2',
  time: DateTime.utc(2025, 3, 6).toISO(),
};

describe('Block store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a Block', async () => {
    await BlockTable.create(testBlock1);
  });

  it('Successfully finds all Blocks', async () => {
    await Promise.all([
      BlockTable.create(testBlock1),
      BlockTable.create(testBlock2),
    ]);

    const blocks: BlockFromDatabase[] = await BlockTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(blocks.length).toEqual(2);
    expect(blocks[0]).toEqual(expect.objectContaining(testBlock1));
    expect(blocks[1]).toEqual(expect.objectContaining(testBlock2));
  });

  it('Successfully finds Block with block height', async () => {
    await Promise.all([
      BlockTable.create(testBlock1),
      BlockTable.create(testBlock2),
    ]);

    const blocks: BlockFromDatabase[] = await BlockTable.findAll(
      {
        blockHeight: ['1'],
      },
      [],
      { readReplica: true },
    );

    expect(blocks.length).toEqual(1);
    expect(blocks[0]).toEqual(expect.objectContaining(testBlock1));
  });

  it('Successfully finds a Block', async () => {
    await BlockTable.create(testBlock1);

    const block: BlockFromDatabase | undefined = await
    BlockTable.findByBlockHeight(
      testBlock1.blockHeight,
    );

    expect(block).toEqual(expect.objectContaining(testBlock1));
  });

  it('Unable finds a Block', async () => {
    const block: BlockFromDatabase | undefined = await
    BlockTable.findByBlockHeight(
      testBlock1.blockHeight,
    );
    expect(block).toEqual(undefined);
  });

  it('Successfully gets latest Block', async () => {
    await Promise.all([
      BlockTable.create(testBlock1),
      BlockTable.create(testBlock2),
    ]);

    const block: BlockFromDatabase = await BlockTable.getLatest();
    expect(block).toEqual(expect.objectContaining(testBlock2));
  });

  it('Unable to find latest Block', async () => {
    await expect(BlockTable.getLatest()).rejects.toEqual(new Error('Unable to find latest block'));
  });

  it('Successfully finds first block created on or after timestamp', async () => {
    await Promise.all([
      BlockTable.create(testBlock1),
      BlockTable.create(testBlock2),
    ]);

    const block: BlockFromDatabase | undefined = await BlockTable.findBlockByCreatedOnOrAfter(
      DateTime.utc(2025, 3, 5).toISO(),
    );

    expect(block).toBeDefined();
    expect(block).toEqual(expect.objectContaining(testBlock1));
  });

  it('Successfully finds first block when querying with later timestamp', async () => {
    await Promise.all([
      BlockTable.create(testBlock1),
      BlockTable.create(testBlock2),
    ]);

    const block: BlockFromDatabase | undefined = await BlockTable.findBlockByCreatedOnOrAfter(
      DateTime.utc(2025, 3, 6).toISO(),
    );

    expect(block).toBeDefined();
    expect(block).toEqual(expect.objectContaining(testBlock2));
  });

  it('Returns undefined when no blocks found after timestamp', async () => {
    await Promise.all([
      BlockTable.create(testBlock1),
      BlockTable.create(testBlock2),
    ]);

    const block: BlockFromDatabase | undefined = await BlockTable.findBlockByCreatedOnOrAfter(
      DateTime.utc(2025, 3, 7).toISO(),
    );

    expect(block).toBeUndefined();
  });
});
