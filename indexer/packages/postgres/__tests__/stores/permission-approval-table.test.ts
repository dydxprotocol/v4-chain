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

  const defaultPermissionApproval1 = {
    suborg_id: 'suborg-1',
    arbitrum_approval: '0x1234567890abcdef1234567890abcdef12345678',
    base_approval: '0x2345678901bcdef12345678901bcdef123456789',
    avalanche_approval: '0x3456789012cdef123456789012cdef1234567890',
    optimism_approval: '0x4567890123def1234567890123def12345678901',
    ethereum_approval: '0x567890123def1234567890123def123456789012',
  };

  const defaultPermissionApproval2 = {
    suborg_id: 'suborg-2',
    arbitrum_approval: '0xabcdef1234567890abcdef1234567890abcdef12',
    base_approval: '0xbcdef12345678901bcdef12345678901bcdef123',
    avalanche_approval: '0xcdef123456789012cdef123456789012cdef1234',
    optimism_approval: '0xdef1234567890123def1234567890123def12345',
    ethereum_approval: '0xef12345678901234ef12345678901234ef123456',
  };

  const partialPermissionApproval = {
    suborg_id: 'suborg-3',
    arbitrum_approval: '0x1111111111111111111111111111111111111111',
    ethereum_approval: '0x2222222222222222222222222222222222222222',
    // base_approval, avalanche_approval, optimism_approval are undefined
  };

  describe('create', () => {
    it('Successfully creates a PermissionApproval', async () => {
      const createdApproval = await PermissionApprovalTable.create(defaultPermissionApproval1);
      expect(createdApproval).toEqual(expect.objectContaining(defaultPermissionApproval1));
    });

    it('Successfully creates a PermissionApproval with partial data', async () => {
      const createdApproval = await PermissionApprovalTable.create(partialPermissionApproval);
      expect(createdApproval).toEqual(expect.objectContaining({
        ...partialPermissionApproval,
        base_approval: null,
        avalanche_approval: null,
        optimism_approval: null,
      }));
    });
  });

  describe('findBySuborgId', () => {
    it('Successfully finds a PermissionApproval by suborg_id', async () => {
      await PermissionApprovalTable.create(defaultPermissionApproval1);
      await PermissionApprovalTable.create(defaultPermissionApproval2);

      const approval = await PermissionApprovalTable.findBySuborgId(
        defaultPermissionApproval1.suborg_id,
      );

      expect(approval).toEqual(expect.objectContaining(defaultPermissionApproval1));
    });

    it('Returns undefined when PermissionApproval not found by suborg_id', async () => {
      await PermissionApprovalTable.create(defaultPermissionApproval1);

      const approval = await PermissionApprovalTable.findBySuborgId(
        'nonexistent-suborg',
      );

      expect(approval).toBeUndefined();
    });
  });

  describe('getArbitrumApprovalForSuborg', () => {
    it('Successfully gets arbitrum approval for existing suborg', async () => {
      await PermissionApprovalTable.create(defaultPermissionApproval1);

      const approval = await PermissionApprovalTable.getArbitrumApprovalForSuborg(
        defaultPermissionApproval1.suborg_id,
      );

      expect(approval).toBe(defaultPermissionApproval1.arbitrum_approval);
    });

    it('Returns undefined when suborg does not exist', async () => {
      const approval = await PermissionApprovalTable.getArbitrumApprovalForSuborg('nonexistent-suborg');

      expect(approval).toBeUndefined();
    });

    it('Returns undefined when arbitrum_approval is null', async () => {
      const approvalWithNullArbitrum = {
        suborg_id: 'suborg-null-arbitrum',
        base_approval: '0x1234567890abcdef1234567890abcdef12345678',
      };
      await PermissionApprovalTable.create(approvalWithNullArbitrum);

      const approval = await PermissionApprovalTable.getArbitrumApprovalForSuborg(
        approvalWithNullArbitrum.suborg_id,
      );

      expect(approval).toBeNull();
    });
  });

  describe('getBaseApprovalForSuborg', () => {
    it('Successfully gets base approval for existing suborg', async () => {
      await PermissionApprovalTable.create(defaultPermissionApproval1);

      const approval = await PermissionApprovalTable.getBaseApprovalForSuborg(
        defaultPermissionApproval1.suborg_id,
      );

      expect(approval).toBe(defaultPermissionApproval1.base_approval);
    });

    it('Returns undefined when suborg does not exist', async () => {
      const approval = await PermissionApprovalTable.getBaseApprovalForSuborg('nonexistent-suborg');

      expect(approval).toBeUndefined();
    });

    it('Returns undefined when base_approval is null', async () => {
      const approvalWithNullBase = {
        suborg_id: 'suborg-null-base',
        arbitrum_approval: '0x1234567890abcdef1234567890abcdef12345678',
      };
      await PermissionApprovalTable.create(approvalWithNullBase);

      const approval = await PermissionApprovalTable.getBaseApprovalForSuborg(
        approvalWithNullBase.suborg_id,
      );

      expect(approval).toBeNull();
    });
  });

  describe('getAvalancheApprovalForSuborg', () => {
    it('Successfully gets avalanche approval for existing suborg', async () => {
      await PermissionApprovalTable.create(defaultPermissionApproval1);

      const approval = await PermissionApprovalTable.getAvalancheApprovalForSuborg(
        defaultPermissionApproval1.suborg_id,
      );

      expect(approval).toBe(defaultPermissionApproval1.avalanche_approval);
    });

    it('Returns undefined when suborg does not exist', async () => {
      const approval = await PermissionApprovalTable.getAvalancheApprovalForSuborg('nonexistent-suborg');

      expect(approval).toBeUndefined();
    });

    it('Returns undefined when avalanche_approval is null', async () => {
      const approvalWithNullAvalanche = {
        suborg_id: 'suborg-null-avalanche',
        arbitrum_approval: '0x1234567890abcdef1234567890abcdef12345678',
      };
      await PermissionApprovalTable.create(approvalWithNullAvalanche);

      const approval = await PermissionApprovalTable.getAvalancheApprovalForSuborg(
        approvalWithNullAvalanche.suborg_id,
      );

      expect(approval).toBeNull();
    });
  });

  describe('getOptimismApprovalForSuborg', () => {
    it('Successfully gets optimism approval for existing suborg', async () => {
      await PermissionApprovalTable.create(defaultPermissionApproval1);

      const approval = await PermissionApprovalTable.getOptimismApprovalForSuborg(
        defaultPermissionApproval1.suborg_id,
      );

      expect(approval).toBe(defaultPermissionApproval1.optimism_approval);
    });

    it('Returns undefined when suborg does not exist', async () => {
      const approval = await PermissionApprovalTable.getOptimismApprovalForSuborg('nonexistent-suborg');

      expect(approval).toBeUndefined();
    });

    it('Returns undefined when optimism_approval is null', async () => {
      const approvalWithNullOptimism = {
        suborg_id: 'suborg-null-optimism',
        arbitrum_approval: '0x1234567890abcdef1234567890abcdef12345678',
      };
      await PermissionApprovalTable.create(approvalWithNullOptimism);

      const approval = await PermissionApprovalTable.getOptimismApprovalForSuborg(
        approvalWithNullOptimism.suborg_id,
      );

      expect(approval).toBeNull();
    });
  });

  describe('getEthereumApprovalForSuborg', () => {
    it('Successfully gets ethereum approval for existing suborg', async () => {
      await PermissionApprovalTable.create(defaultPermissionApproval1);

      const approval = await PermissionApprovalTable.getEthereumApprovalForSuborg(
        defaultPermissionApproval1.suborg_id,
      );

      expect(approval).toBe(defaultPermissionApproval1.ethereum_approval);
    });

    it('Returns undefined when suborg does not exist', async () => {
      const approval = await PermissionApprovalTable.getEthereumApprovalForSuborg('nonexistent-suborg');

      expect(approval).toBeUndefined();
    });

    it('Returns undefined when ethereum_approval is null', async () => {
      const approvalWithNullEthereum = {
        suborg_id: 'suborg-null-ethereum',
        arbitrum_approval: '0x1234567890abcdef1234567890abcdef12345678',
      };
      await PermissionApprovalTable.create(approvalWithNullEthereum);

      const approval = await PermissionApprovalTable.getEthereumApprovalForSuborg(
        approvalWithNullEthereum.suborg_id,
      );

      expect(approval).toBeNull();
    });
  });

  describe('update', () => {
    it('Successfully updates a PermissionApproval', async () => {
      await PermissionApprovalTable.create(defaultPermissionApproval1);

      const updatedFields = {
        suborg_id: defaultPermissionApproval1.suborg_id,
        arbitrum_approval: '0xnewArbitrumApproval123456789012345678901',
        base_approval: '0xnewBaseApproval1234567890123456789012345',
      };

      const updatedApproval = await PermissionApprovalTable.update(updatedFields);

      expect(updatedApproval).toEqual(expect.objectContaining({
        suborg_id: defaultPermissionApproval1.suborg_id,
        arbitrum_approval: updatedFields.arbitrum_approval,
        base_approval: updatedFields.base_approval,
        avalanche_approval: defaultPermissionApproval1.avalanche_approval,
        optimism_approval: defaultPermissionApproval1.optimism_approval,
        ethereum_approval: defaultPermissionApproval1.ethereum_approval,
      }));
    });

    it('Returns undefined when trying to update non-existent PermissionApproval', async () => {
      const updatedFields = {
        suborg_id: 'nonexistent-suborg',
        arbitrum_approval: '0xnewApproval123456789012345678901234567',
      };

      const result = await PermissionApprovalTable.update(updatedFields);

      expect(result).toBeUndefined();
    });
  });

  describe('Edge cases', () => {
    it('Handles empty string suborg_id searches', async () => {
      await PermissionApprovalTable.create(defaultPermissionApproval1);

      const approval = await PermissionApprovalTable.findBySuborgId('');
      const arbitrumApproval = await PermissionApprovalTable.getArbitrumApprovalForSuborg('');
      const baseApproval = await PermissionApprovalTable.getBaseApprovalForSuborg('');
      const avalancheApproval = await PermissionApprovalTable.getAvalancheApprovalForSuborg('');
      const optimismApproval = await PermissionApprovalTable.getOptimismApprovalForSuborg('');
      const ethereumApproval = await PermissionApprovalTable.getEthereumApprovalForSuborg('');

      expect(approval).toBeUndefined();
      expect(arbitrumApproval).toBeUndefined();
      expect(baseApproval).toBeUndefined();
      expect(avalancheApproval).toBeUndefined();
      expect(optimismApproval).toBeUndefined();
      expect(ethereumApproval).toBeUndefined();
    });

    it('Handles special characters in suborg_id', async () => {
      const specialSuborgApproval = {
        suborg_id: 'suborg-with-special-chars-!@#$%^&*()',
        arbitrum_approval: '0x1234567890abcdef1234567890abcdef12345678',
      };

      await PermissionApprovalTable.create(specialSuborgApproval);

      const approval = await PermissionApprovalTable.findBySuborgId(
        specialSuborgApproval.suborg_id,
      );
      const arbitrumApproval = await PermissionApprovalTable.getArbitrumApprovalForSuborg(
        specialSuborgApproval.suborg_id,
      );

      expect(approval).toEqual(expect.objectContaining(specialSuborgApproval));
      expect(arbitrumApproval).toBe(specialSuborgApproval.arbitrum_approval);
    });

    it('Handles very long approval addresses', async () => {
      const longApprovalAddress = `0x${'a'.repeat(64)}`; // Very long hex address
      const longApprovalData = {
        suborg_id: 'suborg-long-approval',
        arbitrum_approval: longApprovalAddress,
      };

      await PermissionApprovalTable.create(longApprovalData);

      const approval = await PermissionApprovalTable.getArbitrumApprovalForSuborg(
        longApprovalData.suborg_id,
      );

      expect(approval).toBe(longApprovalAddress);
    });

    it('Handles multiple chain approvals for same suborg', async () => {
      await PermissionApprovalTable.create(defaultPermissionApproval1);

      // Get all approvals for the same suborg
      const arbitrumApproval = await PermissionApprovalTable.getArbitrumApprovalForSuborg(
        defaultPermissionApproval1.suborg_id,
      );
      const baseApproval = await PermissionApprovalTable.getBaseApprovalForSuborg(
        defaultPermissionApproval1.suborg_id,
      );
      const avalancheApproval = await PermissionApprovalTable.getAvalancheApprovalForSuborg(
        defaultPermissionApproval1.suborg_id,
      );
      const optimismApproval = await PermissionApprovalTable.getOptimismApprovalForSuborg(
        defaultPermissionApproval1.suborg_id,
      );
      const ethereumApproval = await PermissionApprovalTable.getEthereumApprovalForSuborg(
        defaultPermissionApproval1.suborg_id,
      );

      expect(arbitrumApproval).toBe(defaultPermissionApproval1.arbitrum_approval);
      expect(baseApproval).toBe(defaultPermissionApproval1.base_approval);
      expect(avalancheApproval).toBe(defaultPermissionApproval1.avalanche_approval);
      expect(optimismApproval).toBe(defaultPermissionApproval1.optimism_approval);
      expect(ethereumApproval).toBe(defaultPermissionApproval1.ethereum_approval);
    });
  });
});
