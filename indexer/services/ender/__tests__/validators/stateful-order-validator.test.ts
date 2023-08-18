import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  IndexerOrder_Side,
  StatefulOrderEventV1,
  OrderRemovalReason,
} from '@dydxprotocol-indexer/v4-protos';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { StatefulOrderValidator } from '../../src/validators/stateful-order-validator';
import {
  defaultHeight,
  defaultMakerOrder,
  defaultOrderId,
  defaultStatefulOrderPlacementEvent,
  defaultStatefulOrderRemovalEvent,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import { binaryToBase64String, createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import { expectDidntLogError, expectLoggedParseMessageError } from '../helpers/validator-helpers';

describe('stateful-order-validator', () => {
  beforeEach(() => {
    jest.spyOn(logger, 'error');
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('validate', () => {
    it.each([
      ['stateful order placement', defaultStatefulOrderPlacementEvent],
      ['stateful order removal', defaultStatefulOrderRemovalEvent],
    ])('does not throw error on valid %s', (_message: string, event: StatefulOrderEventV1) => {
      const validator: StatefulOrderValidator = new StatefulOrderValidator(
        event,
        createBlock(event),
      );

      validator.validate();
      expectDidntLogError();
    });

    it.each([
      // Base Errors
      [
        'does not contain any event',
        {},
        'One of orderPlace or orderRemoval must be defined in StatefulOrderEvent',
      ],

      // Order Placement Validations
      [
        'does not contain orderId',
        {
          orderPlace: {
            order: { ...defaultMakerOrder, orderId: undefined },
          },
        },
        'StatefulOrderEvent placement: Order must contain an orderId',
      ],
      [
        'does not contain a subaccountId',
        {
          orderPlace: {
            order: {
              ...defaultMakerOrder,
              orderId: { ...defaultOrderId, subaccountId: undefined },
            },
          },
        },
        'StatefulOrderEvent placement: OrderId must contain a subaccountId',
      ],
      [
        'does not contain a specified order side',
        {
          orderPlace: {
            order: {
              ...defaultMakerOrder,
              side: IndexerOrder_Side.SIDE_UNSPECIFIED,
            },
          },
        },
        'StatefulOrderEvent placement:  Order must specify an order side',
      ],
      [
        'does not contain a defined goodTilOneof',
        {
          orderPlace: {
            order: {
              ...defaultMakerOrder,
              goodTilBlock: undefined,
              goodTilBlockTime: undefined,
            },
          },
        },
        'StatefulOrderEvent placement: Order must contain a defined goodTilOneof',
      ],
      [
        'does not contain a defined goodTilBlockTime',
        {
          orderPlace: {
            order: defaultMakerOrder,
          },
        },
        'StatefulOrderEvent placement: order must have goodTilBlockTime',
      ],

      // Order Removal Validations
      [
        'Stateful order removal does not contain orderId',
        {
          orderRemoval: {
            removedOrderId: undefined,
            reason: OrderRemovalReason.ORDER_REMOVAL_REASON_REPLACED,
          },
        },
        'StatefulOrderEvent removal must contain an orderId',
      ],
      [
        'Stateful order removal contains invalid reason',
        {
          orderRemoval: {
            removedOrderId: defaultOrderId,
            reason: OrderRemovalReason.ORDER_REMOVAL_REASON_UNSPECIFIED,
          },
        },
        'StatefulOrderEvent removal must contain a valid reason',
      ],
    ])('throws error if event %s', (
      _message: string,
      event: StatefulOrderEventV1,
      message: string,
    ) => {
      const validator: StatefulOrderValidator = new StatefulOrderValidator(
        event,
        createBlock(event),
      );

      expect(() => validator.validate()).toThrow(new ParseMessageError(message));
      expectLoggedParseMessageError(
        StatefulOrderValidator.name,
        message,
        { event },
      );
    });
  });
});

function createBlock(
  statefulOrderEvent: StatefulOrderEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.STATEFUL_ORDER,
    binaryToBase64String(
      StatefulOrderEventV1.encode(statefulOrderEvent).finish(),
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
