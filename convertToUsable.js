const data = require("./convertjson.json");

let total = [];

let cur = {};

// Converts the texts from https://www.gottwein.de/Lat/mart/ausw01.php to a valid json the tabel is converted to a json using https://www.convertjson.com/html-table-to-json.htm

const romanNumeralRegex =
  /^(M{0,3})(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3})$/;

data.forEach((element) => {
  if (romanNumeralRegex.test(element.FIELD2)) {
    total.push(cur);

    cur = {};

    cur["text"] = "";
    cur["name"] = element.FIELD1;
  } else {
    cur["text"] += " " + element.FIELD2;
  }
});

console.log(JSON.stringify(total));
