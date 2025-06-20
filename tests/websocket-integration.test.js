// WebSocket Integration Test
const WebSocket = require('ws');

const integrationTest = {
  testWebSocketConnection: async () => {
    console.log('TEST: WebSocket connection integration');
    
    return new Promise((resolve) => {
      const ws = new WebSocket('ws://localhost:8080/ws?token=');
      let connected = false;
      
      const timeout = setTimeout(() => {
        if (!connected) {
          console.log('  ✗ Connection timeout');
          ws.close();
          resolve(false);
        }
      }, 5000);
      
      ws.on('open', () => {
        console.log('  ✓ WebSocket connected successfully');
        connected = true;
        clearTimeout(timeout);
        ws.close();
        resolve(true);
      });
      
      ws.on('error', (error) => {
        console.log('  ✗ WebSocket connection failed:', error.message);
        clearTimeout(timeout);
        resolve(false);
      });
      
      ws.on('close', (code, reason) => {
        console.log('  ✓ WebSocket closed:', code, reason.toString());
      });
    });
  },

  runTest: async () => {
    console.log('=== WEBSOCKET INTEGRATION TEST ===');
    const result = await integrationTest.testWebSocketConnection();
    console.log(`Result: ${result ? '✅ PASS' : '❌ FAIL'}`);
    return result;
  }
};

if (require.main === module) {
  integrationTest.runTest().then(result => {
    process.exit(result ? 0 : 1);
  });
}

module.exports = integrationTest;