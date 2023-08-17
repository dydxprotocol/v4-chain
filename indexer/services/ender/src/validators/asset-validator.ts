import { marketRefresher } from '@dydxprotocol-indexer/postgres';
import { AssetCreateEventV1, IndexerTendermintEvent } from '@dydxprotocol-indexer/v4-protos';

import { AssetCreationHandler } from '../handlers/asset-handler';
import { Handler } from '../handlers/handler';
import { Validator } from './validator';

export class AssetValidator extends Validator<AssetCreateEventV1> {
  public validate(): void {
    if (this.event.hasMarket) {
      marketRefresher.getMarketFromId(
        this.event.marketId,
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
  ): Handler<AssetCreateEventV1>[] {
    const handler: Handler<AssetCreateEventV1> = new AssetCreationHandler(
      this.block,
      indexerTendermintEvent,
      txId,
      this.event,
    );

    return [handler];
  }
}
