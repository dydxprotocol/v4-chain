use cosmwasm_std::StdError;
use cw_utils::Expiration;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized")]
    Unauthorized {},

    #[error("Escrow expired (expiration: {expiration:?})")]
    Expired { expiration: Expiration },

    #[error("Escrow not expired")]
    NotExpired {},
}
