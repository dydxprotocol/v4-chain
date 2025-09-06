import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import * as PermissionApprovalTable from '../../src/stores/permission-approval-table';

describe('PermissionApproval store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  const defaultPermissionApprovals = [
    {
      suborg_id: 'suborg-1',
      chain_id: 'arbitrum',
      approval: '0x1234567890abcdef1234567890abcdef12345678',
    },
    {
      suborg_id: 'suborg-1',
      chain_id: 'base',
      approval: '0x2345678901bcdef12345678901bcdef123456789',
    },
    {
      suborg_id: 'suborg-1',
      chain_id: 'avalanche',
      approval: '0x3456789012cdef123456789012cdef1234567890',
    },
    {
      suborg_id: 'suborg-1',
      chain_id: 'optimism',
      approval: '0x4567890123def1234567890123def12345678901',
    },
    {
      suborg_id: 'suborg-1',
      chain_id: 'ethereum',
      approval: '0x567890123def1234567890123def123456789012',
    },
  ];

  const suborg2Approvals = [
    {
      suborg_id: 'suborg-2',
      chain_id: 'arbitrum',
      approval: '0xabcdef1234567890abcdef1234567890abcdef12',
    },
    {
      suborg_id: 'suborg-2',
      chain_id: 'base',
      approval: '0xbcdef12345678901bcdef12345678901bcdef123',
    },
  ];

  const partialPermissionApprovals = [
    {
      suborg_id: 'suborg-3',
      chain_id: 'arbitrum',
      approval: '0x1111111111111111111111111111111111111111',
    },
    {
      suborg_id: 'suborg-3',
      chain_id: 'ethereum',
      approval: '0x2222222222222222222222222222222222222222',
    },
  ];

  describe('create', () => {
    it('Successfully creates a PermissionApproval', async () => {
      const createdApproval = await PermissionApprovalTable.create(defaultPermissionApprovals[0]);
      expect(createdApproval).toEqual(expect.objectContaining(defaultPermissionApprovals[0]));
    });

    it('Successfully creates multiple PermissionApprovals for same suborg', async () => {
      for (const approval of defaultPermissionApprovals) {
        const createdApproval = await PermissionApprovalTable.create(approval);
        expect(createdApproval).toEqual(expect.objectContaining(approval));
      }
    });
  });

  describe('findBySuborgId', () => {
    it('Successfully finds all PermissionApprovals by suborg_id', async () => {
      for (const approval of defaultPermissionApprovals) {
        await PermissionApprovalTable.create(approval);
      }
      for (const approval of suborg2Approvals) {
        await PermissionApprovalTable.create(approval);
      }

      const approvals = await PermissionApprovalTable.findBySuborgId('suborg-1');

      expect(approvals).toHaveLength(5);
      expect(approvals).toEqual(expect.arrayContaining(defaultPermissionApprovals));
    });

    it('Returns empty array when PermissionApproval not found by suborg_id', async () => {
      for (const approval of defaultPermissionApprovals) {
        await PermissionApprovalTable.create(approval);
      }

      const approvals = await PermissionApprovalTable.findBySuborgId('nonexistent-suborg');

      expect(approvals).toEqual([]);
    });
  });

  describe('findBySuborgIdAndChainId', () => {
    it('Successfully finds a PermissionApproval by suborg_id and chain_id', async () => {
      for (const approval of defaultPermissionApprovals) {
        await PermissionApprovalTable.create(approval);
      }

      const approval = await PermissionApprovalTable.findBySuborgIdAndChainId('suborg-1', 'arbitrum');

      expect(approval).toEqual(expect.objectContaining(defaultPermissionApprovals[0]));
    });

    it('Returns undefined when PermissionApproval not found by suborg_id and chain_id', async () => {
      for (const approval of defaultPermissionApprovals) {
        await PermissionApprovalTable.create(approval);
      }

      const approval = await PermissionApprovalTable.findBySuborgIdAndChainId('nonexistent-suborg', 'arbitrum');

      expect(approval).toBeUndefined();
    });
  });

  describe('getApprovalForSuborgAndChain', () => {
    it('Successfully gets approval for existing suborg and chain', async () => {
      const arbitrumApproval = defaultPermissionApprovals[0];
      await PermissionApprovalTable.create(arbitrumApproval);

      const approval = await PermissionApprovalTable.getApprovalForSuborgAndChain(arbitrumApproval.suborg_id, 'arbitrum');

      expect(approval).toBe(arbitrumApproval.approval);
    });

    it('Returns undefined when suborg does not exist', async () => {
      const approval = await PermissionApprovalTable.getApprovalForSuborgAndChain('nonexistent-suborg', 'arbitrum');

      expect(approval).toBeUndefined();
    });

    it('Returns undefined when chain approval does not exist for suborg', async () => {
      // Create only base approval for this suborg
      const baseApproval = {
        suborg_id: 'suborg-no-arbitrum',
        chain_id: 'base',
        approval: '0x1234567890abcdef1234567890abcdef12345678',
      };
      await PermissionApprovalTable.create(baseApproval);

      const approval = await PermissionApprovalTable.getApprovalForSuborgAndChain(baseApproval.suborg_id, 'arbitrum');

      expect(approval).toBeUndefined();
    });

    it('Successfully gets approvals for different chains', async () => {
      for (const approvalData of defaultPermissionApprovals) {
        await PermissionApprovalTable.create(approvalData);
      }

      const arbitrumApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('suborg-1', 'arbitrum');
      const baseApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('suborg-1', 'base');
      const avalancheApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('suborg-1', 'avalanche');
      const optimismApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('suborg-1', 'optimism');
      const ethereumApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('suborg-1', 'ethereum');

      expect(arbitrumApproval).toBe(defaultPermissionApprovals[0].approval);
      expect(baseApproval).toBe(defaultPermissionApprovals[1].approval);
      expect(avalancheApproval).toBe(defaultPermissionApprovals[2].approval);
      expect(optimismApproval).toBe(defaultPermissionApprovals[3].approval);
      expect(ethereumApproval).toBe(defaultPermissionApprovals[4].approval);
    });
  });

  describe('update', () => {
    it('Successfully updates a PermissionApproval', async () => {
      const originalApproval = defaultPermissionApprovals[0];
      await PermissionApprovalTable.create(originalApproval);

      const updatedFields = {
        suborg_id: originalApproval.suborg_id,
        chain_id: originalApproval.chain_id,
        approval: '0xnewArbitrumApproval123456789012345678901',
      };

      const updatedApproval = await PermissionApprovalTable.update(updatedFields);

      expect(updatedApproval).toEqual(expect.objectContaining(updatedFields));
    });

    it('Returns undefined when trying to update non-existent PermissionApproval', async () => {
      const updatedFields = {
        suborg_id: 'nonexistent-suborg',
        chain_id: 'arbitrum',
        approval: '0xnewApproval123456789012345678901234567',
      };

      const result = await PermissionApprovalTable.update(updatedFields);

      expect(result).toBeUndefined();
    });
  });

  describe('upsert', () => {
    it('Successfully creates a new PermissionApproval via upsert', async () => {
      const newApproval = {
        suborg_id: 'new-suborg',
        chain_id: 'arbitrum',
        approval: '0xnewApproval123456789012345678901234567',
      };

      const result = await PermissionApprovalTable.upsert(newApproval);

      expect(result).toEqual(expect.objectContaining(newApproval));
    });

    it('Successfully updates existing PermissionApproval via upsert', async () => {
      const originalApproval = defaultPermissionApprovals[0];
      await PermissionApprovalTable.create(originalApproval);

      const updatedApproval = {
        suborg_id: originalApproval.suborg_id,
        chain_id: originalApproval.chain_id,
        approval: '0xupdatedApproval123456789012345678901234567',
      };

      const result = await PermissionApprovalTable.upsert(updatedApproval);

      expect(result).toEqual(expect.objectContaining(updatedApproval));
    });
  });

  describe('Edge cases', () => {
    it('Handles empty string suborg_id searches', async () => {
      for (const approval of defaultPermissionApprovals) {
        await PermissionApprovalTable.create(approval);
      }

      const approvals = await PermissionApprovalTable.findBySuborgId('');
      const arbitrumApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('', 'arbitrum');
      const baseApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('', 'base');
      const avalancheApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('', 'avalanche');
      const optimismApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('', 'optimism');
      const ethereumApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('', 'ethereum');

      expect(approvals).toEqual([]);
      expect(arbitrumApproval).toBeUndefined();
      expect(baseApproval).toBeUndefined();
      expect(avalancheApproval).toBeUndefined();
      expect(optimismApproval).toBeUndefined();
      expect(ethereumApproval).toBeUndefined();
    });

    it('Handles special characters in suborg_id', async () => {
      const specialSuborgApproval = {
        suborg_id: 'suborg-with-special-chars-!@#$%^&*()',
        chain_id: 'arbitrum',
        approval: '0x1234567890abcdef1234567890abcdef12345678',
      };

      await PermissionApprovalTable.create(specialSuborgApproval);

      const approvals = await PermissionApprovalTable.findBySuborgId(
        specialSuborgApproval.suborg_id,
      );
      const arbitrumApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain(specialSuborgApproval.suborg_id, 'arbitrum');

      expect(approvals).toHaveLength(1);
      expect(approvals[0]).toEqual(expect.objectContaining(specialSuborgApproval));
      expect(arbitrumApproval).toBe(specialSuborgApproval.approval);
    });

    it('Handles very long approval addresses', async () => {
      const longApprovalAddress = `0x${'a'.repeat(64)}`; // Very long hex address
      const longApprovalData = {
        suborg_id: 'suborg-long-approval',
        chain_id: 'arbitrum',
        approval: longApprovalAddress,
      };

      await PermissionApprovalTable.create(longApprovalData);

      const approval = await PermissionApprovalTable.getApprovalForSuborgAndChain(longApprovalData.suborg_id, 'arbitrum');

      expect(approval).toBe(longApprovalAddress);
    });

    it('Handles multiple chain approvals for same suborg', async () => {
      for (const approval of defaultPermissionApprovals) {
        await PermissionApprovalTable.create(approval);
      }

      // Get all approvals for the same suborg using the generic function
      const arbitrumApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('suborg-1', 'arbitrum');
      const baseApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('suborg-1', 'base');
      const avalancheApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('suborg-1', 'avalanche');
      const optimismApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('suborg-1', 'optimism');
      const ethereumApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('suborg-1', 'ethereum');

      expect(arbitrumApproval).toBe(defaultPermissionApprovals[0].approval);
      expect(baseApproval).toBe(defaultPermissionApprovals[1].approval);
      expect(avalancheApproval).toBe(defaultPermissionApprovals[2].approval);
      expect(optimismApproval).toBe(defaultPermissionApprovals[3].approval);
      expect(ethereumApproval).toBe(defaultPermissionApprovals[4].approval);
    });

    it('Handles partial approvals for suborg (only some chains)', async () => {
      for (const approval of partialPermissionApprovals) {
        await PermissionApprovalTable.create(approval);
      }

      const approvals = await PermissionApprovalTable.findBySuborgId('suborg-3');
      const arbitrumApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('suborg-3', 'arbitrum');
      const baseApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('suborg-3', 'base');
      const ethereumApproval = await PermissionApprovalTable.getApprovalForSuborgAndChain('suborg-3', 'ethereum');

      expect(approvals).toHaveLength(2);
      expect(arbitrumApproval).toBe(partialPermissionApprovals[0].approval);
      expect(baseApproval).toBeUndefined(); // No base approval for this suborg
      expect(ethereumApproval).toBe(partialPermissionApprovals[1].approval);
    });
  });
});
