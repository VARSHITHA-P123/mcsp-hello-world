const http = require('http');

const PORT = process.env.PORT || 8080;
const HOST = '0.0.0.0';

const NAMESPACE = process.env.NAMESPACE || 'learning-workspace';
const SCENARIO = process.env.SCENARIO || '1 of 3';

const html = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>MCSP Hello World</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: 'Segoe UI', sans-serif;
      background: #0f0f1a;
      color: #fff;
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .container {
      text-align: center;
      padding: 2rem;
    }
    .badge {
      display: inline-block;
      background: #e00;
      color: white;
      font-size: 0.75rem;
      font-weight: 700;
      letter-spacing: 2px;
      padding: 4px 14px;
      border-radius: 20px;
      margin-bottom: 1.5rem;
      text-transform: uppercase;
    }
    h1 {
      font-size: 3rem;
      font-weight: 700;
      margin-bottom: 0.5rem;
      background: linear-gradient(90deg, #fff 0%, #a0a0cc 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }
    .subtitle {
      color: #888;
      font-size: 1rem;
      margin-bottom: 2.5rem;
    }
    .cards {
      display: flex;
      gap: 1rem;
      justify-content: center;
      flex-wrap: wrap;
      margin-bottom: 2.5rem;
    }
    .card {
      background: #1a1a2e;
      border: 1px solid #2a2a4a;
      border-radius: 12px;
      padding: 1.2rem 1.8rem;
      min-width: 160px;
    }
    .card-label {
      font-size: 0.7rem;
      letter-spacing: 1.5px;
      text-transform: uppercase;
      color: #666;
      margin-bottom: 0.4rem;
    }
    .card-value {
      font-size: 1rem;
      font-weight: 600;
      color: #a0a0ff;
    }
    .status {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      background: #0a2a0a;
      border: 1px solid #1a4a1a;
      border-radius: 8px;
      padding: 0.6rem 1.2rem;
      font-size: 0.9rem;
      color: #4caf50;
    }
    .dot {
      width: 8px;
      height: 8px;
      background: #4caf50;
      border-radius: 50%;
      animation: pulse 2s infinite;
    }
    @keyframes pulse {
      0%, 100% { opacity: 1; }
      50% { opacity: 0.3; }
    }
    .operators {
      margin-top: 2.5rem;
      display: flex;
      gap: 0.6rem;
      justify-content: center;
      flex-wrap: wrap;
    }
    .op-tag {
      background: #1a1a2e;
      border: 1px solid #2a2a4a;
      border-radius: 6px;
      padding: 4px 12px;
      font-size: 0.75rem;
      color: #888;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="badge">MultiCloud SaaS Platform</div>
    <h1>Hello World</h1>
    <p class="subtitle">Deployed on Red Hat OpenShift via GitOps</p>
    <div class="cards">
      <div class="card">
        <div class="card-label">Version</div>
        <div class="card-value">1.0.1</div>
      </div>
      <div class="card">
        <div class="card-label">Namespace</div>
        <div class="card-value">${NAMESPACE}</div>
      </div>
      <div class="card">
        <div class="card-label">Platform</div>
        <div class="card-value">OpenShift 4.20</div>
      </div>
      <div class="card">
        <div class="card-label">Scenario</div>
        <div class="card-value">${SCENARIO}</div>
      </div>
    </div>
    <div class="status">
      <div class="dot"></div>
      Application Running
    </div>
    <div class="operators">
      <span class="op-tag">RHACM</span>
      <span class="op-tag">MCE</span>
      <span class="op-tag">GitOps</span>
      <span class="op-tag">Pipelines</span>
      <span class="op-tag">Cert Manager</span>
      <span class="op-tag">External Secrets</span>
    </div>
  </div>
</body>
</html>`;

const server = http.createServer((req, res) => {
  if (req.url === '/health') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({
      message: 'Hello World from MCSP!',
      version: '1.0.1',
      namespace: process.env.NAMESPACE || 'learning-workspace',
      timestamp: new Date().toISOString()
    }));
    return;
  }
  res.writeHead(200, { 'Content-Type': 'text/html' });
  res.end(html);
});

server.listen(PORT, HOST, () => {
  console.log(`MCSP Hello World app running on ${HOST}:${PORT}`);
});
