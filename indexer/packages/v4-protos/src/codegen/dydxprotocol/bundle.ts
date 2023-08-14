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
import * as _22 from "./clob/equity_tier_limit_config";
import * as _23 from "./clob/genesis";
import * as _24 from "./clob/liquidations_config";
import * as _25 from "./clob/liquidations";
import * as _26 from "./clob/matches";
import * as _27 from "./clob/mev";
import * as _28 from "./clob/operation";
import * as _29 from "./clob/order_removals";
import * as _30 from "./clob/order";
import * as _31 from "./clob/process_proposer_matches_events";
import * as _32 from "./clob/query";
import * as _33 from "./clob/tx";
import * as _34 from "./daemons/bridge/bridge";
import * as _35 from "./daemons/liquidation/liquidation";
import * as _36 from "./daemons/pricefeed/price_feed";
import * as _37 from "./epochs/epoch_info";
import * as _38 from "./epochs/genesis";
import * as _39 from "./epochs/query";
import * as _40 from "./feetiers/genesis";
import * as _41 from "./feetiers/params";
import * as _42 from "./feetiers/query";
import * as _43 from "./feetiers/tx";
import * as _44 from "./indexer/events/events";
import * as _45 from "./indexer/indexer_manager/event";
import * as _46 from "./indexer/off_chain_updates/off_chain_updates";
import * as _47 from "./indexer/protocol/v1/clob";
import * as _48 from "./indexer/protocol/v1/subaccount";
import * as _49 from "./indexer/redis/redis_order";
import * as _50 from "./indexer/shared/removal_reason";
import * as _51 from "./indexer/socks/messages";
import * as _52 from "./perpetuals/genesis";
import * as _53 from "./perpetuals/params";
import * as _54 from "./perpetuals/perpetual";
import * as _55 from "./perpetuals/query";
import * as _56 from "./perpetuals/tx";
import * as _57 from "./prices/genesis";
import * as _58 from "./prices/market_param";
import * as _59 from "./prices/market_price";
import * as _60 from "./prices/query";
import * as _61 from "./prices/tx";
import * as _62 from "./rewards/genesis";
import * as _63 from "./rewards/params";
import * as _64 from "./rewards/query";
import * as _65 from "./rewards/reward_share";
import * as _66 from "./rewards/tx";
import * as _67 from "./sending/genesis";
import * as _68 from "./sending/query";
import * as _69 from "./sending/transfer";
import * as _70 from "./sending/tx";
import * as _71 from "./stats/genesis";
import * as _72 from "./stats/params";
import * as _73 from "./stats/query";
import * as _74 from "./stats/stats";
import * as _75 from "./stats/tx";
import * as _76 from "./subaccounts/asset_position";
import * as _77 from "./subaccounts/genesis";
import * as _78 from "./subaccounts/perpetual_position";
import * as _79 from "./subaccounts/query";
import * as _80 from "./subaccounts/subaccount";
import * as _81 from "./vest/genesis";
import * as _82 from "./vest/query";
import * as _83 from "./vest/tx";
import * as _84 from "./vest/vest_entry";
import * as _91 from "./blocktime/query.lcd";
import * as _92 from "./bridge/query.lcd";
import * as _93 from "./clob/query.lcd";
import * as _94 from "./epochs/query.lcd";
import * as _95 from "./feetiers/query.lcd";
import * as _96 from "./perpetuals/query.lcd";
import * as _97 from "./prices/query.lcd";
import * as _98 from "./rewards/query.lcd";
import * as _99 from "./stats/query.lcd";
import * as _100 from "./subaccounts/query.lcd";
import * as _101 from "./vest/query.lcd";
import * as _102 from "./assets/query.rpc.Query";
import * as _103 from "./blocktime/query.rpc.Query";
import * as _104 from "./bridge/query.rpc.Query";
import * as _105 from "./clob/query.rpc.Query";
import * as _106 from "./epochs/query.rpc.Query";
import * as _107 from "./feetiers/query.rpc.Query";
import * as _108 from "./perpetuals/query.rpc.Query";
import * as _109 from "./prices/query.rpc.Query";
import * as _110 from "./rewards/query.rpc.Query";
import * as _111 from "./sending/query.rpc.Query";
import * as _112 from "./stats/query.rpc.Query";
import * as _113 from "./subaccounts/query.rpc.Query";
import * as _114 from "./vest/query.rpc.Query";
import * as _115 from "./blocktime/tx.rpc.msg";
import * as _116 from "./bridge/tx.rpc.msg";
import * as _117 from "./clob/tx.rpc.msg";
import * as _118 from "./feetiers/tx.rpc.msg";
import * as _119 from "./perpetuals/tx.rpc.msg";
import * as _120 from "./prices/tx.rpc.msg";
import * as _121 from "./rewards/tx.rpc.msg";
import * as _122 from "./sending/tx.rpc.msg";
import * as _123 from "./stats/tx.rpc.msg";
import * as _124 from "./vest/tx.rpc.msg";
import * as _125 from "./lcd";
import * as _126 from "./rpc.query";
import * as _127 from "./rpc.tx";
export namespace dydxprotocol {
  export const assets = { ..._5,
    ..._6,
    ..._7,
    ..._8,
    ..._102
  };
  export const blocktime = { ..._9,
    ..._10,
    ..._11,
    ..._12,
    ..._13,
    ..._91,
    ..._103,
    ..._115
  };
  export const bridge = { ..._14,
    ..._15,
    ..._16,
    ..._17,
    ..._18,
    ..._19,
    ..._92,
    ..._104,
    ..._116
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
    ..._32,
    ..._33,
    ..._93,
    ..._105,
    ..._117
  };
  export namespace daemons {
    export const bridge = { ..._34
    };
    export const liquidation = { ..._35
    };
    export const pricefeed = { ..._36
    };
  }
  export const epochs = { ..._37,
    ..._38,
    ..._39,
    ..._94,
    ..._106
  };
  export const feetiers = { ..._40,
    ..._41,
    ..._42,
    ..._43,
    ..._95,
    ..._107,
    ..._118
  };
  export namespace indexer {
    export const events = { ..._44
    };
    export const indexer_manager = { ..._45
    };
    export const off_chain_updates = { ..._46
    };
    export namespace protocol {
      export const v1 = { ..._47,
        ..._48
      };
    }
    export const redis = { ..._49
    };
    export const shared = { ..._50
    };
    export const socks = { ..._51
    };
  }
  export const perpetuals = { ..._52,
    ..._53,
    ..._54,
    ..._55,
    ..._56,
    ..._96,
    ..._108,
    ..._119
  };
  export const prices = { ..._57,
    ..._58,
    ..._59,
    ..._60,
    ..._61,
    ..._97,
    ..._109,
    ..._120
  };
  export const rewards = { ..._62,
    ..._63,
    ..._64,
    ..._65,
    ..._66,
    ..._98,
    ..._110,
    ..._121
  };
  export const sending = { ..._67,
    ..._68,
    ..._69,
    ..._70,
    ..._111,
    ..._122
  };
  export const stats = { ..._71,
    ..._72,
    ..._73,
    ..._74,
    ..._75,
    ..._99,
    ..._112,
    ..._123
  };
  export const subaccounts = { ..._76,
    ..._77,
    ..._78,
    ..._79,
    ..._80,
    ..._100,
    ..._113
  };
  export const vest = { ..._81,
    ..._82,
    ..._83,
    ..._84,
    ..._101,
    ..._114,
    ..._124
  };
  export const ClientFactory = { ..._125,
    ..._126,
    ..._127
  };
}