var pdf = require("pdf-creator-node");
var fs = require("fs");

function countWordsInArray(array) {
  if (!array) {
    return null;
  }

  const wordCountMap = {};

  array.forEach((word) => {
    if (!wordCountMap[word]) {
      wordCountMap[word] = 1;
    } else {
      wordCountMap[word]++;
    }
  });

  const result = [];

  for (const word in wordCountMap) {
    result.push({ word, count: wordCountMap[word] });
  }

  return result;
}

var html = fs.readFileSync("template.html", "utf-8");
var options = {
  format: "A3",
  orientation: "portrait",
  footer: {
    height: "14mm",
    contents: {
      default: '<span style="color: #444;">{{page}}</span>', // fallback value
    },
  },
};
var texts = require("./test.json").texts;

texts.forEach((text, i) => {
  texts[i].words = countWordsInArray(text.words);
});

var document = {
  html: html,
  data: {
    text: texts,
  },
  path: "./output.pdf",
  type: "",
};
console.log(texts);

pdf
  .create(document, options)
  .then((res) => {
    console.log(res);
  })
  .catch((error) => {
    console.error(error);
  });
