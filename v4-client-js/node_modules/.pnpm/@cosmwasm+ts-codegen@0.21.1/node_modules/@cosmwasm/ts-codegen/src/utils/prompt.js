import { filter } from 'fuzzy';
import { prompt as inquirerer } from 'inquirerer';

export const getFuzzySearch = (list) => {
  return (answers, input) => {
    input = input || '';
    return new Promise(function (resolve) {
      setTimeout(function () {
        const fuzzyResult = filter(input, list);
        resolve(
          fuzzyResult.map(function (el) {
            return el.original;
          })
        );
      }, 25);
    });
  };
};

export const getFuzzySearchNames = (nameValueItemList) => {
  const list = nameValueItemList.map(({ name, value }) => name);
  return (answers, input) => {
    input = input || '';
    return new Promise(function (resolve) {
      setTimeout(function () {
        const fuzzyResult = filter(input, list);
        resolve(
          fuzzyResult.map(function (el) {
            return nameValueItemList.find(
              ({ name, value }) => el.original == name
            );
          })
        );
      }, 25);
    });
  };
};
const transform = (questions) => {
  return questions.map((q) => {
    if (q.type === 'fuzzy') {
      const choices = q.choices;
      delete q.choices;
      return {
        ...q,
        type: 'autocomplete',
        source: getFuzzySearch(choices)
      };
    } else if (q.type === 'fuzzy:objects') {
      const choices = q.choices;
      delete q.choices;
      return {
        ...q,
        type: 'autocomplete',
        source: getFuzzySearchNames(choices)
      };
    } else {
      return q;
    }
  });
};

export const prompt = async (questions = [], argv = {}) => {
  questions = transform(questions);
  return await inquirerer(questions, argv);
};
