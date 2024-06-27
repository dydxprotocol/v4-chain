'use strict';

/*
 * TODO
 * rename most instances of "key" to "member"
 * rename most instances of "value" to "score"
 * rename private _foo methods to #foo
 */

const intersect = require('./intersect.js');

const slice = Array.prototype.slice;
const P = 1 / Math.E;

class SortedSet {
  constructor(options = {}) {
    this._unique = !!options.unique;
    this.empty();
  }

  static intersect () {
    return intersect.call(SortedSet, slice.call(arguments));
  }

  add(key, value) {
    let current;

    if (value == null) {
      return this.rem(key);
    }

    current = this._map.get(key);

    if (current !== undefined) {
      if (value === current) {
        return current;
      }
      this._remove(key, current);
    }

    let node = this._insert(key, value);
    if (!node) {
      current === undefined || this._insert(key, current);
      // TODO: can we defer _remove until after insert?
      throw new Error('unique constraint violated');
    }

    this._map.set(key, value);
    return current === undefined ? null : current;
  }

  card() {
    // Returns the sorted set cardinality (number of elements)
    if (this.length) {
      return this.length;
    }
    return 0;
  }

  count(min, max) {
    if (!this.length) {
      return 0;
    }

    if (min == null) {
      min = -Infinity;
    }
    if (max == null) {
      max = Infinity;
    }

    if (min <= this._head.next[0].next.value && max >= this._tail.value) {
      return this.length;
    }

    if (max < min || min > this._tail.value || max < this._head.next[0].next.value) {
      return 0;
    }

    let i;
    let node = this._first(min);
    let count = 0;

    if (!node) {
      return 0;
    }

    for (i = node.next.length - 1; i >= 0; i -= 1) {
      while (node.next[i].next && node.next[i].next.value <= max) {
        count += node.next[i].span;
        node = node.next[i].next;
      }
    }

    // feels hacky and error prone
    return count && count + 1;
  }

  del(key) {
    // Alias for .rem
    return this.rem(key);
  }

  empty() {
    this.length = 0;
    this._level = 1;
    this._map = new Map();
    this._head = new Node(32, null, 0);
    this._tail = null;

    for (let i = 0; i < 32; i += 1) {
      // hrm
      this._head.next[i] = new Level(null, 0);
    }
  }

  get(key) {
    // Alias for
    return this.score(key);
  }

  has(key) {
    return this._map.has(key)
  }

  incrBy(increment, key) {
    // Increases the score of the member specified by key.
    // If member does not exist, a new member is created with
    // the increment as its score.
    //
    // Parameters:
    //   increment
    //     number
    //   key
    //     string
    //
    // Return
    //   number, the new score of the member
    //
    // TODO: Shortcut, could optimize to half

    let score = this.score(key);

    if (score) {
      this.add(key, score + increment);
      return score + increment;
    }

    this.add(key, increment);
    return increment;
  }

  intersect() {
    // intersect values
    let maps = slice.call(arguments);
    maps.unshift(this);
    return intersect.call(this, maps);
  }

  intersectKeys() {
    let maps = slice.call(arguments);
    maps.unshift(this);
    return intersectKeys.call(this, maps);
  }

  keys() {
    if (!this.length) {
      return [];
    }

    let i;
    let array = new Array(this.length);
    let node = this._head.next[0].next;

    for (i = 0; node; node = node.next[0].next) {
      array[i] = node.key;
      i += 1;
    }

    return array;
  }

  range (start, stop, options) {
    // Parameters:
    //   start
    //     inclusive
    //   stop
    //     inclusive
    //   options (optional)
    //     withScores (optional, default to false)
    //
    // Return:
    //   an array

    if (this.length === 0) {
      return [];
    }

    if (start == null) {
      start = 0;
    } else if (start < 0) {
      start = Math.max(this.length + start, 0);
    }

    if (stop == null) {
      stop = this.length - 1;
    } else if (stop < 0) {
      stop = this.length + stop;
    }

    if (start > stop || start >= this.length) {
      return [];
    }

    if (stop >= this.length) {
      stop = this.length - 1;
    }

    if (typeof options !== 'object') {
      options = {
        withScores: false,
      };
    }

    let i = 0;
    let length = stop - start + 1;
    let result
    try {
      result = new Array(length);
    } catch(e) {
      console.log('start', start);
      console.log('stop', stop);
      console.log('Invalid length', length);
      throw e;
    }

    let node = start > 0 ? this._get(start) : this._head.next[0].next;

    if (options.withScores) {
      for (; length--; node = node.next[0].next) {
        result[i] = [node.key, node.value];
        i += 1;
      }
    } else {
      for (; length--; node = node.next[0].next) {
        result[i] = node.key;
        i += 1;
      }
    }

    return result;
  };

  rangeByScore(min, max, options) {
    // Return members with score within inclusive range [min, max].
    //
    // Parameters:
    //   min (number)
    //   max (number)
    //   options (object, optional)
    //     withScores (bool, optional, default false)

    if (!this.length) {
      return [];
    }

    if (typeof options !== 'object') {
      options = {
        withScores: false,
      };
    }

    if (min == null) {
      min = -Infinity;
    }
    if (max == null) {
      max = Infinity;
    }

    if (min <= this._head.next[0].next.value && max >= this._tail.value) {
      return this.toArray({ withScores: options.withScores });
    }

    if (max < min || min > this._tail.value || max < this._head.next[0].next.value) {
      return [];
    }

    let node = this._first(min);
    let result = [];

    if (options.withScores) {
      for (; node && node.value <= max; node = node.next[0].next) {
        result.push([node.key, node.value]);
      }
    } else {
      for (; node && node.value <= max; node = node.next[0].next) {
        result.push(node.key);
      }
    }

    return result;
  }

  rank(key) {
    // Rank of key, ordered by value.
    //
    // Return
    //   integer
    //     if member exists
    ///  null
    //     if member does not exist

    let value = this._map.get(key);

    if (value === undefined) {
      return null;
    }

    let i;
    let node = this._head;
    let next = null;
    let rank = -1;

    for (i = this._level - 1; i >= 0; i -= 1) {
      while ((next = node.next[i].next) && (next.value < value || (next.value === value && next.key <= key))) {
        rank += node.next[i].span;
        node = next;
      }
      if (node.key && node.key === key) {
        return rank;
      }
    }

    return null;
  }

  rem(key) {
    // Remove single member by key.
    //
    // Return
    //   value of the removed key
    //   or null if key does not exist.

    let value = this._map.get(key);
    if (value !== undefined) {
      this._remove(key, value);
      this._map.delete(key);
      return value;
    }
    return null;
  }

  remRangeByRank (start, end) {
    // Parameters:
    //   start
    //     inclusive
    //   end
    //     exclusive
    //
    // Return
    //   positive integer, the number of removed keys.

    let len = this.length;

    if (!len) {
      return 0;
    }

    if (start == null) {
      start = 0;
    } else if (start < 0) {
      start = Math.max(len + start, 0);
    }

    if (end == null) {
      end = len;
    } else if (end < 0) {
      end = len + end;
    }

    if (start > end || start >= len) {
      return 0;
    }
    if (end > len) {
      end = len;
    }

    if (start === 0 && end === len) {
      this.empty();
      return len;
    }

    let node = this._head;
    let update = new Array(32)
    let result, i, next;
    let traversed = -1;

    for (i = this._level - 1; i >= 0; i -= 1) {
      while ((next = node.next[i].next) && (traversed + node.next[i].span) < start) {
        traversed += node.next[i].span;
        node = next;
      }
      update[i] = node;
    }

    let removed = 0;
    traversed += 1;
    node = node.next[0].next;

    while (node && traversed < end) {
      next = node.next[0].next;
      this._removeNode(node, update);
      this._map.delete(node.key);
      removed += 1;
      traversed += 1;
      node = next;
    }

    this.length -= removed;
    return removed;
  }

  remRangeByScore(min, max) {
    // Remove members with value between min and max (inclusive).
    //
    // Return
    //   positive integer, the number of removed elements.

    let result;
    let removed = 0;

    if (!this.length) {
      return 0;
    }

    if (min == null) {
      min = -Infinity;
    }
    if (max == null) {
      max = Infinity;
    }

    if (min <= this._head.next[0].next.value && max >= this._tail.value) {
      removed = this.length;
      this.empty();
      return removed;
    }

    let next, i;
    let node = this._head;
    let update = new Array(32);

    for (i = this._level - 1; i >= 0; i -= 1) {
      while ((next = node.next[i].next) && next.value < min) {
        node = next;
      }
      update[i] = node;
    }
    node = node.next[0].next;

    while (node && node.value <= max) {
      next = node.next[0].next;
      this._removeNode(node, update);
      this._map.delete(node.key);
      removed += 1;
      node = next;
    }

    this.length -= removed;
    return removed;
  }

  score (member) {
    // Return
    //   number, the score of member in the sorted set.
    //   null, if member does not exist in the sorted set.
    let score = this._map.get(member);
    return score === undefined ? null : score;
  }

  set(key, value) {
    // Alias for
    return this.add(key, value);
  }

  slice(start, end, options) {
    // Almost alias for range. Only difference is that
    // the end is exclusive i.e. not included in the range.
    if (typeof end === 'number' && end !== 0) {
      end -= 1;
    }
    return this.range(start, end, options);
  }

  toArray(options) {
    // The whole set, ordered from smallest to largest.
    //
    // Parameters
    //   options (optional)
    //     withScores (optional, default false)
    //       bool

    if (!this.length) {
      return [];
    }

    if (typeof options !== 'object') {
      options = {
        withScores: false,
      };
    }

    let i;
    let array = new Array(this.length);
    let node = this._head.next[0].next;

    if (options.withScores) {
      for (i = 0; node; node = node.next[0].next) {
        array[i] = [node.key, node.value];
        i += 1;
      }
    } else {
      for (i = 0; node; node = node.next[0].next) {
        array[i] = node.key;
        i += 1;
      }
    }

    return array;
  }

  values() {
    // Return values as an array, the smallest value first.

    if (!this.length) {
      return [];
    }

    let i;
    let array = new Array(this.length);
    let node = this._head.next[0].next;

    for (i = 0; node; node = node.next[0].next) {
      array[i] = node.value;
      i += 1;
    }

    return array;
  }

  _first(min) {
    let node = this._tail;

    if (!node || node.value < min) {
      return null;
    }

    node = this._head;
    for (let next = null, i = this._level - 1; i >= 0; i -= 1) {
      while ((next = node.next[i].next) && next.value < min) {
        node = next;
      }
    }

    return node.next[0].next;
  }

  _get(index) {
    // Find and return the node at index.
    // Return null if not found.
    //
    // TODO: optimize when index is less than log(N) from the end
    let i;
    let node = this._head;
    let distance = -1;

    for (i = this._level - 1; i >= 0; i -= 1) {
      while (node.next[i].next && (distance + node.next[i].span) <= index) {
        distance += node.next[i].span;
        node = node.next[i].next;
      }
      if (distance === index) {
        return node;
      }
    }
    return null;
  }

  _insert(key, value) {
    // precondition: does not already have key
    // in unique mode, returns null if the value already exists
    let update = new Array(32);
    let rank = new Array(32);
    let node = this._head;
    let next = null;
    let i;

    for (i = this._level - 1; i >= 0; i -= 1) {
      rank[i] = (i === (this._level - 1) ? 0 : rank[i + 1]);
      // TODO: optimize some more?
      while ((next = node.next[i].next) && next.value <= value) {
        if (next.value === value) {
          if (this._unique) {
            return null;
          }
          if (next.key >= key) {
            break;
          }
        }
        rank[i] += node.next[i].span;
        node = next;
      }
      if (this._unique && node.value === value) {
        return null;
      }
      update[i] = node;
    }

    if (this._unique && node.value === value) {
      return null;
    }

    let level = randomLevel();
    if (level > this._level) {
      // TODO: optimize
      for (i = this._level; i < level; i += 1) {
        rank[i] = 0;
        update[i] = this._head;
        update[i].next[i].span = this.length;
      }
      this._level = level;
    }

    node = new Node(level, key, value);
    for (i = 0; i < level; i += 1) {
      node.next[i] = new Level(update[i].next[i].next, update[i].next[i].span - (rank[0] - rank[i]));
      update[i].next[i].next = node;
      update[i].next[i].span = (rank[0] - rank[i]) + 1;
    }

    for (i = level; i < this._level; i += 1) {
      update[i].next[i].span++;
    }

    node.prev = (update[0] === this._head) ? null : update[0];
    if (node.next[0].next) {
      node.next[0].next.prev = node;
    } else {
      this._tail = node;
    }

    this.length += 1;
    return node;
  }

  _next(value, node) {
    // find node after node when value >= specified value
    //
    let next, i;

    if (!this._tail || this._tail.value < value) {
      return null;
    }

    // search upwards
    for (next = null; (next = node.next[node.next.length - 1].next) && next.value < value; ) {
      node = next;
    }
    if (node.value === value) {
      return node;
    }

    // search downwards
    for (i = node.next.length - 1; i >= 0; i -= 1) {
      while ((next = node.next[i].next) && next.value < value) {
        node = next;
      }

      if (node.value === value) {
        return node;
      }
    }
    return node.next[0].next;
  }

  _remove(key, value) {
    let update = new Array(32);
    let node = this._head;
    let i, next;

    for (i = this._level - 1; i >= 0; i -= 1) {
      while ((next = node.next[i].next) && (next.value < value || (next.value === value && next.key < key))) {
        node = next;
      }
      update[i] = node;
    }

    node = node.next[0].next;

    if (!node || value !== node.value || node.key !== key) {
      return false;
    }

    // delete
    this._removeNode(node, update);
    this.length -= 1;
  }

  _removeNode(node, update) {
    let next = null;
    let i = 0;
    let n = this._level;

    for (; i < n; i += 1) {
      if (update[i].next[i].next === node) {
        update[i].next[i].span += node.next[i].span - 1;
        update[i].next[i].next = node.next[i].next;
      } else {
        update[i].next[i].span -= 1;
      }
    }
    if (next = node.next[0].next) {
      next.prev = node.prev;
    } else {
      this._tail = node.prev;
    }

    while (this._level > 1 && !this._head.next[this._level - 1].next) {
      this._level -= 1;
    }
  }
}

function randomLevel() {
  let level = 1;
  while (Math.random() < P) {
    level += 1;
  }
  return level < 32 ? level : 32;
}

function Level(next, span) {
  this.next = next;
  this.span = span;
}

// value is score, sorted
// key is obj, unique
function Node(level, key, value) {
  this.key = key;
  this.value = value;
  this.next = new Array(level);
  this.prev = null;
}

function Entry(key, value) {
  this.key = key;
  this.value = value;
}

module.exports = SortedSet;
