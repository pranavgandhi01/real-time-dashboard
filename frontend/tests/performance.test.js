// Performance Tests for Large Datasets
const performanceTests = {
  // Test large dataset handling
  testLargeDataset: () => {
    console.log('TEST: Should handle large flight datasets efficiently');
    
    const sizes = [100, 1000, 5000];
    const results = [];
    
    sizes.forEach(size => {
      const startTime = Date.now();
      
      // Simulate processing large dataset
      const flights = Array.from({ length: size }, (_, i) => ({
        icao24: `TEST${i.toString().padStart(6, '0')}`,
        callsign: `FL${i}`,
        latitude: (Math.random() - 0.5) * 180,
        longitude: (Math.random() - 0.5) * 360,
        on_ground: Math.random() > 0.5
      }));
      
      // Filter valid coordinates (performance test)
      const validFlights = flights.filter(f => 
        f.latitude >= -90 && f.latitude <= 90 && 
        f.longitude >= -180 && f.longitude <= 180
      );
      
      const processingTime = Date.now() - startTime;
      const throughput = Math.round(validFlights.length / processingTime * 1000);
      
      results.push({ size, time: processingTime, throughput });
      console.log(`  ✓ ${size} flights: ${processingTime}ms (${throughput} flights/sec)`);
    });
    
    // Performance should be reasonable (>1000 flights/sec)
    const acceptable = results.every(r => r.throughput > 1000);
    return { acceptable, results };
  },

  // Test WebSocket message processing
  testWebSocketPerformance: () => {
    console.log('TEST: WebSocket message processing performance');
    
    const messageCount = 1000;
    const startTime = Date.now();
    
    // Simulate processing multiple WebSocket messages
    for (let i = 0; i < messageCount; i++) {
      const mockMessage = JSON.stringify([{
        icao24: `MSG${i}`,
        callsign: `TEST${i}`,
        latitude: 40 + Math.random(),
        longitude: -74 + Math.random()
      }]);
      
      // Simulate parsing and validation
      const parsed = JSON.parse(mockMessage);
      const valid = parsed.every(f => 
        f.latitude >= -90 && f.latitude <= 90
      );
    }
    
    const totalTime = Date.now() - startTime;
    const messagesPerSec = Math.round(messageCount / totalTime * 1000);
    
    console.log(`  ✓ ${messageCount} messages: ${totalTime}ms (${messagesPerSec} msg/sec)`);
    
    return messagesPerSec > 5000; // Should process >5000 messages/sec
  },

  runTests: () => {
    console.log('=== PERFORMANCE TESTS ===');
    
    const datasetTest = performanceTests.testLargeDataset();
    const websocketTest = performanceTests.testWebSocketPerformance();
    
    const allPassed = datasetTest.acceptable && websocketTest;
    
    console.log('\n=== PERFORMANCE RESULTS ===');
    console.log(`Dataset Processing: ${datasetTest.acceptable ? 'PASS' : 'FAIL'}`);
    console.log(`WebSocket Processing: ${websocketTest ? 'PASS' : 'FAIL'}`);
    console.log(`Overall: ${allPassed ? '✅ PERFORMANCE OK' : '❌ NEEDS OPTIMIZATION'}`);
    
    return allPassed;
  }
};

module.exports = performanceTests;

if (require.main === module) {
  performanceTests.runTests();
}