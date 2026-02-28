const http = require('http');

const PORT = process.env.PORT || 8080;
const HOST = '0.0.0.0';

const server = http.createServer((req, res) => {
  res.writeHead(200, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify({
    message: 'Hello World from MCSP!',
    version: '1.0.0',
    namespace: process.env.NAMESPACE || 'learning-workspace',
    timestamp: new Date().toISOString()
  }));
});

server.listen(PORT, HOST, () => {
  console.log(`MCSP Hello World app running on ${HOST}:${PORT}`);
});