const magicSplit = /^[a-zà-öø-ÿ]+|[A-ZÀ-ÖØ-ß][a-zà-öø-ÿ]+|[a-zà-öø-ÿ]+|[0-9]+|[A-ZÀ-ÖØ-ß]+(?![a-zà-öø-ÿ])/g;
const spaceSplit = /\S+/g;
function getPartsAndIndexes(string, splitRegex) {
  const result = { parts: [], prefixes: [] };
  const matches = string.matchAll(splitRegex);
  let lastWordEndIndex = 0;
  for (const match of matches) {
    if (typeof match.index !== "number")
      continue;
    const word = match[0];
    result.parts.push(word);
    const prefix = string.slice(lastWordEndIndex, match.index).trim();
    result.prefixes.push(prefix);
    lastWordEndIndex = match.index + word.length;
  }
  const tail = string.slice(lastWordEndIndex).trim();
  if (tail) {
    result.parts.push("");
    result.prefixes.push(tail);
  }
  return result;
}
function splitAndPrefix(string, options) {
  const { keepSpecialCharacters = false, keep, prefix = "" } = options || {};
  const normalString = string.trim().normalize("NFC");
  const hasSpaces = normalString.includes(" ");
  const split = hasSpaces ? spaceSplit : magicSplit;
  const partsAndIndexes = getPartsAndIndexes(normalString, split);
  return partsAndIndexes.parts.map((_part, i) => {
    let foundPrefix = partsAndIndexes.prefixes[i] || "";
    let part = _part;
    if (keepSpecialCharacters === false) {
      if (keep) {
        part = part.normalize("NFD").replace(new RegExp(`[^a-zA-Z\xD8\xDF\xF80-9${keep.join("")}]`, "g"), "");
      }
      if (!keep) {
        part = part.normalize("NFD").replace(/[^a-zA-ZØßø0-9]/g, "");
        foundPrefix = "";
      }
    }
    if (keep) {
      foundPrefix = foundPrefix.replace(new RegExp(`[^${keep.join("")}]`, "g"), "");
    }
    if (i === 0) {
      return foundPrefix + part;
    }
    if (!foundPrefix && !part)
      return "";
    if (!hasSpaces) {
      return (foundPrefix || prefix) + part;
    }
    if (!foundPrefix && prefix.match(/\s/)) {
      return " " + part;
    }
    return (foundPrefix || prefix) + part;
  }).filter(Boolean);
}
function capitaliseWord(string) {
  const match = string.matchAll(magicSplit).next().value;
  const firstLetterIndex = match ? match.index : 0;
  return string.slice(0, firstLetterIndex + 1).toUpperCase() + string.slice(firstLetterIndex + 1).toLowerCase();
}

function camelCase(string, options) {
  return splitAndPrefix(string, options).reduce((result, word, index) => {
    return index === 0 || !(word[0] || "").match(magicSplit) ? result + word.toLowerCase() : result + capitaliseWord(word);
  }, "");
}
function pascalCase(string, options) {
  return splitAndPrefix(string, options).reduce((result, word) => {
    return result + capitaliseWord(word);
  }, "");
}
const upperCamelCase = pascalCase;
function kebabCase(string, options) {
  return splitAndPrefix(string, { ...options, prefix: "-" }).join("").toLowerCase();
}
function snakeCase(string, options) {
  return splitAndPrefix(string, { ...options, prefix: "_" }).join("").toLowerCase();
}
function constantCase(string, options) {
  return splitAndPrefix(string, { ...options, prefix: "_" }).join("").toUpperCase();
}
function trainCase(string, options) {
  return splitAndPrefix(string, { ...options, prefix: "-" }).map((word) => capitaliseWord(word)).join("");
}
function adaCase(string, options) {
  return splitAndPrefix(string, { ...options, prefix: "_" }).map((part) => capitaliseWord(part)).join("");
}
function cobolCase(string, options) {
  return splitAndPrefix(string, { ...options, prefix: "-" }).join("").toUpperCase();
}
function dotNotation(string, options) {
  return splitAndPrefix(string, { ...options, prefix: "." }).join("");
}
function pathCase(string, options = { keepSpecialCharacters: true }) {
  return splitAndPrefix(string, options).reduce((result, word, i) => {
    const prefix = i === 0 || word[0] === "/" ? "" : "/";
    return result + prefix + word;
  }, "");
}
function spaceCase(string, options = { keepSpecialCharacters: true }) {
  return splitAndPrefix(string, { ...options, prefix: " " }).join("");
}
function capitalCase(string, options = { keepSpecialCharacters: true }) {
  return splitAndPrefix(string, { ...options, prefix: " " }).reduce((result, word) => {
    return result + capitaliseWord(word);
  }, "");
}
function lowerCase(string, options = { keepSpecialCharacters: true }) {
  return splitAndPrefix(string, { ...options, prefix: " " }).join("").toLowerCase();
}
function upperCase(string, options = { keepSpecialCharacters: true }) {
  return splitAndPrefix(string, { ...options, prefix: " " }).join("").toUpperCase();
}

export { adaCase, camelCase, capitalCase, cobolCase, constantCase, dotNotation, kebabCase, lowerCase, pascalCase, pathCase, snakeCase, spaceCase, trainCase, upperCamelCase, upperCase };
