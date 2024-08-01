import { Model } from 'objection';

import {
  NonNegativeNumericPattern,
} from '../lib/validators';
import {
  CandleResolution,
  IsoString,
} from '../types';

export default class CandleModel extends Model {
  static get tableName() {
    return 'candles';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'startedAt',
        'ticker',
        'resolution',
        'low',
        'high',
        'open',
        'close',
        'baseTokenVolume',
        'trades',
        'usdVolume',
        'startingOpenInterest',
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        startedAt: { type: 'string', format: 'date-time' },
        ticker: { type: 'string' },
        resolution: { type: 'string', enum: [...Object.values(CandleResolution)] },
        low: { type: 'string', pattern: NonNegativeNumericPattern },
        high: { type: 'string', pattern: NonNegativeNumericPattern },
        open: { type: 'string', pattern: NonNegativeNumericPattern },
        close: { type: 'string', pattern: NonNegativeNumericPattern },
        baseTokenVolume: { type: 'string', pattern: NonNegativeNumericPattern },
        usdVolume: { type: 'string', pattern: NonNegativeNumericPattern },
        trades: { type: 'integer' },
        startingOpenInterest: { type: 'string', pattern: NonNegativeNumericPattern },
        orderbookMidPriceOpen: { type: ['string', 'null'], pattern: NonNegativeNumericPattern },
        orderbookMidPriceClose: { type: ['string', 'null'], pattern: NonNegativeNumericPattern },
      },
    };
  }

  id!: string;

  startedAt!: IsoString;

  ticker!: string;

  resolution!: CandleResolution;

  low!: string;

  high!: string;

  open!: string;

  close!: string;

  baseTokenVolume!: string;

  usdVolume!: string;

  trades!: number;

  startingOpenInterest!: string;

  orderbookMidPriceOpen?: string;

  orderbookMidPriceClose?: string;
}
