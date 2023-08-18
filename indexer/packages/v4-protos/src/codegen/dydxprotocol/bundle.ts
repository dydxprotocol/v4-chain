import * as _5 from "./assets/asset";
import * as _6 from "./assets/genesis";
import * as _7 from "./assets/query";
import * as _8 from "./assets/tx";
import * as _9 from "./blocktime/blocktime";
import * as _10 from "./blocktime/genesis";
import * as _11 from "./blocktime/params";
import * as _12 from "./blocktime/query";
import * as _13 from "./blocktime/tx";
import * as _14 from "./bridge/bridge_event_info";
import * as _15 from "./bridge/bridge_event";
import * as _16 from "./bridge/genesis";
import * as _17 from "./bridge/params";
import * as _18 from "./bridge/query";
import * as _19 from "./bridge/tx";
import * as _20 from "./clob/block_rate_limit_config";
import * as _21 from "./clob/clob_pair";
import * as _22 from "./clob/genesis";
import * as _23 from "./clob/liquidations_config";
import * as _24 from "./clob/liquidations";
import * as _25 from "./clob/matches";
import * as _26 from "./clob/operation";
import * as _27 from "./clob/order_removals";
import * as _28 from "./clob/order";
import * as _29 from "./clob/process_proposer_matches_events";
import * as _30 from "./clob/query";
import * as _31 from "./clob/tx";
import * as _32 from "./daemons/bridge/bridge";
import * as _33 from "./daemons/liquidation/liquidation";
import * as _34 from "./daemons/pricefeed/price_feed";
import * as _35 from "./epochs/epoch_info";
import * as _36 from "./epochs/genesis";
import * as _37 from "./epochs/query";
import * as _38 from "./feetiers/genesis";
import * as _39 from "./feetiers/params";
import * as _40 from "./feetiers/query";
import * as _41 from "./feetiers/tx";
import * as _42 from "./indexer/events/events";
import * as _43 from "./indexer/indexer_manager/event";
import * as _44 from "./indexer/off_chain_updates/off_chain_updates";
import * as _45 from "./indexer/protocol/v1/clob";
import * as _46 from "./indexer/protocol/v1/subaccount";
import * as _47 from "./indexer/redis/redis_order";
import * as _48 from "./indexer/shared/removal_reason";
import * as _49 from "./indexer/socks/messages";
import * as _50 from "./perpetuals/genesis";
import * as _51 from "./perpetuals/params";
import * as _52 from "./perpetuals/perpetual";
import * as _53 from "./perpetuals/query";
import * as _54 from "./perpetuals/tx";
import * as _55 from "./prices/genesis";
import * as _56 from "./prices/market_param";
import * as _57 from "./prices/market_price";
import * as _58 from "./prices/query";
import * as _59 from "./prices/tx";
import * as _60 from "./rewards/genesis";
import * as _61 from "./rewards/params";
import * as _62 from "./rewards/query";
import * as _63 from "./rewards/reward_share";
import * as _64 from "./rewards/tx";
import * as _65 from "./sending/genesis";
import * as _66 from "./sending/query";
import * as _67 from "./sending/transfer";
import * as _68 from "./sending/tx";
import * as _69 from "./stats/genesis";
import * as _70 from "./stats/params";
import * as _71 from "./stats/query";
import * as _72 from "./stats/stats";
import * as _73 from "./stats/tx";
import * as _74 from "./subaccounts/asset_position";
import * as _75 from "./subaccounts/genesis";
import * as _76 from "./subaccounts/perpetual_position";
import * as _77 from "./subaccounts/query";
import * as _78 from "./subaccounts/subaccount";
import * as _79 from "./vesting/genesis";
import * as _80 from "./vesting/query";
import * as _81 from "./vesting/tx";
import * as _82 from "./vesting/vesting_entry";
import * as _89 from "./blocktime/query.lcd";
import * as _90 from "./bridge/query.lcd";
import * as _91 from "./clob/query.lcd";
import * as _92 from "./epochs/query.lcd";
import * as _93 from "./feetiers/query.lcd";
import * as _94 from "./perpetuals/query.lcd";
import * as _95 from "./prices/query.lcd";
import * as _96 from "./rewards/query.lcd";
import * as _97 from "./stats/query.lcd";
import * as _98 from "./subaccounts/query.lcd";
import * as _99 from "./vesting/query.lcd";
import * as _100 from "./assets/query.rpc.Query";
import * as _101 from "./blocktime/query.rpc.Query";
import * as _102 from "./bridge/query.rpc.Query";
import * as _103 from "./clob/query.rpc.Query";
import * as _104 from "./epochs/query.rpc.Query";
import * as _105 from "./feetiers/query.rpc.Query";
import * as _106 from "./perpetuals/query.rpc.Query";
import * as _107 from "./prices/query.rpc.Query";
import * as _108 from "./rewards/query.rpc.Query";
import * as _109 from "./sending/query.rpc.Query";
import * as _110 from "./stats/query.rpc.Query";
import * as _111 from "./subaccounts/query.rpc.Query";
import * as _112 from "./vesting/query.rpc.Query";
import * as _113 from "./blocktime/tx.rpc.msg";
import * as _114 from "./bridge/tx.rpc.msg";
import * as _115 from "./clob/tx.rpc.msg";
import * as _116 from "./feetiers/tx.rpc.msg";
import * as _117 from "./perpetuals/tx.rpc.msg";
import * as _118 from "./prices/tx.rpc.msg";
import * as _119 from "./rewards/tx.rpc.msg";
import * as _120 from "./sending/tx.rpc.msg";
import * as _121 from "./stats/tx.rpc.msg";
import * as _122 from "./vesting/tx.rpc.msg";
import * as _123 from "./lcd";
import * as _124 from "./rpc.query";
import * as _125 from "./rpc.tx";
export namespace dydxprotocol {
  export const assets = { ..._5,
    ..._6,
    ..._7,
    ..._8,
    ..._100
  };
  export const blocktime = { ..._9,
    ..._10,
    ..._11,
    ..._12,
    ..._13,
    ..._89,
    ..._101,
    ..._113
  };
  export const bridge = { ..._14,
    ..._15,
    ..._16,
    ..._17,
    ..._18,
    ..._19,
    ..._90,
    ..._102,
    ..._114
  };
  export const clob = { ..._20,
    ..._21,
    ..._22,
    ..._23,
    ..._24,
    ..._25,
    ..._26,
    ..._27,
    ..._28,
    ..._29,
    ..._30,
    ..._31,
    ..._91,
    ..._103,
    ..._115
  };
  export namespace daemons {
    export const bridge = { ..._32
    };
    export const liquidation = { ..._33
    };
    export const pricefeed = { ..._34
    };
  }
  export const epochs = { ..._35,
    ..._36,
    ..._37,
    ..._92,
    ..._104
  };
  export const feetiers = { ..._38,
    ..._39,
    ..._40,
    ..._41,
    ..._93,
    ..._105,
    ..._116
  };
  export namespace indexer {
    export const events = { ..._42
    };
    export const indexer_manager = { ..._43
    };
    export const off_chain_updates = { ..._44
    };
    export namespace protocol {
      export const v1 = { ..._45,
        ..._46
      };
    }
    export const redis = { ..._47
    };
    export const shared = { ..._48
    };
    export const socks = { ..._49
    };
  }
  export const perpetuals = { ..._50,
    ..._51,
    ..._52,
    ..._53,
    ..._54,
    ..._94,
    ..._106,
    ..._117
  };
  export const prices = { ..._55,
    ..._56,
    ..._57,
    ..._58,
    ..._59,
    ..._95,
    ..._107,
    ..._118
  };
  export const rewards = { ..._60,
    ..._61,
    ..._62,
    ..._63,
    ..._64,
    ..._96,
    ..._108,
    ..._119
  };
  export const sending = { ..._65,
    ..._66,
    ..._67,
    ..._68,
    ..._109,
    ..._120
  };
  export const stats = { ..._69,
    ..._70,
    ..._71,
    ..._72,
    ..._73,
    ..._97,
    ..._110,
    ..._121
  };
  export const subaccounts = { ..._74,
    ..._75,
    ..._76,
    ..._77,
    ..._78,
    ..._98,
    ..._111
  };
  export const vesting = { ..._79,
    ..._80,
    ..._81,
    ..._82,
    ..._99,
    ..._112,
    ..._122
  };
  export const ClientFactory = { ..._123,
    ..._124,
    ..._125
  };
}