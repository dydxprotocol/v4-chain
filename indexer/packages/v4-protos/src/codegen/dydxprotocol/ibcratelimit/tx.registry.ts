import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import { MsgSetLimitParams, MsgDeleteLimitParams } from "./tx";
export const registry: ReadonlyArray<[string, GeneratedType]> = [["/dydxprotocol.ibcratelimit.MsgSetLimitParams", MsgSetLimitParams], ["/dydxprotocol.ibcratelimit.MsgDeleteLimitParams", MsgDeleteLimitParams]];
export const load = (protoRegistry: Registry) => {
  registry.forEach(([typeUrl, mod]) => {
    protoRegistry.register(typeUrl, mod);
  });
};
export const MessageComposer = {
  encoded: {
    setLimitParams(value: MsgSetLimitParams) {
      return {
        typeUrl: "/dydxprotocol.ibcratelimit.MsgSetLimitParams",
        value: MsgSetLimitParams.encode(value).finish()
      };
    },
    deleteLimitParams(value: MsgDeleteLimitParams) {
      return {
        typeUrl: "/dydxprotocol.ibcratelimit.MsgDeleteLimitParams",
        value: MsgDeleteLimitParams.encode(value).finish()
      };
    }
  },
  withTypeUrl: {
    setLimitParams(value: MsgSetLimitParams) {
      return {
        typeUrl: "/dydxprotocol.ibcratelimit.MsgSetLimitParams",
        value
      };
    },
    deleteLimitParams(value: MsgDeleteLimitParams) {
      return {
        typeUrl: "/dydxprotocol.ibcratelimit.MsgDeleteLimitParams",
        value
      };
    }
  },
  fromPartial: {
    setLimitParams(value: MsgSetLimitParams) {
      return {
        typeUrl: "/dydxprotocol.ibcratelimit.MsgSetLimitParams",
        value: MsgSetLimitParams.fromPartial(value)
      };
    },
    deleteLimitParams(value: MsgDeleteLimitParams) {
      return {
        typeUrl: "/dydxprotocol.ibcratelimit.MsgDeleteLimitParams",
        value: MsgDeleteLimitParams.fromPartial(value)
      };
    }
  }
};