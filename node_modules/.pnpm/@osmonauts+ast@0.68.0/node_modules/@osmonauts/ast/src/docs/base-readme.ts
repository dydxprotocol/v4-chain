interface ReadmeArgs {
    imgSrc: string;
    description: string;
    libName: string;
    libPrettyName: string;
    baseModule: string;
    exampleAddr: string;
    signingBaseClient: string;
    chainName: string;
    denom: string;
}

const replaceChars = (str: string) => {
    return str.split(' ').map(s => {
        return s.replace(/\W/g, '')
    }).join('-').toLowerCase();
};

export const readme = ({ imgSrc, description, libName, libPrettyName, baseModule, exampleAddr, signingBaseClient, chainName, denom }: ReadmeArgs) => {
    return `
# ${libPrettyName}

<p align="center">
    <img src="${imgSrc}" width="80"><br />
    ${description}
</p>


## install

\`\`\`sh
npm install ${libName}
\`\`\`
## Table of contents

- [${libPrettyName}](#${replaceChars(libName)})
    - [Install](#install)
    - [Table of contents](#table-of-contents)
- [Usage](#usage)
    - [RPC Clients](#rpc-clients)
    - [Composing Messages](#composing-messages)
        - ${baseModule}
        - Cosmos, CosmWasm, and IBC
            - [CosmWasm](#cosmwasm-messages)
            - [IBC](#ibc-messages)
            - [Cosmos](#cosmos-messages)
- [Wallets and Signers](#connecting-with-wallets-and-signing-messages)
    - [Stargate Client](#initializing-the-stargate-client)
    - [Creating Signers](#creating-signers)
    - [Broadcasting Messages](#broadcasting-messages)
- [Advanced Usage](#advanced-usage)
- [Developing](#developing)
- [Credits](#credits)

## Usage

### RPC Clients

\`\`\`js
import { ${baseModule} } from '${libName}';

const { createRPCQueryClient } = ${baseModule}.ClientFactory; 
const client = await createRPCQueryClient({ rpcEndpoint: RPC_ENDPOINT });

// now you can query the cosmos modules
const balance = await client.cosmos.bank.v1beta1
    .allBalances({ address: '${exampleAddr}' });

// you can also query the ${baseModule} modules
const balances = await client.${baseModule}.exchange.v1beta1
    .exchangeBalances()
\`\`\`

### Composing Messages

Import the \`${baseModule}\` object from \`${libName}\`. 

\`\`\`js
import { ${baseModule} } from '${libName}';

const {
    createSpotLimitOrder,
    createSpotMarketOrder,
    deposit
} = ${baseModule}.exchange.v1beta1.MessageComposer.withTypeUrl;
\`\`\`

#### Auction Messages

\`\`\`js
const {
    bid
} = ${baseModule}.auction.v1beta1.MessageComposer.withTypeUrl;
\`\`\`

#### CosmWasm Messages

\`\`\`js
import { cosmwasm } from "${libName}";

const {
    clearAdmin,
    executeContract,
    instantiateContract,
    migrateContract,
    storeCode,
    updateAdmin
} = cosmwasm.wasm.v1.MessageComposer.withTypeUrl;
\`\`\`

#### IBC Messages

\`\`\`js
import { ibc } from '${libName}';

const {
    transfer
} = ibc.applications.transfer.v1.MessageComposer.withTypeUrl
\`\`\`

#### Cosmos Messages

\`\`\`js
import { cosmos } from '${libName}';

const {
    fundCommunityPool,
    setWithdrawAddress,
    withdrawDelegatorReward,
    withdrawValidatorCommission
} = cosmos.distribution.v1beta1.MessageComposer.fromPartial;

const {
    multiSend,
    send
} = cosmos.bank.v1beta1.MessageComposer.fromPartial;

const {
    beginRedelegate,
    createValidator,
    delegate,
    editValidator,
    undelegate
} = cosmos.staking.v1beta1.MessageComposer.fromPartial;

const {
    deposit,
    submitProposal,
    vote,
    voteWeighted
} = cosmos.gov.v1beta1.MessageComposer.fromPartial;
\`\`\`

## Connecting with Wallets and Signing Messages

⚡️ For web interfaces, we recommend using [cosmos-kit](https://github.com/cosmology-tech/cosmos-kit). Continue below to see how to manually construct signers and clients.

Here are the docs on [creating signers](https://github.com/cosmology-tech/cosmos-kit/tree/main/packages/react#signing-clients) in cosmos-kit that can be used with Keplr and other wallets.

### Initializing the Stargate Client

Use \`${signingBaseClient}\` to get your \`SigningStargateClient\`, with the proto/amino messages full-loaded. No need to manually add amino types, just require and initialize the client:

\`\`\`js
import { ${signingBaseClient} } from '${libName}';

const stargateClient = await ${signingBaseClient}({
    rpcEndpoint,
    signer // OfflineSigner
});
\`\`\`
### Creating Signers

To broadcast messages, you can create signers with a variety of options:

* [cosmos-kit](https://github.com/cosmology-tech/cosmos-kit/tree/main/packages/react#signing-clients) (recommended)
* [keplr](https://docs.keplr.app/api/cosmjs.html)
* [cosmjs](https://gist.github.com/webmaster128/8444d42a7eceeda2544c8a59fbd7e1d9)
### Amino Signer

Likely you'll want to use the Amino, so unless you need proto, you should use this one:

\`\`\`js
import { getOfflineSignerAmino as getOfflineSigner } from 'cosmjs-utils';
\`\`\`
### Proto Signer

\`\`\`js
import { getOfflineSignerProto as getOfflineSigner } from 'cosmjs-utils';
\`\`\`

WARNING: NOT RECOMMENDED TO USE PLAIN-TEXT MNEMONICS. Please take care of your security and use best practices such as AES encryption and/or methods from 12factor applications.

\`\`\`js
import { chains } from 'chain-registry';

const mnemonic =
    'unfold client turtle either pilot stock floor glow toward bullet car science';
    const chain = chains.find(({ chain_name }) => chain_name === '${chainName}');
    const signer = await getOfflineSigner({
    mnemonic,
    chain
    });
\`\`\`
### Broadcasting Messages

Now that you have your \`stargateClient\`, you can broadcast messages:

\`\`\`js
const { send } = cosmos.bank.v1beta1.MessageComposer.withTypeUrl;

const msg = send({
    amount: [
    {
        denom: '${denom}',
        amount: '1000'
    }
    ],
    toAddress: address,
    fromAddress: address
});

const fee: StdFee = {
    amount: [
    {
        denom: '${denom}',
        amount: '864'
    }
    ],
    gas: '86364' // this may need to be adjusted
};
const response = await stargateClient.signAndBroadcast(address, [msg], fee);
\`\`\`

## Advanced Usage


If you want to manually construct a stargate client

\`\`\`js
import { OfflineSigner, GeneratedType, Registry } from "@cosmjs/proto-signing";
import { AminoTypes, SigningStargateClient } from "@cosmjs/stargate";

import { 
    cosmosAminoConverters,
    cosmosProtoRegistry,
    cosmwasmAminoConverters,
    cosmwasmProtoRegistry,
    ibcProtoRegistry,
    ibcAminoConverters,
    ${baseModule}AminoConverters,
    ${baseModule}ProtoRegistry
} from '${libName}';

const signer: OfflineSigner = /* create your signer (see above)  */
const rpcEndpoint = 'https://rpc.cosmos.directory/${baseModule}'; // or another URL

const protoRegistry: ReadonlyArray<[string, GeneratedType]> = [
    ...cosmosProtoRegistry,
    ...cosmwasmProtoRegistry,
    ...ibcProtoRegistry,
    ...${baseModule}ProtoRegistry
];

const aminoConverters = {
    ...cosmosAminoConverters,
    ...cosmwasmAminoConverters,
    ...ibcAminoConverters,
    ...${baseModule}AminoConverters
};

const registry = new Registry(protoRegistry);
const aminoTypes = new AminoTypes(aminoConverters);

const stargateClient = await SigningStargateClient.connectWithSigner(rpcEndpoint, signer, {
    registry,
    aminoTypes
});
\`\`\`

## Developing

When first cloning the repo:

\`\`\`
yarn
yarn build
\`\`\`

### Codegen

Contract schemas live in \`./contracts\`, and protos in \`./proto\`. Look inside of \`scripts/codegen.js\` and configure the settings for bundling your SDK and contracts into \`${libName}\`:

\`\`\`
yarn codegen
\`\`\`

### Publishing

Build the types and then publish:

\`\`\`
yarn build:ts
yarn publish
\`\`\`

## Credits

🛠 Built by Cosmology — if you like our tools, please consider delegating to [our validator ⚛️](https://cosmology.tech/validator)

Code built with the help of these related projects:

* [@cosmwasm/ts-codegen](https://github.com/CosmWasm/ts-codegen) for generated CosmWasm contract Typescript classes
* [@osmonauts/telescope](https://github.com/osmosis-labs/telescope) a "babel for the Cosmos", Telescope is a TypeScript Transpiler for Cosmos Protobufs.
* [cosmos-kit](https://github.com/cosmology-tech/cosmos-kit) A wallet connector for the Cosmos ⚛️

    `;
}