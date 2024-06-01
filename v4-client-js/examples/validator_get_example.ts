import { Network } from '../src/clients/constants';
import { ValidatorClient } from '../src/clients/validator-client';
import { DYDX_TEST_ADDRESS } from './constants';

async function test(): Promise<void> {
  const client = await ValidatorClient.connect(Network.testnet().validatorConfig);

  try {
    const account = await client.get.getAccount(DYDX_TEST_ADDRESS);
    console.log('Account');
    console.log(JSON.stringify(account));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const coins = await client.get.getAccountBalances(DYDX_TEST_ADDRESS);
    console.log('Balances');
    console.log(JSON.stringify(coins));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const subaccounts = await client.get.getSubaccounts();
    console.log('Subaccounts');
    console.log(JSON.stringify(subaccounts));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const subaccount = await client.get.getSubaccount(DYDX_TEST_ADDRESS, 0);
    console.log('Subaccount 0');
    console.log(JSON.stringify(subaccount));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const subaccount = await client.get.getSubaccount(DYDX_TEST_ADDRESS, 1);
    console.log('Subaccount 1');
    console.log(JSON.stringify(subaccount));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const clobpairs = await client.get.getAllClobPairs();
    console.log('Clobs');
    console.log(JSON.stringify(clobpairs));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const prices = await client.get.getAllPrices();
    console.log('Prices');
    console.log(JSON.stringify(prices));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const feeTiers = await client.get.getFeeTiers();
    console.log('Fee Tiers');
    console.log(JSON.stringify(feeTiers));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const userfeeTier = await client.get.getUserFeeTier(DYDX_TEST_ADDRESS);
    console.log('User Fee Tiers');
    console.log(JSON.stringify(userfeeTier));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const userStats = await client.get.getUserStats(DYDX_TEST_ADDRESS);
    console.log('User Fee Tiers');
    console.log(JSON.stringify(userStats));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const equityTierLimitConfiguration = await client.get.getEquityTierLimitConfiguration();
    console.log('Equity Tier Limit Configuration');
    console.log(JSON.stringify(equityTierLimitConfiguration));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const rewardsParams = await client.get.getRewardsParams();
    console.log('Rewards Params');
    console.log(JSON.stringify(rewardsParams));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const delegations = await client.get.getDelegatorDelegations(DYDX_TEST_ADDRESS);
    console.log('Delegations');
    console.log(JSON.stringify(delegations));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const unbondingDelegations = await client
      .get.getDelegatorUnbondingDelegations(DYDX_TEST_ADDRESS);
    console.log('Unbonding Delegationss');
    console.log(JSON.stringify(unbondingDelegations));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const getAllBondedValidators = await client.get.getAllValidators('BOND_STATUS_BONDED');
    console.log('All Validators');
    console.log(JSON.stringify(getAllBondedValidators));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const delayedCompleteBridgeMessagesParams = await client
      .get.getDelayedCompleteBridgeMessages();
    console.log('All delayed complete bridge messages');
    console.log(JSON.stringify(delayedCompleteBridgeMessagesParams));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }

  try {
    const delayedCompleteBridgeMessagesParams = await client
      .get.getDelayedCompleteBridgeMessages(DYDX_TEST_ADDRESS);
    console.log('Filtered delayed complete bridge messages');
    console.log(JSON.stringify(delayedCompleteBridgeMessagesParams));
  } catch (error) {
    console.log(JSON.stringify(error.message));
  }
}

test().then(() => {
}).catch((error) => {
  console.log(error.message);
});
