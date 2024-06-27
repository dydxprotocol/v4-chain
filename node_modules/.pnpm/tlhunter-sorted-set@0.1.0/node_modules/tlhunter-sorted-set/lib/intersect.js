function binaryIntersect(a, b) {
  let lookup = Object.create(null), result = [];
  for (; a; a = a.next[0].next)
    lookup[a.key] = true;
  for (; b; b = b.next[0].next)
    if (lookup[b.key])
      result.push(b.key);
  return result;
}

function ternaryIntersect(a, b, c) {
  let lookup = Object.create(null), result = [];
  for (; a; a = a.next[0].next)
    lookup[a.key] = 0;
  for (; b; b = b.next[0].next)
    if (lookup[b.key] === 0)
      lookup[b.key] = 1;
  for (; c; c = c.next[0].next)
    if (lookup[c.key] === 1)
      result.push(c.key);
  return result;
}

function intersect(nodes) {
  let result, node, lookup, x, i, j, n;
  if (!nodes.length)
    return [];
  for (i = nodes.length - 1; i >= 0; i--) {
    if (!nodes[i].length) // abort
      return [];
    nodes[i] = nodes[i]._head.next[0].next;
  }
  if (nodes.length === 1)
    return nodes[0].toArray({field: 'key'});
  if (nodes.length === 2)
    return binaryIntersect(nodes[0], nodes[1]);
  if (nodes.length === 3)
    return ternaryIntersect(nodes[0], nodes[1], nodes[2]);
    /*return nodes[0].length <= nodes[1].length ?
      binaryIntersect(nodes[0], nodes[1]) :
      binaryIntersect(nodes[1], nodes[0]);*/
  lookup = Object.create(null);
  for (node = nodes.shift(); node; node = node.next[0].next)
    lookup[node.key] = 0;
  for (i = 0, n = nodes.length - 1; i < n; i++) {
    x = 0;
    j = i + 1;
    for (node = nodes[i]; node; node = node.next[0].next) {
      if (lookup[node.key] === i) {
        lookup[node.key] = j;
        x++;
      }
    }
    if (!x) // useful?
      return [];
  }
  result = [];
  for (node = nodes[i]; node; node = node.next[0].next)
    if (lookup[node.key] === i)
      result.push(node.key);
  return result;
}

module.exports = intersect;
