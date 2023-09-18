import { SubaccountUpdateEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { SubaccountUpdate } from '../lib/translated-types';

export function subaccountUpdateEventV1ToSubaccountUpdate(
  event: SubaccountUpdateEventV1,
): SubaccountUpdate {
  return {
    subaccountId: event.subaccountId,
    updatedPerpetualPositions: event.updatedPerpetualPositions,
    updatedAssetPositions: event.updatedAssetPositions,
  };
}
