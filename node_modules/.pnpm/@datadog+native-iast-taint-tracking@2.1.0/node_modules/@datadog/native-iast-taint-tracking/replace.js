const isSpecialRegex = /(\$\$)|(\$&)|(\$`)|(\$')|(\$\d)/

function getReplace (addon) {
  function isSpecialReplacement (replacer) {
    return replacer.indexOf('$') > -1 && !!replacer.match(isSpecialRegex)
  }
  function shouldBePropagated (transactionId, thisArg, replacer) {
    return addon.isTainted(transactionId, thisArg, replacer) && !isSpecialReplacement(replacer)
  }
  if (addon.replace) {
    return addon.replace
  }
  return function replace (transactionId, result, thisArg, matcher, replacer) {
    if (transactionId && typeof thisArg === 'string' && typeof replacer === 'string') {
      if (typeof matcher === 'string') {
        if (shouldBePropagated(transactionId, thisArg, replacer)) {
          const index = thisArg.indexOf(matcher)
          if (index > -1) {
            return addon.replaceStringByString(transactionId, result, thisArg, matcher, replacer, index)
          }
        }
      } else if (matcher instanceof RegExp) {
        if (shouldBePropagated(transactionId, thisArg, replacer)) {
          const replacements = []
          let lastIndex = -1
          for (let match = matcher.exec(thisArg), i = 0; match != null; match = matcher.exec(thisArg), i++) {
            const index = match.index
            if (index !== lastIndex) {
              replacements.push([index, match[0]])
              lastIndex = index
            } else {
              break
            }
          }
          return addon.replaceStringByStringUsingRegex(transactionId, result, thisArg, matcher,
            replacer, replacements)
        }
      }
    }
    return result
  }
}

module.exports = getReplace
