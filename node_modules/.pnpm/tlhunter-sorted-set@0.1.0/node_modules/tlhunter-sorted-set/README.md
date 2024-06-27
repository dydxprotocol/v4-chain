# tlhunter-sorted-set

A JavaScript implementation of Redis' [Sorted Sets](https://redis.io/commands#sorted_set). Keeps a collection of "members" in order based on their score. Uses skip lists under the hood, [like Redis does](http://stackoverflow.com/a/9626334/638546).

This is a fork of the brilliant but abandoned [redis-sorted-map](https://www.npmjs.com/package/redis-sorted-map) package by [Akseli Palén](https://github.com/axelpale) which was itself a fork of the brilliant but abandoned [sorted-map](https://www.npmjs.com/package/sorted-map) package by [Eli Skeggs](https://github.com/skeggse).

![A Skip List](/doc/skip-list.png?raw=true)

Image: The skip list data structure allows search, insert, and removal in O(log(n)) time in average.

## Install

```sh
$ npm install tlhunter-sorted-set
```


## Test

Run any of the following:

```sh
$ npm test
```

_Note:_ remember to `npm install`!


## API

The API mostly follows Redis' [Sorted Set Commands](https://redis.io/commands#sorted_set), with a few additional methods such as `.has(member)`.

Members can be strings, symbols, objects, or really any primitive value.

```js
const SortedSet = require('tlhunter-sorted-set');

const z = new SortedSet();

// average O(log(N))
z.add('Terminator', 8.0); // => null
z.add('District 9', 8.0); // => null
z.add('Ex Machina', 0.7); // => null
z.add('Ex Machina', 7.7); // => 0.7
// alias
z.set('The Matrix', 8.7); // => null

// average O(1)
z.has('Terminator'); // => true
z.has('Blade Runner'); // => false

// average O(1)
z.score('Ex Machina'); // => 7.7
z.score('Blade Runner'); // => null
// alias
z.get('The Matrix'); // => 8.7

// average O(log(N))
z.rem('Ex Machina'); // => 7.7
// average O(1)
z.rem('Ex Machina'); // => null
// alias
z.del('Ex Machina'); // => null

// average O(log(N)+M) where M is the number of elements between min and max
z.rangeByScore(7, 8);
// => ['Ex Machina', 'District 9', 'Terminator']
z.rangeByScore(8); // [8.0-∞)
// => ['District 9', 'Terminator', 'The Matrix']
z.rangeByScore(8, null, { withScores: true });
// => [['District 9', 8.0], ['Terminator', 8.0], ['The Matrix', 8.7]]

// average O(log(N)+log(M)) where M as in rangeByScore
z.count(7, 8); // => 3

// average O(log(N))
z.rank('Ex Machina'); // => 0
z.rank('Terminator'); // => 2
z.rank('Blade Runner'); // => null

// average O(log(N)+M) where M as in range
z.range(0, 2);
// => ['Ex Machina', 'District 9', 'Terminator']
z.range(0, 2, { withScores: true });
// => [['Ex Machina', 7.7],
//     ['District 9', 8],
//     ['Terminator', 8]]
z.range(-1); // => ['The Matrix']
// almost alias
z.slice(0, 3);
// => ['Ex Machina', 'District 9', 'Terminator']

// Set cardinality (number of elements)
// average O(1)
z.card(); // => 4
// alias
z.length // => 4

```


## Intersection

```js
const a = new SortedSet(), b = new SortedSet();

a.add('5a600e10', 16);
a.add('5a600e12', 10);
a.add('5a600e14', 9);
a.add('5a600e15', 14);
a.add('5a600e17', 20);
a.add('5a600e18', 13);
a.add('5a600e19', 15);
a.add('5a600e1a', 19);
a.add('5a600e1b', 7);
a.add('5a600e1c', 13);
a.add('5a600e1e', 10);

b.add('5a600e10', 0);
b.add('5a600e11', 15);
b.add('5a600e13', 5);
b.add('5a600e14', 3);
b.add('5a600e15', 14);
b.add('5a600e17', 12);
b.add('5a600e19', 12);
b.add('5a600e1b', 16);
b.add('5a600e1c', 12);
b.add('5a600e1d', 17);
b.add('5a600e1f', 3);

SortedSet.intersect(a, b);
// => ['5a600e10', '5a600e14', '5a600e17', '5a600e19', '5a600e1c', '5a600e15', '5a600e1b']

SortedSet.intersect(b, a);
// => ['5a600e1b', '5a600e14', '5a600e1c', '5a600e15', '5a600e19', '5a600e10', '5a600e17']

// works, but not preferred
a.intersect(b);
// => ['5a600e10', '5a600e14', '5a600e17', '5a600e19', '5a600e1c', '5a600e15', '5a600e1b']

const c = new SortedSet();

c.add('5a600e10', 7);
c.add('5a600e12', 20);
c.add('5a600e13', 9);
c.add('5a600e14', 19);
c.add('5a600e16', 19);
c.add('5a600e17', 1);
c.add('5a600e18', 18);
c.add('5a600e1a', 6);
c.add('5a600e1c', 15);
c.add('5a600e1f', 4);

// for best performance, the smallest set should be first
SortedSet.intersect(c, a, b);
// => ['5a600e10', '5a600e14', '5a600e17', '5a600e1c']
```


## Unique

You can enable unique values with the unique option, which causes `set` to throw an error if the value provided already belongs to a different key.

```js
const z = new SortedSet({unique: true});

z.add('5a600e10', 16);
z.add('5a600e11', 6);
z.add('5a600e12', 17);
z.add('5a600e13', 11);
z.add('5a600e14', 14);
z.add('5a600e15', 19);
z.add('5a600e16', 3);
z.add('5a600e17', 12);
z.add('5a600e18', 10);

// currently O(log(N)) because it needs to attempt to insert the value
z.add('5a600e19', 11); // throws
z.add('5a600e14', 14); // => 14
```


## Licence

MIT
