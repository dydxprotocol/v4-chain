import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  IndexerOrder_Side,
  StatefulOrderEventV1,
  OrderRemovalReason,
  IndexerOrder_ConditionType,
} from '@dydxprotocol-indexer/v4-protos';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { StatefulOrderValidator } from '../../src/validators/stateful-order-validator';
import {
  defaultConditionalOrderPlacementEvent,
  defaultConditionalOrderTriggeredEvent,
  defaultHeight,
  defaultLongTermOrderPlacementEvent,
  defaultMakerOrder,
  defaultOrderId,
  defaultStatefulOrderPlacementEvent,
  defaultStatefulOrderRemovalEvent,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import { expectDidntLogError, expectLoggedParseMessageError } from '../helpers/validator-helpers';
import { ORDER_FLAG_CONDITIONAL, ORDER_FLAG_LONG_TERM, ORDER_FLAG_SHORT_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import Long from 'long';

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
      ['conditional order placement', defaultConditionalOrderPlacementEvent],
      ['conditional order triggered', defaultConditionalOrderTriggeredEvent],
      ['long term order placement', defaultLongTermOrderPlacementEvent],
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
        'One of orderPlace, orderRemoval, conditionalOrderPlacement, ' +
        'conditionalOrderTriggered, longTermOrderPlacement must be defined in StatefulOrderEvent',
      ],

      // TODO(IND-334): Remove tests after deprecating StatefulOrderPlacement events
      // Order Placement Validations
      [
        'does not contain orderId',
        {
          orderPlace: {
            order: { ...defaultMakerOrder, orderId: undefined },
          },
        },
        'StatefulOrderEvent stateful order: Order must contain an orderId',
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
        'StatefulOrderEvent stateful order: OrderId must contain a subaccountId',
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
        'StatefulOrderEvent stateful order:  Order must specify an order side',
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
        'StatefulOrderEvent stateful order: Order must contain a defined goodTilOneof',
      ],
      [
        'does not contain a defined goodTilBlockTime',
        {
          orderPlace: {
            order: defaultMakerOrder,
          },
        },
        'StatefulOrderEvent stateful order: order must have goodTilBlockTime',
      ],

      // Long term Order Placement Validations
      [
        'does not contain orderId',
        {
          longTermOrderPlacement: {
            order: { ...defaultMakerOrder, orderId: undefined },
          },
        },
        'StatefulOrderEvent stateful order: Order must contain an orderId',
      ],
      [
        'does not contain a subaccountId',
        {
          longTermOrderPlacement: {
            order: {
              ...defaultMakerOrder,
              orderId: { ...defaultOrderId, subaccountId: undefined },
            },
          },
        },
        'StatefulOrderEvent stateful order: OrderId must contain a subaccountId',
      ],
      [
        'does not contain a specified order side',
        {
          longTermOrderPlacement: {
            order: {
              ...defaultMakerOrder,
              side: IndexerOrder_Side.SIDE_UNSPECIFIED,
            },
          },
        },
        'StatefulOrderEvent stateful order:  Order must specify an order side',
      ],
      [
        'does not contain a defined goodTilOneof',
        {
          longTermOrderPlacement: {
            order: {
              ...defaultMakerOrder,
              goodTilBlock: undefined,
              goodTilBlockTime: undefined,
            },
          },
        },
        'StatefulOrderEvent stateful order: Order must contain a defined goodTilOneof',
      ],
      [
        'does not contain a defined goodTilBlockTime',
        {
          longTermOrderPlacement: {
            order: defaultMakerOrder,
          },
        },
        'StatefulOrderEvent stateful order: order must have goodTilBlockTime',
      ],
      [
        'does not contain the correct order flag',
        {
          longTermOrderPlacement: {
            order: {
              ...defaultMakerOrder,
              orderId: {
                ...defaultMakerOrder.orderId!,
                orderFlags: ORDER_FLAG_SHORT_TERM,
              },
              goodTilBlockTime: 123,
            },
          },
        },
        `StatefulOrderEvent long term order must have order flag ${ORDER_FLAG_LONG_TERM}`,
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

      // Conditional Order Placement Validations
      [
        'conditional order placement does not contain orderId',
        {
          conditionalOrderPlacement: {
            order: { ...defaultMakerOrder, orderId: undefined },
          },
        },
        'StatefulOrderEvent stateful order: Order must contain an orderId',
      ],
      [
        'conditional order placement does not contain subaccountId',
        {
          conditionalOrderPlacement: {
            order: {
              ...defaultMakerOrder,
              orderId: {
                ...defaultOrderId,
                subaccountId: undefined,
              },
            },
          },
        },
        'StatefulOrderEvent stateful order: OrderId must contain a subaccountId',
      ],
      [
        'conditional order placement does not contain a specified order side',
        {
          conditionalOrderPlacement: {
            order: {
              ...defaultMakerOrder,
              side: IndexerOrder_Side.SIDE_UNSPECIFIED,
            },
          },
        },
        'StatefulOrderEvent stateful order:  Order must specify an order side',
      ],
      [
        'conditional order placement does not contain a defined goodTilOneof',
        {
          conditionalOrderPlacement: {
            order: {
              ...defaultMakerOrder,
              goodTilBlock: undefined,
              goodTilBlockTime: undefined,
            },
          },
        },
        'StatefulOrderEvent stateful order: Order must contain a defined goodTilOneof',
      ],
      [
        'conditional order placement does not contain the correct order flag',
        {
          conditionalOrderPlacement: {
            order: {
              ...defaultMakerOrder,
              orderId: {
                ...defaultMakerOrder.orderId!,
                orderFlags: ORDER_FLAG_SHORT_TERM,
              },
              goodTilBlockTime: 123,
            },
          },
        },
        `StatefulOrderEvent conditional order must have order flag ${ORDER_FLAG_CONDITIONAL}`,
      ],
      [
        'conditional order placement does not contain a trigger subticks greater than zero',
        {
          conditionalOrderPlacement: {
            order: {
              ...defaultMakerOrder,
              orderId: {
                ...defaultMakerOrder.orderId!,
                orderFlags: ORDER_FLAG_CONDITIONAL,
              },
              goodTilBlockTime: 123,
            },
          },
        },
        'StatefulOrderEvent conditional order must have trigger price > 0',
      ],
      [
        'conditional order placement does not contain a valid condition type',
        {
          conditionalOrderPlacement: {
            order: {
              ...defaultMakerOrder,
              orderId: {
                ...defaultMakerOrder.orderId!,
                orderFlags: ORDER_FLAG_CONDITIONAL,
              },
              goodTilBlockTime: 123,
              conditionalOrderTriggerSubticks: Long.fromValue(1000000, true),
              conditionType: IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED,
            },
          },
        },
        'StatefulOrderEvent conditional order must have valid condition type',
      ],

      // Conditional order triggered Validations
      [
        'conditional order triggered does not contain orderId',
        {
          conditionalOrderTriggered: {
            triggeredOrderId: undefined,
          },
        },
        'StatefulOrderEvent conditional order triggered must contain an orderId',
      ],
      [
        'conditional order triggered does not contain the correct order flag',
        {
          conditionalOrderTriggered: {
            triggeredOrderId: {
              ...defaultOrderId,
              orderFlags: ORDER_FLAG_SHORT_TERM,
            },
          },
        },
        `StatefulOrderEvent conditional order triggered must have order flag ${ORDER_FLAG_CONDITIONAL}`,
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
    StatefulOrderEventV1.encode(statefulOrderEvent).finish(),
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
