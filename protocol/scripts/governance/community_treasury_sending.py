import csv
import decimal
import json
import sys

# Note: value assumes that the token exponent is 18 (i.e. 1 full token = 1 * 10^18 native token denom)
TOKEN_MULTIPLE = decimal.Decimal(1000000000000000000)


# Usage:
# 1. Update below section with appropriate values
# 2. Run `python3 community_treasury_sending.py <csv file>`.
#    Each line should have two fields: recipient address and number of full tokens to send.
# 3. Use generated `.json` file to submit proposal.

########################################################################## 
### BEGIN: Required proposal fields.                                   ###
### TODO: update below fields as needed                                ###
########################################################################## 
TITLE="TODO: Fill in proposal title" 
NATIVE_TOKEN_DENOM="adv4tnt" # TODO: Replace with production token 
PROPOSAL_BODY="""
TODO: Include a brief summary of the proposal and link to relevant governance forum discussion
"""
########################################################################## 
### END: Required proposal fields.                                     ###
### TODO: update above fields as needed                                ###
########################################################################## 

GOV_MODULE_ADDRESS="dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky" # Governance module account
COMMUNITY_TREASURY_MODULE_NAME = "community_treasury"
OUTPUT_FILE="proposal_community_treasury_sending.json"
DEPOSIT = "{}{}".format(str(decimal.Decimal(10000)*TOKEN_MULTIPLE), NATIVE_TOKEN_DENOM)

def get_single_send_message(recipient_address, token_amount):
    return {
        "@type": "/dydxprotocol.sending.MsgSendFromModuleToAccount",
        "authority": GOV_MODULE_ADDRESS,
        "sender_module_name": COMMUNITY_TREASURY_MODULE_NAME,
        "recipient": recipient_address,
        "coin": {
            "amount": str(int(TOKEN_MULTIPLE * token_amount)),
            "denom":  NATIVE_TOKEN_DENOM
        }
    }

if __name__ == "__main__":
    if len(sys.argv) != 2:
        sys.exit("Usage: python3 community_treasury_sending.py <csv file>")

    proposal_template = {
        "title": TITLE,
        "deposit": DEPOSIT,
        "summary": PROPOSAL_BODY,
        "messages": [],
    }

    with open(sys.argv[1], newline='') as csvfile:
        r = csv.reader(csvfile)
        for row in r:
            proposal_template["messages"].append(get_single_send_message(row[0], decimal.Decimal(row[1])))

    with open(OUTPUT_FILE, 'w') as file:
        json.dump(proposal_template, file, indent=4)
        print(f"Output written to {OUTPUT_FILE}")
