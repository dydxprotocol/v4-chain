/// <reference types="long" />
import { Data, DataSDKType, Commit, CommitSDKType, BlockID, BlockIDSDKType } from "../../../../tendermint/types/types";
import { EvidenceList, EvidenceListSDKType } from "../../../../tendermint/types/evidence";
import { Consensus, ConsensusSDKType } from "../../../../tendermint/version/types";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../../../helpers";
/**
 * Block is tendermint type Block, with the Header proposer address
 * field converted to bech32 string.
 */
export interface Block {
    header?: Header;
    data?: Data;
    evidence?: EvidenceList;
    lastCommit?: Commit;
}
/**
 * Block is tendermint type Block, with the Header proposer address
 * field converted to bech32 string.
 */
export interface BlockSDKType {
    header?: HeaderSDKType;
    data?: DataSDKType;
    evidence?: EvidenceListSDKType;
    last_commit?: CommitSDKType;
}
/** Header defines the structure of a Tendermint block header. */
export interface Header {
    /** basic block info */
    version?: Consensus;
    chainId: string;
    height: Long;
    time?: Date;
    /** prev block info */
    lastBlockId?: BlockID;
    /** hashes of block data */
    lastCommitHash: Uint8Array;
    dataHash: Uint8Array;
    /** hashes from the app output from the prev block */
    validatorsHash: Uint8Array;
    /** validators for the next block */
    nextValidatorsHash: Uint8Array;
    /** consensus params for current block */
    consensusHash: Uint8Array;
    /** state after txs from the previous block */
    appHash: Uint8Array;
    lastResultsHash: Uint8Array;
    /** consensus info */
    evidenceHash: Uint8Array;
    /**
     * proposer_address is the original block proposer address, formatted as a Bech32 string.
     * In Tendermint, this type is `bytes`, but in the SDK, we convert it to a Bech32 string
     * for better UX.
     */
    proposerAddress: string;
}
/** Header defines the structure of a Tendermint block header. */
export interface HeaderSDKType {
    version?: ConsensusSDKType;
    chain_id: string;
    height: Long;
    time?: Date;
    last_block_id?: BlockIDSDKType;
    last_commit_hash: Uint8Array;
    data_hash: Uint8Array;
    validators_hash: Uint8Array;
    next_validators_hash: Uint8Array;
    consensus_hash: Uint8Array;
    app_hash: Uint8Array;
    last_results_hash: Uint8Array;
    evidence_hash: Uint8Array;
    proposer_address: string;
}
export declare const Block: {
    encode(message: Block, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Block;
    fromPartial(object: DeepPartial<Block>): Block;
};
export declare const Header: {
    encode(message: Header, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Header;
    fromPartial(object: DeepPartial<Header>): Header;
};
