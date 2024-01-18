#[cfg(not(feature = "library"))]
use cosmwasm_std::{
    entry_point, to_binary, Addr, BankMsg, Binary, Coin, Deps, DepsMut, Env, MessageInfo, Response,
    StdResult
};

use crate::error::ContractError;
use crate::msg::{ArbiterResponse, ExecuteMsg, InstantiateMsg, QueryMsg};
use crate::state::{Config, CONFIG};
use cw2::set_contract_version;
use dydx_cosmwasm::{SendingMsg, SubaccountId};

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
    let config = Config {
        arbiter: deps.api.addr_validate(&msg.arbiter)?,
        recipient: deps.api.addr_validate(&msg.recipient)?,
        source: info.sender,
        expiration: msg.expiration,
    };

    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;

    if let Some(expiration) = msg.expiration {
        if expiration.is_expired(&env.block) {
            return Err(ContractError::Expired { expiration });
        }
    }
    CONFIG.save(deps.storage, &config)?;
    Ok(Response::default())
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response<SendingMsg>, ContractError> {
    match msg {
        ExecuteMsg::Approve { quantity } => execute_approve(deps, env, info, quantity),
        ExecuteMsg::Refund {} => execute_refund(deps, env, info),
    }
}

fn execute_approve(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    quantity: Option<u64>,
    //quantity: Option<Vec<Coin>>,
) -> Result<Response<SendingMsg>, ContractError> {
    let config = CONFIG.load(deps.storage)?;
    if info.sender != config.arbiter {
        return Err(ContractError::Unauthorized {});
    }

    // throws error if the contract is expired
    if let Some(expiration) = config.expiration {
        if expiration.is_expired(&env.block) {
            return Err(ContractError::Expired { expiration });
        }
    }

    let balance = deps.querier.query_balance(&env.contract.address, "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5".to_string())?.amount.u128() as u64;
    //let balance = deps.querier.query_all_balances(&env.contract.address)?;

    let amount = if let Some(quantity) = quantity {
        quantity
    } else {
        // release everything
        // Querier guarantees to return up-to-date data, including funds sent in this handle message
        // https://github.com/CosmWasm/wasmd/blob/master/x/wasm/internal/keeper/keeper.go#L185-L192
        balance
    };
    Ok(send_tokens(env.contract.address, config.recipient, amount, "approve"))
}

fn execute_refund(deps: DepsMut, env: Env, _info: MessageInfo) -> Result<Response<SendingMsg>, ContractError> {
    let config = CONFIG.load(deps.storage)?;
    // anyone can try to refund, as long as the contract is expired
    if let Some(expiration) = config.expiration {
        if !expiration.is_expired(&env.block) {
            return Err(ContractError::NotExpired {});
        }
    } else {
        return Err(ContractError::NotExpired {});
    }

    // Querier guarantees to return up-to-date data, including funds sent in this handle message
    // https://github.com/CosmWasm/wasmd/blob/master/x/wasm/internal/keeper/keeper.go#L185-L192
    let balance = deps.querier.query_balance(&env.contract.address, "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5".to_string())?.amount.u128() as u64;
    Ok(send_tokens(env.contract.address, config.source, balance, "refund"))
}

// this is a helper to move the tokens, so the business logic is easy to read
fn send_tokens(from_address: Addr, to_address: Addr, amount: u64, action: &str) -> Response<SendingMsg> {
    let deposit = SendingMsg::DepositToSubaccount {
        sender: from_address.clone().into(),
        recipient: SubaccountId {
            owner: to_address.clone().into(),
            number: 0,
        },
        asset_id: 0,
        quantums: amount,
    };
    Response::new()
        .add_message(deposit)
        .add_attribute("action", action)
        .add_attribute("to", to_address)
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::Arbiter {} => to_binary(&query_arbiter(deps)?),
    }
}

fn query_arbiter(deps: Deps) -> StdResult<ArbiterResponse> {
    let config = CONFIG.load(deps.storage)?;
    let addr = config.arbiter;
    Ok(ArbiterResponse { arbiter: addr })
}
