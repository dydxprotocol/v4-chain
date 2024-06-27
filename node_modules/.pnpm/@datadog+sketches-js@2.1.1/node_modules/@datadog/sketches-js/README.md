# sketches-js

![Continuous Integration](https://github.com/DataDog/sketches-js/workflows/Continuous%20Integration/badge.svg) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This repo contains the TypeScript implementation of the distributed quantile sketch algorithm [DDSketch](http://www.vldb.org/pvldb/vol12/p2195-masson.pdf). DDSketch is mergeable, meaning that multiple sketches from distributed systems can be combined in a central node.

## Installation

The package is under [@datadog/sketches-js](https://www.npmjs.com/package/@datadog/sketches-js) and can be installed through NPM or Yarn:

```sh
# NPM
npm install @datadog/sketches-js

# Yarn
yarn add @datadog/sketches-js
```

When using Protobuf serialization, the [protobufjs](https://www.npmjs.com/package/protobufjs) module must also be installed manually:

```sh
# NPM
npm install protobufjs

# Yarn
yarn add protobufjs
```

## Usage

### Initialize a sketch

To initialize a sketch with the default parameters:

```js
import { DDSketch } from '@datadog/sketches-js'; // or const { DDSketch } = require('@datadog/sketches-js');
const sketch = new DDSketch();
```

#### Modify the `relativeAccuracy`

If you want more granular control over how accurate the sketch's results will be, you can pass a `relativeAccuracy` parameter when initializing a sketch.

Whereas other histograms use _rank error_ guarantees (i.e. retrieving the p99 of the histogram will give you a value between p98.9 and p99.1), DDSketch uses a _relative error_ guarantee (if the actual value at p99 is 100, the value will be between 99 and 101 for a `relativeAccuracy` of 0.01).

This property makes DDSketch especially useful for long-tailed distributions of data, like measurements of latency.

```js
import { DDSketch } from '@datadog/sketches-js';

const sketch = new DDSketch({
  relativeAccuracy: 0.01, // `relativeAccuracy` must be between 0 and 1
});
```

### Add values to a sketch

To add a number to a sketch, call `sketch.accept(value)`. Both positive and negative numbers are supported.

```js
const measurementOne = 1607374726;
const measurementTwo = 0;
const measurementThree = -3.1415;

sketch.accept(measurementOne);
sketch.accept(measurementTwo);
sketch.accept(measurementThree);
```

### Retrieve measurements from the sketch

To retrieve measurements from a sketch, use `sketch.getValueAtQuantile(quantile)`. Any number between 0 and 1 (inclusive) can be used as a quantile.

Additionally, common summary statistics are available such as `sketch.min`, `sketch.max`, `sketch.sum`, and `sketch.count`:

```js
const measurementOne = 1607374726;
const measurementTwo = 0;
const measurementThree = -3.1415;

sketch.accept(measurementOne);
sketch.accept(measurementTwo);
sketch.accept(measurementThree);

sketch.getValueAtQuantile(0)     // -3.1415
sketch.getValueAtQuantile(0.5)   // 0
sketch.getValueAtQuantile(0.99)  // 1607374726
sketch.getValueAtQuantile(1)     // 1607374726

sketch.min                       // -3.1415
sketch.max                       // 1607374726
sketch.count                     // 3
sketch.sum                       // 1607374722.86
```

### Merge multiple sketches

Independent sketches can be merged together, provided that they were initialized with the same `relativeAccuracy`. This allows collecting and transmitting measurements in a distributed manner, and merging their results together while preserving the `relativeAccuracy` guarantee.

```js
import { DDSketch } from '@datadog/sketches-js';

const sketch1 = new DDSketch();
const sketch2 = new DDSketch();

[1,2,3,4,5].forEach(value => sketch1.accept(value));
[6,7,8,9,10].forEach(value => sketch2.accept(value));

// `sketch2` is merged into `sketch1`, without modifying `sketch2`
sketch1.merge(sketch2);

sketch1.getValueAtQuantile(1) // 10
```

## References
* [DDSketch: A Fast and Fully-Mergeable Quantile Sketch with Relative-Error Guarantees](http://www.vldb.org/pvldb/vol12/p2195-masson.pdf). Charles Masson, Jee E. Rim and Homin K. Lee. 2019.
* Java implementation: [https://github.com/DataDog/sketches-java](https://github.com/DataDog/sketches-java)
* Go implementation: [https://github.com/DataDog/sketches-go](https://github.com/DataDog/sketches-go)
* Python implementation: [https://github.com/DataDog/sketches-py](https://github.com/DataDog/sketches-py)
