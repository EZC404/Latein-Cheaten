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

var html = fs.readFileSync("./pdfGen/template.html", "utf-8");
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
var texts = require("../test.json").texts;

texts.forEach((text, i) => {
  texts[i].words = countWordsInArray(text.words);
  texts[i].erkannteWords = text.erkannteWords.toString().substring(0, 5);
  texts[i].totalWordCount = text.text.split(" ").length;
});

let sorted = texts
  .slice()
  .sort((a, b) => (a.erkannteWords < b.erkannteWords ? 1 : -1));

var document = {
  html: html,
  data: {
    text: texts,
    sorted: sorted,
  },
  path: "./output.pdf",
  type: "",
};

pdf
  .create(document, options)
  .then((res) => {
    console.log(res);
  })
  .catch((error) => {
    console.error(error);
  });
