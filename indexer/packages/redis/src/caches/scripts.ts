import { createHash } from 'crypto';
import { readFileSync } from 'fs';
import path from 'path';

import { Callback, RedisClient } from 'redis';

import { LuaScript } from '../types';

/**
 * Loads the specified script on the supplied redis client.
 *
 * @param script The script to load.
 * @param client The redis client to load the script on.
 */
export function loadScript(script: LuaScript, client: RedisClient): Promise<void> {
  return new Promise((resolve, reject) => {
    const callback: Callback<string> = (
      err: Error | null,
      result: string,
    ) => {
      if (err) {
        return reject(err);
      } else if (script.hash !== result) {
        return reject(new Error(
          `SHA1 hash does not match for script ${script.name}. Expected ${script.hash}, received ${result}.`));
      }
      return resolve();
    };
    client.script(
      'load',
      script.script,
      callback,
    );
  });
}

/**
 * Loads a named script from the specified path.
 *
 * @param name The name of the script.
 * @param scriptPath The path to the script.
 * @returns The created script object
 */
function newLuaScript(name: string, scriptPath: string): LuaScript {
  const script: string = readFileSync(path.resolve(__dirname, scriptPath)).toString();
  const hash: string = createHash('sha1').update(script, 'utf8').digest('hex');
  return {
    name,
    script,
    hash,
  };
}

// Lua Scripts for deleting zero price levels
export const deleteZeroPriceLevelScript: LuaScript = newLuaScript('deleteZeroPriceLevel', '../scripts/delete_zero_level.lua');
export const deleteStalePriceLevelScript: LuaScript = newLuaScript('deleteStalePriceLevel', '../scripts/delete_stale_price_level.lua');
// Lua Scripts for updating/retrieving the orderbook levels, keeping the lastUpdated cache in sync
export const incrementOrderbookLevelScript: LuaScript = newLuaScript('incrementOrderbookLevel', '../scripts/increment_orderbook_level.lua');
export const updateOrderScript: LuaScript = newLuaScript('updateOrder', '../scripts/update_order.lua');
export const placeOrderScript: LuaScript = newLuaScript('placeOrder', '../scripts/place_order.lua');
export const removeOrderScript: LuaScript = newLuaScript('removeOrder', '../scripts/remove_order.lua');
export const addCanceledOrderIdScript: LuaScript = newLuaScript('addCanceledOrderId', '../scripts/add_canceled_order_id.lua');
export const addStatefulOrderUpdateScript: LuaScript = newLuaScript('addStatefulOrderUpdate', '../scripts/add_stateful_order_update.lua');
export const removeStatefulOrderUpdateScript: LuaScript = newLuaScript('removeStatefulOrderUpdate', '../scripts/remove_stateful_order_update.lua');
export const addOrderbookMidPricesScript: LuaScript = newLuaScript('addOrderbookMidPrices', '../scripts/add_orderbook_mid_prices.lua');
export const getOrderbookMidPricesScript: LuaScript = newLuaScript('getOrderbookMidPrices', '../scripts/get_orderbook_mid_prices.lua');

export const allLuaScripts: LuaScript[] = [
  deleteZeroPriceLevelScript,
  deleteStalePriceLevelScript,
  incrementOrderbookLevelScript,
  updateOrderScript,
  placeOrderScript,
  removeOrderScript,
  addCanceledOrderIdScript,
  addStatefulOrderUpdateScript,
  removeStatefulOrderUpdateScript,
  addOrderbookMidPricesScript,
  getOrderbookMidPricesScript,
];
