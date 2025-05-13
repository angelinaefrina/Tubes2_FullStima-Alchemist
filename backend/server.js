const express = require('express');
const path = require('path');
const app = express();

const cors = require('cors');
app.use(cors());

// Pastikan ini menunjuk ke folder PUBLIC, bukan langsung ke svgs
app.use('/svgs', express.static(path.join(__dirname, 'public/svgs')));
app.use('/svgs', (req, res, next) => {
    console.log('SVG request:', req.url);
    next();
  });

app.listen(8080, () => {
  console.log('Server running at http://localhost:8080');
});
