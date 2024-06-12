use cosmwasm_std::to_json_binary;
#[cfg(not(feature = "library"))]
use cosmwasm_std::{
    entry_point, to_binary, Addr, BankMsg, Binary, Coin, Deps, DepsMut, Env, MessageInfo, Response,
    StdResult
};

use crate::error::ContractError;
use crate::msg::{ArbiterResponse, ExecuteMsg, InstantiateMsg};
use crate::state::{Config, CONFIG};
use cw2::set_contract_version;
use dydx_cosmwasm::DydxMsg;
use dydx_cosmwasm::{DydxQuerier, DydxQueryWrapper, SubaccountId};
use dydx_cosmwasm::DydxQuery;

// version info for migration info
const CONTRACT_NAME: &str = "crates.io:dydx-messages-example";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;
    Ok(Response::default())
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response<DydxMsg>, ContractError> {
    match msg {
        ExecuteMsg::DydxMsg(dydx_msg) => {
            match dydx_msg{
                DydxMsg::DepositToSubaccount { sender, recipient, asset_id, quantums } => {
                    let deposit = DydxMsg::DepositToSubaccount {
                        sender: sender,
                        recipient: recipient,
                        asset_id: asset_id,
                        quantums: quantums,
                    };
                    Ok(Response::new().add_message(deposit))
                }
                DydxMsg::WithdrawFromSubaccount { sender, recipient, asset_id, quantums } => {
                    let withdraw = DydxMsg::WithdrawFromSubaccount {
                        sender: sender,
                        recipient: recipient,
                        asset_id: asset_id,
                        quantums: quantums,
                    };
                    Ok(Response::new().add_message(withdraw))
                }
                DydxMsg::PlaceOrder { order } => {
                    let place_order = DydxMsg::PlaceOrder {
                        order: order,
                    };
                    Ok(Response::new().add_message(place_order))
                }
                DydxMsg::CancelOrder { order_id, good_til_oneof } => {
                    let cancel_order = DydxMsg::CancelOrder {
                        order_id: order_id,
                        good_til_oneof: good_til_oneof,
                    };
                    Ok(Response::new().add_message(cancel_order))
                }
                _ => Err(ContractError::InvalidDydxMsg{}),
            }
        }
    }
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps<DydxQueryWrapper>, _env: Env, msg: DydxQuery) -> StdResult<Binary> {
    let dydx_querier = DydxQuerier::new(&deps.querier);

    match msg {
        DydxQuery::MarketPrice { id } => to_json_binary(&dydx_querier.query_market_price(id)?),
        DydxQuery::Subaccount { owner, number } => to_json_binary(&dydx_querier.query_subaccount(owner, number)?),
        DydxQuery::PerpetualClobDetails { id } => to_json_binary(&dydx_querier.query_perpetual_clob_details(id)?),
    }
}