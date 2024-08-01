import {
  IndexerTendermintEvent,
  UpdateYieldParamsEventV1,
} from '@dydxprotocol-indexer/v4-protos';

import { YieldParamsHandler } from '../handlers/yield-params-handler';
import { Handler } from '../handlers/handler';
import { Validator } from './validator';

export class YieldParamsValidator extends Validator<UpdateYieldParamsEventV1> {
  public validate(): void {


    if (this.event.assetYieldIndex === undefined || this.event.assetYieldIndex === '') {
        return this.logAndThrowParseMessageError(
          'UpdateYieldParamsEvent must have an assetYieldIndex that is defined and non-empty',
          { event: this.event },
        );
    }

    if (this.event.sdaiPrice === undefined || this.event.sdaiPrice === '') {
        return this.logAndThrowParseMessageError(
          'UpdateYieldParamsEvent must have an sDAIPrice that is defined and non-empty',
          { event: this.event },
        );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    _: string,
  ): Handler<UpdateYieldParamsEventV1>[] {
    return [
        new YieldParamsHandler(
          this.block,
          this.blockEventIndex,
          indexerTendermintEvent,
          txId,
          this.event,
        ),
    ];
  }
}
