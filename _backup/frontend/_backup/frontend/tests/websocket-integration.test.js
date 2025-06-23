const WebSocket = require('ws');

class WebSocketIntegrationTest {
  constructor() {
    this.wsUrl = process.env.WS_URL || 'ws://localhost:8080/ws';
    this.timeout = 5000;
  }

  async testConnection() {
    console.log('🔌 Testing WebSocket connection...');
    
    return new Promise((resolve) => {
      const ws = new WebSocket(this.wsUrl);
      let connected = false;
      
      const timer = setTimeout(() => {
        if (!connected) {
          console.log('  ❌ Connection timeout');
          ws.close();
          resolve(false);
        }
      }, this.timeout);
      
      ws.on('open', () => {
        console.log('  ✅ Connected successfully');
        connected = true;
        clearTimeout(timer);
        ws.close();
        resolve(true);
      });
      
      ws.on('error', (error) => {
        console.log('  ❌ Connection failed:', error.message);
        clearTimeout(timer);
        resolve(false);
      });
    });
  }

  async testDataFlow() {
    console.log('📊 Testing data flow...');
    
    return new Promise((resolve) => {
      const ws = new WebSocket(this.wsUrl);
      let dataReceived = false;
      
      const timer = setTimeout(() => {
        console.log('  ❌ No data received within timeout');
        ws.close();
        resolve(false);
      }, this.timeout);
      
      ws.on('message', (data) => {
        try {
          const parsed = JSON.parse(data);
          if (parsed && Array.isArray(parsed)) {
            console.log(`  ✅ Received ${parsed.length} flights`);
            dataReceived = true;
            clearTimeout(timer);
            ws.close();
            resolve(true);
          }
        } catch (e) {
          console.log('  ⚠️ Invalid JSON received');
        }
      });
      
      ws.on('error', () => {
        clearTimeout(timer);
        resolve(false);
      });
    });
  }

  async run() {
    console.log('🧪 WebSocket Integration Tests');
    console.log(`URL: ${this.wsUrl}`);
    
    const tests = [
      { name: 'Connection', test: () => this.testConnection() },
      { name: 'Data Flow', test: () => this.testDataFlow() }
    ];
    
    let passed = 0;
    for (const { name, test } of tests) {
      const result = await test();
      if (result) passed++;
    }
    
    const success = passed === tests.length;
    console.log(`\n${success ? '✅' : '❌'} ${passed}/${tests.length} tests passed`);
    return success;
  }
}

if (require.main === module) {
  new WebSocketIntegrationTest().run().then(result => {
    process.exit(result ? 0 : 1);
  });
}

module.exports = WebSocketIntegrationTest;