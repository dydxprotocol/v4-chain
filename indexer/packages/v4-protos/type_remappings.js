const fs = require('fs');
const path = require('path');

const directory = path.join(__dirname, 'src/codegen/dydxprotocol');

function replaceBytes(filePath) {
  let data = fs.readFileSync(filePath, 'utf8');

  // Replace 'bytes' with 'Uint8Array'
  data = data.replace(/\bbytes\b/g, 'Uint8Array');

  // Replace 'bytes.fromAmino(value)' with 'value'
  data = data.replace(
    /Uint8Array\.fromAmino\(([^)]+)\)/g,
    'bytesFromBase64(value)',
  );

  // Replace 'bytes.toAmino(v)' with 'v'
  data = data.replace(/Uint8Array\.toAmino\(([^)]+)\)/g, '$1');

  // Replace 'Uint8Array.fromPartial(value)' with 'value'
  data = data.replace(/Uint8Array\.fromPartial\(([^)]+)\)/g, '$1');

  fs.writeFileSync(filePath, data, 'utf8');
}

function processDirectory(dir) {
  fs.readdirSync(dir).forEach((file) => {
    const fullPath = path.join(dir, file);
    if (fs.lstatSync(fullPath).isDirectory()) {
      processDirectory(fullPath);
    } else if (fullPath.endsWith('.ts')) {
      replaceBytes(fullPath);
    }
  });
}

processDirectory(directory);
