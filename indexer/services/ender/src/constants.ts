export const MILLIS_IN_NANOS: number = 1_000_000;
export const SECONDS_IN_MILLIS: number = 1_000;

// Type sourced from protocol https://github.com/dydxprotocol/v4/blob/main/lib/constants.go#L6
export const QUOTE_CURRENCY_ATOMIC_RESOLUTION: number = -6;

// ============= PARALLELIZATION ID EVENT TYPES =============
// Used to prepend to parallelization ids to ensure that event types from different handlers are
// processed chronologically

// SubaccountUpdate and OrderFill events for the same subaccount are processed chronologically.
export const SUBACCOUNT_ORDER_FILL_EVENT_TYPE: string = 'subaccount_order_fill';

// StatefulOrder and OrderFill events for the same order are processed chronologically.
export const STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE: string = 'stateful_order_order_fill';
