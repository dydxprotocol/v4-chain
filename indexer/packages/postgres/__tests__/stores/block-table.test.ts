import { BlockFromDatabase } from 'packages/postgres/src/types';
import * as BlockTable from '../../src/stores/block-table';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import { defaultBlock, defaultBlock2 } from '../helpers/constants';

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
    await BlockTable.create(defaultBlock);
  });

  it('Successfully finds all Blocks', async () => {
    await Promise.all([
      BlockTable.create(defaultBlock),
      BlockTable.create(defaultBlock2),
    ]);

    const blocks: BlockFromDatabase[] = await BlockTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(blocks.length).toEqual(2);
    expect(blocks[0]).toEqual(expect.objectContaining(defaultBlock));
    expect(blocks[1]).toEqual(expect.objectContaining(defaultBlock2));
  });

  it('Successfully finds Block with block height', async () => {
    await Promise.all([
      BlockTable.create(defaultBlock),
      BlockTable.create(defaultBlock2),
    ]);

    const blocks: BlockFromDatabase[] = await BlockTable.findAll(
      {
        blockHeight: ['1'],
      },
      [],
      { readReplica: true },
    );

    expect(blocks.length).toEqual(1);
    expect(blocks[0]).toEqual(expect.objectContaining(defaultBlock));
  });

  it('Successfully finds a Block', async () => {
    await BlockTable.create(defaultBlock);

    const block: BlockFromDatabase | undefined = await
    BlockTable.findByBlockHeight(
      defaultBlock.blockHeight,
    );

    expect(block).toEqual(expect.objectContaining(defaultBlock));
  });

  it('Unable finds a Block', async () => {
    const block: BlockFromDatabase | undefined = await
    BlockTable.findByBlockHeight(
      defaultBlock.blockHeight,
    );
    expect(block).toEqual(undefined);
  });

  it('Successfully gets latest Block', async () => {
    await Promise.all([
      BlockTable.create(defaultBlock),
      BlockTable.create(defaultBlock2),
    ]);

    const block: BlockFromDatabase | undefined = await BlockTable.getLatest();
    expect(block).toEqual(expect.objectContaining(defaultBlock2));
  });

  it('Unable to find latest Block', async () => {
    const block: BlockFromDatabase | undefined = await BlockTable.getLatest();
    expect(block).toEqual(undefined);
  });
});
