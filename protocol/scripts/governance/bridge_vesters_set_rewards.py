import json
from dateutil import parser
from datetime import datetime, timedelta
import pytz
import sys

NINE_ZEROS="000000000"


# Usage:
# 1. Update below section with appropriate values
# 2. Run `python3 bridge_vesters_set_rewards.py`
# 3. Use generated `.json` file submitting proposal. This should be done within a few hours after
#    the script is generated so that the estimated delay_blocks are accurate.


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
# TODO: update as needed. Can be checked in Mintscan by visiting a future block: 
# `www.mintscan.io/chain_name/block/9000000`
# This can fluctuate; value range for reference: 1.2 <= X <= 1.3
AVG_BLOCK_TIME_FOR_ESTIMATE=1.25
########################################################################## 
### END: Required proposal fields.                                     ###
### TODO: update above fields as needed                                ###
########################################################################## 

##########################################################################
### Network specific constants                                         ###
### Only change if used on a non-prod network                          ### 
##########################################################################
VOTING_PERIOD_DAYS = 4

########################################################################## 
### Script contants - do not change                                    ###
########################################################################## 
GOV_MODULE_ADDRESS="dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky" # Governance module account
COMMUNITY_VESTER_ADDRESS="dydx1wxje320an3karyc6mjw4zghs300dmrjkwn7xtk"
REWARDS_VESTER_ADDRESS="dydx1ltyc6y4skclzafvpznpt2qjwmfwgsndp458rmp"
DELAY_MSG_MODULE_ADDRESS="dydx1mkkvp26dngu6n8rmalaxyp3gwkjuzztq5zx6tr"
DEPOSIT=f"10000{NINE_ZEROS}{NINE_ZEROS}{NATIVE_TOKEN_DENOM}" # 10,000 native tokens
OUTPUT_FILE='proposal_bridge_vesters_set_rewards.json'

# Helper function to estimate number of blocks between two times.
def estimate_blocks_between_timestamps(base_block_timestamp, average_block_time, future_timestamp):
    # Convert ISO 8601 timestamps to datetime objects
    base_block_time = parser.isoparse(base_block_timestamp)
    future_time = parser.isoparse(future_timestamp)

    # Calculate the time difference in seconds
    time_difference = (future_time - base_block_time).total_seconds()

    # Estimate the number of blocks that will be added in this time
    estimated_blocks = float(time_difference) / average_block_time

    return int(round(estimated_blocks))

# Get current time in UTC
current_utc_time = datetime.now(pytz.utc)
# Add voting period so we can use the estimated time for gov proposal execution.
estimated_proposal_pass_time = current_utc_time + timedelta(days=VOTING_PERIOD_DAYS)
formatted_estimated_proposal_pass_time = estimated_proposal_pass_time.strftime("%Y-%m-%dT%H:%M:%S+00:00")

# Delayed update block heights are expected be after the proposal pass time.
delay_blocks_update_1 = estimate_blocks_between_timestamps(
    formatted_estimated_proposal_pass_time, 
    AVG_BLOCK_TIME_FOR_ESTIMATE,
    UPDATE_1_TIME_UTC,
)
delay_blocks_update_2 = estimate_blocks_between_timestamps(
    formatted_estimated_proposal_pass_time, 
    AVG_BLOCK_TIME_FOR_ESTIMATE,
    UPDATE_2_TIME_UTC,
)

if delay_blocks_update_1 <= 0 or delay_blocks_update_2 <= 0:
    sys.exit(f"Estimated delay_blocks <= 0: {delay_blocks_update_1}, {delay_blocks_update_2}")

print(f"Estimated proposal pass time ({VOTING_PERIOD_DAYS} days from now) = {formatted_estimated_proposal_pass_time}")
print(f"Estimated block delay for **first** update @ {UPDATE_1_TIME_UTC} = {delay_blocks_update_1}")
print(f"Estimated block delay for **second** update @ {UPDATE_2_TIME_UTC} = {delay_blocks_update_2}")

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
                "treasury_account": "rewards_treasury",
                "denom": NATIVE_TOKEN_DENOM,
                "denom_exponent": -18,
                "market_id": 1000001,
                "fee_multiplier_ppm": REWARDS_MULTIPLIER
            }
        },
    ]
}

# Add delayed messages
for delayed_block_number, new_fee_multiplier in [
    (delay_blocks_update_1, REWARDS_MULTIPLIER_UPDATE_1), 
    (delay_blocks_update_2, REWARDS_MULTIPLIER_UPDATE_2),
]:
    proposal_template["messages"].append(
        {
            "@type": "/dydxprotocol.delaymsg.MsgDelayMessage",
            "authority": GOV_MODULE_ADDRESS,
            "msg": {
                "@type": "/dydxprotocol.rewards.MsgUpdateParams",
                "authority": DELAY_MSG_MODULE_ADDRESS,
                "params": {
                    "treasury_account": "rewards_treasury",
                    "denom": NATIVE_TOKEN_DENOM,
                    "denom_exponent": -18,
                    "market_id": 1000001,
                    "fee_multiplier_ppm": new_fee_multiplier
                }
            },
            "delay_blocks": delayed_block_number
        }
    )

with open(OUTPUT_FILE, 'w') as file:
    json.dump(proposal_template, file, indent=4)
    print(f"Output written to {OUTPUT_FILE}")
