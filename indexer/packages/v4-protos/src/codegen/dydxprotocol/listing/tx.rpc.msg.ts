import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { MsgSetMarketsHardCap, MsgSetMarketsHardCapResponse, MsgCreateMarketPermissionless, MsgCreateMarketPermissionlessResponse, MsgSetListingVaultDepositParams, MsgSetListingVaultDepositParamsResponse, MsgUpgradeIsolatedPerpetualToCross, MsgUpgradeIsolatedPerpetualToCrossResponse } from "./tx";
/** Msg defines the Msg service. */

export interface Msg {
  /** SetMarketsHardCap sets a hard cap on the number of markets listed */
  setMarketsHardCap(request: MsgSetMarketsHardCap): Promise<MsgSetMarketsHardCapResponse>;
  /** CreateMarketPermissionless creates a new market without going through x/gov */

  createMarketPermissionless(request: MsgCreateMarketPermissionless): Promise<MsgCreateMarketPermissionlessResponse>;
  /** SetListingVaultDepositParams sets PML megavault deposit params */

  setListingVaultDepositParams(request: MsgSetListingVaultDepositParams): Promise<MsgSetListingVaultDepositParamsResponse>;
  /**
   * UpgradeIsolatedPerpetualToCross upgrades a perpetual from isolated to cross
   * margin
   */

  upgradeIsolatedPerpetualToCross(request: MsgUpgradeIsolatedPerpetualToCross): Promise<MsgUpgradeIsolatedPerpetualToCrossResponse>;
}
export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.setMarketsHardCap = this.setMarketsHardCap.bind(this);
    this.createMarketPermissionless = this.createMarketPermissionless.bind(this);
    this.setListingVaultDepositParams = this.setListingVaultDepositParams.bind(this);
    this.upgradeIsolatedPerpetualToCross = this.upgradeIsolatedPerpetualToCross.bind(this);
  }

  setMarketsHardCap(request: MsgSetMarketsHardCap): Promise<MsgSetMarketsHardCapResponse> {
    const data = MsgSetMarketsHardCap.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.listing.Msg", "SetMarketsHardCap", data);
    return promise.then(data => MsgSetMarketsHardCapResponse.decode(new _m0.Reader(data)));
  }

  createMarketPermissionless(request: MsgCreateMarketPermissionless): Promise<MsgCreateMarketPermissionlessResponse> {
    const data = MsgCreateMarketPermissionless.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.listing.Msg", "CreateMarketPermissionless", data);
    return promise.then(data => MsgCreateMarketPermissionlessResponse.decode(new _m0.Reader(data)));
  }

  setListingVaultDepositParams(request: MsgSetListingVaultDepositParams): Promise<MsgSetListingVaultDepositParamsResponse> {
    const data = MsgSetListingVaultDepositParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.listing.Msg", "SetListingVaultDepositParams", data);
    return promise.then(data => MsgSetListingVaultDepositParamsResponse.decode(new _m0.Reader(data)));
  }

  upgradeIsolatedPerpetualToCross(request: MsgUpgradeIsolatedPerpetualToCross): Promise<MsgUpgradeIsolatedPerpetualToCrossResponse> {
    const data = MsgUpgradeIsolatedPerpetualToCross.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.listing.Msg", "UpgradeIsolatedPerpetualToCross", data);
    return promise.then(data => MsgUpgradeIsolatedPerpetualToCrossResponse.decode(new _m0.Reader(data)));
  }

}