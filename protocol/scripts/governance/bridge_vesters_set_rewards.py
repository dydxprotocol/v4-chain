import json
from dateutil import parser

NINE_ZEROS="000000000"


# Usage:
# 1. Update below section with appropriate values
# 2. Run `python3 bridge_vesters_set_rewards.py`


########################################################################## 
### BEGIN: Required proposal fields.                                   ###
### TODO: update below fields as needed                                ###
########################################################################## 
TITLE="TODO: Fill in proposal title" 
NATIVE_TOKEN_DENOM="adv4tnt" # TODO: Replace with production token 
PROPOSAL_BODY="""
TODO: Fill in proposal content. 
Include detailed description of the proposal, links to governance discussion forums (if applicable).
Any signal proposal in text should also be added here.
""" 
# TODO: update the amount of tokens to transfer from bridge to vesters. 
# Note: value assumes that the token exponent is 18 (i.e. 1 full token = 1 * 10^18 native token denom
TRANSFER_AMOUNT_TO_COMMUNITY_VESTER=f"50000000123000000{NINE_ZEROS}" # example given: 50_000_000.123 tokens
TRANSFER_AMOUNT_TO_REWARDS_VESTER=f"20000000123000000{NINE_ZEROS}" # example given: 20_000_000.123 tokens
# Rewards multiplier right after proposal passes
REWARDS_MULTIPLIER=330000 # in parts-per-million (example given: 0.33)
# First update to the rewards multiplier. 
REWARDS_MULTIPLIER_UPDATE_1=660000 # in parts-per-million (example given: 0.66)
# TODO: update time
UPDATE_1_TIME_UTC="2024-12-01T15:00:00+00:00"
# Second update to the rewards multipler. Fill in value and time.
REWARDS_MULTIPLIER_UPDATE_2=900000 # in parts-per-million (example given: 0.9)
# TODO: update time
UPDATE_2_TIME_UTC="2025-12-01T15:00:00+00:00"
# TODO: update as needed. Can be checked in Mintscan. Value range for reference: 1.2 < X < 1.35
AVG_BLOCK_TIME_FOR_ESTIMATE=1.3
########################################################################## 
### END: Required proposal fields.                                     ###
### TODO: update above fields as needed                                ###
########################################################################## 


########################################################################## 
### Script contants - do not change                                    ###
########################################################################## 
GOV_MODULE_ADDRESS="dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky" # Governance module account
COMMUNITY_VESTER_ADDRESS="dydx1wxje320an3karyc6mjw4zghs300dmrjkwn7xtk"
REWARDS_VESTER_ADDRESS="dydx1ltyc6y4skclzafvpznpt2qjwmfwgsndp458rmp"
DELAY_MSG_MODULE_ADDRESS="dydx1mkkvp26dngu6n8rmalaxyp3gwkjuzztq5zx6tr"
DEPOSIT=f"10000{NINE_ZEROS}{NINE_ZEROS}{NATIVE_TOKEN_DENOM}" # 10,000 native tokens
OUTPUT_FILE='proposal_bridge_vesters_set_rewards.json'
REF_BLOCK_HEIGHT_FOR_ESTIMATE=1793496 
REF_BLOCK_TIME_FOR_ESTIMATE="2023-11-23T17:10:19+00:00"

# Helper function to estimate future block time.
def estimate_block_height(current_block_height, current_block_timestamp, average_block_time, future_timestamp):
    """
    Estimates the block height at a future timestamp.

    :param current_block_height: Current height of the blockchain (in blocks)
    :param current_block_timestamp: UTC Timestamp of the current block (ISO 8601 format)
    :param average_block_time: Average time it takes to mine a block (in seconds, can be a float)
    :param future_timestamp: UTC Future timestamp for which block height is to be estimated (ISO 8601 format)
    :return: Estimated block height at the future timestamp
    """
    # Convert ISO 8601 timestamps to datetime objects
    current_block_time = parser.isoparse(current_block_timestamp)
    future_time = parser.isoparse(future_timestamp)

    # Calculate the time difference in seconds
    time_difference = (future_time - current_block_time).total_seconds()

    # Estimate the number of blocks that will be added in this time
    estimated_blocks = float(time_difference) / average_block_time

    # Calculate the estimated future block height
    estimated_future_block_height = current_block_height + estimated_blocks

    return int(round(estimated_future_block_height))

rewards_update_1_block = estimate_block_height(
    REF_BLOCK_HEIGHT_FOR_ESTIMATE, 
    REF_BLOCK_TIME_FOR_ESTIMATE,
    AVG_BLOCK_TIME_FOR_ESTIMATE,
    UPDATE_1_TIME_UTC,
)
rewards_update_2_block = estimate_block_height(
    REF_BLOCK_HEIGHT_FOR_ESTIMATE, 
    REF_BLOCK_TIME_FOR_ESTIMATE,
    AVG_BLOCK_TIME_FOR_ESTIMATE,
    UPDATE_2_TIME_UTC,
)

print(f"Estimated block height for rewards multiplier update 1 {UPDATE_1_TIME_UTC} = {rewards_update_1_block}")
print(f"Estimated block height for rewards multiplier update 2 {UPDATE_2_TIME_UTC} = {rewards_update_2_block}")

proposal_template = {
    "title": TITLE,
    "deposit": DEPOSIT,
    "summary": PROPOSAL_BODY,
    "messages": [
        {
            "@type": "/dydxprotocol.sending.MsgSendFromModuleToAccount",
            "authority": GOV_MODULE_ADDRESS,
            "sender_module_name": "bridge",
            "recipient": COMMUNITY_VESTER_ADDRESS,
            "coin": {
                "amount": TRANSFER_AMOUNT_TO_COMMUNITY_VESTER,
                "denom": NATIVE_TOKEN_DENOM
            }
        },
        {
            "@type": "/dydxprotocol.sending.MsgSendFromModuleToAccount",
            "authority": GOV_MODULE_ADDRESS,
            "sender_module_name": "bridge",
            "recipient": REWARDS_VESTER_ADDRESS,
            "coin": {
                "amount": TRANSFER_AMOUNT_TO_REWARDS_VESTER,
                "denom":  NATIVE_TOKEN_DENOM
            }
        },
        {
            "@type": "/dydxprotocol.rewards.MsgUpdateParams",
            "authority": GOV_MODULE_ADDRESS,
            "params": {
                "treasuryAccount": "rewards_treasury",
                "denom": NATIVE_TOKEN_DENOM,
                "denomExponent": -18,
                "marketId": 1000001,
                "fee_multiplier_ppm": REWARDS_MULTIPLIER
            }
        },
    ]
}

# Add delayed messages
for delayed_block_number, new_fee_multiplier in [
    (rewards_update_1_block, REWARDS_MULTIPLIER_UPDATE_1), 
    (rewards_update_2_block, REWARDS_MULTIPLIER_UPDATE_2),
]:
    proposal_template["messages"].append(
        {
            "@type": "/dydxprotocol.delaymsg.MsgDelayMessage",
            "authority": GOV_MODULE_ADDRESS,
            "msg": {
                "@type": "/dydxprotocol.rewards.MsgUpdateParams",
                "authority": DELAY_MSG_MODULE_ADDRESS,
                "params": {
                    "treasuryAccount": "rewards_treasury",
                    "denom": NATIVE_TOKEN_DENOM,
                    "denomExponent": -18,
                    "marketId": 1000001,
                    "fee_multiplier_ppm": new_fee_multiplier
                }
            },
            "delay_blocks": delayed_block_number
        }
    )

with open(OUTPUT_FILE, 'w') as file:
    json.dump(proposal_template, file, indent=4)
    print(f"Output written to {OUTPUT_FILE}")
