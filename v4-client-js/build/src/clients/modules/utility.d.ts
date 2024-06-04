import type { ComplianceResponse, HeightResponse, TimeResponse } from '../types';
import RestClient from './rest';
export default class UtilityClient extends RestClient {
    /**
     * @description Get the current time of the Indexer
     * @returns {TimeResponse} isoString and epoch
     */
    getTime(): Promise<TimeResponse>;
    /**
     * @description Get the block height of the most recent block processed by the Indexer
     * @returns {HeightResponse} block height and time
     */
    getHeight(): Promise<HeightResponse>;
    /**
     * @description Screen an address to see if it is restricted
     * @param {string} address evm or dydx address
     * @returns {ComplianceResponse} whether the specified address is restricted
     */
    screen(address: string): Promise<ComplianceResponse>;
}
