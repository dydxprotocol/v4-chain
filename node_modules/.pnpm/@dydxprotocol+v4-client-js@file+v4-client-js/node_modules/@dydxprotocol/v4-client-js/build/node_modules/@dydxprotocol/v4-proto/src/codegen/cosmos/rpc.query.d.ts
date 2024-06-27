import { HttpEndpoint } from "@cosmjs/tendermint-rpc";
export declare const createRPCQueryClient: ({ rpcEndpoint }: {
    rpcEndpoint: string | HttpEndpoint;
}) => Promise<{
    cosmos: {
        app: {
            v1alpha1: {
                config(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/app/v1alpha1/query").QueryConfigRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/app/v1alpha1/query").QueryConfigResponse>;
            };
        };
        auth: {
            v1beta1: {
                accounts(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryAccountsRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryAccountsResponse>;
                account(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryAccountRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryAccountResponse>;
                accountAddressByID(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryAccountAddressByIDRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryAccountAddressByIDResponse>;
                params(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryParamsRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryParamsResponse>;
                moduleAccounts(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryModuleAccountsRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryModuleAccountsResponse>;
                moduleAccountByName(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryModuleAccountByNameRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryModuleAccountByNameResponse>;
                bech32Prefix(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").Bech32PrefixRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").Bech32PrefixResponse>;
                addressBytesToString(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").AddressBytesToStringRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").AddressBytesToStringResponse>;
                addressStringToBytes(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").AddressStringToBytesRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").AddressStringToBytesResponse>;
                accountInfo(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryAccountInfoRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/auth/v1beta1/query").QueryAccountInfoResponse>;
            };
        };
        authz: {
            v1beta1: {
                grants(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/authz/v1beta1/query").QueryGrantsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/authz/v1beta1/query").QueryGrantsResponse>;
                granterGrants(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/authz/v1beta1/query").QueryGranterGrantsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/authz/v1beta1/query").QueryGranterGrantsResponse>;
                granteeGrants(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/authz/v1beta1/query").QueryGranteeGrantsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/authz/v1beta1/query").QueryGranteeGrantsResponse>;
            };
        };
        bank: {
            v1beta1: {
                balance(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryBalanceRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryBalanceResponse>;
                allBalances(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryAllBalancesRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryAllBalancesResponse>;
                spendableBalances(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QuerySpendableBalancesRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QuerySpendableBalancesResponse>;
                spendableBalanceByDenom(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QuerySpendableBalanceByDenomRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QuerySpendableBalanceByDenomResponse>;
                totalSupply(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryTotalSupplyRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryTotalSupplyResponse>;
                supplyOf(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QuerySupplyOfRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QuerySupplyOfResponse>;
                params(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryParamsRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryParamsResponse>;
                denomMetadata(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryDenomMetadataRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryDenomMetadataResponse>;
                denomMetadataByQueryString(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryDenomMetadataByQueryStringRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryDenomMetadataByQueryStringResponse>;
                denomsMetadata(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryDenomsMetadataRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryDenomsMetadataResponse>;
                denomOwners(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryDenomOwnersRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QueryDenomOwnersResponse>;
                sendEnabled(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QuerySendEnabledRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/bank/v1beta1/query").QuerySendEnabledResponse>;
            };
        };
        base: {
            node: {
                v1beta1: {
                    config(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/node/v1beta1/query").ConfigRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/node/v1beta1/query").ConfigResponse>;
                    status(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/node/v1beta1/query").StatusRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/node/v1beta1/query").StatusResponse>;
                };
            };
            tendermint: {
                v1beta1: {
                    getNodeInfo(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").GetNodeInfoRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").GetNodeInfoResponse>;
                    getSyncing(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").GetSyncingRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").GetSyncingResponse>;
                    getLatestBlock(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").GetLatestBlockRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").GetLatestBlockResponse>;
                    getBlockByHeight(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").GetBlockByHeightRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").GetBlockByHeightResponse>;
                    getLatestValidatorSet(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").GetLatestValidatorSetRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").GetLatestValidatorSetResponse>;
                    getValidatorSetByHeight(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").GetValidatorSetByHeightRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").GetValidatorSetByHeightResponse>;
                    aBCIQuery(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").ABCIQueryRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/base/tendermint/v1beta1/query").ABCIQueryResponse>;
                };
            };
        };
        circuit: {
            v1: {
                account(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/circuit/v1/query").QueryAccountRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/circuit/v1/query").AccountResponse>;
                accounts(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/circuit/v1/query").QueryAccountsRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/circuit/v1/query").AccountsResponse>;
                disabledList(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/circuit/v1/query").QueryDisabledListRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/circuit/v1/query").DisabledListResponse>;
            };
        };
        consensus: {
            v1: {
                params(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/consensus/v1/query").QueryParamsRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/consensus/v1/query").QueryParamsResponse>;
            };
        };
        distribution: {
            v1beta1: {
                params(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryParamsRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryParamsResponse>;
                validatorDistributionInfo(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryValidatorDistributionInfoRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryValidatorDistributionInfoResponse>;
                validatorOutstandingRewards(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryValidatorOutstandingRewardsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryValidatorOutstandingRewardsResponse>;
                validatorCommission(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryValidatorCommissionRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryValidatorCommissionResponse>;
                validatorSlashes(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryValidatorSlashesRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryValidatorSlashesResponse>;
                delegationRewards(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryDelegationRewardsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryDelegationRewardsResponse>;
                delegationTotalRewards(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryDelegationTotalRewardsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryDelegationTotalRewardsResponse>;
                delegatorValidators(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryDelegatorValidatorsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryDelegatorValidatorsResponse>;
                delegatorWithdrawAddress(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryDelegatorWithdrawAddressRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryDelegatorWithdrawAddressResponse>;
                communityPool(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryCommunityPoolRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/distribution/v1beta1/query").QueryCommunityPoolResponse>;
            };
        };
        evidence: {
            v1beta1: {
                evidence(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/evidence/v1beta1/query").QueryEvidenceRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/evidence/v1beta1/query").QueryEvidenceResponse>;
                allEvidence(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/evidence/v1beta1/query").QueryAllEvidenceRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/evidence/v1beta1/query").QueryAllEvidenceResponse>;
            };
        };
        feegrant: {
            v1beta1: {
                allowance(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/feegrant/v1beta1/query").QueryAllowanceRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/feegrant/v1beta1/query").QueryAllowanceResponse>;
                allowances(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/feegrant/v1beta1/query").QueryAllowancesRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/feegrant/v1beta1/query").QueryAllowancesResponse>;
                allowancesByGranter(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/feegrant/v1beta1/query").QueryAllowancesByGranterRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/feegrant/v1beta1/query").QueryAllowancesByGranterResponse>;
            };
        };
        gov: {
            v1: {
                constitution(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryConstitutionRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryConstitutionResponse>;
                proposal(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryProposalRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryProposalResponse>;
                proposals(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryProposalsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryProposalsResponse>;
                vote(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryVoteRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryVoteResponse>;
                votes(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryVotesRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryVotesResponse>;
                params(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryParamsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryParamsResponse>;
                deposit(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryDepositRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryDepositResponse>;
                deposits(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryDepositsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryDepositsResponse>;
                tallyResult(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryTallyResultRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1/query").QueryTallyResultResponse>;
            };
            v1beta1: {
                proposal(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryProposalRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryProposalResponse>;
                proposals(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryProposalsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryProposalsResponse>;
                vote(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryVoteRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryVoteResponse>;
                votes(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryVotesRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryVotesResponse>;
                params(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryParamsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryParamsResponse>;
                deposit(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryDepositRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryDepositResponse>;
                deposits(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryDepositsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryDepositsResponse>;
                tallyResult(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryTallyResultRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/gov/v1beta1/query").QueryTallyResultResponse>;
            };
        };
        group: {
            v1: {
                groupInfo(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupInfoRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupInfoResponse>;
                groupPolicyInfo(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupPolicyInfoRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupPolicyInfoResponse>;
                groupMembers(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupMembersRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupMembersResponse>;
                groupsByAdmin(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupsByAdminRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupsByAdminResponse>;
                groupPoliciesByGroup(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupPoliciesByGroupRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupPoliciesByGroupResponse>;
                groupPoliciesByAdmin(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupPoliciesByAdminRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupPoliciesByAdminResponse>;
                proposal(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryProposalRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryProposalResponse>;
                proposalsByGroupPolicy(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryProposalsByGroupPolicyRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryProposalsByGroupPolicyResponse>;
                voteByProposalVoter(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryVoteByProposalVoterRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryVoteByProposalVoterResponse>;
                votesByProposal(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryVotesByProposalRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryVotesByProposalResponse>;
                votesByVoter(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryVotesByVoterRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryVotesByVoterResponse>;
                groupsByMember(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupsByMemberRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupsByMemberResponse>;
                tallyResult(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryTallyResultRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryTallyResultResponse>;
                groups(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupsRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/group/v1/query").QueryGroupsResponse>;
            };
        };
        mint: {
            v1beta1: {
                params(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/mint/v1beta1/query").QueryParamsRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/mint/v1beta1/query").QueryParamsResponse>;
                inflation(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/mint/v1beta1/query").QueryInflationRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/mint/v1beta1/query").QueryInflationResponse>;
                annualProvisions(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/mint/v1beta1/query").QueryAnnualProvisionsRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/mint/v1beta1/query").QueryAnnualProvisionsResponse>;
            };
        };
        nft: {
            v1beta1: {
                balance(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QueryBalanceRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QueryBalanceResponse>;
                owner(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QueryOwnerRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QueryOwnerResponse>;
                supply(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QuerySupplyRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QuerySupplyResponse>;
                nFTs(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QueryNFTsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QueryNFTsResponse>;
                nFT(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QueryNFTRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QueryNFTResponse>;
                class(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QueryClassRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QueryClassResponse>;
                classes(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QueryClassesRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/nft/v1beta1/query").QueryClassesResponse>;
            };
        };
        orm: {
            query: {
                v1alpha1: {
                    get(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/orm/query/v1alpha1/query").GetRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/orm/query/v1alpha1/query").GetResponse>;
                    list(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/orm/query/v1alpha1/query").ListRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/orm/query/v1alpha1/query").ListResponse>;
                };
            };
        };
        params: {
            v1beta1: {
                params(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/params/v1beta1/query").QueryParamsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/params/v1beta1/query").QueryParamsResponse>;
                subspaces(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/params/v1beta1/query").QuerySubspacesRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/params/v1beta1/query").QuerySubspacesResponse>;
            };
        };
        slashing: {
            v1beta1: {
                params(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/slashing/v1beta1/query").QueryParamsRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/slashing/v1beta1/query").QueryParamsResponse>;
                signingInfo(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/slashing/v1beta1/query").QuerySigningInfoRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/slashing/v1beta1/query").QuerySigningInfoResponse>;
                signingInfos(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/slashing/v1beta1/query").QuerySigningInfosRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/slashing/v1beta1/query").QuerySigningInfosResponse>;
            };
        };
        staking: {
            v1beta1: {
                validators(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryValidatorsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryValidatorsResponse>;
                validator(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryValidatorRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryValidatorResponse>;
                validatorDelegations(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryValidatorDelegationsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryValidatorDelegationsResponse>;
                validatorUnbondingDelegations(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryValidatorUnbondingDelegationsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryValidatorUnbondingDelegationsResponse>;
                delegation(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryDelegationRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryDelegationResponse>;
                unbondingDelegation(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryUnbondingDelegationRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryUnbondingDelegationResponse>;
                delegatorDelegations(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryDelegatorDelegationsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryDelegatorDelegationsResponse>;
                delegatorUnbondingDelegations(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryDelegatorUnbondingDelegationsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryDelegatorUnbondingDelegationsResponse>;
                redelegations(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryRedelegationsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryRedelegationsResponse>;
                delegatorValidators(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryDelegatorValidatorsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryDelegatorValidatorsResponse>;
                delegatorValidator(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryDelegatorValidatorRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryDelegatorValidatorResponse>;
                historicalInfo(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryHistoricalInfoRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryHistoricalInfoResponse>;
                pool(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryPoolRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryPoolResponse>;
                params(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryParamsRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/staking/v1beta1/query").QueryParamsResponse>;
            };
        };
        tx: {
            v1beta1: {
                simulate(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").SimulateRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").SimulateResponse>;
                getTx(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").GetTxRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").GetTxResponse>;
                broadcastTx(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").BroadcastTxRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").BroadcastTxResponse>;
                getTxsEvent(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").GetTxsEventRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").GetTxsEventResponse>;
                getBlockWithTxs(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").GetBlockWithTxsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").GetBlockWithTxsResponse>;
                txDecode(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").TxDecodeRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").TxDecodeResponse>;
                txEncode(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").TxEncodeRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").TxEncodeResponse>;
                txEncodeAmino(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").TxEncodeAminoRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").TxEncodeAminoResponse>;
                txDecodeAmino(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").TxDecodeAminoRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/tx/v1beta1/service").TxDecodeAminoResponse>;
            };
        };
        upgrade: {
            v1beta1: {
                currentPlan(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/upgrade/v1beta1/query").QueryCurrentPlanRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/upgrade/v1beta1/query").QueryCurrentPlanResponse>;
                appliedPlan(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/upgrade/v1beta1/query").QueryAppliedPlanRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/upgrade/v1beta1/query").QueryAppliedPlanResponse>;
                upgradedConsensusState(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/upgrade/v1beta1/query").QueryUpgradedConsensusStateRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/upgrade/v1beta1/query").QueryUpgradedConsensusStateResponse>;
                moduleVersions(request: import("@dydxprotocol/v4-proto/src/codegen/cosmos/upgrade/v1beta1/query").QueryModuleVersionsRequest): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/upgrade/v1beta1/query").QueryModuleVersionsResponse>;
                authority(request?: import("@dydxprotocol/v4-proto/src/codegen/cosmos/upgrade/v1beta1/query").QueryAuthorityRequest | undefined): Promise<import("@dydxprotocol/v4-proto/src/codegen/cosmos/upgrade/v1beta1/query").QueryAuthorityResponse>;
            };
        };
    };
}>;
