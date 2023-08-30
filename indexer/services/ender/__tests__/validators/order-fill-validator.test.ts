import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  OrderFillEventV1,
  IndexerOrder_Side,
} from '@dydxprotocol-indexer/v4-protos';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { OrderFillValidator } from '../../src/validators/order-fill-validator';
import {
  defaultHeight,
  defaultLiquidationEvent,
  defaultLiquidationOrder,
  defaultMakerOrder,
  defaultOrderEvent,
  defaultOrderId,
  defaultTakerOrder,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import { binaryToBase64String, createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import { expectDidntLogError, expectLoggedParseMessageError } from '../helpers/validator-helpers';

describe('order-fill-validator', () => {
  beforeEach(() => {
    jest.spyOn(logger, 'error');
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('validate', () => {
    it.each([
      ['order event', defaultOrderEvent],
      ['liquidation event', defaultLiquidationEvent],
    ])('does not throw error on valid %s', (_message: string, event: OrderFillEventV1) => {
      const validator: OrderFillValidator = new OrderFillValidator(
        event,
        createBlock(event),
      );

      validator.validate();
      expectDidntLogError();
    });

    it.each([
      // Base Errors
      [
        'does not contain maker order',
        {
          ...defaultOrderEvent,
          makerOrder: undefined,
        },
        'OrderFillEvent must contain a maker order',
      ],

      // Maker Order Validations
      [
        'does not contain orderId',
        {
          ...defaultOrderEvent,
          makerOrder: { ...defaultMakerOrder, orderId: undefined },
        },
        'OrderFillEvent must contain a makerOrder: Order must contain an orderId',
      ],
      [
        'does not contain a subaccountId',
        {
          ...defaultOrderEvent,
          makerOrder: {
            ...defaultMakerOrder,
            orderId: {
              ...defaultOrderId,
              subaccountId: undefined,
            },
          },
        },
        'OrderFillEvent must contain a makerOrder: OrderId must contain a subaccountId',
      ],
      [
        'does not contain a specified order side',
        {
          ...defaultOrderEvent,
          makerOrder: { ...defaultMakerOrder, side: IndexerOrder_Side.SIDE_UNSPECIFIED },
        },
        'OrderFillEvent must contain a makerOrder:  Order must specify an order side',
      ],
      [
        'does not contain a defined goodTilOneof',
        {
          ...defaultOrderEvent,
          makerOrder: {
            ...defaultMakerOrder,
            goodTilBlock: undefined,
            goodTilBlockTime: undefined,
          },
        },
        'OrderFillEvent must contain a makerOrder: Order must contain a defined goodTilOneof',
      ],

      // Taker Order Event Validations
      [
        'does not contain orderId',
        {
          ...defaultOrderEvent,
          order: { ...defaultTakerOrder, orderId: undefined },
        } as OrderFillEventV1,
        'OrderFillEvent must contain a takerOrder: Order must contain an orderId',
      ],
      [
        'does not contain a subaccountId',
        {
          ...defaultOrderEvent,
          order: {
            ...defaultTakerOrder,
            orderId: {
              ...defaultOrderId,
              subaccountId: undefined,
            },
          },
        } as OrderFillEventV1,
        'OrderFillEvent must contain a takerOrder: OrderId must contain a subaccountId',
      ],
      [
        'does not contain a specified order side',
        {
          ...defaultOrderEvent,
          order: { ...defaultTakerOrder, side: IndexerOrder_Side.SIDE_UNSPECIFIED },
        } as OrderFillEventV1,
        'OrderFillEvent must contain a takerOrder:  Order must specify an order side',
      ],
      [
        'does not contain a defined goodTilOneof',
        {
          ...defaultOrderEvent,
          order: { ...defaultTakerOrder, goodTilBlock: undefined, goodTilBlockTime: undefined },
        } as OrderFillEventV1,
        'OrderFillEvent must contain a takerOrder: Order must contain a defined goodTilOneof',
      ],

      // Taker Liquidation Event validations
      [
        'does not contain liquidated subaccountId',
        {
          ...defaultLiquidationEvent,
          liquidationOrder: {
            ...defaultLiquidationOrder,
            liquidated: undefined,
          },
        } as OrderFillEventV1,
        'LiquidationOrder must contain a liquidated subaccountId',
      ],
    ])('throws error if event %s', (_message: string, event: OrderFillEventV1, message: string) => {
      const validator: OrderFillValidator = new OrderFillValidator(
        event,
        createBlock(event),
      );

      expect(() => validator.validate()).toThrow(new ParseMessageError(message));
      expectLoggedParseMessageError(
        OrderFillValidator.name,
        message,
        { event },
      );
    });
  });
});

function createBlock(
  orderFillEvent: OrderFillEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.ORDER_FILL,
    binaryToBase64String(
      OrderFillEventV1.encode(orderFillEvent).finish(),
    ),
    0,
    0,
  );

  return createIndexerTendermintBlock(
    defaultHeight,
    defaultTime,
    [event],
    [defaultTxHash],
  );
}
