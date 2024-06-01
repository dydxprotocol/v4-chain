const { code } = require("./build");

const b = code`if (true) { logTrue(); } else { logFalse(); }`;
console.log(b.toString());
