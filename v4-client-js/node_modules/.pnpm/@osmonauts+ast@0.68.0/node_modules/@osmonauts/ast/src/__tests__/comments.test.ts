import { expectCode } from "../../test-utils";

it('AminoConverter', () => {
    expectCode({
        "type": "ExportNamedDeclaration",
        "exportKind": "type",
        "specifiers": [],
        "source": null,
        "declaration": {
            "type": "TSInterfaceDeclaration",
            "id": {
                "type": "Identifier",
                "name": "Service"
            },
            "body": {
                "type": "TSInterfaceBody",
                "body": [
                    {
                        "type": "TSMethodSignature",
                        "key": {
                            "type": "Identifier",
                            "name": "Simulate"
                        },
                        "computed": false,
                        "parameters": [
                            {
                                "type": "Identifier",
                                "name": "request",
                                "typeAnnotation": {
                                    "type": "TSTypeAnnotation",
                                    "typeAnnotation": {
                                        "type": "TSTypeReference",
                                        "typeName": {
                                            "type": "Identifier",
                                            "name": "SimulateRequest"
                                        }
                                    }
                                }
                            }
                        ],
                        "typeAnnotation": {
                            "type": "TSTypeAnnotation",
                            "typeAnnotation": {
                                "type": "TSTypeReference",
                                "typeName": {
                                    "type": "Identifier",
                                    "name": "Promise"
                                },
                                "typeParameters": {
                                    "type": "TSTypeParameterInstantiation",
                                    "params": [
                                        {
                                            "type": "TSTypeReference",
                                            "typeName": {
                                                "type": "Identifier",
                                                "name": "SimulateResponse"
                                            }
                                        }
                                    ]
                                }
                            }
                        },
                        "kind": "method",
                        "trailingComments": [
                            {
                                "type": "CommentBlock",
                                "value": "* GetTx fetches a tx by hash. "
                            }
                        ],
                        "leadingComments": [
                            {
                                "type": "CommentBlock",
                                "value": "* Simulate simulates executing a transaction for estimating gas usage. "
                            }
                        ]
                    },
                    {
                        "type": "TSMethodSignature",
                        "key": {
                            "type": "Identifier",
                            "name": "GetTx"
                        },
                        "computed": false,
                        "parameters": [
                            {
                                "type": "Identifier",
                                "name": "request",
                                "typeAnnotation": {
                                    "type": "TSTypeAnnotation",
                                    "typeAnnotation": {
                                        "type": "TSTypeReference",
                                        "typeName": {
                                            "type": "Identifier",
                                            "name": "GetTxRequest"
                                        }
                                    }
                                }
                            }
                        ],
                        "typeAnnotation": {
                            "type": "TSTypeAnnotation",
                            "typeAnnotation": {
                                "type": "TSTypeReference",
                                "typeName": {
                                    "type": "Identifier",
                                    "name": "Promise"
                                },
                                "typeParameters": {
                                    "type": "TSTypeParameterInstantiation",
                                    "params": [
                                        {
                                            "type": "TSTypeReference",
                                            "typeName": {
                                                "type": "Identifier",
                                                "name": "GetTxResponse"
                                            }
                                        }
                                    ]
                                }
                            }
                        },
                        "kind": "method",
                        "trailingComments": [
                            {
                                "type": "CommentBlock",
                                "value": "* BroadcastTx broadcast transaction. "
                            }
                        ],
                        "leadingComments": [
                            {
                                "type": "CommentBlock",
                                "value": "* GetTx fetches a tx by hash. "
                            }
                        ]
                    },
                    {
                        "type": "TSMethodSignature",
                        "key": {
                            "type": "Identifier",
                            "name": "BroadcastTx"
                        },
                        "computed": false,
                        "parameters": [
                            {
                                "type": "Identifier",
                                "name": "request",
                                "typeAnnotation": {
                                    "type": "TSTypeAnnotation",
                                    "typeAnnotation": {
                                        "type": "TSTypeReference",
                                        "typeName": {
                                            "type": "Identifier",
                                            "name": "BroadcastTxRequest"
                                        }
                                    }
                                }
                            }
                        ],
                        "typeAnnotation": {
                            "type": "TSTypeAnnotation",
                            "typeAnnotation": {
                                "type": "TSTypeReference",
                                "typeName": {
                                    "type": "Identifier",
                                    "name": "Promise"
                                },
                                "typeParameters": {
                                    "type": "TSTypeParameterInstantiation",
                                    "params": [
                                        {
                                            "type": "TSTypeReference",
                                            "typeName": {
                                                "type": "Identifier",
                                                "name": "BroadcastTxResponse"
                                            }
                                        }
                                    ]
                                }
                            }
                        },
                        "kind": "method",
                        "trailingComments": [
                            {
                                "type": "CommentBlock",
                                "value": "* GetTxsEvent fetches txs by event. "
                            }
                        ],
                        "leadingComments": [
                            {
                                "type": "CommentBlock",
                                "value": "* BroadcastTx broadcast transaction. "
                            }
                        ]
                    },
                    {
                        "type": "TSMethodSignature",
                        "key": {
                            "type": "Identifier",
                            "name": "GetTxsEvent"
                        },
                        "computed": false,
                        "parameters": [
                            {
                                "type": "Identifier",
                                "name": "request",
                                "typeAnnotation": {
                                    "type": "TSTypeAnnotation",
                                    "typeAnnotation": {
                                        "type": "TSTypeReference",
                                        "typeName": {
                                            "type": "Identifier",
                                            "name": "GetTxsEventRequest"
                                        }
                                    }
                                }
                            }
                        ],
                        "typeAnnotation": {
                            "type": "TSTypeAnnotation",
                            "typeAnnotation": {
                                "type": "TSTypeReference",
                                "typeName": {
                                    "type": "Identifier",
                                    "name": "Promise"
                                },
                                "typeParameters": {
                                    "type": "TSTypeParameterInstantiation",
                                    "params": [
                                        {
                                            "type": "TSTypeReference",
                                            "typeName": {
                                                "type": "Identifier",
                                                "name": "GetTxsEventResponse"
                                            }
                                        }
                                    ]
                                }
                            }
                        },
                        "kind": "method",
                        "trailingComments": [
                            {
                                "type": "CommentBlock",
                                "value": "*\n     * GetBlockWithTxs fetches a block with decoded txs.\n     *\n     * Since: cosmos-sdk 0.45.2\n     "
                            }
                        ],
                        "leadingComments": [
                            {
                                "type": "CommentBlock",
                                "value": "* GetTxsEvent fetches txs by event. "
                            }
                        ]
                    },
                    {
                        "type": "TSMethodSignature",
                        "key": {
                            "type": "Identifier",
                            "name": "GetBlockWithTxs"
                        },
                        "computed": false,
                        "parameters": [
                            {
                                "type": "Identifier",
                                "name": "request",
                                "typeAnnotation": {
                                    "type": "TSTypeAnnotation",
                                    "typeAnnotation": {
                                        "type": "TSTypeReference",
                                        "typeName": {
                                            "type": "Identifier",
                                            "name": "GetBlockWithTxsRequest"
                                        }
                                    }
                                }
                            }
                        ],
                        "typeAnnotation": {
                            "type": "TSTypeAnnotation",
                            "typeAnnotation": {
                                "type": "TSTypeReference",
                                "typeName": {
                                    "type": "Identifier",
                                    "name": "Promise"
                                },
                                "typeParameters": {
                                    "type": "TSTypeParameterInstantiation",
                                    "params": [
                                        {
                                            "type": "TSTypeReference",
                                            "typeName": {
                                                "type": "Identifier",
                                                "name": "GetBlockWithTxsResponse"
                                            }
                                        }
                                    ]
                                }
                            }
                        },
                        "kind": "method",
                        "leadingComments": [
                            {
                                "type": "CommentBlock",
                                "value": "*\n     * GetBlockWithTxs fetches a block with decoded txs.\n     *\n     * Since: cosmos-sdk 0.45.2\n     "
                            }
                        ]
                    }
                ]
            }
        },
        "leadingComments": [
            {
                "type": "CommentBlock",
                "value": "* Service defines a gRPC service for interacting with transactions. "
            }
        ]
    })
})
