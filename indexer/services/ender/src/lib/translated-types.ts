import { IndexerAssetPosition, IndexerPerpetualPosition, IndexerSubaccountId } from '@dydxprotocol-indexer/v4-protos';

export interface SubaccountUpdate {
  subaccountId?: IndexerSubaccountId;
  updatedPerpetualPositions: IndexerPerpetualPosition[];
  updatedAssetPositions: IndexerAssetPosition[];
}
