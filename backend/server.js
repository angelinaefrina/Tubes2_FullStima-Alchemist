const express = require('express');
const cors = require('cors');
const path = require('path');

const app = express();
const PORT = 8080;

app.use(cors());

app.use(express.static(path.join(__dirname, 'public')));

app.get('/recipe.json', (req, res) => {
  res.sendFile(path.join(__dirname, 'public', 'recipe.json'));
});

app.listen(PORT, () => {
  console.log(`Server is running at http://localhost:${PORT}`);
});
