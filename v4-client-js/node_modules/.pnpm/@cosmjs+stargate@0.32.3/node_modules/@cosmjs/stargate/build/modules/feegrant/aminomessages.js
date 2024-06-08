"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createFeegrantAminoConverters = void 0;
function createFeegrantAminoConverters() {
    return {
    // For Cosmos SDK < 0.46 the Amino JSON codec was broken on chain and thus inaccessible.
    // Now this can be implemented for 0.46+ chains, see
    // https://github.com/cosmos/cosmjs/issues/1092
    //
    // "/cosmos.feegrant.v1beta1.MsgGrantAllowance": IMPLEMENT_ME,
    // "/cosmos.feegrant.v1beta1.MsgRevokeAllowance": IMPLEMENT_ME,
    };
}
exports.createFeegrantAminoConverters = createFeegrantAminoConverters;
//# sourceMappingURL=aminomessages.js.map