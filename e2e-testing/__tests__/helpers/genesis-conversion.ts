import * as fs from 'fs';

const jsonFilePath = '../../genesis.json'; // Replace with the path to your JSON file

fs.readFile(jsonFilePath, 'utf8', (err, data) => {
  if (err) {
    console.error('Error reading JSON file:', err);
    return;
  }

  try {
    const jsonObject = JSON.parse(data);
    const jsonString = JSON.stringify(jsonObject);
    const escapedString = jsonString.replace(/\\/g, '\\\\').replace(/"/g, '\\"');

    console.log(`"${escapedString}"`);
  } catch (err) {
    console.error('Error parsing JSON:', err);
  }
});
