import { IndexerTendermintEvent, MarketEventV1 } from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';

import { Handler, HandlerInitializer } from '../handlers/handler';
import { MarketCreateHandler } from '../handlers/markets/market-create-handler';
import { MarketModifyHandler } from '../handlers/markets/market-modify-handler';
import { MarketPriceUpdateHandler } from '../handlers/markets/market-price-update-handler';
import { MarketCreateEventMessage, MarketModifyEventMessage, MarketPriceUpdateEventMessage } from '../lib/types';
import { Validator } from './validator';

export class MarketValidator extends Validator<MarketEventV1> {
  public validate(): void {
    if (
      this.event.marketCreate === undefined &&
      this.event.marketModify === undefined &&
      this.event.priceUpdate === undefined
    ) {
      return this.logAndThrowParseMessageError(
        'One of marketCreate, marketModify, or priceUpdate must be defined in MarketEvent',
        { event: this.event },
      );
    }
    if (this.event.marketCreate !== undefined) {
      this.validateMarketCreate();
    } else if (this.event.marketModify !== undefined) {
      this.validateMarketModify();
    } else { // priceUpdate
      this.validatePriceUpdate();
    }
  }

  private validateMarketCreate(): void {
    const marketCreate: MarketCreateEventMessage = this.event as MarketCreateEventMessage;
    if (marketCreate.marketCreate.base === undefined) {
      return this.logAndThrowParseMessageError(
        'Invalid MarketCreate, base field is undefined',
        { event: this.event },
      );
    }

    if (marketCreate.marketCreate.base!.pair === '') {
      return this.logAndThrowParseMessageError(
        'Invalid MarketCreate, pair is empty',
        { event: this.event },
      );
    }

    if (marketCreate.marketCreate.base!.minPriceChangePpm === 0) {
      return this.logAndThrowParseMessageError(
        'Invalid MarketCreate, minPriceChangePpm is 0',
        { event: this.event },
      );
    }
  }

  private validateMarketModify(): void {
    const marketModify: MarketModifyEventMessage = this.event as MarketModifyEventMessage;
    if (marketModify.marketModify.base === undefined) {
      return this.logAndThrowParseMessageError(
        'Invalid MarketModify, base field is undefined',
        { event: this.event },
      );
    }

    if (marketModify.marketModify.base!.pair === '') {
      return this.logAndThrowParseMessageError(
        'Invalid MarketModify, pair is empty',
        { event: this.event },
      );
    }

    if (marketModify.marketModify.base!.minPriceChangePpm === 0) {
      return this.logAndThrowParseMessageError(
        'Invalid MarketModify, minPriceChangePpm is 0',
        { event: this.event },
      );
    }
  }

  private validatePriceUpdate(): void {
    const marketPriceUpdate:
    MarketPriceUpdateEventMessage = this.event as MarketPriceUpdateEventMessage;
    if (
      marketPriceUpdate.priceUpdate.priceWithExponent <= Long.fromValue(0)
    ) {
      return this.logAndThrowParseMessageError(
        'Invalid MarketPriceUpdate, priceWithExponent must be > 0',
        { event: this.event },
      );
    }
  }

  public getHandlerInitializer() : HandlerInitializer | undefined {
    if (this.event.marketCreate !== undefined) {
      return MarketCreateHandler;
    } else if (this.event.marketModify !== undefined) {
      return MarketModifyHandler;
    } else if (this.event.priceUpdate !== undefined) {
      return MarketPriceUpdateHandler;
    }
    return undefined;
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    __: string,
  ): Handler<MarketEventV1>[] {
    const Initializer:
    HandlerInitializer | undefined = this.getHandlerInitializer();
    if (Initializer === undefined) {
      this.logAndThrowParseMessageError(
        'Cannot process event',
        { event: this.event },
      );
    }
    // @ts-ignore
    const handler: Handler<MarketEvent> = new Initializer(
      this.block,
      this.blockEventIndex,
      indexerTendermintEvent,
      txId,
      this.event,
    );

    return [handler];
  }
}
