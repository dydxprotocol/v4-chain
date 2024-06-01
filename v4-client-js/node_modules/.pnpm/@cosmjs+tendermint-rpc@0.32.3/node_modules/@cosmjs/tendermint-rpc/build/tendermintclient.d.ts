import { Comet38Client } from "./comet38";
import { HttpEndpoint } from "./rpcclients";
import { Tendermint34Client } from "./tendermint34";
import { Tendermint37Client } from "./tendermint37";
/**
 * A TendermintClient is either a Tendermint34Client or a Tendermint37Client
 *
 * @deprecated use `CometClient`
 */
export type TendermintClient = Tendermint34Client | Tendermint37Client;
/** A CometClient is either a Tendermint34Client, Tendermint37Client or a Comet38Client */
export type CometClient = Tendermint34Client | Tendermint37Client | Comet38Client;
export declare function isTendermint34Client(client: CometClient): client is Tendermint34Client;
export declare function isTendermint37Client(client: CometClient): client is Tendermint37Client;
export declare function isComet38Client(client: CometClient): client is Comet38Client;
/**
 * Auto-detects the version of the backend and uses a suitable client.
 */
export declare function connectComet(endpoint: string | HttpEndpoint): Promise<CometClient>;
