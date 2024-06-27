/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
"use strict";
var $protobuf = require("protobufjs/minimal");
// Common aliases
var $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;
// Exported root namespace
var $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});
$root.DDSketch = (function () {
    /**
     * Properties of a DDSketch.
     * @exports IDDSketch
     * @interface IDDSketch
     * @property {IIndexMapping|null} [mapping] DDSketch mapping
     * @property {IStore|null} [positiveValues] DDSketch positiveValues
     * @property {IStore|null} [negativeValues] DDSketch negativeValues
     * @property {number|null} [zeroCount] DDSketch zeroCount
     */
    /**
     * Constructs a new DDSketch.
     * @exports DDSketch
     * @classdesc Represents a DDSketch.
     * @implements IDDSketch
     * @constructor
     * @param {IDDSketch=} [properties] Properties to set
     */
    function DDSketch(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }
    /**
     * DDSketch mapping.
     * @member {IIndexMapping|null|undefined} mapping
     * @memberof DDSketch
     * @instance
     */
    DDSketch.prototype.mapping = null;
    /**
     * DDSketch positiveValues.
     * @member {IStore|null|undefined} positiveValues
     * @memberof DDSketch
     * @instance
     */
    DDSketch.prototype.positiveValues = null;
    /**
     * DDSketch negativeValues.
     * @member {IStore|null|undefined} negativeValues
     * @memberof DDSketch
     * @instance
     */
    DDSketch.prototype.negativeValues = null;
    /**
     * DDSketch zeroCount.
     * @member {number} zeroCount
     * @memberof DDSketch
     * @instance
     */
    DDSketch.prototype.zeroCount = 0;
    /**
     * Creates a new DDSketch instance using the specified properties.
     * @function create
     * @memberof DDSketch
     * @static
     * @param {IDDSketch=} [properties] Properties to set
     * @returns {DDSketch} DDSketch instance
     */
    DDSketch.create = function create(properties) {
        return new DDSketch(properties);
    };
    /**
     * Encodes the specified DDSketch message. Does not implicitly {@link DDSketch.verify|verify} messages.
     * @function encode
     * @memberof DDSketch
     * @static
     * @param {IDDSketch} message DDSketch message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    DDSketch.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.mapping != null && Object.hasOwnProperty.call(message, "mapping"))
            $root.IndexMapping.encode(message.mapping, writer.uint32(/* id 1, wireType 2 =*/ 10).fork()).ldelim();
        if (message.positiveValues != null && Object.hasOwnProperty.call(message, "positiveValues"))
            $root.Store.encode(message.positiveValues, writer.uint32(/* id 2, wireType 2 =*/ 18).fork()).ldelim();
        if (message.negativeValues != null && Object.hasOwnProperty.call(message, "negativeValues"))
            $root.Store.encode(message.negativeValues, writer.uint32(/* id 3, wireType 2 =*/ 26).fork()).ldelim();
        if (message.zeroCount != null && Object.hasOwnProperty.call(message, "zeroCount"))
            writer.uint32(/* id 4, wireType 1 =*/ 33).double(message.zeroCount);
        return writer;
    };
    /**
     * Encodes the specified DDSketch message, length delimited. Does not implicitly {@link DDSketch.verify|verify} messages.
     * @function encodeDelimited
     * @memberof DDSketch
     * @static
     * @param {IDDSketch} message DDSketch message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    DDSketch.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };
    /**
     * Decodes a DDSketch message from the specified reader or buffer.
     * @function decode
     * @memberof DDSketch
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {DDSketch} DDSketch
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    DDSketch.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.DDSketch();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.mapping = $root.IndexMapping.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.positiveValues = $root.Store.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.negativeValues = $root.Store.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.zeroCount = reader.double();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    };
    /**
     * Decodes a DDSketch message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof DDSketch
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {DDSketch} DDSketch
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    DDSketch.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };
    /**
     * Verifies a DDSketch message.
     * @function verify
     * @memberof DDSketch
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    DDSketch.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.mapping != null && message.hasOwnProperty("mapping")) {
            var error = $root.IndexMapping.verify(message.mapping);
            if (error)
                return "mapping." + error;
        }
        if (message.positiveValues != null && message.hasOwnProperty("positiveValues")) {
            var error = $root.Store.verify(message.positiveValues);
            if (error)
                return "positiveValues." + error;
        }
        if (message.negativeValues != null && message.hasOwnProperty("negativeValues")) {
            var error = $root.Store.verify(message.negativeValues);
            if (error)
                return "negativeValues." + error;
        }
        if (message.zeroCount != null && message.hasOwnProperty("zeroCount"))
            if (typeof message.zeroCount !== "number")
                return "zeroCount: number expected";
        return null;
    };
    /**
     * Creates a DDSketch message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof DDSketch
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {DDSketch} DDSketch
     */
    DDSketch.fromObject = function fromObject(object) {
        if (object instanceof $root.DDSketch)
            return object;
        var message = new $root.DDSketch();
        if (object.mapping != null) {
            if (typeof object.mapping !== "object")
                throw TypeError(".DDSketch.mapping: object expected");
            message.mapping = $root.IndexMapping.fromObject(object.mapping);
        }
        if (object.positiveValues != null) {
            if (typeof object.positiveValues !== "object")
                throw TypeError(".DDSketch.positiveValues: object expected");
            message.positiveValues = $root.Store.fromObject(object.positiveValues);
        }
        if (object.negativeValues != null) {
            if (typeof object.negativeValues !== "object")
                throw TypeError(".DDSketch.negativeValues: object expected");
            message.negativeValues = $root.Store.fromObject(object.negativeValues);
        }
        if (object.zeroCount != null)
            message.zeroCount = Number(object.zeroCount);
        return message;
    };
    /**
     * Creates a plain object from a DDSketch message. Also converts values to other types if specified.
     * @function toObject
     * @memberof DDSketch
     * @static
     * @param {DDSketch} message DDSketch
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    DDSketch.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.mapping = null;
            object.positiveValues = null;
            object.negativeValues = null;
            object.zeroCount = 0;
        }
        if (message.mapping != null && message.hasOwnProperty("mapping"))
            object.mapping = $root.IndexMapping.toObject(message.mapping, options);
        if (message.positiveValues != null && message.hasOwnProperty("positiveValues"))
            object.positiveValues = $root.Store.toObject(message.positiveValues, options);
        if (message.negativeValues != null && message.hasOwnProperty("negativeValues"))
            object.negativeValues = $root.Store.toObject(message.negativeValues, options);
        if (message.zeroCount != null && message.hasOwnProperty("zeroCount"))
            object.zeroCount = options.json && !isFinite(message.zeroCount) ? String(message.zeroCount) : message.zeroCount;
        return object;
    };
    /**
     * Converts this DDSketch to JSON.
     * @function toJSON
     * @memberof DDSketch
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    DDSketch.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };
    return DDSketch;
})();
$root.IndexMapping = (function () {
    /**
     * Properties of an IndexMapping.
     * @exports IIndexMapping
     * @interface IIndexMapping
     * @property {number|null} [gamma] IndexMapping gamma
     * @property {number|null} [indexOffset] IndexMapping indexOffset
     * @property {IndexMapping.Interpolation|null} [interpolation] IndexMapping interpolation
     */
    /**
     * Constructs a new IndexMapping.
     * @exports IndexMapping
     * @classdesc Represents an IndexMapping.
     * @implements IIndexMapping
     * @constructor
     * @param {IIndexMapping=} [properties] Properties to set
     */
    function IndexMapping(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }
    /**
     * IndexMapping gamma.
     * @member {number} gamma
     * @memberof IndexMapping
     * @instance
     */
    IndexMapping.prototype.gamma = 0;
    /**
     * IndexMapping indexOffset.
     * @member {number} indexOffset
     * @memberof IndexMapping
     * @instance
     */
    IndexMapping.prototype.indexOffset = 0;
    /**
     * IndexMapping interpolation.
     * @member {IndexMapping.Interpolation} interpolation
     * @memberof IndexMapping
     * @instance
     */
    IndexMapping.prototype.interpolation = 0;
    /**
     * Creates a new IndexMapping instance using the specified properties.
     * @function create
     * @memberof IndexMapping
     * @static
     * @param {IIndexMapping=} [properties] Properties to set
     * @returns {IndexMapping} IndexMapping instance
     */
    IndexMapping.create = function create(properties) {
        return new IndexMapping(properties);
    };
    /**
     * Encodes the specified IndexMapping message. Does not implicitly {@link IndexMapping.verify|verify} messages.
     * @function encode
     * @memberof IndexMapping
     * @static
     * @param {IIndexMapping} message IndexMapping message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    IndexMapping.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.gamma != null && Object.hasOwnProperty.call(message, "gamma"))
            writer.uint32(/* id 1, wireType 1 =*/ 9).double(message.gamma);
        if (message.indexOffset != null && Object.hasOwnProperty.call(message, "indexOffset"))
            writer.uint32(/* id 2, wireType 1 =*/ 17).double(message.indexOffset);
        if (message.interpolation != null && Object.hasOwnProperty.call(message, "interpolation"))
            writer.uint32(/* id 3, wireType 0 =*/ 24).int32(message.interpolation);
        return writer;
    };
    /**
     * Encodes the specified IndexMapping message, length delimited. Does not implicitly {@link IndexMapping.verify|verify} messages.
     * @function encodeDelimited
     * @memberof IndexMapping
     * @static
     * @param {IIndexMapping} message IndexMapping message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    IndexMapping.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };
    /**
     * Decodes an IndexMapping message from the specified reader or buffer.
     * @function decode
     * @memberof IndexMapping
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {IndexMapping} IndexMapping
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    IndexMapping.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.IndexMapping();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.gamma = reader.double();
                    break;
                case 2:
                    message.indexOffset = reader.double();
                    break;
                case 3:
                    message.interpolation = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    };
    /**
     * Decodes an IndexMapping message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof IndexMapping
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {IndexMapping} IndexMapping
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    IndexMapping.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };
    /**
     * Verifies an IndexMapping message.
     * @function verify
     * @memberof IndexMapping
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    IndexMapping.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.gamma != null && message.hasOwnProperty("gamma"))
            if (typeof message.gamma !== "number")
                return "gamma: number expected";
        if (message.indexOffset != null && message.hasOwnProperty("indexOffset"))
            if (typeof message.indexOffset !== "number")
                return "indexOffset: number expected";
        if (message.interpolation != null && message.hasOwnProperty("interpolation"))
            switch (message.interpolation) {
                default:
                    return "interpolation: enum value expected";
                case 0:
                case 1:
                case 2:
                case 3:
                    break;
            }
        return null;
    };
    /**
     * Creates an IndexMapping message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof IndexMapping
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {IndexMapping} IndexMapping
     */
    IndexMapping.fromObject = function fromObject(object) {
        if (object instanceof $root.IndexMapping)
            return object;
        var message = new $root.IndexMapping();
        if (object.gamma != null)
            message.gamma = Number(object.gamma);
        if (object.indexOffset != null)
            message.indexOffset = Number(object.indexOffset);
        switch (object.interpolation) {
            case "NONE":
            case 0:
                message.interpolation = 0;
                break;
            case "LINEAR":
            case 1:
                message.interpolation = 1;
                break;
            case "QUADRATIC":
            case 2:
                message.interpolation = 2;
                break;
            case "CUBIC":
            case 3:
                message.interpolation = 3;
                break;
        }
        return message;
    };
    /**
     * Creates a plain object from an IndexMapping message. Also converts values to other types if specified.
     * @function toObject
     * @memberof IndexMapping
     * @static
     * @param {IndexMapping} message IndexMapping
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    IndexMapping.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.gamma = 0;
            object.indexOffset = 0;
            object.interpolation = options.enums === String ? "NONE" : 0;
        }
        if (message.gamma != null && message.hasOwnProperty("gamma"))
            object.gamma = options.json && !isFinite(message.gamma) ? String(message.gamma) : message.gamma;
        if (message.indexOffset != null && message.hasOwnProperty("indexOffset"))
            object.indexOffset = options.json && !isFinite(message.indexOffset) ? String(message.indexOffset) : message.indexOffset;
        if (message.interpolation != null && message.hasOwnProperty("interpolation"))
            object.interpolation = options.enums === String ? $root.IndexMapping.Interpolation[message.interpolation] : message.interpolation;
        return object;
    };
    /**
     * Converts this IndexMapping to JSON.
     * @function toJSON
     * @memberof IndexMapping
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    IndexMapping.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };
    /**
     * Interpolation enum.
     * @name IndexMapping.Interpolation
     * @enum {number}
     * @property {number} NONE=0 NONE value
     * @property {number} LINEAR=1 LINEAR value
     * @property {number} QUADRATIC=2 QUADRATIC value
     * @property {number} CUBIC=3 CUBIC value
     */
    IndexMapping.Interpolation = (function () {
        var valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "NONE"] = 0;
        values[valuesById[1] = "LINEAR"] = 1;
        values[valuesById[2] = "QUADRATIC"] = 2;
        values[valuesById[3] = "CUBIC"] = 3;
        return values;
    })();
    return IndexMapping;
})();
$root.Store = (function () {
    /**
     * Properties of a Store.
     * @exports IStore
     * @interface IStore
     * @property {Object.<string,number>|null} [binCounts] Store binCounts
     * @property {Array.<number>|null} [contiguousBinCounts] Store contiguousBinCounts
     * @property {number|null} [contiguousBinIndexOffset] Store contiguousBinIndexOffset
     */
    /**
     * Constructs a new Store.
     * @exports Store
     * @classdesc Represents a Store.
     * @implements IStore
     * @constructor
     * @param {IStore=} [properties] Properties to set
     */
    function Store(properties) {
        this.binCounts = {};
        this.contiguousBinCounts = [];
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }
    /**
     * Store binCounts.
     * @member {Object.<string,number>} binCounts
     * @memberof Store
     * @instance
     */
    Store.prototype.binCounts = $util.emptyObject;
    /**
     * Store contiguousBinCounts.
     * @member {Array.<number>} contiguousBinCounts
     * @memberof Store
     * @instance
     */
    Store.prototype.contiguousBinCounts = $util.emptyArray;
    /**
     * Store contiguousBinIndexOffset.
     * @member {number} contiguousBinIndexOffset
     * @memberof Store
     * @instance
     */
    Store.prototype.contiguousBinIndexOffset = 0;
    /**
     * Creates a new Store instance using the specified properties.
     * @function create
     * @memberof Store
     * @static
     * @param {IStore=} [properties] Properties to set
     * @returns {Store} Store instance
     */
    Store.create = function create(properties) {
        return new Store(properties);
    };
    /**
     * Encodes the specified Store message. Does not implicitly {@link Store.verify|verify} messages.
     * @function encode
     * @memberof Store
     * @static
     * @param {IStore} message Store message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Store.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.binCounts != null && Object.hasOwnProperty.call(message, "binCounts"))
            for (var keys = Object.keys(message.binCounts), i = 0; i < keys.length; ++i)
                writer.uint32(/* id 1, wireType 2 =*/ 10).fork().uint32(/* id 1, wireType 0 =*/ 8).sint32(keys[i]).uint32(/* id 2, wireType 1 =*/ 17).double(message.binCounts[keys[i]]).ldelim();
        if (message.contiguousBinCounts != null && message.contiguousBinCounts.length) {
            writer.uint32(/* id 2, wireType 2 =*/ 18).fork();
            for (var i = 0; i < message.contiguousBinCounts.length; ++i)
                writer.double(message.contiguousBinCounts[i]);
            writer.ldelim();
        }
        if (message.contiguousBinIndexOffset != null && Object.hasOwnProperty.call(message, "contiguousBinIndexOffset"))
            writer.uint32(/* id 3, wireType 0 =*/ 24).sint32(message.contiguousBinIndexOffset);
        return writer;
    };
    /**
     * Encodes the specified Store message, length delimited. Does not implicitly {@link Store.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Store
     * @static
     * @param {IStore} message Store message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Store.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };
    /**
     * Decodes a Store message from the specified reader or buffer.
     * @function decode
     * @memberof Store
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Store} Store
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Store.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.Store(), key, value;
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    if (message.binCounts === $util.emptyObject)
                        message.binCounts = {};
                    var end2 = reader.uint32() + reader.pos;
                    key = 0;
                    value = 0;
                    while (reader.pos < end2) {
                        var tag2 = reader.uint32();
                        switch (tag2 >>> 3) {
                            case 1:
                                key = reader.sint32();
                                break;
                            case 2:
                                value = reader.double();
                                break;
                            default:
                                reader.skipType(tag2 & 7);
                                break;
                        }
                    }
                    message.binCounts[key] = value;
                    break;
                case 2:
                    if (!(message.contiguousBinCounts && message.contiguousBinCounts.length))
                        message.contiguousBinCounts = [];
                    if ((tag & 7) === 2) {
                        var end2 = reader.uint32() + reader.pos;
                        while (reader.pos < end2)
                            message.contiguousBinCounts.push(reader.double());
                    }
                    else
                        message.contiguousBinCounts.push(reader.double());
                    break;
                case 3:
                    message.contiguousBinIndexOffset = reader.sint32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    };
    /**
     * Decodes a Store message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Store
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Store} Store
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Store.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };
    /**
     * Verifies a Store message.
     * @function verify
     * @memberof Store
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Store.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.binCounts != null && message.hasOwnProperty("binCounts")) {
            if (!$util.isObject(message.binCounts))
                return "binCounts: object expected";
            var key = Object.keys(message.binCounts);
            for (var i = 0; i < key.length; ++i) {
                if (!$util.key32Re.test(key[i]))
                    return "binCounts: integer key{k:sint32} expected";
                if (typeof message.binCounts[key[i]] !== "number")
                    return "binCounts: number{k:sint32} expected";
            }
        }
        if (message.contiguousBinCounts != null && message.hasOwnProperty("contiguousBinCounts")) {
            if (!Array.isArray(message.contiguousBinCounts))
                return "contiguousBinCounts: array expected";
            for (var i = 0; i < message.contiguousBinCounts.length; ++i)
                if (typeof message.contiguousBinCounts[i] !== "number")
                    return "contiguousBinCounts: number[] expected";
        }
        if (message.contiguousBinIndexOffset != null && message.hasOwnProperty("contiguousBinIndexOffset"))
            if (!$util.isInteger(message.contiguousBinIndexOffset))
                return "contiguousBinIndexOffset: integer expected";
        return null;
    };
    /**
     * Creates a Store message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Store
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Store} Store
     */
    Store.fromObject = function fromObject(object) {
        if (object instanceof $root.Store)
            return object;
        var message = new $root.Store();
        if (object.binCounts) {
            if (typeof object.binCounts !== "object")
                throw TypeError(".Store.binCounts: object expected");
            message.binCounts = {};
            for (var keys = Object.keys(object.binCounts), i = 0; i < keys.length; ++i)
                message.binCounts[keys[i]] = Number(object.binCounts[keys[i]]);
        }
        if (object.contiguousBinCounts) {
            if (!Array.isArray(object.contiguousBinCounts))
                throw TypeError(".Store.contiguousBinCounts: array expected");
            message.contiguousBinCounts = [];
            for (var i = 0; i < object.contiguousBinCounts.length; ++i)
                message.contiguousBinCounts[i] = Number(object.contiguousBinCounts[i]);
        }
        if (object.contiguousBinIndexOffset != null)
            message.contiguousBinIndexOffset = object.contiguousBinIndexOffset | 0;
        return message;
    };
    /**
     * Creates a plain object from a Store message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Store
     * @static
     * @param {Store} message Store
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Store.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.arrays || options.defaults)
            object.contiguousBinCounts = [];
        if (options.objects || options.defaults)
            object.binCounts = {};
        if (options.defaults)
            object.contiguousBinIndexOffset = 0;
        var keys2;
        if (message.binCounts && (keys2 = Object.keys(message.binCounts)).length) {
            object.binCounts = {};
            for (var j = 0; j < keys2.length; ++j)
                object.binCounts[keys2[j]] = options.json && !isFinite(message.binCounts[keys2[j]]) ? String(message.binCounts[keys2[j]]) : message.binCounts[keys2[j]];
        }
        if (message.contiguousBinCounts && message.contiguousBinCounts.length) {
            object.contiguousBinCounts = [];
            for (var j = 0; j < message.contiguousBinCounts.length; ++j)
                object.contiguousBinCounts[j] = options.json && !isFinite(message.contiguousBinCounts[j]) ? String(message.contiguousBinCounts[j]) : message.contiguousBinCounts[j];
        }
        if (message.contiguousBinIndexOffset != null && message.hasOwnProperty("contiguousBinIndexOffset"))
            object.contiguousBinIndexOffset = message.contiguousBinIndexOffset;
        return object;
    };
    /**
     * Converts this Store to JSON.
     * @function toJSON
     * @memberof Store
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Store.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };
    return Store;
})();
module.exports = $root;
